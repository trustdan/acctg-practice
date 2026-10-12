package engine

// Known error tags for Stream-M03 (Module 3: Cash Flows & Financial Statements).
const (
	TagM03OCFPayableWrongDirection         = "m03_ocf_payable_wrong_direction"
	TagM03OCFPrepaidWrongDirection         = "m03_ocf_prepaid_wrong_direction"
	TagM03OCFUnearnedRevenueWrongDirection = "m03_ocf_unearned_revenue_wrong_direction"
	TagM03QuickRatioWrongDirection         = "m03_quick_ratio_wrong_direction"
	TagM03InventoryTreatedAsQuickAsset     = "m03_inventory_treated_as_quick_asset"
	TagM03WorkingCapitalMiscalculated      = "m03_working_capital_miscalculated"
)

// Supported semantic event family IDs for Stream-M03.
const (
	FamilyM03OCFAccountsPayableIncrease = "m03_ocf_ap_increase"
	FamilyM03OCFPrepaidIncrease         = "m03_ocf_prepaid_increase"
	FamilyM03OCFUnearnedRevenueIncrease = "m03_ocf_unearned_revenue_increase"
	FamilyM03OCFQuickRatio              = "m03_ocf_quick_ratio"
)
