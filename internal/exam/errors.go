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
	default:
		clean := strings.ReplaceAll(conceptID, "_", " ")
		return strings.Title(clean)
	}
}
