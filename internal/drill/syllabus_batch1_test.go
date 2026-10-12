package drill_test

import (
	"testing"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
)

func TestSyllabusBatch1StageBuilders(t *testing.T) {
	_, _, catalog := setupTestDrill(t)
	eng := engine.NewEngine(catalog)
	gen := drill.NewGenerator(catalog, eng)

	testCases := []struct {
		name     string
		question bank.QuestionJSON
		amount   int64
	}{
		{
			name: "Module 5: Bad Debt Allowance",
			question: bank.QuestionJSON{
				ID:               "test_bad_debt",
				Version:          1,
				FamilyID:         bank.FamilyBadDebtAllowance,
				RuleVersion:      1,
				Status:           bank.StatusDraft,
				ScenarioTemplate: "Company estimates ${amount_dollars} of uncollectible accounts.",
				Parameters:       map[string][]int64{"amount_minor_units": {500000}},
				Concepts:         []string{"bad_debt_allowance"},
				Module:           5,
			},
			amount: 500000,
		},
		{
			name: "Module 7: Issue Bonds at Par",
			question: bank.QuestionJSON{
				ID:               "test_issue_bonds",
				Version:          1,
				FamilyID:         bank.FamilyIssueBondsPar,
				RuleVersion:      1,
				Status:           bank.StatusDraft,
				ScenarioTemplate: "Company issues ${amount_dollars} of 10-year bonds at par for cash.",
				Parameters:       map[string][]int64{"amount_minor_units": {10000000}},
				Concepts:         []string{"bonds_payable_classification"},
				Module:           7,
			},
			amount: 10000000,
		},
		{
			name: "Module 6: Acquire Goodwill",
			question: bank.QuestionJSON{
				ID:               "test_goodwill",
				Version:          1,
				FamilyID:         bank.FamilyAcquireGoodwill,
				RuleVersion:      1,
				Status:           bank.StatusDraft,
				ScenarioTemplate: "Company acquires a business unit for ${amount_dollars} cash recognizing goodwill.",
				Parameters:       map[string][]int64{"amount_minor_units": {20000000}},
				Concepts:         []string{"goodwill_classification"},
				Module:           6,
			},
			amount: 20000000,
		},
		{
			name: "Module 3: Credit Sales Cash Flow Analysis",
			question: bank.QuestionJSON{
				ID:               "test_credit_cf",
				Version:          1,
				FamilyID:         bank.FamilyCreditSalesCashFlow,
				RuleVersion:      1,
				Status:           bank.StatusDraft,
				ScenarioTemplate: "Company performs ${amount_dollars} of services on credit.",
				Parameters:       map[string][]int64{"amount_minor_units": {1500000}},
				Concepts:         []string{"statement_of_cash_flows_indirect"},
				Module:           3,
			},
			amount: 1500000,
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
