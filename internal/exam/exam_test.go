package exam

import (
	"bytes"
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/curriculum"
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
)

func setupTestExam(t *testing.T, seed int64, totalQ int, timeLimit time.Duration) (*ExamRunner, *domain.AccountCatalog, *engine.Engine, *drill.Generator, []bank.QuestionJSON) {
	t.Helper()
	catalog, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	if err != nil {
		t.Fatalf("failed loading accounts: %v", err)
	}
	eng := engine.NewEngine(catalog)
	gen := drill.NewGenerator(catalog, eng)

	bankObj, err := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), catalog)
	if err != nil {
		t.Fatalf("failed loading seed questions: %v", err)
	}

	runner, err := NewExamRunner(ExamConfig{
		SessionID:      "exam-test-1",
		Questions:      bankObj.Questions,
		Catalog:        catalog,
		Engine:         eng,
		Generator:      gen,
		TotalQuestions: totalQ,
		TimeLimit:      timeLimit,
		Seed:           seed,
		Scaffold:       domain.ScaffoldFaded,
		ClockNow:       time.Date(2026, 10, 2, 10, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("failed creating exam runner: %v", err)
	}

	return runner, catalog, eng, gen, bankObj.Questions
}

func TestAnswersWithheldThroughoutSession(t *testing.T) {
	runner, _, _, _, _ := setupTestExam(t, 42, 5, 0)

	// Invariant 1: Requesting hint must be strictly suppressed
	if err := runner.RequestHint(); err != ErrHintsSuppressedInExam {
		t.Errorf("expected ErrHintsSuppressedInExam, got %v", err)
	}

	// Invariant 2: Requesting explanation must be strictly suppressed
	if err := runner.RequestExplanation(); err != ErrExplanationsSuppressedInExam {
		t.Errorf("expected ErrExplanationsSuppressedInExam, got %v", err)
	}

	// Invariant 3: Reference cheatsheet must be strictly suppressed
	if err := runner.RecordReferenceUse(); err != ErrReferenceSuppressedInExam {
		t.Errorf("expected ErrReferenceSuppressedInExam, got %v", err)
	}

	// Invariant 4: Submitting an answer records response without returning feedback,
	// correctness, hints, or explanations
	eq, stageAns, err := runner.CurrentQuestion()
	if err != nil {
		t.Fatalf("unexpected error getting current question: %v", err)
	}

	// Submit an intentionally wrong option (first option that is not correct)
	wrongOptID := ""
	for _, opt := range stageAns.Options {
		if opt.ID != stageAns.CorrectOptionID {
			wrongOptID = opt.ID
			break
		}
	}
	if wrongOptID == "" {
		t.Fatalf("could not find distractor option")
	}

	err = runner.SubmitOption(wrongOptID, time.Now().UTC())
	if err != nil {
		t.Fatalf("SubmitOption failed: %v", err)
	}

	// Check that runner advanced immediately without retry
	eqNext, stageAnsNext, err := runner.CurrentQuestion()
	if err != nil {
		t.Fatalf("failed getting next question stage: %v", err)
	}

	if eqNext.QuestionIndex == eq.QuestionIndex && stageAnsNext.Stage == stageAns.Stage {
		t.Errorf("expected runner to advance to next stage without retry")
	}
}

func TestResultsReproducibilityWithDeterministicSeed(t *testing.T) {
	seed := int64(987654321)

	runner1, _, _, _, _ := setupTestExam(t, seed, 6, 0)
	runner2, _, _, _, _ := setupTestExam(t, seed, 6, 0)

	if len(runner1.Questions) != len(runner2.Questions) {
		t.Fatalf("question count mismatch: %d vs %d", len(runner1.Questions), len(runner2.Questions))
	}

	for i := 0; i < len(runner1.Questions); i++ {
		q1 := runner1.Questions[i]
		q2 := runner2.Questions[i]

		if q1.Instance.QuestionID != q2.Instance.QuestionID {
			t.Errorf("question %d ID mismatch: %s vs %s", i, q1.Instance.QuestionID, q2.Instance.QuestionID)
		}
		if q1.Instance.PromptText != q2.Instance.PromptText {
			t.Errorf("question %d prompt mismatch", i)
		}
		if q1.Instance.Parameters["amount_minor_units"] != q2.Instance.Parameters["amount_minor_units"] {
			t.Errorf("question %d parameter mismatch", i)
		}

		for _, st := range q1.Stages {
			st1 := q1.Instance.StageAnswers[st]
			st2 := q2.Instance.StageAnswers[st]

			if st1.CorrectOptionID != st2.CorrectOptionID {
				t.Errorf("stage %s correct option mismatch", st)
			}
			if len(st1.Options) != len(st2.Options) {
				t.Errorf("stage %s option count mismatch", st)
			}
			for optIdx := range st1.Options {
				if st1.Options[optIdx].ID != st2.Options[optIdx].ID {
					t.Errorf("stage %s option %d ID mismatch", st, optIdx)
				}
			}
		}
	}
}

func TestInterruptedSessionResumesExplicitly(t *testing.T) {
	seed := int64(12345)
	runner, catalog, eng, gen, allQ := setupTestExam(t, seed, 5, 10*time.Minute)

	// Answer 2 stages
	for i := 0; i < 2; i++ {
		_, st, err := runner.CurrentQuestion()
		if err != nil {
			t.Fatalf("error getting stage: %v", err)
		}
		if err := runner.SubmitOption(st.CorrectOptionID, time.Now().UTC()); err != nil {
			t.Fatalf("SubmitOption error: %v", err)
		}
	}
	runner.Tick(120 * time.Second) // 2 minutes elapsed

	// Interrupt session
	record := runner.Interrupt()
	if record.Status != ExamStatusInterrupted {
		t.Errorf("expected ExamStatusInterrupted, got %s", record.Status)
	}
	if record.ElapsedSeconds != 120 {
		t.Errorf("expected 120 elapsed seconds, got %d", record.ElapsedSeconds)
	}

	attempts := runner.Attempts
	if len(attempts) != 2 {
		t.Fatalf("expected 2 attempts recorded, got %d", len(attempts))
	}

	// Now resume runner
	resumedRunner, err := ResumeExamRunner(record, attempts, allQ, catalog, eng, gen)
	if err != nil {
		t.Fatalf("ResumeExamRunner failed: %v", err)
	}

	if resumedRunner.ElapsedSeconds != 120 {
		t.Errorf("resumed elapsed mismatch: %d", resumedRunner.ElapsedSeconds)
	}
	if resumedRunner.TotalQuestions != 5 {
		t.Errorf("resumed total questions mismatch: %d", resumedRunner.TotalQuestions)
	}
	if len(resumedRunner.Attempts) != 2 {
		t.Errorf("resumed attempts mismatch: %d", len(resumedRunner.Attempts))
	}

	// Verify resumed runner is at question 2 (first question had 2 stages completed)
	eqCur, _, err := resumedRunner.CurrentQuestion()
	if err != nil {
		t.Fatalf("resumed CurrentQuestion error: %v", err)
	}
	if eqCur.QuestionIndex != 1 {
		t.Errorf("expected resumed runner to be at question index 1, got %d", eqCur.QuestionIndex)
	}
}

func TestExamTimerTickAndTimeout(t *testing.T) {
	runner, _, _, _, _ := setupTestExam(t, 555, 4, 60*time.Second)

	// Tick 30 seconds
	expired := runner.Tick(30 * time.Second)
	if expired {
		t.Errorf("did not expect timer to expire at 30s")
	}
	if runner.Remaining() != 30*time.Second {
		t.Errorf("expected 30s remaining, got %v", runner.Remaining())
	}

	// Tick another 35 seconds (exceeding 60s limit)
	expired = runner.Tick(35 * time.Second)
	if !expired {
		t.Errorf("expected timer to expire")
	}
	if !runner.IsCompleted() {
		t.Errorf("expected exam to be completed on timeout")
	}
	if !runner.TimedOut {
		t.Errorf("expected TimedOut flag to be true")
	}
}

func TestSummariesBySkillAndErrorPattern(t *testing.T) {
	runner, _, _, _, _ := setupTestExam(t, 777, 4, 0)

	// Answer all questions, deliberately picking 1 distractor with an error tag
	errorTagTriggered := ""
	for !runner.IsCompleted() {
		_, st, err := runner.CurrentQuestion()
		if err != nil {
			break
		}
		if errorTagTriggered == "" {
			// Pick a distractor option that has an error tag
			for _, opt := range st.Options {
				if opt.ID != st.CorrectOptionID && opt.ErrorTag != "" {
					errorTagTriggered = opt.ErrorTag
					_ = runner.SubmitOption(opt.ID, time.Now().UTC())
					break
				}
			}
		} else {
			// Answer correctly
			_ = runner.SubmitOption(st.CorrectOptionID, time.Now().UTC())
		}
	}

	report, err := runner.Finish(time.Now().UTC())
	if err != nil {
		t.Fatalf("Finish failed: %v", err)
	}

	if report.TotalItems == 0 {
		t.Fatalf("expected non-zero total items")
	}

	// Verify skills summary
	if len(report.Skills) == 0 {
		t.Errorf("expected skills breakdown, got empty")
	}
	for _, sk := range report.Skills {
		if sk.ConceptName == "" {
			t.Errorf("concept %s missing display name", sk.ConceptID)
		}
		if sk.Total == 0 {
			t.Errorf("concept %s has 0 total items", sk.ConceptID)
		}
	}

	// Verify error pattern summary
	if errorTagTriggered != "" {
		if len(report.Errors) == 0 {
			t.Errorf("expected error summary for triggered tag %s, got 0", errorTagTriggered)
		}
		found := false
		for _, errItem := range report.Errors {
			if errItem.Tag == errorTagTriggered {
				found = true
				if errItem.Misconception == "" || errItem.Remediation == "" {
					t.Errorf("incomplete error explanation for %s", errorTagTriggered)
				}
			}
		}
		if !found {
			t.Errorf("triggered error tag %s not found in report errors", errorTagTriggered)
		}
	}

	// Verify unmasked question reviews exist after finish
	if len(report.Reviews) != report.TotalItems {
		t.Errorf("expected %d review items, got %d", report.TotalItems, len(report.Reviews))
	}

	// Verify text formatting produces clean output
	text := report.FormatText()
	if len(text) == 0 {
		t.Errorf("FormatText produced empty string")
	}
}
