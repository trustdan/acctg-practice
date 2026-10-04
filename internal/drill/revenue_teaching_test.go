package drill_test

import (
	"strings"
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// The most tempting balanced misconception must still be wrong when the
// learner starts directly at entry assembly, without the category stages.
func TestMisconceptionsWithFadedScaffold(t *testing.T) {
	gen, questions, _ := setupTestDrill(t)
	cases := []struct {
		id, wrong, tag string
	}{
		{"cash_service_basic", "opt_dr_ar_cr_rev", engine.TagWrongAccount},
		{"customer_advance_basic", "opt_dr_cash_cr_rev", engine.TagRevenueRecognizedPrematurely},
		{"service_on_credit_basic", "opt_dr_cash_cr_rev", engine.TagCashRecordedWhenUncollected},
		{"collect_receivable_basic", "opt_dr_cash_cr_rev", engine.TagDuplicateRevenueOnCollection},
		{"earn_advance_basic", "opt_dr_cash_cr_rev", engine.TagCashRecordedOnEarningAdvance},
		{"cash_rent_basic", "opt_dr_prepaid_cr_cash", engine.TagWrongAccount},
		{"borrow_cash_basic", "opt_dr_cash_cr_rev", engine.TagRevenueRecordedOnBorrowing},
		{"issue_shares_basic", "opt_dr_cash_cr_rev", engine.TagRevenueRecordedOnShareIssue},
		{"prepaid_insurance_retail", "opt_dr_expense_cr_cash", engine.TagExpenseRecordedOnPrepaidPurchase},
		{"prepaid_consumption_retail", "opt_no_entry", engine.TagPrepaidNotExpensedOnConsumption},
		{"equipment_cash_rental", "opt_dr_expense_cr_cash", engine.TagExpenseRecordedOnEquipmentPurchase},
		{"repay_principal_logistics", "opt_dr_expense_cr_cash", engine.TagExpenseRecordedOnLoanRepayment},
		{"dividend_cash_retail", "opt_dr_expense_cr_cash", engine.TagExpenseRecordedOnDividend},
	}
	for _, tc := range cases {
		t.Run(tc.id, func(t *testing.T) {
			q := getQuestionByID(questions, tc.id)
			for _, amount := range q.Parameters["amount_minor_units"] {
				inst, err := gen.GenerateInstanceWithScaffold(*q, 42, map[string]int64{"amount_minor_units": amount}, domain.ScaffoldFaded)
				if err != nil {
					t.Fatal(err)
				}
				session := drill.NewSession("faded-revenue", inst)
				stage, err := session.CurrentStage()
				if err != nil || stage.Stage != domain.StageBalancedEntry {
					t.Fatalf("faded scaffold must begin with entry assembly: %v", err)
				}
				feedback, err := session.SubmitOption(tc.wrong, time.Unix(0, 0).UTC())
				if err != nil {
					t.Fatal(err)
				}
				if feedback.IsCorrect || feedback.AdvanceStage || feedback.ErrorTag != tc.tag {
					t.Fatalf("balanced misconception incorrectly graded: %+v", feedback)
				}
				if !strings.Contains(feedback.Hint, "?") || strings.Contains(feedback.Hint, "Debit ") || strings.Contains(feedback.Hint, "Credit ") {
					t.Fatalf("hint must prompt reasoning without giving the entry: %q", feedback.Hint)
				}
				feedback, err = session.SubmitOption(stage.CorrectOptionID, time.Unix(1, 0).UTC())
				if err != nil || !feedback.IsCorrect || feedback.AssistanceLevel != domain.AssistanceRetry {
					t.Fatalf("assisted correction must remain retry evidence: %+v, %v", feedback, err)
				}
				if strings.Contains(feedback.Explanation, "${amount}") || strings.Contains(feedback.Explanation, "%s") {
					t.Fatalf("unrendered teaching amount: %q", feedback.Explanation)
				}
				if !session.Recap().IsBalanced {
					t.Fatal("canonical entry must remain balanced")
				}
			}
		})
	}
}

func TestTeachingDoesNotRevealLaterStageAnswers(t *testing.T) {
	gen, questions, _ := setupTestDrill(t)
	for _, q := range questions.Questions {
		if q.Status == "retired" {
			continue
		}
		id := q.ID
		inst, err := gen.GenerateInstance(q, 42, map[string]int64{"amount_minor_units": q.Parameters["amount_minor_units"][0]})
		if err != nil {
			t.Fatal(err)
		}
		first := inst.StageAnswers[domain.StageIdentifyAccount]
		for _, word := range []string{"debit", "credit"} {
			if strings.Contains(strings.ToLower(first.Explanation), word) {
				t.Errorf("%s identifies the side before the side stage: %s", id, first.Explanation)
			}
		}
		for _, stage := range inst.StageAnswers {
			if !strings.Contains(stage.CausalHint, "?") {
				t.Errorf("%s %s hint must ask a question", id, stage.Stage)
			}
			for _, option := range stage.Options {
				if option.ID == stage.CorrectOptionID && option.ErrorTag != "" {
					t.Errorf("%s %s correct option has error tag %q", id, stage.Stage, option.ErrorTag)
				}
				if strings.Contains(option.Text, "Revenue / Equity") || strings.Contains(option.Text, "memo entry only") {
					t.Errorf("%s retains misleading option: %s", id, option.Text)
				}
			}
		}
	}
}
