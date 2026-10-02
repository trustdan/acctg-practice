package storage_test

import (
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/exam"
	"github.com/trustdan/acctg-practice/internal/storage"
)

func TestExamPersistenceAndEvidenceSeparation(t *testing.T) {
	db, err := storage.Open(":memory:")
	if err != nil {
		t.Fatalf("failed opening in-memory db: %v", err)
	}
	defer db.Close()

	// 1. Record a learning attempt in the regular attempts table
	learnSessID := "sess-learn-1"
	_ = db.SaveSession(storage.SessionRecord{
		ID:        learnSessID,
		Mode:      "drill",
		StartedAt: time.Now().UTC(),
	})
	qInst := &domain.QuestionInstance{
		InstanceID:    "inst-learn-1",
		QuestionID:    "q-learn-1",
		Version:       1,
		FamilyID:      "cash_service",
		RuleVersion:   1,
		PromptText:    "Service rendered for cash",
		Parameters:    map[string]int64{"amount_minor_units": 10000},
		RandomSeed:    42,
		Entry:         domain.Entry{},
		StageAnswers:  map[domain.DrillStage]domain.StageAnswer{},
		ScaffoldLevel: domain.ScaffoldFull,
	}
	_ = db.SaveQuestionInstance(qInst, learnSessID)

	learnAttempt := domain.Attempt{
		AttemptID:        "att-learn-1",
		SessionID:        learnSessID,
		InstanceID:       "inst-learn-1",
		QuestionID:       "q-learn-1",
		QuestionVersion:  1,
		Stage:            domain.StageBalancedEntry,
		ConceptID:        "cash_vs_revenue",
		SelectedOptionID: "opt-1",
		IsCorrect:        true,
		Assistance:       domain.AssistanceNone,
		GradingVersion:   1,
		AnsweredAt:       time.Now().UTC(),
	}
	if err := db.RecordAttempt(learnAttempt); err != nil {
		t.Fatalf("failed recording learning attempt: %v", err)
	}

	// 2. Create an exam session and record exam attempts in dedicated tables
	examSessID := "exam-test-100"
	examSess := exam.ExamSessionRecord{
		ID:               examSessID,
		TotalQuestions:   5,
		TimeLimitSeconds: 600,
		ElapsedSeconds:   120,
		Status:           exam.ExamStatusInProgress,
		StartedAt:        time.Now().UTC().Add(-2 * time.Minute),
		Score:            0,
		CorrectCount:     0,
		TotalAttempts:    0,
		Seed:             999,
	}
	if err := db.SaveExamSession(examSess); err != nil {
		t.Fatalf("failed saving exam session: %v", err)
	}

	examAtts := []exam.ExamAttemptRecord{
		{
			ID:               "exam-att-1",
			ExamSessionID:    examSessID,
			InstanceID:       "inst-learn-1",
			QuestionID:       "q-learn-1",
			QuestionIndex:    0,
			Stage:            domain.StageBalancedEntry,
			ConceptID:        "cash_vs_revenue",
			SelectedOptionID: "opt-1",
			CorrectOptionID:  "opt-1",
			IsCorrect:        true,
			GradingVersion:   1,
			AnsweredAt:       time.Now().UTC(),
		},
		{
			ID:               "exam-att-2",
			ExamSessionID:    examSessID,
			InstanceID:       "inst-learn-1",
			QuestionID:       "q-learn-1",
			QuestionIndex:    0,
			Stage:            domain.StageEquationEffect,
			ConceptID:        "accounting_equation_invariants",
			SelectedOptionID: "opt-wrong",
			CorrectOptionID:  "opt-right",
			IsCorrect:        false,
			ErrorTag:         "unbalanced_entry",
			GradingVersion:   1,
			AnsweredAt:       time.Now().UTC(),
		},
	}

	if err := db.RecordExamAttemptsBatch(examAtts); err != nil {
		t.Fatalf("failed recording batch exam attempts: %v", err)
	}

	// 3. Verify CRITICAL INVARIANT: db.GetAllAttempts() must contain ONLY learning attempts
	allLearningAttempts, err := db.GetAllAttempts()
	if err != nil {
		t.Fatalf("failed querying all attempts: %v", err)
	}
	if len(allLearningAttempts) != 1 {
		t.Fatalf("expected exactly 1 learning attempt in attempts table, got %d", len(allLearningAttempts))
	}
	if allLearningAttempts[0].AttemptID != "att-learn-1" {
		t.Errorf("unexpected attempt in learning table: %s", allLearningAttempts[0].AttemptID)
	}

	// 4. Verify exam attempts are retrieved from dedicated table
	retrievedExamAtts, err := db.GetExamAttempts(examSessID)
	if err != nil {
		t.Fatalf("failed retrieving exam attempts: %v", err)
	}
	if len(retrievedExamAtts) != 2 {
		t.Fatalf("expected 2 exam attempts, got %d", len(retrievedExamAtts))
	}
	if retrievedExamAtts[1].ErrorTag != "unbalanced_entry" {
		t.Errorf("expected error tag unbalanced_entry, got %s", retrievedExamAtts[1].ErrorTag)
	}

	// 5. Test interrupted session detection and resumption
	interrupted, err := db.GetLatestInterruptedExamSession()
	if err != nil {
		t.Fatalf("failed checking interrupted exam session: %v", err)
	}
	if interrupted == nil || interrupted.ID != examSessID {
		t.Fatalf("expected to detect interrupted exam session %s, got %+v", examSessID, interrupted)
	}

	// 6. Test session completion
	compTime := time.Now().UTC()
	err = db.CompleteExamSession(examSessID, compTime, 50.0, 1, 2, 180)
	if err != nil {
		t.Fatalf("failed completing exam session: %v", err)
	}

	completedSess, err := db.GetExamSession(examSessID)
	if err != nil {
		t.Fatalf("failed retrieving completed exam session: %v", err)
	}
	if completedSess.Status != exam.ExamStatusCompleted {
		t.Errorf("expected status completed, got %s", completedSess.Status)
	}
	if completedSess.Score != 50.0 {
		t.Errorf("expected score 50.0, got %f", completedSess.Score)
	}
	if completedSess.ElapsedSeconds != 180 {
		t.Errorf("expected elapsed seconds 180, got %d", completedSess.ElapsedSeconds)
	}

	// 7. Verify no interrupted session remains
	interruptedAfter, err := db.GetLatestInterruptedExamSession()
	if err != nil {
		t.Fatalf("failed checking interrupted session: %v", err)
	}
	if interruptedAfter != nil {
		t.Errorf("expected no interrupted session after completion, got %+v", interruptedAfter)
	}

	// 8. Test list exam sessions
	allSessions, err := db.ListExamSessions()
	if err != nil {
		t.Fatalf("failed listing exam sessions: %v", err)
	}
	if len(allSessions) != 1 {
		t.Errorf("expected 1 exam session, got %d", len(allSessions))
	}
}
