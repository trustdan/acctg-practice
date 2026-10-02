package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/trustdan/acctg-practice/internal/bank"
)

var retireSeq uint64

// RecordApprovalEvent persists an audit record of an approval, rejection, repair, or retirement action.
func (d *DB) RecordApprovalEvent(event bank.ApprovalEvent) error {
	query := `
INSERT INTO approval_events (
    id, candidate_id, question_id, action, reviewer, rule_version, source, notes, created_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

	_, err := d.db.Exec(query,
		event.ID, event.CandidateID, event.QuestionID, string(event.Action),
		event.Reviewer, event.RuleVersion, event.Source, event.Notes, event.CreatedAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("failed to record approval event: %w", err)
	}
	return nil
}

// GetApprovalEvents retrieves approval events matching a candidate or question ID, newest first.
func (d *DB) GetApprovalEvents(targetID string) ([]bank.ApprovalEvent, error) {
	query := `
SELECT id, candidate_id, question_id, action, reviewer, rule_version, source, notes, created_at
FROM approval_events
WHERE candidate_id = ? OR question_id = ?
ORDER BY created_at DESC;`

	rows, err := d.db.Query(query, targetID, targetID)
	if err != nil {
		return nil, fmt.Errorf("failed to query approval events: %w", err)
	}
	defer rows.Close()

	return scanApprovalEvents(rows)
}

// ListApprovalEvents retrieves all approval audit events, newest first.
func (d *DB) ListApprovalEvents() ([]bank.ApprovalEvent, error) {
	query := `
SELECT id, candidate_id, question_id, action, reviewer, rule_version, source, notes, created_at
FROM approval_events
ORDER BY created_at DESC;`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list approval events: %w", err)
	}
	defer rows.Close()

	return scanApprovalEvents(rows)
}

// PublishQuestion atomically records an approved question in published_questions,
// logs the approval audit event in approval_events, and updates candidate status.
func (d *DB) PublishQuestion(q bank.QuestionJSON, candidateID string, event bank.ApprovalEvent) error {
	paramsJSON, err := json.Marshal(q.Parameters)
	if err != nil {
		return fmt.Errorf("failed to marshal parameters: %w", err)
	}

	conceptsJSON, err := json.Marshal(q.Concepts)
	if err != nil {
		return fmt.Errorf("failed to marshal concepts: %w", err)
	}

	teachingJSON, err := json.Marshal(q.Teaching)
	if err != nil {
		return fmt.Errorf("failed to marshal teaching: %w", err)
	}

	fixtureJSON, err := json.Marshal(q.ExpectedFixture)
	if err != nil {
		return fmt.Errorf("failed to marshal fixture: %w", err)
	}

	reviewer := ""
	if q.Review.Reviewer != nil {
		reviewer = *q.Review.Reviewer
	}
	approvedAt := time.Now().UTC()
	if q.Review.ApprovedAt != nil {
		t, err := time.Parse(time.RFC3339, *q.Review.ApprovedAt)
		if err == nil {
			approvedAt = t
		}
	}

	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin publish transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Record approval event
	eventSQL := `
INSERT INTO approval_events (
    id, candidate_id, question_id, action, reviewer, rule_version, source, notes, created_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

	_, err = tx.Exec(eventSQL,
		event.ID, event.CandidateID, event.QuestionID, string(event.Action),
		event.Reviewer, event.RuleVersion, event.Source, event.Notes, event.CreatedAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("failed inserting approval event in transaction: %w", err)
	}

	// 2. Insert or update published question
	pubSQL := `
INSERT INTO published_questions (
    id, version, family_id, rule_version, status, scenario_template,
    parameters_json, concepts_json, fixture_json, reviewer, approved_at, source, created_at, teaching_json
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    version = excluded.version,
    family_id = excluded.family_id,
    rule_version = excluded.rule_version,
    status = excluded.status,
    scenario_template = excluded.scenario_template,
    parameters_json = excluded.parameters_json,
    concepts_json = excluded.concepts_json,
    fixture_json = excluded.fixture_json,
    reviewer = excluded.reviewer,
    approved_at = excluded.approved_at,
    source = excluded.source,
    teaching_json = excluded.teaching_json;`

	_, err = tx.Exec(pubSQL,
		q.ID, q.Version, q.FamilyID, q.RuleVersion, q.Status, q.ScenarioTemplate,
		string(paramsJSON), string(conceptsJSON), string(fixtureJSON),
		reviewer, approvedAt, q.Review.Source, time.Now().UTC(), string(teachingJSON),
	)
	if err != nil {
		return fmt.Errorf("failed inserting published question: %w", err)
	}

	// 3. If candidateID is provided, update candidate question status
	if candidateID != "" {
		candSQL := `UPDATE candidate_questions SET status = ? WHERE id = ?;`
		_, err = tx.Exec(candSQL, bank.StatusApprovedActive, candidateID)
		if err != nil {
			return fmt.Errorf("failed updating candidate status: %w", err)
		}
	}

	return tx.Commit()
}

// GetPublishedQuestion retrieves a single published question by ID.
func (d *DB) GetPublishedQuestion(id string) (*bank.QuestionJSON, error) {
	query := `
SELECT id, version, family_id, rule_version, status, scenario_template,
       parameters_json, concepts_json, fixture_json, reviewer, approved_at, source, teaching_json
FROM published_questions WHERE id = ?;`

	row := d.db.QueryRow(query, id)
	return scanPublishedQuestion(row)
}

// ListPublishedQuestions retrieves published questions, optionally filtering for active practice status.
func (d *DB) ListPublishedQuestions(onlyActive bool) ([]bank.QuestionJSON, error) {
	query := `
SELECT id, version, family_id, rule_version, status, scenario_template,
       parameters_json, concepts_json, fixture_json, reviewer, approved_at, source, teaching_json
FROM published_questions`

	if onlyActive {
		query += fmt.Sprintf(" WHERE status IN ('%s', '%s', '%s')",
			bank.StatusActive, bank.StatusApprovedActive, bank.StatusSeedPendingReview)
	}
	query += " ORDER BY id ASC;"

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query published questions: %w", err)
	}
	defer rows.Close()

	var questions []bank.QuestionJSON
	for rows.Next() {
		q, err := scanPublishedQuestion(rows)
		if err != nil {
			return nil, err
		}
		questions = append(questions, *q)
	}
	return questions, rows.Err()
}

// RetireQuestion marks a question as retired, logs an approval audit event, and excludes it from future practice.
// Historical attempts remain completely preserved and replayable.
func (d *DB) RetireQuestion(questionID string, reviewer string, reason string, clock time.Time) (*bank.ApprovalEvent, error) {
	if clock.IsZero() {
		clock = time.Now().UTC()
	}

	seq := atomic.AddUint64(&retireSeq, 1)
	apprID := fmt.Sprintf("appr_retire_%s_%s_%06d_%d", questionID, clock.UTC().Format("20060102_150405"), (clock.Nanosecond()/1000)%1000000, seq)
	event := bank.ApprovalEvent{
		ID:          apprID,
		CandidateID: questionID,
		QuestionID:  questionID,
		Action:      bank.ActionRetire,
		Reviewer:    reviewer,
		RuleVersion: 1,
		Source:      "manual_review",
		Notes:       reason,
		CreatedAt:   clock.UTC(),
	}

	tx, err := d.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin retirement transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Insert retirement event
	eventSQL := `
INSERT INTO approval_events (
    id, candidate_id, question_id, action, reviewer, rule_version, source, notes, created_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

	_, err = tx.Exec(eventSQL,
		event.ID, event.CandidateID, event.QuestionID, string(event.Action),
		event.Reviewer, event.RuleVersion, event.Source, event.Notes, event.CreatedAt.UTC(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed inserting retirement approval event: %w", err)
	}

	// 2. Update published_questions status if it exists there
	pubSQL := `UPDATE published_questions SET status = ? WHERE id = ?;`
	_, _ = tx.Exec(pubSQL, bank.StatusRetired, questionID)

	// 3. Update candidate_questions status if it exists there
	candSQL := `UPDATE candidate_questions SET status = ? WHERE id = ?;`
	_, _ = tx.Exec(candSQL, bank.StatusRetired, questionID)

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit retirement: %w", err)
	}

	return &event, nil
}

// IsQuestionRetired checks whether a question ID has been retired either in published_questions or approval_events.
func (d *DB) IsQuestionRetired(questionID string) bool {
	// Check approval_events for retirement action
	var count int
	err := d.db.QueryRow(
		"SELECT COUNT(*) FROM approval_events WHERE question_id = ? AND action = ?",
		questionID, string(bank.ActionRetire),
	).Scan(&count)
	if err == nil && count > 0 {
		return true
	}

	// Check published_questions status
	var status string
	err = d.db.QueryRow(
		"SELECT status FROM published_questions WHERE id = ?",
		questionID,
	).Scan(&status)
	if err == nil && status == bank.StatusRetired {
		return true
	}

	return false
}

// GetActiveBankQuestions merges seed questions and published database questions,
// strictly excluding any retired questions and verifying active practice eligibility.
//
// Invariants (AGENT-CONTRACT):
// 1. Candidates never auto-promote based on use or score (candidates in candidate_questions are never loaded).
// 2. Only approved, versioned content enters the active bank.
// 3. Retired questions are strictly excluded from practice.
func (d *DB) GetActiveBankQuestions(seeds []bank.QuestionJSON) ([]bank.QuestionJSON, error) {
	// Query all retired question IDs
	retiredRows, err := d.db.Query("SELECT DISTINCT question_id FROM approval_events WHERE action = ?", string(bank.ActionRetire))
	if err != nil {
		return nil, fmt.Errorf("failed querying retired questions: %w", err)
	}
	defer retiredRows.Close()

	retiredMap := make(map[string]bool)
	for retiredRows.Next() {
		var qID string
		if err := retiredRows.Scan(&qID); err == nil {
			retiredMap[qID] = true
		}
	}

	var active []bank.QuestionJSON

	// 1. Process seed questions
	for _, seed := range seeds {
		if retiredMap[seed.ID] {
			seed.Status = bank.StatusRetired
		}
		if bank.IsActiveForPractice(seed.Status) {
			active = append(active, seed)
		}
	}

	// 2. Process published database questions
	pubQuestions, err := d.ListPublishedQuestions(false)
	if err != nil {
		return nil, fmt.Errorf("failed listing published questions: %w", err)
	}

	for _, pub := range pubQuestions {
		if retiredMap[pub.ID] {
			pub.Status = bank.StatusRetired
		}
		if bank.IsActiveForPractice(pub.Status) {
			active = append(active, pub)
		}
	}

	return active, nil
}

func scanApprovalEvents(rows *sql.Rows) ([]bank.ApprovalEvent, error) {
	var events []bank.ApprovalEvent
	for rows.Next() {
		var e bank.ApprovalEvent
		var actionStr string
		var notesVal sql.NullString

		err := rows.Scan(
			&e.ID, &e.CandidateID, &e.QuestionID, &actionStr,
			&e.Reviewer, &e.RuleVersion, &e.Source, &notesVal, &e.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed scanning approval event: %w", err)
		}
		e.Action = bank.ApprovalAction(actionStr)
		if notesVal.Valid {
			e.Notes = notesVal.String
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func scanPublishedQuestion(s rowScanner) (*bank.QuestionJSON, error) {
	var q bank.QuestionJSON
	var paramsJSON, conceptsJSON, fixtureJSON, teachingJSON string
	var reviewer, source string
	var approvedAt time.Time

	err := s.Scan(
		&q.ID, &q.Version, &q.FamilyID, &q.RuleVersion, &q.Status, &q.ScenarioTemplate,
		&paramsJSON, &conceptsJSON, &fixtureJSON, &reviewer, &approvedAt, &source, &teachingJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("published question not found")
		}
		return nil, fmt.Errorf("failed scanning published question: %w", err)
	}

	if err := json.Unmarshal([]byte(paramsJSON), &q.Parameters); err != nil {
		return nil, fmt.Errorf("failed unmarshaling parameters: %w", err)
	}
	if err := json.Unmarshal([]byte(conceptsJSON), &q.Concepts); err != nil {
		return nil, fmt.Errorf("failed unmarshaling concepts: %w", err)
	}
	if err := json.Unmarshal([]byte(fixtureJSON), &q.ExpectedFixture); err != nil {
		return nil, fmt.Errorf("failed unmarshaling fixture: %w", err)
	}

	approvedAtStr := approvedAt.UTC().Format(time.RFC3339)
	q.Review = bank.ReviewJSON{
		Reviewer:   &reviewer,
		ApprovedAt: &approvedAtStr,
		Source:     source,
	}

	if err := json.Unmarshal([]byte(teachingJSON), &q.Teaching); err != nil {
		return nil, fmt.Errorf("failed unmarshaling teaching: %w", err)
	}
	return &q, nil
}
