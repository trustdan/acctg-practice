package drill_test

import (
	"testing"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
)

func TestSyllabusBatch4StageBuilders(t *testing.T) {
	_, _, catalog := setupTestDrill(t)
	eng := engine.NewEngine(catalog)
	gen := drill.NewGenerator(catalog, eng)

	testCases := []struct {
		name     string
		question bank.QuestionJSON
		amount   int64
	}{
		{
			name: "Module 6: Goodwill Impairment Testing",
			question: bank.QuestionJSON{
				ID:               "test_goodwill_impairment",
				Version:          1,
				FamilyID:         bank.FamilyGoodwillImpairment,
				RuleVersion:      1,
				Status:           bank.StatusDraft,
				ScenarioTemplate: "Company records ${amount_dollars} goodwill impairment write-down.",
				Parameters:       map[string][]int64{"amount_minor_units": {5000000}},
				Concepts:         []string{"goodwill_impairment"},
				Module:           6,
			},
			amount: 5000000,
		},
		{
			name: "Module 8: Digital Asset Accounting (Tesla Bitcoin)",
			question: bank.QuestionJSON{
				ID:               "test_digital_assets",
				Version:          1,
				FamilyID:         bank.FamilyDigitalAssetAccounting,
				RuleVersion:      1,
				Status:           bank.StatusDraft,
				ScenarioTemplate: "Corporation purchases ${amount_dollars} of cryptocurrency for cash.",
				Parameters:       map[string][]int64{"amount_minor_units": {150000000}},
				Concepts:         []string{"digital_asset_accounting"},
				Module:           8,
			},
			amount: 150000000,
		},
		{
			name: "Module 9: Reporting Strategy & Conservative Revenue Policy",
			question: bank.QuestionJSON{
				ID:               "test_reporting_policy",
				Version:          1,
				FamilyID:         bank.FamilyReportingPolicyChoice,
				RuleVersion:      1,
				Status:           bank.StatusDraft,
				ScenarioTemplate: "Company receives ${amount_dollars} cash for service contract with pending obligations.",
				Parameters:       map[string][]int64{"amount_minor_units": {2500000}},
				Concepts:         []string{"reporting_strategy"},
				Module:           9,
			},
			amount: 2500000,
		},
		{
			name: "Module 10: Financial Forecasting & Pro-Forma Modeling",
			question: bank.QuestionJSON{
				ID:               "test_financial_forecasting",
				Version:          1,
				FamilyID:         bank.FamilyFinancialForecasting,
				RuleVersion:      1,
				Status:           bank.StatusDraft,
				ScenarioTemplate: "Analyst models ${amount_dollars} projected credit sales using disaggregated unit drivers.",
				Parameters:       map[string][]int64{"amount_minor_units": {4000000}},
				Concepts:         []string{"financial_forecasting"},
				Module:           10,
			},
			amount: 4000000,
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
