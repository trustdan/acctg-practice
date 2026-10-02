package storage

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"time"
)

type Migration struct {
	Version     int
	Description string
	SQL         string
}

// migrations defines the ordered sequence of schema migrations.
var migrations = []Migration{
	{
		Version:     1,
		Description: "Initialize core tables: sessions, question_instances, attempts, and mastery_projections",
		SQL: `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at TIMESTAMP NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    mode TEXT NOT NULL,
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS question_instances (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    question_id TEXT NOT NULL,
    question_version INTEGER NOT NULL,
    family_id TEXT NOT NULL,
    rule_version INTEGER NOT NULL,
    prompt_text TEXT NOT NULL,
    parameters_json TEXT NOT NULL,
    random_seed INTEGER NOT NULL,
    stage_answers_json TEXT NOT NULL,
    entry_json TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_instances_session ON question_instances(session_id);
CREATE INDEX IF NOT EXISTS idx_instances_family ON question_instances(family_id);

CREATE TABLE IF NOT EXISTS attempts (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    instance_id TEXT NOT NULL REFERENCES question_instances(id) ON DELETE CASCADE,
    stage TEXT NOT NULL,
    concept_id TEXT NOT NULL,
    selected_option_id TEXT NOT NULL,
    is_correct INTEGER NOT NULL,
    assistance TEXT NOT NULL,
    error_tag TEXT,
    grading_version INTEGER NOT NULL,
    answered_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_attempts_instance ON attempts(instance_id);
CREATE INDEX IF NOT EXISTS idx_attempts_concept ON attempts(concept_id);
CREATE INDEX IF NOT EXISTS idx_attempts_session ON attempts(session_id);

CREATE TABLE IF NOT EXISTS mastery_projections (
    concept_id TEXT PRIMARY KEY,
    alpha REAL NOT NULL,
    beta REAL NOT NULL,
    last_attempt_at TIMESTAMP,
    updated_at TIMESTAMP NOT NULL
);
`,
	},
	{
		Version:     2,
		Description: "Add reference_used to attempts and scaffold_level to question_instances",
		SQL: `
ALTER TABLE attempts ADD COLUMN reference_used INTEGER NOT NULL DEFAULT 0;
ALTER TABLE question_instances ADD COLUMN scaffold_level INTEGER NOT NULL DEFAULT 0;
`,
	},
	{
		Version:     3,
		Description: "Add candidate_questions table for Stage 12 creative candidate generation",
		SQL: `
CREATE TABLE IF NOT EXISTS candidate_questions (
    id TEXT PRIMARY KEY,
    family_id TEXT NOT NULL,
    rule_version INTEGER NOT NULL,
    status TEXT NOT NULL,
    scenario_template TEXT NOT NULL,
    parameters_json TEXT NOT NULL,
    concepts_json TEXT NOT NULL,
    derived_postings_json TEXT NOT NULL,
    derived_equation_json TEXT NOT NULL,
    explanation_json TEXT NOT NULL,
    source TEXT NOT NULL,
    model TEXT,
    prompt_text TEXT,
    target_family TEXT,
    target_concept TEXT,
    created_at TIMESTAMP NOT NULL,
    validation_status TEXT NOT NULL,
    rejection_reason TEXT
);

CREATE INDEX IF NOT EXISTS idx_candidates_family ON candidate_questions(family_id);
CREATE INDEX IF NOT EXISTS idx_candidates_status ON candidate_questions(status);
CREATE INDEX IF NOT EXISTS idx_candidates_validation ON candidate_questions(validation_status);
`,
	},
	{
		Version:     4,
		Description: "Add approval_events and published_questions tables for Stage 13 review and promotion workflow",
		SQL: `
CREATE TABLE IF NOT EXISTS approval_events (
    id TEXT PRIMARY KEY,
    candidate_id TEXT NOT NULL,
    question_id TEXT NOT NULL,
    action TEXT NOT NULL,
    reviewer TEXT NOT NULL,
    rule_version INTEGER NOT NULL,
    source TEXT NOT NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_approval_candidate ON approval_events(candidate_id);
CREATE INDEX IF NOT EXISTS idx_approval_question ON approval_events(question_id);
CREATE INDEX IF NOT EXISTS idx_approval_action ON approval_events(action);
CREATE INDEX IF NOT EXISTS idx_approval_created ON approval_events(created_at);

CREATE TABLE IF NOT EXISTS published_questions (
    id TEXT PRIMARY KEY,
    version INTEGER NOT NULL,
    family_id TEXT NOT NULL,
    rule_version INTEGER NOT NULL,
    status TEXT NOT NULL,
    scenario_template TEXT NOT NULL,
    parameters_json TEXT NOT NULL,
    concepts_json TEXT NOT NULL,
    fixture_json TEXT NOT NULL,
    reviewer TEXT NOT NULL,
    approved_at TIMESTAMP NOT NULL,
    source TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_published_family ON published_questions(family_id);
CREATE INDEX IF NOT EXISTS idx_published_status ON published_questions(status);
`,
	},
	{
		Version:     5,
		Description: "Add exam_sessions and exam_attempts tables for Stage 16 exam mode evidence separation",
		SQL: `
CREATE TABLE IF NOT EXISTS exam_sessions (
    id TEXT PRIMARY KEY,
    total_questions INTEGER NOT NULL,
    time_limit_seconds INTEGER NOT NULL DEFAULT 0,
    elapsed_seconds INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL,
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    score REAL,
    correct_count INTEGER NOT NULL DEFAULT 0,
    total_attempts INTEGER NOT NULL DEFAULT 0,
    seed INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_exam_sessions_status ON exam_sessions(status);
CREATE INDEX IF NOT EXISTS idx_exam_sessions_started ON exam_sessions(started_at);

CREATE TABLE IF NOT EXISTS exam_attempts (
    id TEXT PRIMARY KEY,
    exam_session_id TEXT NOT NULL REFERENCES exam_sessions(id) ON DELETE CASCADE,
    instance_id TEXT NOT NULL REFERENCES question_instances(id) ON DELETE CASCADE,
    question_id TEXT NOT NULL,
    question_index INTEGER NOT NULL,
    stage TEXT NOT NULL,
    concept_id TEXT NOT NULL,
    selected_option_id TEXT NOT NULL,
    correct_option_id TEXT NOT NULL,
    is_correct INTEGER NOT NULL,
    error_tag TEXT,
    grading_version INTEGER NOT NULL,
    answered_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_exam_attempts_session ON exam_attempts(exam_session_id);
CREATE INDEX IF NOT EXISTS idx_exam_attempts_concept ON exam_attempts(concept_id);
CREATE INDEX IF NOT EXISTS idx_exam_attempts_tag ON exam_attempts(error_tag);
`,
	},
}

// ApplyMigrations applies all pending migrations in order.
// If the database already exists on disk, it creates a timestamped backup before migrating.
// If any migration fails, it rolls back and preserves the prior database state.
func ApplyMigrations(db *sql.DB, dbPath string) error {
	// Create backup if database file already exists and is non-empty
	var backupPath string
	if dbPath != "" && dbPath != ":memory:" {
		info, err := os.Stat(dbPath)
		if err == nil && info.Size() > 0 {
			backupPath = fmt.Sprintf("%s.bak.%s", dbPath, time.Now().UTC().Format("20060102-150405"))
			if err := copyFile(dbPath, backupPath); err != nil {
				return fmt.Errorf("failed to create pre-migration backup %s: %w", backupPath, err)
			}
		}
	}

	// Ensure schema_migrations table exists
	initTableSQL := `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at TIMESTAMP NOT NULL,
    description TEXT NOT NULL
);`
	if _, err := db.Exec(initTableSQL); err != nil {
		return fmt.Errorf("failed to initialize schema_migrations table: %w", err)
	}

	currentVersion, err := getCurrentMigrationVersion(db)
	if err != nil {
		return fmt.Errorf("failed to query current migration version: %w", err)
	}

	for _, m := range migrations {
		if m.Version > currentVersion {
			if err := applySingleMigration(db, m); err != nil {
				errMsg := fmt.Sprintf("migration %d (%s) failed: %v", m.Version, m.Description, err)
				if backupPath != "" {
					errMsg += fmt.Sprintf(". Prior database preserved at backup: %s", backupPath)
				}
				return fmt.Errorf("%s", errMsg)
			}
		}
	}

	return nil
}

func getCurrentMigrationVersion(db *sql.DB) (int, error) {
	row := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations")
	var ver int
	if err := row.Scan(&ver); err != nil {
		return 0, err
	}
	return ver, nil
}

func applySingleMigration(db *sql.DB, m Migration) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not start migration transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(m.SQL); err != nil {
		return fmt.Errorf("failed executing migration SQL: %w", err)
	}

	recordSQL := "INSERT INTO schema_migrations (version, applied_at, description) VALUES (?, ?, ?)"
	if _, err := tx.Exec(recordSQL, m.Version, time.Now().UTC(), m.Description); err != nil {
		return fmt.Errorf("failed recording migration version: %w", err)
	}

	return tx.Commit()
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
