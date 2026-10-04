package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"

	"github.com/trustdan/acctg-practice/internal/domain"
)

// SessionRecord represents a practice session in durable storage.
type SessionRecord struct {
	ID          string     `json:"id"`
	Mode        string     `json:"mode"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// DB manages database operations with SQLite.
type DB struct {
	db     *sql.DB
	dbPath string
}

// Open initializes SQLite connection, configures WAL pragmas, and applies migrations.
func Open(dbPath string) (*DB, error) {
	if dbPath != ":memory:" {
		dir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dir, 0700); err != nil {
			return nil, fmt.Errorf("failed to create db directory %s: %w", dir, err)
		}
	}

	dsn := dbPath
	if dbPath != ":memory:" {
		dsn = fmt.Sprintf("%s?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)", dbPath)
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if err := ApplyMigrations(db, dbPath); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return &DB{
		db:     db,
		dbPath: dbPath,
	}, nil
}

// Close safely shuts down the SQLite connection.
func (d *DB) Close() error {
	return d.db.Close()
}

// SaveSession creates or updates a practice session.
func (d *DB) SaveSession(sess SessionRecord) error {
	query := `
INSERT INTO sessions (id, mode, started_at, completed_at)
VALUES (?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    completed_at = excluded.completed_at;`

	_, err := d.db.Exec(query, sess.ID, sess.Mode, sess.StartedAt.UTC(), sess.CompletedAt)
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}
	return nil
}

// CompleteSession marks a session completed with a timestamp.
func (d *DB) CompleteSession(sessionID string, completedAt time.Time) error {
	query := `UPDATE sessions SET completed_at = ? WHERE id = ?;`
	res, err := d.db.Exec(query, completedAt.UTC(), sessionID)
	if err != nil {
		return fmt.Errorf("failed to complete session: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("session %s not found", sessionID)
	}
	return nil
}

// GetSession retrieves a session record by ID.
func (d *DB) GetSession(sessionID string) (*SessionRecord, error) {
	row := d.db.QueryRow("SELECT id, mode, started_at, completed_at FROM sessions WHERE id = ?", sessionID)
	var s SessionRecord
	var completedAt sql.NullTime
	if err := row.Scan(&s.ID, &s.Mode, &s.StartedAt, &completedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session %s not found", sessionID)
		}
		return nil, err
	}
	if completedAt.Valid {
		t := completedAt.Time
		s.CompletedAt = &t
	}
	return &s, nil
}

// SaveQuestionInstance persists an instantiated question with its exact prompt, parameters, entry, and stage answers for replay.
func (d *DB) SaveQuestionInstance(inst *domain.QuestionInstance, sessionID string) error {
	paramsJSON, err := json.Marshal(inst.Parameters)
	if err != nil {
		return fmt.Errorf("failed to marshal parameters: %w", err)
	}

	stageAnswersJSON, err := json.Marshal(inst.StageAnswers)
	if err != nil {
		return fmt.Errorf("failed to marshal stage answers: %w", err)
	}

	entryJSON, err := json.Marshal(inst.Entry)
	if err != nil {
		return fmt.Errorf("failed to marshal entry: %w", err)
	}

	pedagogyJSON, err := json.Marshal(inst.Pedagogy)
	if err != nil {
		return fmt.Errorf("failed to marshal pedagogy: %w", err)
	}
	query := `
INSERT INTO question_instances (
    id, session_id, question_id, question_version, family_id, rule_version,
    prompt_text, parameters_json, random_seed, stage_answers_json, entry_json, created_at, scaffold_level, pedagogy_json
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO NOTHING;`

	_, err = d.db.Exec(query,
		inst.InstanceID, sessionID, inst.QuestionID, inst.Version, inst.FamilyID, inst.RuleVersion,
		inst.PromptText, string(paramsJSON), inst.RandomSeed, string(stageAnswersJSON), string(entryJSON), time.Now().UTC(),
		int(inst.ScaffoldLevel), string(pedagogyJSON),
	)
	if err != nil {
		return fmt.Errorf("failed to save question instance: %w", err)
	}
	return nil
}

// GetQuestionInstance retrieves an instantiated question snapshot for replay.
func (d *DB) GetQuestionInstance(instanceID string) (*domain.QuestionInstance, error) {
	query := `
SELECT id, question_id, question_version, family_id, rule_version,
       prompt_text, parameters_json, random_seed, stage_answers_json, entry_json, scaffold_level, pedagogy_json
FROM question_instances WHERE id = ?;`

	row := d.db.QueryRow(query, instanceID)
	var inst domain.QuestionInstance
	var paramsJSON, stageAnswersJSON, entryJSON, pedagogyJSON string
	var scaffoldLevelInt int

	err := row.Scan(
		&inst.InstanceID, &inst.QuestionID, &inst.Version, &inst.FamilyID, &inst.RuleVersion,
		&inst.PromptText, &paramsJSON, &inst.RandomSeed, &stageAnswersJSON, &entryJSON, &scaffoldLevelInt, &pedagogyJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("question instance %s not found", instanceID)
		}
		return nil, err
	}
	inst.ScaffoldLevel = domain.ScaffoldLevel(scaffoldLevelInt)
	if err := json.Unmarshal([]byte(pedagogyJSON), &inst.Pedagogy); err != nil {
		return nil, fmt.Errorf("invalid pedagogy snapshot: %w", err)
	}

	if err := json.Unmarshal([]byte(paramsJSON), &inst.Parameters); err != nil {
		return nil, fmt.Errorf("failed to unmarshal parameters: %w", err)
	}
	if err := json.Unmarshal([]byte(stageAnswersJSON), &inst.StageAnswers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stage answers: %w", err)
	}
	if err := json.Unmarshal([]byte(entryJSON), &inst.Entry); err != nil {
		return nil, fmt.Errorf("failed to unmarshal entry: %w", err)
	}

	return &inst, nil
}

// RecordAttempt persists an attempt with idempotent write semantics (INSERT OR IGNORE).
// Calling this multiple times with the same attempt ID will NOT create duplicates or inflate history.
func (d *DB) RecordAttempt(att domain.Attempt) error {
	if err := att.Validate(); err != nil {
		return fmt.Errorf("invalid attempt: %w", err)
	}

	isCorrectInt := 0
	if att.IsCorrect {
		isCorrectInt = 1
	}

	refUsedInt := 0
	if att.ReferenceUsed {
		refUsedInt = 1
	}

	query := `
INSERT INTO attempts (
    id, session_id, instance_id, stage, concept_id,
    selected_option_id, is_correct, assistance, error_tag,
    grading_version, answered_at, reference_used, setting_group, pedagogy_version
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO NOTHING;`

	_, err := d.db.Exec(query,
		att.AttemptID, att.SessionID, att.InstanceID, string(att.Stage), att.ConceptID,
		att.SelectedOptionID, isCorrectInt, string(att.Assistance), att.ErrorTag,
		att.GradingVersion, att.AnsweredAt.UTC(), refUsedInt, att.SettingGroup, att.PedagogyVersion,
	)
	if err != nil {
		return fmt.Errorf("failed to record attempt: %w", err)
	}
	return nil
}

// RecordAttemptsBatch saves multiple attempts inside a single atomic transaction.
func (d *DB) RecordAttemptsBatch(attempts []domain.Attempt) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
INSERT INTO attempts (
    id, session_id, instance_id, stage, concept_id,
    selected_option_id, is_correct, assistance, error_tag,
    grading_version, answered_at, reference_used, setting_group, pedagogy_version
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO NOTHING;`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, att := range attempts {
		if err := att.Validate(); err != nil {
			return fmt.Errorf("attempt %s invalid: %w", att.AttemptID, err)
		}
		isCorrectInt := 0
		if att.IsCorrect {
			isCorrectInt = 1
		}
		refUsedInt := 0
		if att.ReferenceUsed {
			refUsedInt = 1
		}
		_, err := stmt.Exec(
			att.AttemptID, att.SessionID, att.InstanceID, string(att.Stage), att.ConceptID,
			att.SelectedOptionID, isCorrectInt, string(att.Assistance), att.ErrorTag,
			att.GradingVersion, att.AnsweredAt.UTC(), refUsedInt, att.SettingGroup, att.PedagogyVersion,
		)
		if err != nil {
			return fmt.Errorf("failed executing batch insert for %s: %w", att.AttemptID, err)
		}
	}

	return tx.Commit()
}

// GetAttemptsForInstance retrieves all recorded attempts for a question instance in chronological order.
func (d *DB) GetAttemptsForInstance(instanceID string) ([]domain.Attempt, error) {
	return d.queryAttempts("SELECT id, session_id, instance_id, stage, concept_id, selected_option_id, is_correct, assistance, error_tag, grading_version, answered_at, reference_used, setting_group, pedagogy_version FROM attempts WHERE instance_id = ? ORDER BY answered_at ASC", instanceID)
}

// GetAttemptsForSession retrieves all recorded attempts for a session in chronological order.
func (d *DB) GetAttemptsForSession(sessionID string) ([]domain.Attempt, error) {
	return d.queryAttempts("SELECT id, session_id, instance_id, stage, concept_id, selected_option_id, is_correct, assistance, error_tag, grading_version, answered_at, reference_used, setting_group, pedagogy_version FROM attempts WHERE session_id = ? ORDER BY answered_at ASC", sessionID)
}

// GetAllAttempts retrieves all recorded attempts across the entire database.
func (d *DB) GetAllAttempts() ([]domain.Attempt, error) {
	return d.queryAttempts("SELECT id, session_id, instance_id, stage, concept_id, selected_option_id, is_correct, assistance, error_tag, grading_version, answered_at, reference_used, setting_group, pedagogy_version FROM attempts ORDER BY answered_at ASC")
}

// GetAttemptsForConcept retrieves attempts tied to a specific learning concept.
func (d *DB) GetAttemptsForConcept(conceptID string) ([]domain.Attempt, error) {
	return d.queryAttempts("SELECT id, session_id, instance_id, stage, concept_id, selected_option_id, is_correct, assistance, error_tag, grading_version, answered_at, reference_used, setting_group, pedagogy_version FROM attempts WHERE concept_id = ? ORDER BY answered_at ASC", conceptID)
}

func (d *DB) queryAttempts(query string, args ...any) ([]domain.Attempt, error) {
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed querying attempts: %w", err)
	}
	defer rows.Close()

	var attempts []domain.Attempt
	for rows.Next() {
		var a domain.Attempt
		var isCorrectInt int
		var refUsedInt int
		var stageStr, assistStr string
		var errorTag sql.NullString

		err := rows.Scan(
			&a.AttemptID, &a.SessionID, &a.InstanceID, &stageStr, &a.ConceptID,
			&a.SelectedOptionID, &isCorrectInt, &assistStr, &errorTag,
			&a.GradingVersion, &a.AnsweredAt, &refUsedInt, &a.SettingGroup, &a.PedagogyVersion,
		)
		if err != nil {
			return nil, fmt.Errorf("failed scanning attempt row: %w", err)
		}

		a.Stage = domain.DrillStage(stageStr)
		a.Assistance = domain.AssistanceLevel(assistStr)
		a.IsCorrect = (isCorrectInt == 1)
		a.ReferenceUsed = (refUsedInt == 1)
		if errorTag.Valid {
			a.ErrorTag = errorTag.String
		}
		attempts = append(attempts, a)
	}

	return attempts, rows.Err()
}

// Backup creates a clean snapshot copy of the current database file to destPath.
func (d *DB) Backup(destPath string) error {
	if d.dbPath == ":memory:" {
		return fmt.Errorf("cannot backup in-memory database to file path")
	}
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create destination dir: %w", err)
	}
	// Use SQLite VACUUM INTO to create an atomic, safe live backup without corruption
	_, err := d.db.Exec("VACUUM INTO ?", destPath)
	if err != nil {
		return fmt.Errorf("failed to backup database via VACUUM INTO: %w", err)
	}
	return nil
}

// ExportJSON exports all sessions, question instances, and attempts to an io.Writer as clean JSON.
func (d *DB) ExportJSON(w io.Writer) error {
	type ExportData struct {
		ExportedAt time.Time                 `json:"exported_at"`
		Sessions   []SessionRecord           `json:"sessions"`
		Instances  []domain.QuestionInstance `json:"instances"`
		Attempts   []domain.Attempt          `json:"attempts"`
	}

	attempts, err := d.GetAllAttempts()
	if err != nil {
		return fmt.Errorf("failed to get attempts for export: %w", err)
	}

	sessRows, err := d.db.Query("SELECT id, mode, started_at, completed_at FROM sessions ORDER BY started_at ASC")
	if err != nil {
		return fmt.Errorf("failed querying sessions for export: %w", err)
	}
	defer sessRows.Close()

	var sessions []SessionRecord
	for sessRows.Next() {
		var s SessionRecord
		var comp sql.NullTime
		if err := sessRows.Scan(&s.ID, &s.Mode, &s.StartedAt, &comp); err != nil {
			return err
		}
		if comp.Valid {
			t := comp.Time
			s.CompletedAt = &t
		}
		sessions = append(sessions, s)
	}

	instRows, err := d.db.Query("SELECT id FROM question_instances ORDER BY created_at ASC")
	if err != nil {
		return fmt.Errorf("failed querying instances for export: %w", err)
	}
	defer instRows.Close()

	var instances []domain.QuestionInstance
	for instRows.Next() {
		var instID string
		if err := instRows.Scan(&instID); err != nil {
			return err
		}
		inst, err := d.GetQuestionInstance(instID)
		if err != nil {
			return err
		}
		instances = append(instances, *inst)
	}

	data := ExportData{
		ExportedAt: time.Now().UTC(),
		Sessions:   sessions,
		Instances:  instances,
		Attempts:   attempts,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}
