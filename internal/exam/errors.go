package exam

import (
	"strings"

	"github.com/trustdan/acctg-practice/internal/engine"
)

// DistractorExplanation returns the accounting misconception name and pedagogical remediation
// for a diagnostic distractor error tag.
func DistractorExplanation(tag string) (misconception string, remediation string) {
	switch tag {
	case engine.TagRevenueRecognizedPrematurely:
		return "Premature Revenue Recognition",
			"Recognized revenue upon receiving cash before services were delivered. Under accrual accounting and the revenue recognition principle, revenue is recognized when performance obligations are satisfied, not when cash is received. Customer advances are liabilities (Unearned Revenue) until earned."

	case engine.TagDuplicateRevenueOnCollection:
		return "Duplicate Revenue on Collection",
			"Recognized revenue again upon collecting cash for a prior credit sale. Revenue was already recognized when the service was billed (Dr Accounts Receivable, Cr Service Revenue). Cash collection settles the existing receivable (Dr Cash, Cr Accounts Receivable) with no second revenue impact."

	case engine.TagCashRecordedOnEarningAdvance:
		return "Cash Recorded on Fulfilling Advance",
			"Recorded cash inflow when fulfilling previously unearned revenue. Cash was received in an earlier period; fulfilling the obligation satisfies the liability (Dr Unearned Revenue, Cr Service Revenue) with no new cash flow."

	case engine.TagExpenseRecordedOnPrepaidPurchase:
		return "Prepaid Asset Expensed at Acquisition",
			"Expensed multi-period coverage immediately at payment. Under the matching principle, expenditures providing future economic benefits must be capitalized as assets (e.g. Prepaid Insurance or Prepaid Rent) and expensed systematically as benefits are consumed."

	case engine.TagPrepaidNotExpensedOnConsumption:
		return "Prepaid Benefit Not Expensed on Consumption",
			"Failed to record the adjusting expense as prepaid benefits expired. As time passes and coverage is used, an adjusting entry must transfer the consumed cost from Prepaid asset to Operating Expense."

	case engine.TagExpenseRecordedOnEquipmentPurchase:
		return "Capital Asset Expensed Immediately",
			"Expensed durable equipment immediately as an operating expense. Purchasing equipment is a capital expenditure acquiring a productive asset (asset exchange Dr Equipment, Cr Cash), not an immediate expense of current operations."

	case engine.TagRevenueRecordedOnBorrowing:
		return "Bank Borrowing Recorded as Revenue",
			"Treated borrowed bank loan proceeds as earned revenue. Borrowing incurs a contractual obligation to repay the lender (Notes Payable liability), not an earnings generation from operations."

	case engine.TagExpenseRecordedOnLoanRepayment:
		return "Principal Repayment Expensed",
			"Treated repayment of loan principal as an operating expense. Paying down loan principal satisfies an existing liability (Dr Notes Payable, Cr Cash); only interest incurred is an expense."

	case engine.TagRevenueRecordedOnShareIssue:
		return "Owner Capital Recorded as Revenue",
			"Treated owner investment as earned revenue. Contributions from stockholders represent contributed equity capital (Dr Cash, Cr Common Stock), not earnings from business operations."

	case engine.TagExpenseRecordedOnDividend:
		return "Dividends Expensed as Operating Cost",
			"Treated cash dividends as an operating expense. Dividends are distributions of accumulated profits to stockholders that reduce retained equity directly, not expenses incurred to generate revenue."

	case engine.TagCashRecordedOnDividendDeclaration:
		return "Cash Deducted on Dividend Declaration",
			"Recorded cash outflow on the declaration date. Declaring a dividend establishes a legal liability (Dr Dividends, Cr Dividends Payable); cash is not disbursed until the later payment date."

	case engine.TagDividendsDebitedOnPayment:
		return "Dividends Debited Twice on Payment",
			"Debited dividends account when paying a previously declared dividend. Equity was already reduced upon declaration; paying the dividend satisfies the liability (Dr Dividends Payable, Cr Cash)."

	case engine.TagCashRecordedWhenUncollected:
		return "Cash Recorded Before Collection",
			"Recorded cash for services billed on account. No cash arrives when the customer is invoiced; the company records Accounts Receivable and collects cash later."

	case engine.TagCashRecordedOnPrepaidExpiration:
		return "Cash Recorded on Prepaid Expiration",
			"Recorded a cash payment when prepaid coverage expired. Cash left when the policy was bought; expiration moves cost from Prepaid Insurance to Insurance Expense with no cash effect."

	case engine.TagRevenueDeferredWhenEarned:
		return "Earned Revenue Recorded as Unearned",
			"Treated revenue for completed work as an obligation. Once the service is performed, the company owes the customer nothing further and records Service Revenue."

	case engine.TagAdvanceConfusedWithReceivable:
		return "Advance Confused with Receivable",
			"Treated a customer prepayment as a customer debt. When a customer pays before the work, the company owes future service (Unearned Revenue); Accounts Receivable arises only when the customer owes the company."

	case engine.TagPayableRecordedForCashPayment:
		return "Payable Recorded for a Cash Payment",
			"Recorded an amount owed when the company paid cash immediately. Accounts Payable arises only when payment is deferred; an immediate payment credits Cash."

	case engine.TagEquationEffectMissed:
		return "Equation Effect Misread",
			"Described the event as an exchange that leaves totals or equity unchanged when a liability or equity account also changed. Place each account in the entry under Assets, Liabilities, or Equity and total both sides."

	case engine.TagReversedSides:
		return "Debit / Credit Inversion",
			"Reversed the debit and credit sides of the entry. In standard accounting, Debit is always on the Left and Credit is always on the Right. Assets/expenses increase with Debits; liabilities/equity/revenue increase with Credits."

	case engine.TagUnbalancedEntry:
		return "Unbalanced Entry",
			"Debits did not equal credits in the selected entry. The fundamental invariant of double-entry bookkeeping requires total debit dollars to equal total credit dollars for every transaction."

	case engine.TagWrongAccount:
		return "Incorrect Account Selection",
			"Selected an account that does not participate in the economic reality of this transaction."

	case engine.TagWrongAmount:
		return "Incorrect Dollar Amount",
			"The chosen dollar amount does not reconcile to the transaction parameters."

	case engine.TagDirectWriteOffExpense:
		return "Direct Write-Off Violates GAAP Matching",
			"Credited Accounts Receivable directly instead of establishing an allowance. Under the allowance method, estimated losses are matched to the period of sale using the contra-asset Allowance for Doubtful Accounts."

	case engine.TagWriteOffExpensedAgain:
		return "Write-Off Expensed a Second Time",
			"Debited Bad Debt Expense upon specific write-off. Expense was already recognized when estimating uncollectibles; writing off a specific uncollectible account reduces the allowance, not net income again."

	case engine.TagStockIssuedInsteadOfBond:
		return "Debt Security Confused with Equity",
			"Recorded Common Stock when issuing bonds. Bonds are long-term contractual debt obligations requiring repayment of principal and interest (Bonds Payable liability), not equity ownership shares."

	case engine.TagGoodwillExpensed:
		return "Goodwill Expensed at Acquisition",
			"Expensed Goodwill immediately at acquisition. Goodwill purchased in a business combination has future economic value and must be capitalized as an intangible asset, not expensed."

	case engine.TagCashFlowOperatingWrongDirection:
		return "Cash Flow Operating Adjustment Inversion",
			"Treated non-cash revenue as cash flow. Under the indirect Statement of Cash Flows, an increase in Accounts Receivable means revenue was recognized without collecting cash, so it must be deducted from Net Income."

	case engine.TagDiscountRecordedAsLoss:
		return "Bond Discount Expensed as Immediate Loss",
			"Treated the discount on bonds payable as an immediate loss or expense. The discount represents additional cost of debt that must be recorded in a contra-liability account (Discount on Bonds Payable) and amortized over the bond term to Interest Expense."

	case engine.TagExpensingFeasibleSoftware:
		return "Feasible Software Development Costs Expensed",
			"Expensed software development costs after technological feasibility was established. Under US GAAP (ASC 985-20), research costs before feasibility are expensed as R&D, but coding and testing costs after feasibility must be capitalized as an intangible asset."

	case engine.TagCurrentRatioWrongDirection:
		return "Current Ratio Impact Inverted",
			"Assumed that paying off current liabilities lowers liquidity. Paying accounts payable with cash reduces both Current Assets and Current Liabilities by the same amount, leaving Working Capital unchanged; for any company with a Current Ratio > 1.0, this equal reduction mathematically increases the Current Ratio."

	case engine.TagTreasuryStockRecordedAsAsset:
		return "Treasury Stock Recorded as Asset",
			"Recorded Treasury Stock as an asset or investment. A corporation cannot own an equity claim on itself; repurchased shares are contra-equity reducing total stockholders' equity, not an asset."

	case engine.TagIgnoringBondAmortization:
		return "Bond Discount Amortization Omitted",
			"Recorded interest expense equal only to the cash coupon payment, ignoring discount amortization. Under US GAAP, issuing bonds at a discount increases effective interest expense; the discount must be amortized over the bond term."

	case engine.TagRepoTreatedAsSale:
		return "Repurchase Agreement Treated as Asset Sale",
			"Treated the repurchase agreement as a sale of collateral rather than a secured borrowing. Under US GAAP economic substance, short-term repurchase agreements where the transferor retains effective control are collateralized financing transactions (Notes Payable), not asset sales."

	case engine.TagPrematureRevenueRecognition:
		return "Deferred Performance Obligation Recognized as Immediate Revenue",
			"Recognized the entire sales price as immediate revenue at point of sale, ignoring undelivered obligations. Under US GAAP multi-element revenue rules (ASC 606), consideration must be allocated to distinct performance obligations, deferring unearned revenue for future deliverables."

	case engine.TagGoodwillAmortized:
		return "Goodwill Amortized Instead of Impairment Tested",
			"Treated goodwill as subject to periodic amortization. Under US GAAP (ASC 350), goodwill has an indefinite life and cannot be amortized; it must be evaluated at least annually for impairment and written down only if fair value drops below carrying value."

	case engine.TagCryptoRecordedAsCashEquivalent:
		return "Digital Asset Treated as Cash or Cash Equivalent",
			"Treated cryptocurrency as cash or cash equivalents. Under US GAAP (ASC 350), digital assets lack statutory legal tender status and are classified as indefinite-lived intangible assets, tested for impairment without recognizing unrealized holding gains until sold."

	case engine.TagAggressiveRevenueAcceleration:
		return "Aggressive Premature Revenue Recognition",
			"Accelerated revenue recognition before fulfilling performance obligations under an aggressive policy. Prudent and conservative GAAP financial reporting requires deferring revenue to Unearned Revenue until goods or services are delivered, preserving earnings quality."

	case engine.TagAggregatedForecastAssumption:
		return "Crude Aggregated Top-Line Forecast Assumption",
			"Relied on crude aggregated percentage growth rather than disaggregated unit-level drivers. Disaggregated operational drivers (unit volume, pricing power, store count) provide structural predictive accuracy, whereas aggregated top-line growth hides volume vs price dynamics."

	default:
		clean := strings.ReplaceAll(tag, "_", " ")
		clean = strings.Title(clean)
		return clean, "Review the core economic reality and normal account directions for this transaction family."
	}
}

// ConceptDisplayName returns a student-friendly display title for a concept ID.
func ConceptDisplayName(conceptID string) string {
	switch conceptID {
	case "cash_vs_revenue":
		return "Cash Timing vs. Revenue Recognition"
	case "cash_classification":
		return "Cash Classification & Normal Side"
	case "service_revenue_classification":
		return "Service Revenue & Equity Impact"
	case "debit_credit_translation":
		return "Debit/Credit Direction Mechanics"
	case "customer_advance_liability":
		return "Customer Advances & Unearned Revenue"
	case "unearned_revenue_settlement":
		return "Fulfilling Unearned Revenue"
	case "accrued_receivables":
		return "Accounts Receivable on Credit Sales"
	case "receivable_collection":
		return "Collecting Accounts Receivable"
	case "prepaid_asset_capitalization":
		return "Prepaid Asset Capitalization"
	case "prepaid_expense_adjustment":
		return "Prepaid Expense Adjustment"
	case "capital_expenditures":
		return "Equipment & Capital Expenditures"
	case "debt_borrowing":
		return "Bank Loans & Notes Payable"
	case "debt_principal_repayment":
		return "Principal Repayment vs. Expense"
	case "equity_capital_contributions":
		return "Common Stock & Owner Contributions"
	case "dividend_accounting":
		return "Dividends vs. Operating Expenses"
	case "dividend_declaration_vs_payment":
		return "Dividend Declaration vs. Payment"
	case "accounting_equation_invariants":
		return "Accounting Equation Effects (A = L + E)"
	case "allowance_write_off":
		return "Specific Receivable Write-Off (Allowance Method)"
	case "bond_discount_issuance":
		return "Issuing Bonds at a Discount (Contra-Liability)"
	case "software_capitalization":
		return "Software Development Cost Capitalization (ASC 985-20)"
	case "current_ratio_liquidity":
		return "Current Ratio & Working Capital Liquidity Analysis"
	case "treasury_stock_accounting":
		return "Treasury Stock & Contra-Equity Mechanics"
	case "bond_interest_expense":
		return "Bond Interest Expense & Discount Amortization"
	case "repurchase_agreements":
		return "Repurchase Agreements & Collateralized Borrowing (Repo 105)"
	case "multi_element_revenue_recognition":
		return "Multi-Element Revenue Allocation (ASC 606)"
	case "goodwill_impairment":
		return "Goodwill Impairment Testing (ASC 350)"
	case "digital_asset_accounting":
		return "Digital Asset & Cryptocurrency Accounting"
	case "reporting_strategy":
		return "Aggressive vs. Conservative Reporting Strategy"
	case "financial_forecasting":
		return "Disaggregated Financial Statement Forecasting"
	default:
		clean := strings.ReplaceAll(conceptID, "_", " ")
		return strings.Title(clean)
	}
}
