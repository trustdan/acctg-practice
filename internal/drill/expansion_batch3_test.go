package drill_test

import (
	"testing"

	"github.com/trustdan/acctg-practice/internal/domain"
)

func TestExpansionBatch3SemanticsTeachingAndContrasts(t *testing.T) {
	expected := map[string][2]domain.AccountID{
		"equipment_purchase_workshop":    {"equipment", "cash"},
		"equipment_purchase_replacement": {"equipment", "cash"},
		"cash_rent_workshop":             {"rent_expense", "cash"},
		"borrow_cash_future_equipment":   {"cash", "notes_payable"},
		"issue_shares_future_equipment":  {"cash", "common_stock"},
	}
	pairs := [][2]string{
		{"equipment_purchase_workshop", "cash_rent_workshop"},
		{"equipment_purchase_replacement", "cash_rent_workshop"},
		{"borrow_cash_future_equipment", "issue_shares_future_equipment"},
		{"borrow_cash_future_equipment", "equipment_purchase_workshop"},
		{"issue_shares_future_equipment", "equipment_purchase_workshop"},
	}
	verifyPilotBatch(t, "../../docs/EXPANSION-BATCH3.json", expected, pairs)
}

func TestPublishedExpansionBatch3MatchesReviewedPackage(t *testing.T) {
	verifyPublishedPilot(t, "../../docs/EXPANSION-BATCH3.json")
}
