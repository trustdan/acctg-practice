package drill_test

import (
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"strings"
	"testing"
)

func TestExpansionBatch14SemanticsTeachingAndContrasts(t *testing.T) {
	verifyPilotBatch(t, "../../docs/EXPANSION-BATCH14.json", map[string][2]domain.AccountID{
		"cash_service_undeposited_notes":    {"cash", "service_revenue"},
		"service_on_credit_cash_sale_label": {"accounts_receivable", "service_revenue"},
		"collect_receivable_before_due":     {"cash", "accounts_receivable"},
		"customer_advance_customer_books":   {"cash", "unearned_revenue"},
		"earn_advance_separate_unpaid_job":  {"unearned_revenue", "service_revenue"},
	}, [][2]string{
		{"cash_service_undeposited_notes", "issue_shares_customer_investor"},
		{"service_on_credit_cash_sale_label", "cash_service_undeposited_notes"},
		{"collect_receivable_before_due", "customer_advance_customer_books"},
		{"customer_advance_customer_books", "earn_advance_separate_unpaid_job"},
		{"earn_advance_separate_unpaid_job", "service_on_credit_cash_sale_label"},
	})
}
func TestPublishedExpansionBatch14MatchesReviewedPackage(t *testing.T) {
	verifyPublishedPilot(t, "../../docs/EXPANSION-BATCH14.json")
}
func TestExpansionBatch14PaymentApplicationAndRemainingBalances(t *testing.T) {
	gen, _, catalog := setupTestDrill(t)
	pack, err := bank.LoadQuestionBankFile("../../docs/EXPANSION-BATCH14.json", catalog)
	if err != nil {
		t.Fatal(err)
	}
	for _, q := range pack.Questions {
		q.Status = bank.StatusApprovedActive
		for _, amount := range q.Parameters["amount_minor_units"] {
			ledger := domain.NewTAccountLedger(catalog)
			want := map[domain.AccountID]int64{"cash": 2 * amount, "accounts_receivable": amount, "unearned_revenue": amount, "service_revenue": 2 * amount, "common_stock": 3 * amount, "notes_payable": amount}
			for id, n := range want {
				side := domain.SideDebit
				if id == "unearned_revenue" || id == "service_revenue" || id == "common_stock" || id == "notes_payable" {
					side = domain.SideCredit
				}
				ledger.GetOrCreate(id).Post(side, domain.NewMoney(n))
			}
			inst, err := gen.GenerateInstance(q, 101, map[string]int64{"amount_minor_units": amount})
			if err != nil {
				t.Fatal(err)
			}
			for _, p := range inst.Entry.Postings {
				ledger.GetOrCreate(p.AccountID).Post(p.Side, p.Amount)
			}
			switch q.FamilyID {
			case "cash_service":
				want["cash"] += amount
				want["service_revenue"] += amount
			case "service_on_credit":
				want["accounts_receivable"] += amount
				want["service_revenue"] += amount
			case "collect_receivable":
				want["cash"] += amount
				want["accounts_receivable"] = 0
			case "customer_advance":
				want["cash"] += amount
				want["unearned_revenue"] += amount
			case "earn_advance":
				want["unearned_revenue"] = 0
				want["service_revenue"] += amount
			}
			for id, n := range want {
				if ledger.GetOrCreate(id).NetBalance.Cents() != n {
					t.Fatalf("%s %s misapplied receipt, repeated earning, or changed unrelated capital/debt", q.ID, id)
				}
			}
		}
	}
}

func TestExpansionBatch14PrepaidHintsKeepSeparateInvoiceInScope(t *testing.T) {
	gen, _, catalog := setupTestDrill(t)
	pack, err := bank.LoadQuestionBankFile("../../docs/EXPANSION-BATCH14.json", catalog)
	if err != nil {
		t.Fatal(err)
	}
	q := *getQuestionByID(pack, "earn_advance_separate_unpaid_job")
	q.Status = bank.StatusApprovedActive
	inst, err := gen.GenerateInstance(q, 101, map[string]int64{"amount_minor_units": 10000})
	if err != nil {
		t.Fatal(err)
	}
	hint := inst.StageAnswers[domain.StageIdentifyAccount].HintForOption("opt_ar")
	if !strings.Contains(hint, "For this prepaid work") {
		t.Fatalf("hint must distinguish prepaid job from unrelated unpaid invoice: %s", hint)
	}
}
