package engine

// Stream-M06 Error Tags & Misconceptions for Module 6: Property, Plant & Equipment, Depreciation, and Impairments.
const (
	TagM06DirectAssetCredit               = "m06_direct_asset_credit"
	TagM06AccumulatedDeprExpensed         = "m06_accumulated_depreciation_expensed"
	TagM06TaxDepreciationRecordedOnBooks  = "m06_tax_depreciation_recorded_on_books"
	TagM06GainTreatedAsRevenue            = "m06_gain_treated_as_operating_revenue"
	TagM06DisposalAccumDeprRetained       = "m06_disposal_accum_depr_retained"
	TagM06ImpairmentWithoutRecoverability = "m06_impairment_without_recoverability_test"
	TagM06ImpairmentReversedUnderGAAP     = "m06_impairment_reversed_under_gaap"
)

// Stream-M06 Family ID constants for candidate PPE questions.
const (
	FamilyM06StraightLineDepr         = "m06_ppe_straight_line_depr"
	FamilyM06MACRSTaxTiming           = "m06_ppe_macrs_tax_timing"
	FamilyM06AssetDisposalGainLoss    = "m06_ppe_asset_disposal_gain_loss"
	FamilyM06ImpairmentRecoverability = "m06_ppe_impairment_recoverability"
)
