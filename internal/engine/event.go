package engine

import (
	"github.com/trustdan/acctg-practice/internal/domain"
)

// Known error tags for diagnostic feedback on learner mistakes.
const (
	TagCorrect                            = "correct"
	TagDuplicateRevenueOnCollection       = "duplicate_revenue_on_collection"
	TagRevenueRecognizedPrematurely       = "revenue_recognized_prematurely"
	TagCashRecordedOnEarningAdvance       = "cash_recorded_on_earning_advance"
	TagExpenseRecordedOnPrepaidPurchase   = "expense_recorded_on_prepaid_purchase"
	TagPrepaidNotExpensedOnConsumption    = "prepaid_not_expensed_on_consumption"
	TagRevenueRecordedOnBorrowing         = "revenue_recorded_on_borrowing"
	TagRevenueRecordedOnShareIssue        = "revenue_recorded_on_share_issue"
	TagExpenseRecordedOnLoanRepayment     = "expense_recorded_on_loan_repayment"
	TagExpenseRecordedOnEquipmentPurchase = "expense_recorded_on_equipment_purchase"
	TagExpenseRecordedOnDividend          = "expense_recorded_on_dividend"
	TagCashRecordedOnDividendDeclaration  = "cash_recorded_on_dividend_declaration"
	TagDividendsDebitedOnPayment          = "dividends_debited_on_payment"
	TagCashRecordedWhenUncollected        = "cash_recorded_when_uncollected"
	TagCashRecordedOnPrepaidExpiration    = "cash_recorded_on_prepaid_expiration"
	TagRevenueDeferredWhenEarned          = "revenue_deferred_when_earned"
	TagAdvanceConfusedWithReceivable      = "advance_confused_with_receivable"
	TagPayableRecordedForCashPayment      = "payable_recorded_for_cash_payment"
	TagEquationEffectMissed               = "equation_effect_missed"
	TagReversedSides                      = "reversed_sides"
	TagUnbalancedEntry                    = "unbalanced_entry"
	TagWrongAccount                       = "wrong_account"
	TagWrongAmount                        = "wrong_amount"
	TagDirectWriteOffExpense              = "direct_write_off_expense"
	TagWriteOffExpensedAgain              = "write_off_expensed_again"
	TagStockIssuedInsteadOfBond           = "stock_issued_instead_of_bond"
	TagDiscountRecordedAsLoss             = "discount_recorded_as_loss"
	TagPremiumRecordedAsGain              = "premium_recorded_as_gain"
	TagGoodwillExpensed                   = "goodwill_expensed"
	TagTreasuryStockRecordedAsAsset       = "treasury_stock_recorded_as_asset"
	TagExpensingFeasibleSoftware          = "expensing_feasible_software"
	TagCashFlowOperatingWrongDirection    = "cash_flow_operating_wrong_direction"
	TagCurrentRatioWrongDirection         = "current_ratio_wrong_direction"
	TagIgnoringBondAmortization           = "ignoring_bond_amortization"
	TagRepoTreatedAsSale                  = "repo_treated_as_sale"
	TagPrematureRevenueRecognition        = "premature_revenue_recognition"
	TagGoodwillAmortized                  = "goodwill_amortized"
	TagCryptoRecordedAsCashEquivalent     = "crypto_recorded_as_cash_equivalent"
	TagAggressiveRevenueAcceleration      = "aggressive_revenue_acceleration"
	TagAggregatedForecastAssumption       = "aggregated_forecast_assumption"
)

// TransactionEvent describes an economic event to be processed by the engine.
type TransactionEvent struct {
	FamilyID   string           `json:"family_id"`
	Parameters map[string]int64 `json:"parameters"`
}

// ProcessedEvent contains the canonical derived entry, equation effects, and semantic explanation.
type ProcessedEvent struct {
	Event       TransactionEvent     `json:"event"`
	Entry       domain.Entry         `json:"entry"`
	Equation    domain.EquationDelta `json:"equation"`
	Explanation EventExplanation     `json:"explanation"`
}

// EventExplanation provides causal, pedagogical commentary on why this entry is correct.
type EventExplanation struct {
	Summary         string `json:"summary"`
	DebitRationale  string `json:"debit_rationale"`
	CreditRationale string `json:"credit_rationale"`
	EquationSummary string `json:"equation_summary"`
}

// SemanticEvaluation reports the result of grading a candidate entry against canonical event semantics.
type SemanticEvaluation struct {
	IsCorrect bool     `json:"is_correct"`
	ErrorTags []string `json:"error_tags,omitempty"`
	Feedback  string   `json:"feedback"`
}
