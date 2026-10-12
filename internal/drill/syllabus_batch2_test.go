package drill_test

import (
	"testing"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
)

func TestSyllabusBatch2StageBuilders(t *testing.T) {
	_, _, catalog := setupTestDrill(t)
	eng := engine.NewEngine(catalog)
	gen := drill.NewGenerator(catalog, eng)

	testCases := []struct {
		name     string
		question bank.QuestionJSON
		amount   int64
	}{
		{
			name: "Module 5: Customer Write-Off",
			question: bank.QuestionJSON{
				ID:               "test_write_off",
				Version:          1,
				FamilyID:         bank.FamilyWriteOffReceivable,
				RuleVersion:      1,
				Status:           bank.StatusDraft,
				ScenarioTemplate: "Company writes off ${amount_dollars} of uncollectible accounts under allowance method.",
				Parameters:       map[string][]int64{"amount_minor_units": {120000}},
				Concepts:         []string{"allowance_write_off"},
				Module:           5,
			},
			amount: 120000,
		},
		{
			name: "Module 7: Issue Bonds at Discount",
			question: bank.QuestionJSON{
				ID:               "test_bonds_discount",
				Version:          1,
				FamilyID:         bank.FamilyIssueBondsDiscount,
				RuleVersion:      1,
				Status:           bank.StatusDraft,
				ScenarioTemplate: "Company issues ${amount_dollars} of 10-year bonds at a discount.",
				Parameters: map[string][]int64{
					"amount_minor_units":   {10000000},
					"cash_minor_units":     {9500000},
					"discount_minor_units": {500000},
				},
				Concepts: []string{"bond_discount_issuance"},
				Module:   7,
			},
			amount: 10000000,
		},
		{
			name: "Module 9: Capitalize Software Post-Feasibility",
			question: bank.QuestionJSON{
				ID:               "test_cap_software",
				Version:          1,
				FamilyID:         bank.FamilyCapitalizeSoftwareFeasibility,
				RuleVersion:      1,
				Status:           bank.StatusDraft,
				ScenarioTemplate: "Company pays ${amount_dollars} cash for software development post-feasibility.",
				Parameters:       map[string][]int64{"amount_minor_units": {3000000}},
				Concepts:         []string{"software_capitalization"},
				Module:           9,
			},
			amount: 3000000,
		},
		{
			name: "Module 3/8: Current Ratio Liquidity Analysis",
			question: bank.QuestionJSON{
				ID:               "test_current_ratio",
				Version:          1,
				FamilyID:         bank.FamilyCurrentRatioAnalysis,
				RuleVersion:      1,
				Status:           bank.StatusDraft,
				ScenarioTemplate: "Retail company pays ${amount_dollars} cash to settle Accounts Payable.",
				Parameters:       map[string][]int64{"amount_minor_units": {2000000}},
				Concepts:         []string{"current_ratio_liquidity"},
				Module:           3,
			},
			amount: 2000000,
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

			stages := []domain.DrillStage{
				domain.StageIdentifyAccount,
				domain.StageAccountCategory,
				domain.StageDirection,
				domain.StageDebitCredit,
				domain.StageCounterAccount,
				domain.StageBalancedEntry,
				domain.StageEquationEffect,
			}
			for _, st := range stages {
				ans, ok := inst.StageAnswers[st]
				if !ok {
					t.Errorf("missing stage %s", st)
					continue
				}
				if ans.CorrectOptionID == "" {
					t.Errorf("stage %s missing CorrectOptionID", st)
				}
				if len(ans.Options) < 2 {
					t.Errorf("stage %s has fewer than 2 options: %d", st, len(ans.Options))
				}
				if ans.CausalHint == "" || ans.Explanation == "" {
					t.Errorf("stage %s missing hint or explanation", st)
				}
			}
		})
	}
}
