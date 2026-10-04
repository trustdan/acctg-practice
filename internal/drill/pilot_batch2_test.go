package drill_test

import (
	"testing"

	"github.com/trustdan/acctg-practice/internal/domain"
)

func TestPilotBatch2SemanticsTeachingAndContrasts(t *testing.T) {
	expected := map[string][2]domain.AccountID{
		"prepaid_purchase_clinic":    {"prepaid_insurance", "cash"},
		"prepaid_consumption_clinic": {"insurance_expense", "prepaid_insurance"},
		"cash_rent_clinic":           {"rent_expense", "cash"},
		"issue_shares_studio":        {"cash", "common_stock"},
		"repay_principal_studio":     {"notes_payable", "cash"},
		"dividend_cash_studio":       {"dividends", "cash"},
	}
	pairs := [][2]string{
		{"prepaid_purchase_clinic", "prepaid_consumption_clinic"},
		{"prepaid_purchase_clinic", "cash_rent_clinic"},
		{"issue_shares_studio", "borrow_cash_future_payroll"},
		{"repay_principal_studio", "dividend_cash_studio"},
		{"cash_rent_clinic", "repay_principal_studio"},
		{"cash_rent_clinic", "dividend_cash_studio"},
	}
	verifyPilotBatch(t, "../../docs/PILOT-BATCH2.json", expected, pairs)
}

func TestPublishedPilotBatch2MatchesReviewedPackage(t *testing.T) {
	verifyPublishedPilot(t, "../../docs/PILOT-BATCH2.json")
}
