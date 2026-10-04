package drill_test

import (
	"testing"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
)

func TestExpansionBatch6SemanticsTeachingAndContrasts(t *testing.T) {
	expected := map[string][2]domain.AccountID{
		"equipment_cash_existing_note":   {"equipment", "cash"},
		"prepaid_purchase_machine_cover": {"prepaid_insurance", "cash"},
		"prepaid_consumption_no_claim":   {"insurance_expense", "prepaid_insurance"},
		"repay_principal_final_payment":  {"notes_payable", "cash"},
		"dividend_cash_prior_earnings":   {"dividends", "cash"},
	}
	pairs := [][2]string{
		{"equipment_cash_existing_note", "prepaid_purchase_machine_cover"},
		{"equipment_cash_existing_note", "borrow_cash_future_equipment"},
		{"prepaid_purchase_machine_cover", "prepaid_consumption_no_claim"},
		{"repay_principal_final_payment", "dividend_cash_prior_earnings"},
		{"repay_principal_final_payment", "borrow_cash_basic"},
		{"dividend_cash_prior_earnings", "cash_rent_workshop"},
	}
	verifyPilotBatch(t, "../../docs/EXPANSION-BATCH6.json", expected, pairs)
}

func TestPublishedExpansionBatch6MatchesReviewedPackage(t *testing.T) {
	verifyPublishedPilot(t, "../../docs/EXPANSION-BATCH6.json")
}

func TestExpansionBatch6ExistingDebtAndFinalPayoff(t *testing.T) {
	gen, _, catalog := setupTestDrill(t)
	pack, err := bank.LoadQuestionBankFile("../../docs/EXPANSION-BATCH6.json", catalog)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"equipment_cash_existing_note", "repay_principal_final_payment"} {
		q := getQuestionByID(pack, id)
		if q == nil {
			t.Fatalf("missing %s", id)
		}
		preview := *q
		preview.Status = bank.StatusApprovedActive
		for _, amount := range q.Parameters["amount_minor_units"] {
			inst, err := gen.GenerateInstance(preview, 100, map[string]int64{"amount_minor_units": amount})
			if err != nil {
				t.Fatal(err)
			}
			ledger := domain.NewTAccountLedger(catalog)
			ledger.GetOrCreate("cash").Post(domain.SideDebit, domain.NewMoney(amount))
			ledger.GetOrCreate("notes_payable").Post(domain.SideCredit, domain.NewMoney(amount))
			for _, posting := range inst.Entry.Postings {
				ledger.GetOrCreate(posting.AccountID).Post(posting.Side, posting.Amount)
			}
			note := ledger.GetOrCreate("notes_payable")
			if id == "equipment_cash_existing_note" {
				if note.NetBalance.Cents() != amount || note.BalanceSide != domain.SideCredit {
					t.Fatalf("cash acquisition changed the unrelated note: %+v", note)
				}
			} else if note.NetBalance.Cents() != 0 {
				t.Fatalf("final principal was not fully settled: %+v", note)
			}
			if ledger.GetOrCreate("cash").NetBalance.Cents() != 0 {
				t.Fatal("payment did not reduce opening cash by exactly the stated amount")
			}
		}
	}
}
