package drill_test

import (
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"strings"
	"testing"
)

func TestExpansionBatch8SemanticsTeachingAndContrasts(t *testing.T) {
	expected := map[string][2]domain.AccountID{
		"cash_service_unpaid_booking": {"cash", "service_revenue"},
		"service_on_credit_due_today": {"accounts_receivable", "service_revenue"},
		"collect_receivable_overdue":  {"cash", "accounts_receivable"},
		"customer_advance_same_day":   {"cash", "unearned_revenue"},
		"earn_advance_same_day":       {"unearned_revenue", "service_revenue"},
	}
	pairs := [][2]string{
		{"cash_service_unpaid_booking", "service_on_credit_due_today"},
		{"cash_service_unpaid_booking", "collect_receivable_overdue"},
		{"cash_service_unpaid_booking", "customer_advance_same_day"},
		{"customer_advance_same_day", "earn_advance_same_day"},
		{"earn_advance_same_day", "cash_service_unpaid_booking"},
	}
	verifyPilotBatch(t, "../../docs/EXPANSION-BATCH8.json", expected, pairs)
}

func TestPublishedExpansionBatch8MatchesReviewedPackage(t *testing.T) {
	verifyPublishedPilot(t, "../../docs/EXPANSION-BATCH8.json")
}

func TestExpansionBatch8SameDayReceiptThenEarning(t *testing.T) {
	gen, _, catalog := setupTestDrill(t)
	pack, err := bank.LoadQuestionBankFile("../../docs/EXPANSION-BATCH8.json", catalog)
	if err != nil {
		t.Fatal(err)
	}
	morning := *getQuestionByID(pack, "customer_advance_same_day")
	evening := *getQuestionByID(pack, "earn_advance_same_day")
	morning.Status, evening.Status = bank.StatusApprovedActive, bank.StatusApprovedActive
	for _, amount := range morning.Parameters["amount_minor_units"] {
		ledger := domain.NewTAccountLedger(catalog)
		for i, q := range []bank.QuestionJSON{morning, evening} {
			inst, err := gen.GenerateInstance(q, 101, map[string]int64{"amount_minor_units": amount})
			if err != nil {
				t.Fatal(err)
			}
			for _, stage := range inst.StageAnswers {
				text := stage.Prompt + stage.CausalHint + stage.Explanation
				for _, hint := range stage.MistakeHints {
					text += hint
				}
				for _, invalid := range []string{"next month", "last month", "No cash changes hands today", "no cash moved today", "Did money arrive today"} {
					if strings.Contains(text, invalid) {
						t.Fatalf("%s teaching conflicts with same-day timing: %s", q.ID, invalid)
					}
				}
			}
			for _, p := range inst.Entry.Postings {
				ledger.GetOrCreate(p.AccountID).Post(p.Side, p.Amount)
			}
			wantLiability, wantRevenue := amount, int64(0)
			if i == 1 {
				wantLiability, wantRevenue = 0, amount
			}
			if ledger.GetOrCreate("cash").NetBalance.Cents() != amount || ledger.GetOrCreate("unearned_revenue").NetBalance.Cents() != wantLiability || ledger.GetOrCreate("service_revenue").NetBalance.Cents() != wantRevenue {
				t.Fatalf("step %d duplicated receipt or recognized earning at wrong time", i)
			}
		}
	}
}
