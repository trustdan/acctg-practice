package drill_test

import (
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"testing"
)

func TestExpansionBatch12SemanticsTeachingAndContrasts(t *testing.T) {
	verifyPilotBatch(t, "../../docs/EXPANSION-BATCH12.json", map[string][2]domain.AccountID{
		"equipment_cash_additional_machine":   {"equipment", "cash"},
		"repay_principal_bank_landlord":       {"notes_payable", "cash"},
		"dividend_cash_service_provider":      {"dividends", "cash"},
		"prepaid_purchase_broker_own_policy":  {"prepaid_insurance", "cash"},
		"prepaid_consumption_final_remaining": {"insurance_expense", "prepaid_insurance"},
	}, [][2]string{
		{"equipment_cash_additional_machine", "repay_principal_equipment_loan"},
		{"repay_principal_bank_landlord", "cash_rent_stockholder_landlord"},
		{"dividend_cash_service_provider", "repay_principal_bank_landlord"},
		{"prepaid_purchase_broker_own_policy", "cash_service_basic"},
		{"prepaid_purchase_broker_own_policy", "prepaid_consumption_final_remaining"},
	})
}

func TestPublishedExpansionBatch12MatchesReviewedPackage(t *testing.T) {
	verifyPublishedPilot(t, "../../docs/EXPANSION-BATCH12.json")
}

func TestExpansionBatch12PreservesPriorCostsAndClearsOnlyRemainingCoverage(t *testing.T) {
	gen, _, catalog := setupTestDrill(t)
	pack, err := bank.LoadQuestionBankFile("../../docs/EXPANSION-BATCH12.json", catalog)
	if err != nil {
		t.Fatal(err)
	}
	for _, q := range pack.Questions {
		q.Status = bank.StatusApprovedActive
		for _, amount := range q.Parameters["amount_minor_units"] {
			ledger := domain.NewTAccountLedger(catalog)
			want := map[domain.AccountID]int64{"cash": 5 * amount, "equipment": 3 * amount, "notes_payable": 3 * amount, "prepaid_insurance": amount, "rent_expense": 2 * amount, "insurance_expense": 2 * amount, "common_stock": 4 * amount, "dividends": 0}
			for id, n := range want {
				side := domain.SideDebit
				if id == "notes_payable" || id == "common_stock" {
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
			case "equipment_purchase_cash":
				want["equipment"] += amount
				want["cash"] -= amount
			case "repay_note_principal":
				want["notes_payable"] -= amount
				want["cash"] -= amount
			case "dividend_cash":
				want["dividends"] += amount
				want["cash"] -= amount
			case "prepaid_purchase":
				want["prepaid_insurance"] += amount
				want["cash"] -= amount
			case "prepaid_consumption":
				want["prepaid_insurance"] = 0
				want["insurance_expense"] += amount
			}
			for id, n := range want {
				if ledger.GetOrCreate(id).NetBalance.Cents() != n {
					t.Fatalf("%s %s duplicated old cost/payment, altered unrelated balance, or omitted final expiry", q.ID, id)
				}
			}
		}
	}
}
