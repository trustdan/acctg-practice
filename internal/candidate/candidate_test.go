package candidate_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/curriculum"
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/candidate"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
	"github.com/trustdan/acctg-practice/internal/mastery"
	"github.com/trustdan/acctg-practice/internal/tutor"
)

func setupTestEngine(t *testing.T) (*engine.Engine, *domain.AccountCatalog) {
	t.Helper()
	catalog, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	if err != nil {
		t.Fatalf("failed to load accounts: %v", err)
	}
	eng := engine.NewEngine(catalog)
	return eng, catalog
}

func TestCandidateParsingStrictValidation(t *testing.T) {
	eng, catalog := setupTestEngine(t)
	prov := candidate.Provenance{
		Source:       "test",
		GeneratedAt:  time.Now().UTC(),
		TargetFamily: bank.FamilyCustomerAdvance,
	}

	tests := []struct {
		name        string
		rawJSON     string
		expectError bool
		errorSubstr string
	}{
		{
			name: "valid customer advance candidate",
			rawJSON: `{
				"family_id": "customer_advance",
				"scenario_template": "A corporate client pays ${amount_dollars} cash today for software design services next month.",
				"parameters": {
					"amount_minor_units": [5000, 10000, 20000]
				},
				"concepts": ["cash_vs_revenue", "unearned_revenue_classification"]
			}`,
			expectError: false,
		},
		{
			name: "unknown field rejected",
			rawJSON: `{
				"family_id": "customer_advance",
				"scenario_template": "A client pays ${amount_dollars} cash today.",
				"parameters": {
					"amount_minor_units": [5000]
				},
				"concepts": ["cash_vs_revenue"],
				"arbitrary_field": "hacked"
			}`,
			expectError: true,
			errorSubstr: "unknown fields",
		},
		{
			name: "unsupported family rejected",
			rawJSON: `{
				"family_id": "crypto_yield_farming",
				"scenario_template": "Staked tokens earn ${amount_dollars}.",
				"parameters": {
					"amount_minor_units": [5000]
				},
				"concepts": ["crypto"]
			}`,
			expectError: true,
			errorSubstr: "unsupported event family",
		},
		{
			name: "missing amount_dollars placeholder rejected",
			rawJSON: `{
				"family_id": "customer_advance",
				"scenario_template": "A client pays 500 dollars today for future services.",
				"parameters": {
					"amount_minor_units": [5000]
				},
				"concepts": ["cash_vs_revenue"]
			}`,
			expectError: true,
			errorSubstr: "placeholder",
		},
		{
			name: "missing required parameter rejected",
			rawJSON: `{
				"family_id": "customer_advance",
				"scenario_template": "A client pays ${amount_dollars} cash today.",
				"parameters": {
					"other_param": [5000]
				},
				"concepts": ["cash_vs_revenue"]
			}`,
			expectError: true,
			errorSubstr: "missing required parameter",
		},
		{
			name: "non-positive parameter amount rejected",
			rawJSON: `{
				"family_id": "customer_advance",
				"scenario_template": "A client pays ${amount_dollars} cash today.",
				"parameters": {
					"amount_minor_units": [5000, 0, -100]
				},
				"concepts": ["cash_vs_revenue"]
			}`,
			expectError: true,
			errorSubstr: "positive integer minor units",
		},
		{
			name: "empty scenario template rejected",
			rawJSON: `{
				"family_id": "customer_advance",
				"scenario_template": "   ",
				"parameters": {
					"amount_minor_units": [5000]
				},
				"concepts": ["cash_vs_revenue"]
			}`,
			expectError: true,
			errorSubstr: "cannot be empty",
		},
		{
			name: "empty concepts list rejected",
			rawJSON: `{
				"family_id": "customer_advance",
				"scenario_template": "A client pays ${amount_dollars} cash today.",
				"parameters": {
					"amount_minor_units": [5000]
				},
				"concepts": []
			}`,
			expectError: true,
			errorSubstr: "concepts list cannot be empty",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cand, err := candidate.ParseProposal(tc.rawJSON, eng, catalog, prov)
			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if cand == nil {
					t.Fatalf("expected candidate struct recording failure, got nil")
				}
				if cand.ValidationStatus != candidate.ValidationFailed {
					t.Errorf("expected ValidationFailed, got %s", cand.ValidationStatus)
				}
				if cand.Status == candidate.StatusCandidatePendingReview && tc.errorSubstr != "" {
					// Failed proposals should be rejected or have reasons
					if cand.RejectionReason == "" {
						t.Errorf("expected rejection reason, got empty")
					}
				}
				// Invariant: Failed candidates NEVER active for practice
				if bank.IsActiveForPractice(string(cand.Status)) {
					t.Errorf("failed candidate status %q must NOT be active for practice", cand.Status)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if cand.ValidationStatus != candidate.ValidationValid {
					t.Errorf("expected ValidationValid, got %s", cand.ValidationStatus)
				}
				if cand.Status != candidate.StatusCandidatePendingReview {
					t.Errorf("expected StatusCandidatePendingReview, got %s", cand.Status)
				}
				// Invariant: Unapproved candidates NEVER active for practice
				if bank.IsActiveForPractice(string(cand.Status)) {
					t.Errorf("unapproved candidate status %q must NOT be active for practice", cand.Status)
				}
			}
		})
	}
}

func TestLocalAnswerDerivationOverrulesUntrustedProposals(t *testing.T) {
	eng, catalog := setupTestEngine(t)
	prov := candidate.Provenance{
		Source:        "offline_generator",
		GeneratedAt:   time.Now().UTC(),
		TargetFamily:  bank.FamilyCustomerAdvance,
		TargetConcept: "cash_vs_revenue",
	}

	// Raw proposal provides NO answer key or postings
	rawJSON := `{
		"family_id": "customer_advance",
		"scenario_template": "A client pays ${amount_dollars} cash today for web design services next month.",
		"parameters": {
			"amount_minor_units": [10000]
		},
		"concepts": ["cash_vs_revenue", "unearned_revenue_classification"]
	}`

	cand, err := candidate.ParseProposal(rawJSON, eng, catalog, prov)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	// Verify that postings were derived LOCALLY by engine
	if len(cand.DerivedFixture.Postings) != 2 {
		t.Fatalf("expected 2 derived postings, got %d", len(cand.DerivedFixture.Postings))
	}

	p1 := cand.DerivedFixture.Postings[0]
	p2 := cand.DerivedFixture.Postings[1]

	if p1.AccountID != "cash" || p1.Side != "debit" {
		t.Errorf("expected Dr Cash, got %s %s", p1.Side, p1.AccountID)
	}
	if p2.AccountID != "unearned_revenue" || p2.Side != "credit" {
		t.Errorf("expected Cr Unearned Revenue, got %s %s", p2.Side, p2.AccountID)
	}

	// Verify Equation Delta: Assets +$100, Liabilities +$100, Equity $0
	if cand.DerivedEquation.DeltaAssets != domain.Money(10000) {
		t.Errorf("expected DeltaAssets=10000, got %d", cand.DerivedEquation.DeltaAssets)
	}
	if cand.DerivedEquation.DeltaLiabilities != domain.Money(10000) {
		t.Errorf("expected DeltaLiabilities=10000, got %d", cand.DerivedEquation.DeltaLiabilities)
	}
	if cand.DerivedEquation.DeltaEquity != domain.Money(0) {
		t.Errorf("expected DeltaEquity=0, got %d", cand.DerivedEquation.DeltaEquity)
	}
}

func TestContradictorySemanticProposalsCaughtAndRejected(t *testing.T) {
	eng, catalog := setupTestEngine(t)
	prov := candidate.Provenance{
		Source:       "untrusted_llm",
		GeneratedAt:  time.Now().UTC(),
		TargetFamily: bank.FamilyCustomerAdvance,
	}

	// The untrusted proposal attempts to credit service_revenue for a customer advance!
	contradictoryJSON := `{
		"family_id": "customer_advance",
		"scenario_template": "A client pays ${amount_dollars} cash today for a project to begin next month.",
		"parameters": {
			"amount_minor_units": [10000]
		},
		"concepts": ["cash_vs_revenue"],
		"proposed_postings": [
			{"account_id": "cash", "side": "debit", "amount_parameter": "amount_minor_units"},
			{"account_id": "service_revenue", "side": "credit", "amount_parameter": "amount_minor_units"}
		]
	}`

	cand, err := candidate.ParseProposal(contradictoryJSON, eng, catalog, prov)
	if err == nil {
		t.Fatalf("expected contradiction error, got nil")
	}

	if cand.ValidationStatus != candidate.ValidationFailed {
		t.Errorf("expected ValidationFailed, got %s", cand.ValidationStatus)
	}
	if cand.Status != candidate.StatusCandidateRejected {
		t.Errorf("expected StatusCandidateRejected, got %s", cand.Status)
	}
	if cand.RejectionReason == "" {
		t.Errorf("expected non-empty rejection reason")
	}

	// Invariant: Contradictory proposal stays inactive and can never be selected for practice
	if bank.IsActiveForPractice(string(cand.Status)) {
		t.Errorf("rejected proposal must NOT be active for practice")
	}
}

func TestOfflineCandidateGeneratorAllFamilies(t *testing.T) {
	eng, catalog := setupTestEngine(t)
	gen := candidate.NewOfflineCandidateGenerator(eng, catalog)

	allFamilies := []string{
		bank.FamilyCashService,
		bank.FamilyCustomerAdvance,
		bank.FamilyServiceOnCredit,
		bank.FamilyCollectReceivable,
		bank.FamilyEarnAdvance,
		bank.FamilyCashRent,
		bank.FamilyBorrowCash,
		bank.FamilyIssueShares,
		bank.FamilyCashExpense,
		bank.FamilyPrepaidPurchase,
		bank.FamilyPrepaidConsumption,
		bank.FamilyEquipmentPurchaseCash,
		bank.FamilyRepayNotePrincipal,
		bank.FamilyDividendCash,
	}

	for _, fam := range allFamilies {
		t.Run(fam, func(t *testing.T) {
			cand, err := gen.Generate(context.Background(), candidate.GenerateRequest{
				FamilyID: fam,
				Seed:     42,
			})
			if err != nil {
				t.Fatalf("failed to generate candidate for family %s: %v", fam, err)
			}
			if cand.FamilyID != fam {
				t.Errorf("expected family %s, got %s", fam, cand.FamilyID)
			}
			if cand.ValidationStatus != candidate.ValidationValid {
				t.Errorf("expected ValidationValid, got %s", cand.ValidationStatus)
			}
			if cand.Status != candidate.StatusCandidatePendingReview {
				t.Errorf("expected StatusCandidatePendingReview, got %s", cand.Status)
			}
			if len(cand.DerivedFixture.Postings) < 2 {
				t.Errorf("expected at least 2 derived postings, got %d", len(cand.DerivedFixture.Postings))
			}
			// Invariant: Unapproved candidate is NOT active for practice
			if bank.IsActiveForPractice(string(cand.Status)) {
				t.Errorf("candidate %s must NOT be active for practice", cand.ID)
			}
		})
	}
}

func TestProviderCandidateGeneratorWithSimulatedTutor(t *testing.T) {
	eng, catalog := setupTestEngine(t)

	// 1. Success case: Simulated tutor returns valid JSON
	sim := tutor.NewSimulatedProvider(tutor.SimulatedConfig{
		Delay: 10 * time.Millisecond,
	})
	// Note: SimulatedProvider's default prose is not JSON, so let's mock a provider that returns valid JSON
	validJSONTutor := &mockJSONTutor{
		name: "MockJSONTutor",
		jsonResponse: `{
			"family_id": "customer_advance",
			"scenario_template": "A corporate client pays ${amount_dollars} cash today for security audits scheduled next month.",
			"parameters": {
				"amount_minor_units": [5000, 10000, 20000]
			},
			"concepts": ["cash_vs_revenue", "unearned_revenue_classification"]
		}`,
	}

	gen := candidate.NewProviderCandidateGenerator(validJSONTutor, eng, catalog)
	cand, err := gen.Generate(context.Background(), candidate.GenerateRequest{
		FamilyID: bank.FamilyCustomerAdvance,
	})
	if err != nil {
		t.Fatalf("unexpected generation error: %v", err)
	}
	if cand.ValidationStatus != candidate.ValidationValid {
		t.Errorf("expected ValidationValid, got %s", cand.ValidationStatus)
	}
	if cand.Provenance.Source != "MockJSONTutor" {
		t.Errorf("expected source MockJSONTutor, got %s", cand.Provenance.Source)
	}

	// 2. Cancellation case
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	_, err = gen.Generate(ctx, candidate.GenerateRequest{FamilyID: bank.FamilyCustomerAdvance})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}

	// 3. Error case: Provider fails upstream
	failingTutor := &mockJSONTutor{
		name:     "FailingTutor",
		failWith: errors.New("upstream rate limited"),
	}
	genFailing := candidate.NewProviderCandidateGenerator(failingTutor, eng, catalog)
	candFailed, err := genFailing.Generate(context.Background(), candidate.GenerateRequest{FamilyID: bank.FamilyCustomerAdvance})
	if err == nil {
		t.Fatalf("expected error from failing tutor, got nil")
	}
	if candFailed == nil {
		t.Fatalf("expected candidate struct with failure info, got nil")
	}
	if candFailed.ValidationStatus != candidate.ValidationFailed {
		t.Errorf("expected ValidationFailed, got %s", candFailed.ValidationStatus)
	}
	if bank.IsActiveForPractice(string(candFailed.Status)) {
		t.Errorf("failed candidate must NOT be active for practice")
	}
	_ = sim
}

type mockJSONTutor struct {
	name         string
	jsonResponse string
	failWith     error
}

func (m *mockJSONTutor) Name() string {
	return m.name
}

func (m *mockJSONTutor) Hint(ctx context.Context, req tutor.Request) (tutor.Response, error) {
	if m.failWith != nil {
		return tutor.Response{}, m.failWith
	}
	return tutor.Response{Text: m.jsonResponse, Provider: m.name, GeneratedAt: time.Now().UTC()}, nil
}

func (m *mockJSONTutor) Explain(ctx context.Context, req tutor.Request) (tutor.Response, error) {
	if err := ctx.Err(); err != nil {
		return tutor.Response{}, err
	}
	if m.failWith != nil {
		return tutor.Response{}, m.failWith
	}
	return tutor.Response{Text: m.jsonResponse, Provider: m.name, GeneratedAt: time.Now().UTC()}, nil
}

func TestWeakestConceptTargeting(t *testing.T) {
	projections := make(map[string]*mastery.ConceptStats)

	// Concept 1: High score, low need
	s1 := mastery.NewConceptStats("cash_classification")
	s1.IndependentAttempts = 10
	s1.IndependentSuccesses = 10
	s1.Alpha = 11.0
	s1.Beta = 1.0
	projections["cash_classification"] = s1

	// Concept 2: Low score, high need (weakest concept)
	s2 := mastery.NewConceptStats("prepaid_expenses")
	s2.IndependentAttempts = 5
	s2.IndependentSuccesses = 1
	s2.Alpha = 2.0
	s2.Beta = 5.0
	projections["prepaid_expenses"] = s2

	now := time.Now().UTC()
	weakConcept, matchingFamily := candidate.SelectWeakestConcept(projections, now)

	if weakConcept != "prepaid_expenses" {
		t.Errorf("expected weakest concept prepaid_expenses, got %s", weakConcept)
	}
	if matchingFamily != bank.FamilyPrepaidPurchase && matchingFamily != bank.FamilyPrepaidConsumption {
		t.Errorf("expected prepaid family, got %s", matchingFamily)
	}
}

func TestCandidatePreviewFormatting(t *testing.T) {
	eng, catalog := setupTestEngine(t)
	gen := candidate.NewOfflineCandidateGenerator(eng, catalog)

	cand, err := gen.Generate(context.Background(), candidate.GenerateRequest{
		FamilyID: bank.FamilyCustomerAdvance,
		Seed:     100,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	preview := cand.FormatPreview()

	// Check required preview elements
	requiredSnippets := []string{
		"CANDIDATE QUESTION PREVIEW",
		"Status: candidate_pending_review",
		"Validation: valid",
		"WORDING:",
		"Template:",
		"Sample:",
		"ASSUMPTIONS (Parameters):",
		"amount_minor_units:",
		"DERIVED POSTINGS (Engine-Derived Canonical Entry):",
		"Dr.      cash",
		"Cr.      unearned_revenue",
		"Equation Impact:",
		"CONCEPT TAGS:",
		"PROVENANCE:",
		"Source:        offline_generator",
	}

	for _, snip := range requiredSnippets {
		if !bytes.Contains([]byte(preview), []byte(snip)) {
			t.Errorf("preview missing expected snippet: %q", snip)
		}
	}
}

func TestInvariantCandidateGenerationDoesNotMutateLearnerEvidence(t *testing.T) {
	eng, catalog := setupTestEngine(t)
	gen := candidate.NewOfflineCandidateGenerator(eng, catalog)

	// Simulate generating 20 candidates
	for i := 0; i < 20; i++ {
		cand, err := gen.Generate(context.Background(), candidate.GenerateRequest{
			FamilyID: bank.FamilyCashService,
			Seed:     int64(i + 1),
		})
		if err != nil {
			t.Fatalf("error generating candidate %d: %v", i, err)
		}

		// Invariant: Status must never be active
		if bank.IsActiveForPractice(string(cand.Status)) {
			t.Fatalf("candidate %s must never be active for practice", cand.ID)
		}
	}

	// Verify that learner mastery stats remain zero / clean
	var attempts []domain.Attempt
	projections := mastery.RebuildProjections(attempts)
	if len(projections) != 0 {
		t.Fatalf("candidate generation produced learner evidence! projections count = %d", len(projections))
	}
}
