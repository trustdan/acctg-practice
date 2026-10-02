package candidate

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/curriculum"
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/engine"
)

func setupTestCatalogAndEngine(t *testing.T) (*engine.Engine, *bank.AccountsFile) {
	cat, file, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	if err != nil {
		t.Fatalf("failed loading accounts: %v", err)
	}
	eng := engine.NewEngine(cat)
	return eng, file
}

func TestBalancedButWrongAdvanceRevenueCaughtByReview(t *testing.T) {
	eng, _ := setupTestCatalogAndEngine(t)

	// A candidate whose derived entry balances (Dr Cash / Cr Unearned Revenue),
	// but whose scenario text contradicts the family semantics by stating services were fully completed today!
	cand := &CandidateQuestion{
		ID:               "cand_advance_wrong_wording",
		FamilyID:         bank.FamilyCustomerAdvance,
		RuleVersion:      1,
		Status:           StatusCandidatePendingReview,
		ScenarioTemplate: "Acme Web Studio completed a website redesign for client Apex and received ${amount_dollars} cash today upon full delivery and completion.",
		Parameters: map[string][]int64{
			"amount_minor_units": {150000},
		},
		Concepts:         []string{"cash_vs_revenue", "unearned_revenue"},
		ValidationStatus: ValidationValid,
		DerivedFixture: bank.FixtureJSON{
			Postings: []bank.FixturePostingJSON{
				{AccountID: "cash", Side: "debit", AmountParameter: "amount_minor_units"},
				{AccountID: "unearned_revenue", Side: "credit", AmountParameter: "amount_minor_units"},
			},
		},
		Provenance: Provenance{
			Source:      "creative_generator_v1",
			GeneratedAt: time.Now().UTC(),
		},
	}

	// Attempting to approve must fail at the semantic wording gate!
	pub, event, err := ApproveCandidate(cand, "ReviewerDan", "Approved after check", time.Now().UTC())
	if err == nil {
		t.Fatalf("expected balanced-but-wrong candidate to be rejected during review, but got success: pub=%+v, event=%+v", pub, event)
	}

	if !strings.Contains(err.Error(), "semantic wording review failed") {
		t.Errorf("expected error message to cite semantic wording review, got: %v", err)
	}

	// Candidate must remain unpromoted
	if cand.Status == bank.StatusApprovedActive {
		t.Errorf("candidate status must not be approved active after failure")
	}
	_ = eng
}

func TestApproveCandidateSuccess(t *testing.T) {
	now := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	cand := &CandidateQuestion{
		ID:               "cand_advance_valid",
		FamilyID:         bank.FamilyCustomerAdvance,
		RuleVersion:      1,
		Status:           StatusCandidatePendingReview,
		ScenarioTemplate: "Apex Consulting received ${amount_dollars} in cash as an advance retainer from a client for services to be performed next month.",
		Parameters: map[string][]int64{
			"amount_minor_units": {200000},
		},
		Concepts:         []string{"cash_vs_revenue", "unearned_revenue"},
		ValidationStatus: ValidationValid,
		DerivedFixture: bank.FixtureJSON{
			Postings: []bank.FixturePostingJSON{
				{AccountID: "cash", Side: "debit", AmountParameter: "amount_minor_units"},
				{AccountID: "unearned_revenue", Side: "credit", AmountParameter: "amount_minor_units"},
			},
		},
		Provenance: Provenance{
			Source:      "creative_generator_v1",
			GeneratedAt: now,
		},
	}

	pub, event, err := ApproveCandidate(cand, "ProfAccounting", "Vetted wording and timing", now)
	if err != nil {
		t.Fatalf("expected approval to succeed, got error: %v", err)
	}

	if pub == nil || pub.ID != cand.ID {
		t.Fatalf("expected published question ID %s, got %v", cand.ID, pub)
	}
	if pub.Status != bank.StatusApprovedActive {
		t.Errorf("expected published status %s, got %s", bank.StatusApprovedActive, pub.Status)
	}
	if *pub.Review.Reviewer != "ProfAccounting" {
		t.Errorf("expected reviewer ProfAccounting, got %v", *pub.Review.Reviewer)
	}
	if *pub.Review.ApprovedAt != now.Format(time.RFC3339) {
		t.Errorf("expected approved_at %s, got %v", now.Format(time.RFC3339), *pub.Review.ApprovedAt)
	}

	if event == nil || event.Action != bank.ActionApprove {
		t.Fatalf("expected approval event with action %s, got %v", bank.ActionApprove, event)
	}
	if event.Reviewer != "ProfAccounting" {
		t.Errorf("expected event reviewer ProfAccounting, got %s", event.Reviewer)
	}
	if event.RuleVersion != 1 {
		t.Errorf("expected event rule version 1, got %d", event.RuleVersion)
	}
}

func TestApproveCandidateRequiresReviewer(t *testing.T) {
	cand := &CandidateQuestion{
		ID:               "cand_test",
		FamilyID:         bank.FamilyCashService,
		RuleVersion:      1,
		Status:           StatusCandidatePendingReview,
		ScenarioTemplate: "Apex provided services today for ${amount_dollars} cash.",
		Parameters: map[string][]int64{
			"amount_minor_units": {10000},
		},
		Concepts:         []string{"cash_service"},
		ValidationStatus: ValidationValid,
	}

	_, _, err := ApproveCandidate(cand, "   ", "notes", time.Now().UTC())
	if err == nil || !strings.Contains(err.Error(), "reviewer identity is required") {
		t.Fatalf("expected error requiring reviewer, got: %v", err)
	}
}

func TestApproveCandidateRejectsNewUntestedFamily(t *testing.T) {
	cand := &CandidateQuestion{
		ID:               "cand_untested",
		FamilyID:         "crypto_currency_mining",
		RuleVersion:      1,
		Status:           StatusCandidatePendingReview,
		ScenarioTemplate: "Apex mined ${amount_dollars} of crypto today.",
		Parameters: map[string][]int64{
			"amount_minor_units": {10000},
		},
		Concepts:         []string{"crypto"},
		ValidationStatus: ValidationValid,
	}

	_, _, err := ApproveCandidate(cand, "Reviewer1", "notes", time.Now().UTC())
	if err == nil || !strings.Contains(err.Error(), "not a reviewed and tested event family") {
		t.Fatalf("expected error rejecting unreviewed/untested family, got: %v", err)
	}
}

func TestRejectCandidate(t *testing.T) {
	cand := &CandidateQuestion{
		ID:               "cand_flawed",
		FamilyID:         bank.FamilyCashService,
		RuleVersion:      1,
		Status:           StatusCandidatePendingReview,
		ScenarioTemplate: "Apex did something confusing.",
		ValidationStatus: ValidationValid,
		Provenance: Provenance{
			Source: "test_source",
		},
	}

	event, err := RejectCandidate(cand, "ReviewerDan", "Wording is ambiguous and lacks clarity", time.Now().UTC())
	if err != nil {
		t.Fatalf("expected rejection to succeed, got: %v", err)
	}

	if cand.Status != StatusCandidateRejected {
		t.Errorf("expected candidate status to be %s, got %s", StatusCandidateRejected, cand.Status)
	}
	if cand.RejectionReason != "Wording is ambiguous and lacks clarity" {
		t.Errorf("expected rejection reason recorded, got: %s", cand.RejectionReason)
	}
	if event.Action != bank.ActionReject {
		t.Errorf("expected event action %s, got %s", bank.ActionReject, event.Action)
	}
}

func TestRepairCandidate(t *testing.T) {
	eng, file := setupTestCatalogAndEngine(t)
	cat, _, _ := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	_ = file

	// Initially contradictory candidate
	cand := &CandidateQuestion{
		ID:               "cand_needs_repair",
		FamilyID:         bank.FamilyCustomerAdvance,
		RuleVersion:      1,
		Status:           StatusCandidatePendingReview,
		ScenarioTemplate: "Completed website for client and received ${amount_dollars} cash today.",
		Parameters: map[string][]int64{
			"amount_minor_units": {50000},
		},
		Concepts:         []string{"cash_vs_revenue"},
		ValidationStatus: ValidationPending,
		Provenance: Provenance{
			Source: "generator",
		},
	}

	// Repair with clean advance wording
	repairedTemplate := "Client Acme paid ${amount_dollars} cash in advance as a deposit for consulting to be performed next month."
	event, err := RepairCandidate(cand, repairedTemplate, nil, eng, cat, "ProfDan", "Clarified advance timing and future performance", time.Now().UTC())
	if err != nil {
		t.Fatalf("expected repair to succeed, got: %v", err)
	}

	if event.Action != bank.ActionRepair {
		t.Errorf("expected event action %s, got %s", bank.ActionRepair, event.Action)
	}
	if cand.ValidationStatus != ValidationValid {
		t.Errorf("expected candidate validation status %s, got %s (reason: %s)", ValidationValid, cand.ValidationStatus, cand.RejectionReason)
	}
	if cand.ScenarioTemplate != repairedTemplate {
		t.Errorf("expected scenario template updated")
	}

	// Now candidate should pass approval!
	pub, apprEvent, err := ApproveCandidate(cand, "ProfDan", "Approved after repair", time.Now().UTC())
	if err != nil {
		t.Fatalf("expected candidate to approve cleanly after repair, got: %v", err)
	}
	if pub == nil || apprEvent == nil {
		t.Fatalf("expected published question and approval event")
	}
}
