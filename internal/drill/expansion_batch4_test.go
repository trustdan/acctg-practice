package drill_test

import (
	"testing"

	"github.com/trustdan/acctg-practice/internal/domain"
)

func TestExpansionBatch4SemanticsTeachingAndContrasts(t *testing.T) {
	expected := map[string][2]domain.AccountID{
		"prepaid_purchase_renewal":           {"prepaid_insurance", "cash"},
		"prepaid_consumption_partial_policy": {"insurance_expense", "prepaid_insurance"},
		"repay_principal_equipment_loan":     {"notes_payable", "cash"},
		"dividend_cash_working_owners":       {"dividends", "cash"},
		"borrow_cash_stockholder_note":       {"cash", "notes_payable"},
	}
	pairs := [][2]string{
		{"prepaid_purchase_renewal", "prepaid_consumption_partial_policy"},
		{"repay_principal_equipment_loan", "equipment_purchase_workshop"},
		{"repay_principal_equipment_loan", "dividend_cash_working_owners"},
		{"borrow_cash_stockholder_note", "repay_principal_equipment_loan"},
		{"borrow_cash_stockholder_note", "issue_shares_future_equipment"},
		{"dividend_cash_working_owners", "cash_rent_workshop"},
	}
	verifyPilotBatch(t, "../../docs/EXPANSION-BATCH4.json", expected, pairs)
}

func TestPublishedExpansionBatch4MatchesReviewedPackage(t *testing.T) {
	verifyPublishedPilot(t, "../../docs/EXPANSION-BATCH4.json")
}
