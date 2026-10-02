package storage_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/candidate"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
	"github.com/trustdan/acctg-practice/internal/storage"
)

func setupTestQuestionInstance(t *testing.T) *domain.QuestionInstance {
	catalog, _, err := bank.LoadAccountsFile("../../curriculum/accounts.json")
	if err != nil {
		t.Fatalf("failed to load accounts: %v", err)
	}
	eng := engine.NewEngine(catalog)
	gen := drill.NewGenerator(catalog, eng)
	qBank, err := bank.LoadQuestionBankFile("../../curriculum/seed-questions.json", catalog)
	if err != nil {
		t.Fatalf("failed to load questions: %v", err)
	}

	inst, err := gen.GenerateInstance(qBank.Questions[1], 12345, map[string]int64{"amount_minor_units": 10000})
	if err != nil {
		t.Fatalf("failed to generate instance: %v", err)
	}
	return inst
}

func TestRestartRestoresProgress(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "acctg-storage-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "practice.db")

	// Phase 1: Open DB, create session, instance, and attempts, then close
	db, err := storage.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	sessionID := "sess-restart-1"
	startedAt := time.Now().UTC().Add(-10 * time.Minute)
	err = db.SaveSession(storage.SessionRecord{
		ID:        sessionID,
		Mode:      "progressive_drill",
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("failed to save session: %v", err)
	}

	inst := setupTestQuestionInstance(t)
	if err := db.SaveQuestionInstance(inst, sessionID); err != nil {
		t.Fatalf("failed to save question instance: %v", err)
	}

	att1 := domain.Attempt{
		AttemptID:        "att-001",
		SessionID:        sessionID,
		InstanceID:       inst.InstanceID,
		QuestionID:       inst.QuestionID,
		QuestionVersion:  inst.Version,
		Stage:            domain.StageIdentifyAccount,
		ConceptID:        "cash_classification",
		SelectedOptionID: "opt_cash",
		IsCorrect:        true,
		Assistance:       domain.AssistanceNone,
		GradingVersion:   1,
		AnsweredAt:       time.Now().UTC().Add(-5 * time.Minute),
	}
	att2 := domain.Attempt{
		AttemptID:        "att-002",
		SessionID:        sessionID,
		InstanceID:       inst.InstanceID,
		QuestionID:       inst.QuestionID,
		QuestionVersion:  inst.Version,
		Stage:            domain.StageCounterAccount,
		ConceptID:        "cash_vs_revenue",
		SelectedOptionID: "opt_service_rev",
		IsCorrect:        false,
		Assistance:       domain.AssistanceRetry,
		ErrorTag:         "revenue_recognized_prematurely",
		GradingVersion:   1,
		AnsweredAt:       time.Now().UTC().Add(-3 * time.Minute),
	}

	if err := db.RecordAttempt(att1); err != nil {
		t.Fatalf("failed to record attempt 1: %v", err)
	}
	if err := db.RecordAttempt(att2); err != nil {
		t.Fatalf("failed to record attempt 2: %v", err)
	}

	completedAt := time.Now().UTC()
	if err := db.CompleteSession(sessionID, completedAt); err != nil {
		t.Fatalf("failed to complete session: %v", err)
	}

	// Close database
	if err := db.Close(); err != nil {
		t.Fatalf("failed to close database: %v", err)
	}

	// Phase 2: Restart! Reopen database from disk and verify progress is completely restored
	db2, err := storage.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to reopen database: %v", err)
	}
	defer db2.Close()

	sessRestored, err := db2.GetSession(sessionID)
	if err != nil {
		t.Fatalf("failed to retrieve session after restart: %v", err)
	}
	if sessRestored.ID != sessionID || sessRestored.Mode != "progressive_drill" {
		t.Errorf("unexpected restored session: %+v", sessRestored)
	}
	if sessRestored.CompletedAt == nil {
		t.Fatalf("expected completed_at to be restored")
	}

	attempts, err := db2.GetAttemptsForSession(sessionID)
	if err != nil {
		t.Fatalf("failed to retrieve attempts after restart: %v", err)
	}
	if len(attempts) != 2 {
		t.Fatalf("expected 2 restored attempts, got %d", len(attempts))
	}

	if attempts[0].AttemptID != "att-001" || !attempts[0].IsCorrect || attempts[0].Assistance != domain.AssistanceNone {
		t.Errorf("unexpected attempt[0]: %+v", attempts[0])
	}
	if attempts[1].AttemptID != "att-002" || attempts[1].IsCorrect || attempts[1].Assistance != domain.AssistanceRetry || attempts[1].ErrorTag != "revenue_recognized_prematurely" {
		t.Errorf("unexpected attempt[1]: %+v", attempts[1])
	}

	// Verify concept-specific query
	conceptAttempts, err := db2.GetAttemptsForConcept("cash_vs_revenue")
	if err != nil || len(conceptAttempts) != 1 {
		t.Fatalf("unexpected concept attempts query: %v, count=%d", err, len(conceptAttempts))
	}
}

func TestReplayUsesOriginalQuestionVersionAndOptionOrder(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "acctg-storage-replay-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "practice.db")
	db, err := storage.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	_ = db.SaveSession(storage.SessionRecord{ID: "sess-replay", Mode: "machine", StartedAt: time.Now().UTC()})

	originalInst := setupTestQuestionInstance(t)
	if err := db.SaveQuestionInstance(originalInst, "sess-replay"); err != nil {
		t.Fatalf("failed to save question instance: %v", err)
	}

	// Retrieve question instance for replay
	retrievedInst, err := db.GetQuestionInstance(originalInst.InstanceID)
	if err != nil {
		t.Fatalf("failed to get question instance for replay: %v", err)
	}

	// Assert version numbers
	if retrievedInst.Version != originalInst.Version {
		t.Errorf("version mismatch: got %d, want %d", retrievedInst.Version, originalInst.Version)
	}
	if retrievedInst.RuleVersion != originalInst.RuleVersion {
		t.Errorf("rule version mismatch: got %d, want %d", retrievedInst.RuleVersion, originalInst.RuleVersion)
	}
	if retrievedInst.RandomSeed != originalInst.RandomSeed {
		t.Errorf("seed mismatch: got %d, want %d", retrievedInst.RandomSeed, originalInst.RandomSeed)
	}
	if retrievedInst.PromptText != originalInst.PromptText {
		t.Errorf("prompt mismatch: got %q, want %q", retrievedInst.PromptText, originalInst.PromptText)
	}

	// Assert exact option order for every stage
	for stageKey, origStage := range originalInst.StageAnswers {
		retrievedStage, ok := retrievedInst.StageAnswers[stageKey]
		if !ok {
			t.Fatalf("missing stage %s in retrieved instance", stageKey)
		}
		if len(retrievedStage.Options) != len(origStage.Options) {
			t.Fatalf("option count mismatch for stage %s", stageKey)
		}
		for i := range origStage.Options {
			origOpt := origStage.Options[i]
			retOpt := retrievedStage.Options[i]
			if origOpt.ID != retOpt.ID || origOpt.Label != retOpt.Label || origOpt.Text != retOpt.Text {
				t.Fatalf("option order mismatch at index %d for stage %s: orig=%+v, ret=%+v", i, stageKey, origOpt, retOpt)
			}
		}
	}
}

func TestDuplicateWritesCannotInflateHistory(t *testing.T) {
	db, err := storage.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	defer db.Close()

	sessionID := "sess-idempotent"
	_ = db.SaveSession(storage.SessionRecord{ID: sessionID, Mode: "drill", StartedAt: time.Now().UTC()})
	inst := setupTestQuestionInstance(t)
	_ = db.SaveQuestionInstance(inst, sessionID)

	att := domain.Attempt{
		AttemptID:        "att-fixed-uuid",
		SessionID:        sessionID,
		InstanceID:       inst.InstanceID,
		QuestionID:       inst.QuestionID,
		QuestionVersion:  inst.Version,
		Stage:            domain.StageIdentifyAccount,
		ConceptID:        "cash_classification",
		SelectedOptionID: "opt_cash",
		IsCorrect:        true,
		Assistance:       domain.AssistanceNone,
		GradingVersion:   1,
		AnsweredAt:       time.Now().UTC(),
	}

	// Record attempt 1st time
	if err := db.RecordAttempt(att); err != nil {
		t.Fatalf("first record failed: %v", err)
	}

	// Record the exact same attempt 2nd time (must succeed without creating duplicate)
	if err := db.RecordAttempt(att); err != nil {
		t.Fatalf("second record failed: %v", err)
	}

	// Record inside batch insert with same attempt ID
	if err := db.RecordAttemptsBatch([]domain.Attempt{att}); err != nil {
		t.Fatalf("batch insert duplicate failed: %v", err)
	}

	// Query attempts: MUST be exactly 1
	attempts, err := db.GetAllAttempts()
	if err != nil {
		t.Fatalf("failed querying attempts: %v", err)
	}
	if len(attempts) != 1 {
		t.Fatalf("expected exactly 1 attempt despite duplicate writes, got %d", len(attempts))
	}
}

func TestFailedMigrationPreservesPriorDatabase(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "acctg-storage-migration-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "precious.db")

	// Phase 1: Initialize database at version 1 and populate critical learner data
	db, err := storage.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	sessionID := "sess-precious-1"
	_ = db.SaveSession(storage.SessionRecord{ID: sessionID, Mode: "exam", StartedAt: time.Now().UTC()})
	inst := setupTestQuestionInstance(t)
	_ = db.SaveQuestionInstance(inst, sessionID)
	_ = db.RecordAttempt(domain.Attempt{
		AttemptID:        "att-precious-1",
		SessionID:        sessionID,
		InstanceID:       inst.InstanceID,
		QuestionID:       inst.QuestionID,
		QuestionVersion:  1,
		Stage:            domain.StageIdentifyAccount,
		ConceptID:        "cash_classification",
		SelectedOptionID: "opt_cash",
		IsCorrect:        true,
		Assistance:       domain.AssistanceNone,
		GradingVersion:   1,
		AnsweredAt:       time.Now().UTC(),
	})
	_ = db.Close()

	// Phase 2: Simulate a failed future migration
	rawDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open raw sqlite db: %v", err)
	}

	// Try to execute a broken migration directly through a transaction to test rollback
	tx, err := rawDB.Begin()
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	_, err = tx.Exec("THIS IS INVALID SQL SYNTAX THAT MUST FAIL;")
	if err == nil {
		t.Fatalf("expected syntax error")
	}
	_ = tx.Rollback()
	_ = rawDB.Close()

	// Phase 3: Reopen database with standard storage API and verify prior data is 100% intact
	reopenedDB, err := storage.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to reopen database: %v", err)
	}
	defer reopenedDB.Close()

	sess, err := reopenedDB.GetSession(sessionID)
	if err != nil || sess.ID != sessionID {
		t.Fatalf("failed to retrieve session after simulated migration failure: %v", err)
	}

	attempts, err := reopenedDB.GetAttemptsForSession(sessionID)
	if err != nil || len(attempts) != 1 || attempts[0].AttemptID != "att-precious-1" {
		t.Fatalf("prior historical attempts corrupted or lost after simulated failure: %v, count=%d", err, len(attempts))
	}
}

func TestBackupAndJSONExport(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "acctg-storage-backup-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "orig.db")
	backupPath := filepath.Join(tempDir, "backup.db")

	db, err := storage.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	sessionID := "sess-backup"
	_ = db.SaveSession(storage.SessionRecord{ID: sessionID, Mode: "drill", StartedAt: time.Now().UTC()})
	inst := setupTestQuestionInstance(t)
	_ = db.SaveQuestionInstance(inst, sessionID)

	// Test Backup
	if err := db.Backup(backupPath); err != nil {
		t.Fatalf("backup failed: %v", err)
	}

	// Verify backup file exists and can be opened
	backupDB, err := storage.Open(backupPath)
	if err != nil {
		t.Fatalf("failed to open backup database: %v", err)
	}
	defer backupDB.Close()
	sess, err := backupDB.GetSession(sessionID)
	if err != nil || sess.ID != sessionID {
		t.Fatalf("backup verification failed: %v", err)
	}

	// Test ExportJSON
	var buf bytes.Buffer
	if err := db.ExportJSON(&buf); err != nil {
		t.Fatalf("export JSON failed: %v", err)
	}

	var exported map[string]any
	if err := json.Unmarshal(buf.Bytes(), &exported); err != nil {
		t.Fatalf("exported data is not valid JSON: %v", err)
	}
	if _, ok := exported["sessions"]; !ok {
		t.Fatalf("exported JSON missing sessions key")
	}
	if _, ok := exported["instances"]; !ok {
		t.Fatalf("exported JSON missing instances key")
	}
}

func TestMigration2AndReferenceTracking(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "acctg-mig2-test-*")
	if err != nil {
		t.Fatalf("failed creating temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "mig2.db")

	// 1. Initialize DB at migration 1 schema manually
	rawDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed raw open: %v", err)
	}

	initSQL := `
CREATE TABLE schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at TIMESTAMP NOT NULL,
    description TEXT NOT NULL
);
INSERT INTO schema_migrations (version, applied_at, description) VALUES (1, '2026-01-01 00:00:00', 'v1');

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
`
	if _, err := rawDB.Exec(initSQL); err != nil {
		rawDB.Close()
		t.Fatalf("failed initializing v1 schema: %v", err)
	}
	rawDB.Close()

	// 2. Open via storage.Open, which must migrate from v1 to v2
	db, err := storage.Open(dbPath)
	if err != nil {
		t.Fatalf("storage.Open failed on v1 database: %v", err)
	}
	defer db.Close()

	sessID := "sess-mig2"
	if err := db.SaveSession(storage.SessionRecord{ID: sessID, Mode: "drill", StartedAt: time.Now().UTC()}); err != nil {
		t.Fatalf("failed saving session: %v", err)
	}

	inst := setupTestQuestionInstance(t)
	inst.InstanceID = "inst-mig2"
	inst.ScaffoldLevel = domain.ScaffoldIntermediate

	if err := db.SaveQuestionInstance(inst, sessID); err != nil {
		t.Fatalf("failed saving question instance: %v", err)
	}

	loadedInst, err := db.GetQuestionInstance(inst.InstanceID)
	if err != nil {
		t.Fatalf("failed getting question instance: %v", err)
	}
	if loadedInst.ScaffoldLevel != domain.ScaffoldIntermediate {
		t.Fatalf("expected scaffold level %v, got %v", domain.ScaffoldIntermediate, loadedInst.ScaffoldLevel)
	}

	att := domain.Attempt{
		AttemptID:        "att-mig2-1",
		SessionID:        sessID,
		InstanceID:       inst.InstanceID,
		QuestionID:       inst.QuestionID,
		QuestionVersion:  inst.Version,
		Stage:            domain.StageIdentifyAccount,
		ConceptID:        "cash_vs_revenue",
		SelectedOptionID: "opt-1",
		IsCorrect:        true,
		Assistance:       domain.AssistanceReference,
		ReferenceUsed:    true,
		GradingVersion:   1,
		AnsweredAt:       time.Now().UTC(),
	}

	if err := db.RecordAttempt(att); err != nil {
		t.Fatalf("failed recording attempt with reference_used: %v", err)
	}

	attempts, err := db.GetAttemptsForInstance(inst.InstanceID)
	if err != nil {
		t.Fatalf("failed getting attempts: %v", err)
	}
	if len(attempts) != 1 {
		t.Fatalf("expected 1 attempt, got %d", len(attempts))
	}
	if !attempts[0].ReferenceUsed {
		t.Fatalf("expected ReferenceUsed to be true")
	}
	if attempts[0].Assistance != domain.AssistanceReference {
		t.Fatalf("expected assistance level 'reference', got %s", attempts[0].Assistance)
	}
}

func TestMigration3AndCandidatePersistence(t *testing.T) {
	db, err := storage.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	now := time.Now().UTC().Truncate(time.Second)

	// 1. Save valid candidate
	cand1 := candidate.CandidateQuestion{
		ID:               "cand_cust_adv_001",
		FamilyID:         bank.FamilyCustomerAdvance,
		RuleVersion:      1,
		Status:           candidate.StatusCandidatePendingReview,
		ScenarioTemplate: "A client pays ${amount_dollars} cash today for consulting next month.",
		Parameters: map[string][]int64{
			"amount_minor_units": {5000, 10000, 20000},
		},
		Concepts: []string{"cash_vs_revenue", "unearned_revenue_classification"},
		DerivedFixture: bank.FixtureJSON{
			Postings: []bank.FixturePostingJSON{
				{AccountID: "cash", Side: "debit", AmountParameter: "amount_minor_units"},
				{AccountID: "unearned_revenue", Side: "credit", AmountParameter: "amount_minor_units"},
			},
		},
		DerivedEquation: domain.EquationDelta{
			DeltaAssets:      domain.Money(10000),
			DeltaLiabilities: domain.Money(10000),
			DeltaEquity:      domain.Money(0),
		},
		Explanation: engine.EventExplanation{
			Summary: "Customer advance creates unearned revenue liability.",
		},
		Provenance: candidate.Provenance{
			Source:        "offline_generator",
			GeneratedAt:   now,
			TargetFamily:  bank.FamilyCustomerAdvance,
			TargetConcept: "cash_vs_revenue",
		},
		ValidationStatus: candidate.ValidationValid,
	}

	if err := db.SaveCandidate(cand1); err != nil {
		t.Fatalf("failed to save candidate: %v", err)
	}

	// 2. Retrieve candidate by ID
	loaded, err := db.GetCandidate("cand_cust_adv_001")
	if err != nil {
		t.Fatalf("failed to get candidate: %v", err)
	}

	if loaded.ID != cand1.ID {
		t.Errorf("expected ID %s, got %s", cand1.ID, loaded.ID)
	}
	if loaded.FamilyID != cand1.FamilyID {
		t.Errorf("expected family %s, got %s", cand1.FamilyID, loaded.FamilyID)
	}
	if loaded.Status != candidate.StatusCandidatePendingReview {
		t.Errorf("expected status %s, got %s", candidate.StatusCandidatePendingReview, loaded.Status)
	}
	if loaded.ValidationStatus != candidate.ValidationValid {
		t.Errorf("expected validation valid, got %s", loaded.ValidationStatus)
	}
	if len(loaded.DerivedFixture.Postings) != 2 {
		t.Errorf("expected 2 derived postings, got %d", len(loaded.DerivedFixture.Postings))
	}
	if loaded.DerivedEquation.DeltaAssets != domain.Money(10000) {
		t.Errorf("expected DeltaAssets 10000, got %d", loaded.DerivedEquation.DeltaAssets)
	}
	if loaded.Provenance.Source != "offline_generator" {
		t.Errorf("expected source offline_generator, got %s", loaded.Provenance.Source)
	}

	// Invariant check: Candidate must NOT be active for practice
	if bank.IsActiveForPractice(string(loaded.Status)) {
		t.Errorf("candidate question must NOT be active for practice")
	}

	// 3. Save failed candidate proposal with rejection reason
	candFailed := candidate.CandidateQuestion{
		ID:               "cand_failed_002",
		FamilyID:         "unknown_family",
		RuleVersion:      1,
		Status:           candidate.StatusCandidateRejected,
		ScenarioTemplate: "Invalid scenario without placeholder.",
		Parameters: map[string][]int64{
			"amount_minor_units": {1000},
		},
		Concepts:         []string{"unknown"},
		ValidationStatus: candidate.ValidationFailed,
		RejectionReason:  "unsupported event family: \"unknown_family\"",
		Provenance: candidate.Provenance{
			Source:      "untrusted_llm",
			GeneratedAt: now,
		},
	}

	if err := db.SaveCandidate(candFailed); err != nil {
		t.Fatalf("failed to save failed proposal: %v", err)
	}

	loadedFailed, err := db.GetCandidate("cand_failed_002")
	if err != nil {
		t.Fatalf("failed to get failed proposal: %v", err)
	}
	if loadedFailed.ValidationStatus != candidate.ValidationFailed {
		t.Errorf("expected ValidationFailed, got %s", loadedFailed.ValidationStatus)
	}
	if loadedFailed.RejectionReason != candFailed.RejectionReason {
		t.Errorf("expected rejection reason %q, got %q", candFailed.RejectionReason, loadedFailed.RejectionReason)
	}
	if bank.IsActiveForPractice(string(loadedFailed.Status)) {
		t.Errorf("failed proposal must NOT be active for practice")
	}

	// 4. ListCandidates with filters
	allCands, err := db.ListCandidates(storage.CandidateFilter{})
	if err != nil {
		t.Fatalf("failed to list candidates: %v", err)
	}
	if len(allCands) != 2 {
		t.Fatalf("expected 2 candidates in DB, got %d", len(allCands))
	}

	validCands, err := db.ListCandidates(storage.CandidateFilter{
		ValidationStatus: string(candidate.ValidationValid),
	})
	if err != nil {
		t.Fatalf("failed listing valid candidates: %v", err)
	}
	if len(validCands) != 1 || validCands[0].ID != "cand_cust_adv_001" {
		t.Fatalf("expected 1 valid candidate cand_cust_adv_001, got %v", validCands)
	}

	// 5. Delete candidate
	if err := db.DeleteCandidate("cand_failed_002"); err != nil {
		t.Fatalf("failed to delete candidate: %v", err)
	}
	remaining, err := db.ListCandidates(storage.CandidateFilter{})
	if err != nil {
		t.Fatalf("failed listing after delete: %v", err)
	}
	if len(remaining) != 1 {
		t.Fatalf("expected 1 candidate after delete, got %d", len(remaining))
	}
}
