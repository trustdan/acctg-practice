package bank

// Supported semantic event family IDs.
const (
	FamilyCashService           = "cash_service"
	FamilyCustomerAdvance       = "customer_advance"
	FamilyServiceOnCredit       = "service_on_credit"
	FamilyCollectReceivable     = "collect_receivable"
	FamilyEarnAdvance           = "earn_advance"
	FamilyCashRent              = "cash_rent"
	FamilyBorrowCash            = "borrow_cash"
	FamilyIssueShares           = "issue_shares"
	FamilyCashExpense           = "cash_expense"
	FamilyPrepaidPurchase       = "prepaid_purchase"
	FamilyPrepaidConsumption    = "prepaid_consumption"
	FamilyEquipmentPurchaseCash = "equipment_purchase_cash"
	FamilyRepayNotePrincipal    = "repay_note_principal"
	FamilyDividendCash          = "dividend_cash"
	FamilyDeclareDividend       = "declare_dividend"
	FamilyPayDividendPayable    = "pay_dividend_payable"
)

var supportedFamilies = map[string]struct{}{
	FamilyCashService:           {},
	FamilyCustomerAdvance:       {},
	FamilyServiceOnCredit:       {},
	FamilyCollectReceivable:     {},
	FamilyEarnAdvance:           {},
	FamilyCashRent:              {},
	FamilyBorrowCash:            {},
	FamilyIssueShares:           {},
	FamilyCashExpense:           {},
	FamilyPrepaidPurchase:       {},
	FamilyPrepaidConsumption:    {},
	FamilyEquipmentPurchaseCash: {},
	FamilyRepayNotePrincipal:    {},
	FamilyDividendCash:          {},
	FamilyDeclareDividend:       {},
	FamilyPayDividendPayable:    {},
}

// IsSupportedFamily checks whether the family ID has reviewed semantic rules in code.
func IsSupportedFamily(familyID string) bool {
	_, ok := supportedFamilies[familyID]
	return ok
}

// ReviewedFamilyRuleVersion returns the current reviewed rule version for a supported family.
// Returns (version, true) if reviewed and tested, or (0, false) if unreviewed/new.
func ReviewedFamilyRuleVersion(familyID string) (int, bool) {
	if _, ok := supportedFamilies[familyID]; ok {
		// All 13 supported initial families have verified engine rules and regression tests at rule_version 1
		return 1, true
	}
	return 0, false
}

// IsFamilyReviewedAndTested verifies that a family and rule version have reviewed engine rules
// and regression tests in the codebase. New families require explicit engine implementation and tests.
func IsFamilyReviewedAndTested(familyID string, ruleVersion int) bool {
	ver, ok := ReviewedFamilyRuleVersion(familyID)
	if !ok {
		return false
	}
	return ver == ruleVersion
}

// RequiredParameters returns the parameters that must be provided for a given family.
func RequiredParameters(familyID string) []string {
	// All initial families require amount_minor_units
	return []string{"amount_minor_units"}
}
