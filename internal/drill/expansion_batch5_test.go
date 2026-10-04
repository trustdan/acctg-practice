package drill_test

import (
	"testing"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
)

func TestExpansionBatch5SemanticsTeachingAndContrasts(t *testing.T) {
	expected := map[string][2]domain.AccountID{
		"cash_service_separate_job":           {"cash", "service_revenue"},
		"service_on_credit_new_client":        {"accounts_receivable", "service_revenue"},
		"collect_receivable_installment":      {"cash", "accounts_receivable"},
		"customer_advance_additional_receipt": {"cash", "unearned_revenue"},
		"earn_advance_completed_session":      {"unearned_revenue", "service_revenue"},
	}
	pairs := [][2]string{
		{"cash_service_separate_job", "collect_receivable_installment"},
		{"cash_service_separate_job", "service_on_credit_new_client"},
		{"cash_service_separate_job", "customer_advance_additional_receipt"},
		{"customer_advance_additional_receipt", "earn_advance_completed_session"},
		{"service_on_credit_new_client", "earn_advance_completed_session"},
		{"collect_receivable_installment", "customer_advance_additional_receipt"},
	}
	verifyPilotBatch(t, "../../docs/EXPANSION-BATCH5.json", expected, pairs)
}

func TestPublishedExpansionBatch5MatchesReviewedPackage(t *testing.T) {
	verifyPublishedPilot(t, "../../docs/EXPANSION-BATCH5.json")
}

func TestExpansionBatch5PreservesPriorAndRemainingBalances(t *testing.T) {
	gen, _, catalog := setupTestDrill(t)
	pack, err := bank.LoadQuestionBankFile("../../docs/EXPANSION-BATCH5.json", catalog)
	if err != nil {
		t.Fatal(err)
	}
	// Opening entries model amounts explicitly outside today's event. The
	// amount substitutions must affect only the receipt or completed portion.
	cases := []struct {
		id                          string
		openingDebit, openingCredit domain.AccountID
		openingMultiplier           int64
		check                       domain.AccountID
		wantMultiplier              int64
		side                        domain.Side
	}{
		{"cash_service_separate_job", "accounts_receivable", "service_revenue", 3, "accounts_receivable", 3, domain.SideDebit},
		{"service_on_credit_new_client", "cash", "service_revenue", 2, "cash", 2, domain.SideDebit},
		{"collect_receivable_installment", "accounts_receivable", "service_revenue", 3, "accounts_receivable", 2, domain.SideDebit},
		{"earn_advance_completed_session", "cash", "unearned_revenue", 3, "unearned_revenue", 2, domain.SideCredit},
	}
	for _, tc := range cases {
		q := getQuestionByID(pack, tc.id)
		if q == nil {
			t.Fatalf("missing %s", tc.id)
		}
		preview := *q
		preview.Status = bank.StatusApprovedActive
		for _, amount := range q.Parameters["amount_minor_units"] {
			inst, err := gen.GenerateInstance(preview, 100, map[string]int64{"amount_minor_units": amount})
			if err != nil {
				t.Fatal(err)
			}
			ledger := domain.NewTAccountLedger(catalog)
			ledger.GetOrCreate(tc.openingDebit).Post(domain.SideDebit, domain.NewMoney(tc.openingMultiplier*amount))
			ledger.GetOrCreate(tc.openingCredit).Post(domain.SideCredit, domain.NewMoney(tc.openingMultiplier*amount))
			for _, posting := range inst.Entry.Postings {
				ledger.GetOrCreate(posting.AccountID).Post(posting.Side, posting.Amount)
			}
			balance := ledger.GetOrCreate(tc.check)
			if balance.NetBalance.Cents() != tc.wantMultiplier*amount || balance.BalanceSide != tc.side {
				t.Fatalf("%s amount %d changed an unrelated balance or erased a remainder: %+v", tc.id, amount, balance)
			}
		}
	}
	q := getQuestionByID(pack, "customer_advance_additional_receipt")
	if q == nil {
		t.Fatal("missing additional advance")
	}
	preview := *q
	preview.Status = bank.StatusApprovedActive
	for _, amount := range q.Parameters["amount_minor_units"] {
		inst, err := gen.GenerateInstance(preview, 100, map[string]int64{"amount_minor_units": amount})
		if err != nil {
			t.Fatal(err)
		}
		ledger := domain.NewTAccountLedger(catalog)
		for _, id := range []domain.AccountID{"cash", "unearned_revenue"} {
			account := ledger.GetOrCreate(id)
			account.Post(account.NormalSide, domain.NewMoney(amount/2))
		}
		for _, posting := range inst.Entry.Postings {
			ledger.GetOrCreate(posting.AccountID).Post(posting.Side, posting.Amount)
		}
		for _, id := range []domain.AccountID{"cash", "unearned_revenue"} {
			account := ledger.GetOrCreate(id)
			if account.NetBalance.Cents() != amount+amount/2 || account.BalanceSide != account.NormalSide {
				t.Fatalf("additional receipt overwrote or repeated prior %s: %+v", id, account)
			}
		}
	}
}
