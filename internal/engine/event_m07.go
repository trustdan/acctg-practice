package engine

// Stream-M07 Error Tags & Misconceptions for Module 7: Liabilities, Bond Pricing & Effective Interest Amortization.
const (
	TagM07PremiumRecordedAsGain              = "m07_premium_recorded_as_gain"
	TagM07DiscountRecordedAsLoss             = "m07_discount_recorded_as_loss"
	TagM07StraightLineUsedInsteadOfEffective = "m07_straight_line_used_instead_of_effective"
	TagM07IgnoringBondAmortization           = "m07_ignoring_bond_amortization"
	TagM07ContingencyDisclosedWithoutAccrual = "m07_contingency_disclosed_without_accrual"
	TagM07ContingentRecordedAsPayable        = "m07_contingent_recorded_as_trade_payable"
	TagM07RemoteContingencyAccrued           = "m07_remote_contingency_accrued"
)

// Stream-M07 Family ID constants for candidate bond and liability questions.
const (
	FamilyM07BondIssuanceDiscount       = "m07_bnd_issuance_discount"
	FamilyM07BondIssuancePremium        = "m07_bnd_issuance_premium"
	FamilyM07BondEffectiveInterestAmort = "m07_bnd_effective_interest_amort"
	FamilyM07BondContingentLiability    = "m07_bnd_contingent_liability"
)
