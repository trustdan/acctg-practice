package engine

// Stream-M04: Revenue Recognition (ASC 606) Event Family Identifiers.
const (
	FamilyM04VariableConsideration = "m04_variable_consideration"
	FamilyM04ContractAssetAR       = "m04_contract_asset_ar"
	FamilyM04RightOfReturn         = "m04_right_of_return"
	FamilyM04SatisfactionOverTime  = "m04_satisfaction_over_time"
)

// Stream-M04: Diagnostic Error Tags for Misconceptions in ASC 606 Revenue Recognition.
const (
	// TagM04VariableConsiderationUnconstrained is triggered when variable consideration
	// is recognized before meeting the ASC 606 constraint (probable that significant reversal will not occur).
	TagM04VariableConsiderationUnconstrained = "m04_variable_consideration_unconstrained"

	// TagM04UnconditionalRightMisclassifiedAsContractAsset is triggered when an unconditional right to
	// payment (subject only to the passage of time) is misclassified as a conditional contract asset.
	TagM04UnconditionalRightMisclassifiedAsContractAsset = "m04_unconditional_right_misclassified_as_contract_asset"

	// TagM04RevenueRecognizedWithUncertainReturn is triggered when revenue is recognized despite
	// full right of return / refund obligations during an open trial period.
	TagM04RevenueRecognizedWithUncertainReturn = "m04_revenue_recognized_with_uncertain_return"

	// TagM04CashDebitedOnOverTimeSatisfaction is triggered when cash is debited instead of relieving
	// an existing contract liability (Unearned Revenue) when satisfying obligations over time.
	TagM04CashDebitedOnOverTimeSatisfaction = "m04_cash_debited_on_over_time_satisfaction"

	// TagM04ContractLiabilityNotRelieved is triggered when earned revenue is recognized without reducing
	// the associated contract liability (Unearned Revenue).
	TagM04ContractLiabilityNotRelieved = "m04_contract_liability_not_relieved"
)
