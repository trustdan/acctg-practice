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
	FamilyBadDebtAllowance      = "bad_debt_allowance"
	FamilyWriteOffReceivable    = "write_off_receivable"
	FamilyIssueBondsPar         = "issue_bonds_par"
	FamilyIssueBondsDiscount    = "issue_bonds_discount"
	FamilyIssueBondsPremium     = "issue_bonds_premium"
	FamilyPayBondInterest       = "pay_bond_interest"
	FamilyAcquireGoodwill       = "acquire_goodwill"
	FamilyRepurchaseTreasuryStock = "repurchase_treasury_stock"
	FamilyCapitalizeSoftwareFeasibility = "capitalize_software_feasibility"
	FamilyCreditSalesCashFlow   = "credit_sales_cash_flow"
	FamilyCurrentRatioAnalysis  = "current_ratio_analysis"
	FamilyRepo105Analysis       = "repo105_analysis"
	FamilyAppleMultiElementRevenue = "apple_multielement_revenue"
	FamilyGoodwillImpairment    = "goodwill_impairment"
	FamilyDigitalAssetAccounting = "digital_asset_accounting"
	FamilyReportingPolicyChoice = "reporting_policy_choice"
	FamilyFinancialForecasting  = "financial_forecasting"

	// Wave 1: Stream-M03 (Cash Flows & Statement Analysis)
	FamilyM03OCFAccountsPayableIncrease   = "m03_ocf_ap_increase"
	FamilyM03OCFPrepaidIncrease           = "m03_ocf_prepaid_increase"
	FamilyM03OCFUnearnedRevenueIncrease   = "m03_ocf_unearned_revenue_increase"
	FamilyM03OCFQuickRatio                = "m03_ocf_quick_ratio"

	// Wave 1: Stream-M04 (Revenue Recognition ASC 606)
	FamilyM04VariableConsideration = "m04_variable_consideration"
	FamilyM04ContractAssetAR       = "m04_contract_asset_ar"
	FamilyM04RightOfReturn         = "m04_right_of_return"
	FamilyM04SatisfactionOverTime  = "m04_satisfaction_over_time"

	// Wave 1: Stream-M05 (Inventory & Cost of Goods Sold)
	FamilyM05InventoryFIFO         = "m05_inv_fifo_costing"
	FamilyM05InventoryLIFO         = "m05_inv_lifo_costing"
	FamilyM05InventoryNRVWriteDown = "m05_inv_nrv_writedown"
	FamilyM05FreightInCapitalize   = "m05_inv_freight_in_cap"
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
	FamilyBadDebtAllowance:      {},
	FamilyWriteOffReceivable:    {},
	FamilyIssueBondsPar:         {},
	FamilyIssueBondsDiscount:    {},
	FamilyIssueBondsPremium:     {},
	FamilyPayBondInterest:       {},
	FamilyAcquireGoodwill:       {},
	FamilyRepurchaseTreasuryStock: {},
	FamilyCapitalizeSoftwareFeasibility: {},
	FamilyCreditSalesCashFlow:   {},
	FamilyCurrentRatioAnalysis:  {},
	FamilyRepo105Analysis:       {},
	FamilyAppleMultiElementRevenue: {},
	FamilyGoodwillImpairment:    {},
	FamilyDigitalAssetAccounting: {},
	FamilyReportingPolicyChoice: {},
	FamilyFinancialForecasting:  {},

	// Wave 1 Streams
	FamilyM03OCFAccountsPayableIncrease: {},
	FamilyM03OCFPrepaidIncrease:         {},
	FamilyM03OCFUnearnedRevenueIncrease: {},
	FamilyM03OCFQuickRatio:              {},
	FamilyM04VariableConsideration:      {},
	FamilyM04ContractAssetAR:            {},
	FamilyM04RightOfReturn:              {},
	FamilyM04SatisfactionOverTime:       {},
	FamilyM05InventoryFIFO:              {},
	FamilyM05InventoryLIFO:              {},
	FamilyM05InventoryNRVWriteDown:       {},
	FamilyM05FreightInCapitalize:         {},
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
