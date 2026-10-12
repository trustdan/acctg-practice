package drill_test

import (
	"testing"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
)

func TestSyllabusBatch3StageBuilders(t *testing.T) {
	_, _, catalog := setupTestDrill(t)
	eng := engine.NewEngine(catalog)
	gen := drill.NewGenerator(catalog, eng)

	testCases := []struct {
		name     string
		question bank.QuestionJSON
		amount   int64
	}{
		{
			name: "Module 6: Treasury Stock Repurchase",
			question: bank.QuestionJSON{
				ID:               "test_treasury_stock",
				Version:          1,
				FamilyID:         bank.FamilyRepurchaseTreasuryStock,
				RuleVersion:      1,
				Status:           bank.StatusDraft,
				ScenarioTemplate: "Company repurchases ${amount_dollars} of common stock into treasury.",
				Parameters:       map[string][]int64{"amount_minor_units": {2500000}},
				Concepts:         []string{"treasury_stock_accounting"},
				Module:           6,
			},
			amount: 2500000,
		},
		{
			name: "Module 7: Bond Interest & Discount Amortization",
			question: bank.QuestionJSON{
				ID:               "test_bond_interest",
				Version:          1,
				FamilyID:         bank.FamilyPayBondInterest,
				RuleVersion:      1,
				Status:           bank.StatusDraft,
				ScenarioTemplate: "Company records ${amount_dollars} bond interest expense and discount amortization.",
				Parameters: map[string][]int64{
					"amount_minor_units":   {250000},
					"cash_minor_units":     {200000},
					"discount_minor_units": {50000},
				},
				Concepts: []string{"bond_interest_expense"},
				Module:   7,
			},
			amount: 250000,
		},
		{
			name: "Module 8: Lehman Repo 105 Analysis",
			question: bank.QuestionJSON{
				ID:               "test_repo105",
				Version:          1,
				FamilyID:         bank.FamilyRepo105Analysis,
				RuleVersion:      1,
				Status:           bank.StatusDraft,
				ScenarioTemplate: "Bank receives ${amount_dollars} cash under a short-term repo agreement.",
				Parameters:       map[string][]int64{"amount_minor_units": {5000000}},
				Concepts:         []string{"repurchase_agreements"},
				Module:           8,
			},
			amount: 5000000,
		},
		{
			name: "Module 5: Apple Multi-Element Revenue Allocation",
			question: bank.QuestionJSON{
				ID:               "test_apple_revenue",
				Version:          1,
				FamilyID:         bank.FamilyAppleMultiElementRevenue,
				RuleVersion:      1,
				Status:           bank.StatusDraft,
				ScenarioTemplate: "Company sells bundle for ${amount_dollars} cash allocating hardware and software updates.",
				Parameters: map[string][]int64{
					"amount_minor_units":            {100000},
					"delivered_revenue_minor_units": {80000},
					"unearned_revenue_minor_units":  {20000},
				},
				Concepts: []string{"multi_element_revenue_recognition"},
				Module:   5,
			},
			amount: 100000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			seed := int64(42)
			inst, err := gen.GenerateInstance(tc.question, seed, map[string]int64{"amount_minor_units": tc.amount})
			if err != nil {
				t.Fatalf("GenerateInstance failed: %v", err)
			}
			if !inst.Entry.IsBalanced() {
				t.Fatalf("Entry not balanced: %+v", inst.Entry)
			}
			if len(inst.StageAnswers) != 7 {
				t.Fatalf("expected 7 stage answers, got %d", len(inst.StageAnswers))
			}

			// Validate each stage has exactly 1 correct answer and valid options
			for stageKey, ans := range inst.StageAnswers {
				if ans.CorrectOptionID == "" {
					t.Errorf("stage %s missing CorrectOptionID", stageKey)
				}
				correctCount := 0
				for _, opt := range ans.Options {
					if opt.ID == ans.CorrectOptionID {
						correctCount++
						if opt.ErrorTag != "" {
							t.Errorf("stage %s correct option %s has error tag %s", stageKey, opt.ID, opt.ErrorTag)
						}
					}
				}
				if correctCount != 1 {
					t.Errorf("stage %s has %d correct options, expected 1", stageKey, correctCount)
				}
			}

			// Reconcile transaction effects
			recon, err := domain.ReconcileTransaction(inst.Entry, catalog)
			if err != nil || !recon.Reconciled {
				t.Fatalf("Transaction reconciliation failed: %v, recon: %+v", err, recon)
			}
		})
	}
}
