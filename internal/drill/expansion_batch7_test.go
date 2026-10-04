package drill_test

import (
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"testing"
)

func TestExpansionBatch7PriorCapitalRemainsRecorded(t *testing.T) {
	gen, _, catalog := setupTestDrill(t)
	pack, err := bank.LoadQuestionBankFile("../../docs/EXPANSION-BATCH7.json", catalog)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"borrow_cash_prior_investment", "issue_shares_separate_private_sale"} {
		q := *getQuestionByID(pack, id)
		q.Status = bank.StatusApprovedActive
		for _, amount := range q.Parameters["amount_minor_units"] {
			inst, err := gen.GenerateInstance(q, 101, map[string]int64{"amount_minor_units": amount})
			if err != nil {
				t.Fatal(err)
			}
			ledger := domain.NewTAccountLedger(catalog)
			ledger.GetOrCreate("cash").Post(domain.SideDebit, domain.NewMoney(2*amount))
			ledger.GetOrCreate("common_stock").Post(domain.SideCredit, domain.NewMoney(2*amount))
			for _, p := range inst.Entry.Postings {
				ledger.GetOrCreate(p.AccountID).Post(p.Side, p.Amount)
			}
			wantStock, wantNote := 3*amount, int64(0)
			if id == "borrow_cash_prior_investment" {
				wantStock, wantNote = 2*amount, amount
			}
			if ledger.GetOrCreate("common_stock").NetBalance.Cents() != wantStock || ledger.GetOrCreate("notes_payable").NetBalance.Cents() != wantNote || ledger.GetOrCreate("cash").NetBalance.Cents() != 3*amount {
				t.Fatalf("%s repeated prior capital or recorded an unrelated private trade", id)
			}
		}
	}
}

func TestExpansionBatch7SemanticsTeachingAndContrasts(t *testing.T) {
	expected := map[string][2]domain.AccountID{
		"cash_rent_before_opening":           {"rent_expense", "cash"},
		"cash_rent_stockholder_landlord":     {"rent_expense", "cash"},
		"issue_shares_customer_investor":     {"cash", "common_stock"},
		"issue_shares_separate_private_sale": {"cash", "common_stock"},
		"borrow_cash_prior_investment":       {"cash", "notes_payable"},
	}
	pairs := [][2]string{
		{"cash_rent_before_opening", "equipment_purchase_workshop"},
		{"cash_rent_stockholder_landlord", "dividend_cash_working_owners"},
		{"issue_shares_customer_investor", "customer_advance_additional_receipt"},
		{"issue_shares_customer_investor", "cash_service_separate_job"},
		{"issue_shares_separate_private_sale", "borrow_cash_prior_investment"},
	}
	verifyPilotBatch(t, "../../docs/EXPANSION-BATCH7.json", expected, pairs)
}

func TestPublishedExpansionBatch7MatchesReviewedPackage(t *testing.T) {
	verifyPublishedPilot(t, "../../docs/EXPANSION-BATCH7.json")
}
