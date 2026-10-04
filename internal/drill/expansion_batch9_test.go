package drill_test

import (
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"testing"
)

func TestExpansionBatch9SemanticsTeachingAndContrasts(t *testing.T) {
	expected := map[string][2]domain.AccountID{
		"prepaid_purchase_discounted_policy": {"prepaid_insurance", "cash"},
		"prepaid_consumption_closed_month":   {"insurance_expense", "prepaid_insurance"},
		"equipment_cash_secondhand":          {"equipment", "cash"},
		"repay_principal_early_partial":      {"notes_payable", "cash"},
		"dividend_cash_total_distribution":   {"dividends", "cash"},
	}
	pairs := [][2]string{
		{"prepaid_purchase_discounted_policy", "prepaid_consumption_closed_month"},
		{"equipment_cash_secondhand", "prepaid_purchase_discounted_policy"},
		{"equipment_cash_secondhand", "cash_rent_before_opening"},
		{"repay_principal_early_partial", "dividend_cash_total_distribution"},
		{"repay_principal_early_partial", "borrow_cash_basic"},
		{"dividend_cash_total_distribution", "cash_rent_stockholder_landlord"},
	}
	verifyPilotBatch(t, "../../docs/EXPANSION-BATCH9.json", expected, pairs)
}

func TestPublishedExpansionBatch9MatchesReviewedPackage(t *testing.T) {
	verifyPublishedPilot(t, "../../docs/EXPANSION-BATCH9.json")
}

func TestExpansionBatch9PartialBalancesAndTotalDistribution(t *testing.T) {
	gen, _, catalog := setupTestDrill(t)
	pack, err := bank.LoadQuestionBankFile("../../docs/EXPANSION-BATCH9.json", catalog)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"prepaid_consumption_closed_month", "repay_principal_early_partial", "dividend_cash_total_distribution"} {
		q := *getQuestionByID(pack, id)
		q.Status = bank.StatusApprovedActive
		for _, amount := range q.Parameters["amount_minor_units"] {
			ledger := domain.NewTAccountLedger(catalog)
			ledger.GetOrCreate("cash").Post(domain.SideDebit, domain.NewMoney(3*amount))
			ledger.GetOrCreate("prepaid_insurance").Post(domain.SideDebit, domain.NewMoney(3*amount))
			ledger.GetOrCreate("notes_payable").Post(domain.SideCredit, domain.NewMoney(3*amount))
			ledger.GetOrCreate("common_stock").Post(domain.SideCredit, domain.NewMoney(3*amount))
			inst, err := gen.GenerateInstance(q, 101, map[string]int64{"amount_minor_units": amount})
			if err != nil {
				t.Fatal(err)
			}
			for _, p := range inst.Entry.Postings {
				ledger.GetOrCreate(p.AccountID).Post(p.Side, p.Amount)
			}
			cash, prepaid, note, dividends := 3*amount, 3*amount, 3*amount, int64(0)
			switch id {
			case "prepaid_consumption_closed_month":
				prepaid = 2 * amount
			case "repay_principal_early_partial":
				cash, note = 2*amount, 2*amount
			case "dividend_cash_total_distribution":
				cash, dividends = 2*amount, amount
			}
			if ledger.GetOrCreate("cash").NetBalance.Cents() != cash || ledger.GetOrCreate("prepaid_insurance").NetBalance.Cents() != prepaid || ledger.GetOrCreate("notes_payable").NetBalance.Cents() != note || ledger.GetOrCreate("dividends").NetBalance.Cents() != dividends || ledger.GetOrCreate("common_stock").NetBalance.Cents() != 3*amount {
				t.Fatalf("%s cleared unrelated balances, consumed future coverage, or multiplied total payment", id)
			}
		}
	}
}
