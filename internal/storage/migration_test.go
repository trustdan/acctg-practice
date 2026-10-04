package storage_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/internal/exam"
	"github.com/trustdan/acctg-practice/internal/storage"
	_ "modernc.org/sqlite"
)

// TestFreshDatabaseAppliesAllMigrations verifies that a fresh database initializes
// all versioned migrations (v1 through v5) with all expected tables and indexes.
func TestFreshDatabaseAppliesAllMigrations(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "acctg-mig-fresh-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "fresh.db")
	db, err := storage.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open fresh database: %v", err)
	}
	defer db.Close()

	// Verify schema_migrations has all 5 migrations recorded
	rawDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open raw db: %v", err)
	}
	defer rawDB.Close()

	var count int
	if err := rawDB.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count); err != nil {
		t.Fatalf("failed to query schema_migrations: %v", err)
	}
	if count != 9 {
		t.Fatalf("expected 9 migrations, got %d", count)
	}

	var maxVer int
	if err := rawDB.QueryRow("SELECT MAX(version) FROM schema_migrations").Scan(&maxVer); err != nil {
		t.Fatalf("failed to query max version: %v", err)
	}
	if maxVer != 9 {
		t.Fatalf("expected max version 9, got %d", maxVer)
	}

	// Verify all expected tables exist
	expectedTables := []string{
		"schema_migrations",
		"sessions",
		"question_instances",
		"attempts",
		"mastery_projections",
		"candidate_questions",
		"approval_events",
		"published_questions",
		"exam_sessions",
		"exam_attempts",
		"arcade_high_scores",
		"saved_explanations",
	}

	for _, tbl := range expectedTables {
		var tblCount int
		err := rawDB.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", tbl).Scan(&tblCount)
		if err != nil {
			t.Fatalf("failed checking table %s: %v", tbl, err)
		}
		if tblCount != 1 {
			t.Fatalf("expected table %s to exist", tbl)
		}
	}
}

// TestUpgradeRegressionPreservesExistingData simulates an existing database from v1
// containing learner history and verifies that migrating to v5 preserves all historical
// data intact while enabling new features and creating a backup.
func TestUpgradeRegressionPreservesExistingData(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "acctg-mig-upgrade-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "legacy_v1.db")

	// Step 1: Create a v1 legacy database
	rawDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open raw sqlite: %v", err)
	}

	v1InitSQL := `
CREATE TABLE schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at TIMESTAMP NOT NULL,
    description TEXT NOT NULL
);

INSERT INTO schema_migrations (version, applied_at, description)
VALUES (1, datetime('now'), 'Initialize core tables');

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    mode TEXT NOT NULL,
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP
);

CREATE TABLE question_instances (
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

CREATE TABLE attempts (
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

CREATE TABLE mastery_projections (
    concept_id TEXT PRIMARY KEY,
    alpha REAL NOT NULL,
    beta REAL NOT NULL,
    last_attempt_at TIMESTAMP,
    updated_at TIMESTAMP NOT NULL
);

-- Seed legacy learner data
INSERT INTO sessions (id, mode, started_at)
VALUES ('sess-legacy-1', 'progressive_drill', '2026-09-15 10:00:00');

INSERT INTO question_instances (id, session_id, question_id, question_version, family_id, rule_version, prompt_text, parameters_json, random_seed, stage_answers_json, entry_json, created_at)
VALUES ('inst-legacy-1', 'sess-legacy-1', 'customer_advance_basic', 1, 'customer_advance', 1, 'Received $5,000 cash for future work', '{"amount_minor_units": 500000}', 42, '{}', '{}', '2026-09-15 10:01:00');

INSERT INTO attempts (id, session_id, instance_id, stage, concept_id, selected_option_id, is_correct, assistance, grading_version, answered_at)
VALUES ('att-legacy-1', 'sess-legacy-1', 'inst-legacy-1', 'identify_account', 'cash_vs_revenue', 'opt-cash', 1, 'independent', 1, '2026-09-15 10:01:30');

INSERT INTO mastery_projections (concept_id, alpha, beta, updated_at)
VALUES ('cash_vs_revenue', 2.5, 1.0, '2026-09-15 10:02:00');
`
	if _, err := rawDB.Exec(v1InitSQL); err != nil {
		rawDB.Close()
		t.Fatalf("failed initializing v1 schema and data: %v", err)
	}
	rawDB.Close()

	// Step 2: Open through storage.Open to trigger migrations v2 -> v5
	db, err := storage.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open and migrate legacy database: %v", err)
	}
	defer db.Close()

	// Step 3: Verify pre-migration backup was created
	files, err := filepath.Glob(filepath.Join(tempDir, "legacy_v1.db.bak.*"))
	if err != nil {
		t.Fatalf("failed globbing backup files: %v", err)
	}
	if len(files) == 0 {
		t.Fatalf("expected pre-migration backup file to exist")
	}

	// Step 4: Verify legacy data was 100% preserved
	legacySession, err := db.GetSession("sess-legacy-1")
	if err != nil {
		t.Fatalf("failed fetching legacy session: %v", err)
	}
	if legacySession.ID != "sess-legacy-1" || legacySession.Mode != "progressive_drill" {
		t.Fatalf("legacy session corrupted: %+v", legacySession)
	}

	legacyInst, err := db.GetQuestionInstance("inst-legacy-1")
	if err != nil {
		t.Fatalf("failed fetching legacy instance: %v", err)
	}
	if legacyInst.QuestionID != "customer_advance_basic" || legacyInst.RandomSeed != 42 {
		t.Fatalf("legacy question instance corrupted: %+v", legacyInst)
	}
	// Verify v2 column scaffold_level defaults to 0 (ScaffoldFull)
	if legacyInst.ScaffoldLevel != 0 {
		t.Fatalf("expected default scaffold_level 0, got %d", legacyInst.ScaffoldLevel)
	}

	legacyAttempts, err := db.GetAttemptsForInstance("inst-legacy-1")
	if err != nil {
		t.Fatalf("failed fetching legacy attempts: %v", err)
	}
	if len(legacyAttempts) != 1 {
		t.Fatalf("expected 1 legacy attempt, got %d", len(legacyAttempts))
	}
	// Verify v2 column reference_used defaults to false
	if legacyAttempts[0].ReferenceUsed {
		t.Fatalf("expected default reference_used=false on legacy attempt")
	}

	// Step 5: Verify new tables from v3, v4, v5, and v6 can now be written and read
	// v5: exam session
	now := time.Now().UTC()
	err = db.SaveExamSession(exam.ExamSessionRecord{
		ID:               "exam-sess-new-1",
		TotalQuestions:   5,
		TimeLimitSeconds: 600,
		Status:           "completed",
		StartedAt:        now.Add(-5 * time.Minute),
		CompletedAt:      &now,
		Score:            1.0,
		CorrectCount:     5,
		TotalAttempts:    5,
		Seed:             999,
	})
	if err != nil {
		t.Fatalf("failed saving exam session to migrated database: %v", err)
	}

	examSess, err := db.GetExamSession("exam-sess-new-1")
	if err != nil {
		t.Fatalf("failed retrieving exam session: %v", err)
	}
	if examSess.ID != "exam-sess-new-1" || examSess.Score != 1.0 {
		t.Fatalf("retrieved exam session mismatch: %+v", examSess)
	}

	// v6: arcade high score
	err = db.SaveArcadeHighScore(storage.ArcadeHighScore{
		Initials:        "DAN",
		Score:           15000,
		BlastedCount:    35,
		SurvivalSeconds: 120,
	})
	if err != nil {
		t.Fatalf("failed saving arcade high score to migrated database: %v", err)
	}

	topScore, err := db.GetTopArcadeHighScore()
	if err != nil {
		t.Fatalf("failed querying top arcade high score: %v", err)
	}
	if topScore == nil || topScore.Initials != "DAN" || topScore.Score != 15000 {
		t.Fatalf("unexpected top score from migrated database: %+v", topScore)
	}
}

// TestMigrationsIdempotent verifies that running migrations repeatedly on an
// already fully-migrated database is a safe, idempotent no-op.
func TestMigrationsIdempotent(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "acctg-mig-idem-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "idempotent.db")

	// First open
	db1, err := storage.Open(dbPath)
	if err != nil {
		t.Fatalf("first open failed: %v", err)
	}
	if err := db1.Close(); err != nil {
		t.Fatalf("failed closing db1: %v", err)
	}

	// Second open
	db2, err := storage.Open(dbPath)
	if err != nil {
		t.Fatalf("second open failed: %v", err)
	}
	defer db2.Close()

	// Third open via raw sql ApplyMigrations
	rawDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("raw db open failed: %v", err)
	}
	defer rawDB.Close()

	if err := storage.ApplyMigrations(rawDB, dbPath); err != nil {
		t.Fatalf("re-applying migrations failed: %v", err)
	}

	// Verify exactly 9 migration rows exist without duplicates
	var count int
	if err := rawDB.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count); err != nil {
		t.Fatalf("failed counting migrations: %v", err)
	}
	if count != 9 {
		t.Fatalf("expected 9 migration records, got %d", count)
	}
}
