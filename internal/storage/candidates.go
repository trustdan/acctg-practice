package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/trustdan/acctg-practice/internal/candidate"
)

// CandidateFilter provides criteria for querying candidate questions.
type CandidateFilter struct {
	FamilyID         string
	Status           string
	ValidationStatus string
}

// SaveCandidate persists a candidate question with its parameters, derived postings,
// equation effects, and provenance. If a candidate with the same ID exists, it is updated.
func (d *DB) SaveCandidate(cand candidate.CandidateQuestion) error {
	paramsJSON, err := json.Marshal(cand.Parameters)
	if err != nil {
		return fmt.Errorf("failed to marshal parameters: %w", err)
	}

	conceptsJSON, err := json.Marshal(cand.Concepts)
	if err != nil {
		return fmt.Errorf("failed to marshal concepts: %w", err)
	}

	derivedPostingsJSON, err := json.Marshal(cand.DerivedFixture)
	if err != nil {
		return fmt.Errorf("failed to marshal derived fixture: %w", err)
	}

	derivedEquationJSON, err := json.Marshal(cand.DerivedEquation)
	if err != nil {
		return fmt.Errorf("failed to marshal derived equation: %w", err)
	}

	explanationJSON, err := json.Marshal(cand.Explanation)
	if err != nil {
		return fmt.Errorf("failed to marshal explanation: %w", err)
	}

	teachingJSON, err := json.Marshal(cand.Teaching)
	if err != nil {
		return fmt.Errorf("failed to marshal teaching: %w", err)
	}

	query := `
INSERT INTO candidate_questions (
    id, family_id, rule_version, status, scenario_template,
    parameters_json, concepts_json, derived_postings_json,
    derived_equation_json, explanation_json, source, model,
    prompt_text, target_family, target_concept, created_at,
    validation_status, rejection_reason, teaching_json
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    family_id = excluded.family_id,
    rule_version = excluded.rule_version,
    status = excluded.status,
    scenario_template = excluded.scenario_template,
    parameters_json = excluded.parameters_json,
    concepts_json = excluded.concepts_json,
    derived_postings_json = excluded.derived_postings_json,
    derived_equation_json = excluded.derived_equation_json,
    explanation_json = excluded.explanation_json,
    source = excluded.source,
    model = excluded.model,
    prompt_text = excluded.prompt_text,
    target_family = excluded.target_family,
    target_concept = excluded.target_concept,
    created_at = excluded.created_at,
    validation_status = excluded.validation_status,
    rejection_reason = excluded.rejection_reason,
    teaching_json = excluded.teaching_json;`

	_, err = d.db.Exec(query,
		cand.ID, cand.FamilyID, cand.RuleVersion, string(cand.Status), cand.ScenarioTemplate,
		string(paramsJSON), string(conceptsJSON), string(derivedPostingsJSON),
		string(derivedEquationJSON), string(explanationJSON), cand.Provenance.Source, cand.Provenance.Model,
		cand.Provenance.PromptText, cand.Provenance.TargetFamily, cand.Provenance.TargetConcept, cand.Provenance.GeneratedAt.UTC(),
		string(cand.ValidationStatus), cand.RejectionReason, string(teachingJSON),
	)
	if err != nil {
		return fmt.Errorf("failed to save candidate question: %w", err)
	}

	return nil
}

// GetCandidate retrieves a single candidate question by ID.
func (d *DB) GetCandidate(id string) (*candidate.CandidateQuestion, error) {
	query := `
SELECT id, family_id, rule_version, status, scenario_template,
       parameters_json, concepts_json, derived_postings_json,
       derived_equation_json, explanation_json, source, model,
       prompt_text, target_family, target_concept, created_at,
       validation_status, rejection_reason, teaching_json
FROM candidate_questions WHERE id = ?;`

	row := d.db.QueryRow(query, id)
	return scanCandidate(row)
}

// ListCandidates retrieves candidate questions matching the specified filter criteria.
func (d *DB) ListCandidates(filter CandidateFilter) ([]candidate.CandidateQuestion, error) {
	var conditions []string
	var args []any

	if filter.FamilyID != "" {
		conditions = append(conditions, "family_id = ?")
		args = append(args, filter.FamilyID)
	}
	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}
	if filter.ValidationStatus != "" {
		conditions = append(conditions, "validation_status = ?")
		args = append(args, filter.ValidationStatus)
	}

	query := `
SELECT id, family_id, rule_version, status, scenario_template,
       parameters_json, concepts_json, derived_postings_json,
       derived_equation_json, explanation_json, source, model,
       prompt_text, target_family, target_concept, created_at,
       validation_status, rejection_reason, teaching_json
FROM candidate_questions`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY created_at DESC;"

	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query candidates: %w", err)
	}
	defer rows.Close()

	var candidates []candidate.CandidateQuestion
	for rows.Next() {
		cand, err := scanCandidate(rows)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, *cand)
	}

	return candidates, rows.Err()
}

// DeleteCandidate removes a candidate question by ID.
func (d *DB) DeleteCandidate(id string) error {
	_, err := d.db.Exec("DELETE FROM candidate_questions WHERE id = ?", id)
	return err
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanCandidate(s rowScanner) (*candidate.CandidateQuestion, error) {
	var c candidate.CandidateQuestion
	var statusStr, valStatusStr, teachingJSON string
	var paramsJSON, conceptsJSON, derivedPostingsJSON, derivedEqJSON, explJSON string
	var modelVal, promptVal, targetFamilyVal, targetConceptVal, rejReasonVal sql.NullString

	err := s.Scan(
		&c.ID, &c.FamilyID, &c.RuleVersion, &statusStr, &c.ScenarioTemplate,
		&paramsJSON, &conceptsJSON, &derivedPostingsJSON,
		&derivedEqJSON, &explJSON, &c.Provenance.Source, &modelVal,
		&promptVal, &targetFamilyVal, &targetConceptVal, &c.Provenance.GeneratedAt,
		&valStatusStr, &rejReasonVal, &teachingJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("candidate question not found")
		}
		return nil, fmt.Errorf("failed scanning candidate row: %w", err)
	}

	c.Status = candidate.CandidateStatus(statusStr)
	c.ValidationStatus = candidate.ValidationStatus(valStatusStr)
	if modelVal.Valid {
		c.Provenance.Model = modelVal.String
	}
	if promptVal.Valid {
		c.Provenance.PromptText = promptVal.String
	}
	if targetFamilyVal.Valid {
		c.Provenance.TargetFamily = targetFamilyVal.String
	}
	if targetConceptVal.Valid {
		c.Provenance.TargetConcept = targetConceptVal.String
	}
	if rejReasonVal.Valid {
		c.RejectionReason = rejReasonVal.String
	}

	if err := json.Unmarshal([]byte(paramsJSON), &c.Parameters); err != nil {
		return nil, fmt.Errorf("failed to unmarshal parameters: %w", err)
	}
	if err := json.Unmarshal([]byte(conceptsJSON), &c.Concepts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal concepts: %w", err)
	}
	if err := json.Unmarshal([]byte(derivedPostingsJSON), &c.DerivedFixture); err != nil {
		return nil, fmt.Errorf("failed to unmarshal derived fixture: %w", err)
	}
	if err := json.Unmarshal([]byte(derivedEqJSON), &c.DerivedEquation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal derived equation: %w", err)
	}
	if err := json.Unmarshal([]byte(explJSON), &c.Explanation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal explanation: %w", err)
	}

	if err := json.Unmarshal([]byte(teachingJSON), &c.Teaching); err != nil {
		return nil, fmt.Errorf("failed to unmarshal teaching: %w", err)
	}
	return &c, nil
}
