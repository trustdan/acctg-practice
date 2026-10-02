package storage_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/candidate"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/storage"
)

func TestMigration4AndApprovalStorage(t *testing.T) {
	db, err := storage.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}
	defer db.Close()

	now := time.Date(2026, 10, 1, 14, 0, 0, 0, time.UTC)
	candID := "cand_test_adv_1"
	qID := "q_test_adv_1"

	// 1. Record an approval event
	event := bank.ApprovalEvent{
		ID:          "appr_1",
		CandidateID: candID,
		QuestionID:  qID,
		Action:      bank.ActionApprove,
		Reviewer:    "DanReviewer",
		RuleVersion: 1,
		Source:      "offline_generator",
		Notes:       "Passed all checks",
		CreatedAt:   now,
	}

	if err := db.RecordApprovalEvent(event); err != nil {
		t.Fatalf("failed to record approval event: %v", err)
	}

	// 2. Query approval events
	events, err := db.GetApprovalEvents(candID)
	if err != nil {
		t.Fatalf("failed querying approval events: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Reviewer != "DanReviewer" || events[0].Action != bank.ActionApprove {
		t.Errorf("unexpected event: %+v", events[0])
	}
}

func TestPublishQuestionAndActiveBankMerge(t *testing.T) {
	db, err := storage.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}
	defer db.Close()

	now := time.Date(2026, 10, 1, 15, 0, 0, 0, time.UTC)
	candID := "cand_advance_pub_1"

	// First, store candidate in candidate_questions
	cand := candidate.CandidateQuestion{
		ID:               candID,
		FamilyID:         bank.FamilyCustomerAdvance,
		RuleVersion:      1,
		Status:           candidate.StatusCandidatePendingReview,
		ScenarioTemplate: "Client Acme paid ${amount_dollars} in advance for services next month.",
		Parameters: map[string][]int64{
			"amount_minor_units": {10000},
		},
		Concepts: []string{"cash_vs_revenue"},
		DerivedFixture: bank.FixtureJSON{
			Postings: []bank.FixturePostingJSON{
				{AccountID: "cash", Side: "debit", AmountParameter: "amount_minor_units"},
				{AccountID: "unearned_revenue", Side: "credit", AmountParameter: "amount_minor_units"},
			},
		},
		ValidationStatus: candidate.ValidationValid,
		Provenance: candidate.Provenance{
			Source:      "generator",
			GeneratedAt: now,
		},
	}
	if err := db.SaveCandidate(cand); err != nil {
		t.Fatalf("failed saving candidate: %v", err)
	}

	// Create QuestionJSON to publish
	reviewer := "LeadInstructor"
	approvedAtStr := now.Format(time.RFC3339)
	pubQ := bank.QuestionJSON{
		ID:               "q_advance_published_1",
		Version:          1,
		FamilyID:         bank.FamilyCustomerAdvance,
		RuleVersion:      1,
		Status:           bank.StatusApprovedActive,
		ScenarioTemplate: cand.ScenarioTemplate,
		Parameters:       cand.Parameters,
		Concepts:         cand.Concepts,
		ExpectedFixture:  cand.DerivedFixture,
		Review: bank.ReviewJSON{
			Reviewer:   &reviewer,
			ApprovedAt: &approvedAtStr,
			Source:     cand.Provenance.Source,
		},
	}

	apprEvent := bank.ApprovalEvent{
		ID:          "appr_pub_1",
		CandidateID: candID,
		QuestionID:  pubQ.ID,
		Action:      bank.ActionApprove,
		Reviewer:    reviewer,
		RuleVersion: 1,
		Source:      cand.Provenance.Source,
		Notes:       "Promoted to active bank",
		CreatedAt:   now,
	}

	// Publish
	if err := db.PublishQuestion(pubQ, candID, apprEvent); err != nil {
		t.Fatalf("failed publishing question: %v", err)
	}

	// Verify candidate status was updated to approved_active in candidate_questions
	updatedCand, err := db.GetCandidate(candID)
	if err != nil {
		t.Fatalf("failed getting candidate: %v", err)
	}
	if updatedCand.Status != bank.StatusApprovedActive {
		t.Errorf("expected candidate status %s, got %s", bank.StatusApprovedActive, updatedCand.Status)
	}

	// Verify question is in published_questions
	fetchedPub, err := db.GetPublishedQuestion(pubQ.ID)
	if err != nil {
		t.Fatalf("failed getting published question: %v", err)
	}
	if fetchedPub.ID != pubQ.ID || fetchedPub.Status != bank.StatusApprovedActive {
		t.Errorf("unexpected published question: %+v", fetchedPub)
	}

	// Test GetActiveBankQuestions merging seed questions + published question
	seed := bank.QuestionJSON{
		ID:               "seed_1",
		Version:          1,
		FamilyID:         bank.FamilyCashService,
		RuleVersion:      1,
		Status:           bank.StatusApprovedActive,
		ScenarioTemplate: "Apex earned ${amount_dollars} today.",
		Parameters: map[string][]int64{
			"amount_minor_units": {5000},
		},
		Concepts: []string{"cash_service"},
	}

	activeBank, err := db.GetActiveBankQuestions([]bank.QuestionJSON{seed})
	if err != nil {
		t.Fatalf("failed getting active bank: %v", err)
	}
	if len(activeBank) != 2 {
		t.Fatalf("expected 2 active questions (1 seed + 1 published), got %d", len(activeBank))
	}
}

func TestRetireQuestionStrictlyExcludesFromPractice(t *testing.T) {
	db, err := storage.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}
	defer db.Close()

	seedActive := bank.QuestionJSON{
		ID:               "seed_good",
		Version:          1,
		FamilyID:         bank.FamilyCashService,
		RuleVersion:      1,
		Status:           bank.StatusApprovedActive,
		ScenarioTemplate: "Good template ${amount_dollars}",
		Parameters:       map[string][]int64{"amount_minor_units": {1000}},
		Concepts:         []string{"cash"},
	}

	seedToRetire := bank.QuestionJSON{
		ID:               "seed_flawed",
		Version:          1,
		FamilyID:         bank.FamilyCustomerAdvance,
		RuleVersion:      1,
		Status:           bank.StatusApprovedActive,
		ScenarioTemplate: "Flawed template ${amount_dollars}",
		Parameters:       map[string][]int64{"amount_minor_units": {2000}},
		Concepts:         []string{"advance"},
	}

	// Before retirement: both active
	active, err := db.GetActiveBankQuestions([]bank.QuestionJSON{seedActive, seedToRetire})
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	if len(active) != 2 {
		t.Fatalf("expected 2 active questions, got %d", len(active))
	}

	// Retire seed_flawed
	event, err := db.RetireQuestion("seed_flawed", "AuditorAlice", "Discovered ambiguous wording in template", time.Now().UTC())
	if err != nil {
		t.Fatalf("failed retiring question: %v", err)
	}
	if event.Action != bank.ActionRetire {
		t.Errorf("expected retire action, got %s", event.Action)
	}

	// Verify IsQuestionRetired
	if !db.IsQuestionRetired("seed_flawed") {
		t.Errorf("expected seed_flawed to be retired")
	}
	if db.IsQuestionRetired("seed_good") {
		t.Errorf("seed_good should not be retired")
	}

	// After retirement: only seedActive is returned!
	activeAfter, err := db.GetActiveBankQuestions([]bank.QuestionJSON{seedActive, seedToRetire})
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	if len(activeAfter) != 1 {
		t.Fatalf("expected 1 active question after retirement, got %d", len(activeAfter))
	}
	if activeAfter[0].ID != "seed_good" {
		t.Errorf("expected seed_good, got %s", activeAfter[0].ID)
	}
}

func TestInvariantCandidatesNeverAutoPromoteBasedOnUseOrScore(t *testing.T) {
	db, err := storage.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}
	defer db.Close()

	// Store an unapproved candidate
	cand := candidate.CandidateQuestion{
		ID:               "cand_unapproved",
		FamilyID:         bank.FamilyCustomerAdvance,
		RuleVersion:      1,
		Status:           candidate.StatusCandidatePendingReview,
		ScenarioTemplate: "Acme received ${amount_dollars} in advance.",
		Parameters:       map[string][]int64{"amount_minor_units": {50000}},
		Concepts:         []string{"cash_vs_revenue"},
		ValidationStatus: candidate.ValidationValid,
		Provenance: candidate.Provenance{
			Source:      "creative_generator",
			GeneratedAt: time.Now().UTC(),
		},
	}
	if err := db.SaveCandidate(cand); err != nil {
		t.Fatalf("failed saving candidate: %v", err)
	}

	// Simulate learner achieving perfect scores across 50 attempts
	sessID := "sess_learner_perfect"
	_ = db.SaveSession(storage.SessionRecord{
		ID:        sessID,
		Mode:      "standard",
		StartedAt: time.Now().UTC(),
	})

	for i := 0; i < 50; i++ {
		instID := fmt.Sprintf("inst_%d", i)
		inst := &domain.QuestionInstance{
			InstanceID:    instID,
			QuestionID:    "seed_cash_service_1",
			Version:       1,
			FamilyID:      bank.FamilyCashService,
			RuleVersion:   1,
			PromptText:    "Test prompt",
			Parameters:    map[string]int64{"amount_minor_units": 1000},
			RandomSeed:    int64(i),
			StageAnswers:  map[domain.DrillStage]domain.StageAnswer{},
			Entry:         domain.Entry{},
			ScaffoldLevel: domain.ScaffoldFull,
		}
		_ = db.SaveQuestionInstance(inst, sessID)

		attempt := domain.Attempt{
			AttemptID:        fmt.Sprintf("att_%d", i),
			SessionID:        sessID,
			InstanceID:       instID,
			QuestionID:       "seed_cash_service_1",
			QuestionVersion:  1,
			Stage:            domain.StageIdentifyAccount,
			ConceptID:        "cash_vs_revenue",
			SelectedOptionID: "cash",
			IsCorrect:        true,
			Assistance:       domain.AssistanceNone,
			GradingVersion:   1,
			AnsweredAt:       time.Now().UTC(),
		}
		_ = db.RecordAttempt(attempt)
	}

	// Verify unapproved candidate was NEVER auto-promoted!
	storedCand, err := db.GetCandidate("cand_unapproved")
	if err != nil {
		t.Fatalf("failed getting candidate: %v", err)
	}
	if storedCand.Status != candidate.StatusCandidatePendingReview {
		t.Errorf("candidate was unexpectedly mutated: status=%s", storedCand.Status)
	}

	// Candidate is NEVER present in active bank
	active, err := db.GetActiveBankQuestions(nil)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	if len(active) != 0 {
		t.Errorf("expected 0 active bank questions, got %d", len(active))
	}

	// Verify no approval events were created
	events, err := db.ListApprovalEvents()
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected 0 approval events, got %d", len(events))
	}
}

func TestHistoricalAttemptsSurviveQuestionRetirementAndUpdates(t *testing.T) {
	db, err := storage.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}
	defer db.Close()

	// 1. Publish a question version 1
	pubQ := bank.QuestionJSON{
		ID:               "q_dynamic_1",
		Version:          1,
		FamilyID:         bank.FamilyCustomerAdvance,
		RuleVersion:      1,
		Status:           bank.StatusApprovedActive,
		ScenarioTemplate: "Client paid ${amount_dollars} in advance.",
		Parameters:       map[string][]int64{"amount_minor_units": {10000}},
		Concepts:         []string{"cash_vs_revenue"},
		Review: bank.ReviewJSON{
			Source: "initial_review",
		},
	}
	appr := bank.ApprovalEvent{
		ID:          "appr_v1",
		CandidateID: "cand_1",
		QuestionID:  pubQ.ID,
		Action:      bank.ActionApprove,
		Reviewer:    "Reviewer1",
		RuleVersion: 1,
		Source:      "initial_review",
		CreatedAt:   time.Now().UTC(),
	}
	if err := db.PublishQuestion(pubQ, "", appr); err != nil {
		t.Fatalf("failed publishing: %v", err)
	}

	// 2. Learner practices q_dynamic_1 v1
	sessID := "sess_hist_1"
	if err := db.SaveSession(storage.SessionRecord{
		ID:        sessID,
		Mode:      "standard",
		StartedAt: time.Now().UTC(),
	}); err != nil {
		t.Fatalf("failed saving session: %v", err)
	}

	instID := "inst_hist_v1"
	inst := &domain.QuestionInstance{
		InstanceID:  instID,
		QuestionID:  pubQ.ID,
		Version:     1,
		FamilyID:    pubQ.FamilyID,
		RuleVersion: pubQ.RuleVersion,
		PromptText:  "Client paid $100.00 in advance.",
		Parameters:  map[string]int64{"amount_minor_units": 10000},
		RandomSeed:  42,
		StageAnswers: map[domain.DrillStage]domain.StageAnswer{
			domain.StageIdentifyAccount: {CorrectOptionID: "cash"},
		},
		Entry:         domain.Entry{},
		ScaffoldLevel: domain.ScaffoldFull,
	}
	if err := db.SaveQuestionInstance(inst, sessID); err != nil {
		t.Fatalf("failed saving question instance: %v", err)
	}

	attID := "att_hist_1"
	att := domain.Attempt{
		AttemptID:        attID,
		SessionID:        sessID,
		InstanceID:       instID,
		QuestionID:       pubQ.ID,
		QuestionVersion:  1,
		Stage:            domain.StageIdentifyAccount,
		ConceptID:        "cash_vs_revenue",
		SelectedOptionID: "cash",
		IsCorrect:        true,
		Assistance:       domain.AssistanceNone,
		GradingVersion:   1,
		AnsweredAt:       time.Now().UTC(),
	}
	if err := db.RecordAttempt(att); err != nil {
		t.Fatalf("failed saving attempt: %v", err)
	}

	// 3. Later, flaw discovered: question is retired!
	retireEvent, err := db.RetireQuestion(pubQ.ID, "AuditorDan", "Discovered flaw in wording", time.Now().UTC())
	if err != nil {
		t.Fatalf("failed retiring: %v", err)
	}
	if retireEvent.Action != bank.ActionRetire {
		t.Errorf("expected retire action")
	}

	// 4. Verify historical attempts and instance are completely unmodified and retrievable!
	sessionAttempts, err := db.GetAttemptsForSession(sessID)
	if err != nil {
		t.Fatalf("failed getting session attempts: %v", err)
	}
	if len(sessionAttempts) != 1 {
		t.Fatalf("expected 1 historical attempt, got %d", len(sessionAttempts))
	}
	if sessionAttempts[0].AttemptID != attID || !sessionAttempts[0].IsCorrect {
		t.Errorf("historical attempt was corrupted: %+v", sessionAttempts[0])
	}

	allAttempts, err := db.GetAllAttempts()
	if err != nil {
		t.Fatalf("failed getting all attempts: %v", err)
	}
	if len(allAttempts) != 1 {
		t.Fatalf("expected 1 total attempt, got %d", len(allAttempts))
	}

	// And verify the question is retired from future active practice
	if !db.IsQuestionRetired(pubQ.ID) {
		t.Errorf("expected q_dynamic_1 to be marked retired")
	}
	activeQuestions, err := db.GetActiveBankQuestions(nil)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	if len(activeQuestions) != 0 {
		t.Errorf("retired question should not be in active bank, got %d questions", len(activeQuestions))
	}
}
