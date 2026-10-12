package engine

// Stream-M08 Error Tags & Misconceptions for Module 8: Leases (ASC 842) & Stockholders' Equity (Treasury Stock).
const (
	TagM08LeaseExpensedAtCommencement         = "m08_lease_expensed_at_commencement"
	TagM08ROUAssetRecordedAsEquity            = "m08_rou_asset_recorded_as_equity"
	TagM08IgnoringLeaseInterest               = "m08_ignoring_lease_interest"
	TagM08TreasuryStockRecordedAsAsset        = "m08_treasury_stock_recorded_as_asset"
	TagM08GainRecordedOnTreasuryReissuance    = "m08_gain_recorded_on_treasury_reissuance"
	TagM08RetainedEarningsDebitedOnReissuance = "m08_retained_earnings_debited_on_reissuance"
)

// Stream-M08 Family ID constants for candidate lease & equity questions.
const (
	FamilyM08LeaseCommencementROUAsset = "m08_lse_commencement_rou_asset"
	FamilyM08LeasePeriodicPayment      = "m08_lse_periodic_lease_payment"
	FamilyM08TreasuryStockRepurchase   = "m08_lse_treasury_stock_repurchase"
	FamilyM08TreasuryStockReissuance   = "m08_lse_treasury_stock_reissuance"
)
