package engine

import (
	"fmt"

	"github.com/trustdan/acctg-practice/internal/domain"
)

// EventRule defines the contract for deterministic rule implementations.
type EventRule interface {
	DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error)
	Explain(params map[string]int64) EventExplanation
	DiagnoseError(candidate domain.Entry, params map[string]int64) (errorTags []string, feedback string)
}

func getAmount(params map[string]int64) (domain.Money, error) {
	amtVal, ok := params["amount_minor_units"]
	if !ok || amtVal <= 0 {
		return 0, fmt.Errorf("missing or non-positive amount_minor_units parameter")
	}
	return domain.NewMoney(amtVal), nil
}

// 1. Cash Service: Earned and received in cash today.
type CashServiceRule struct{}

func (CashServiceRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (CashServiceRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Completed service and received %s cash today.", amt.FormatDollars()),
		DebitRationale:  "Cash is an Asset that increased; Assets increase by Debit.",
		CreditRationale: "Service Revenue increases Equity; Revenue increases by Credit.",
		EquationSummary: "Assets increase (+Cash); Equity increases (+Service Revenue).",
	}
}

func (CashServiceRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	return nil, ""
}

// 2. Service on Credit: Earned today, billed to customer (unpaid).
type ServiceOnCreditRule struct{}

func (ServiceOnCreditRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "accounts_receivable", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (ServiceOnCreditRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Completed %s service on credit; customer was invoiced.", amt.FormatDollars()),
		DebitRationale:  "Accounts Receivable is an Asset (claim to future cash) that increased; Assets increase by Debit.",
		CreditRationale: "Service Revenue increases Equity as performance is complete; Revenue increases by Credit.",
		EquationSummary: "Assets increase (+Accounts Receivable); Equity increases (+Service Revenue).",
	}
}

func (ServiceOnCreditRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "cash" {
			return []string{TagCashRecordedWhenUncollected}, "No cash was received today; the customer was billed on credit, creating Accounts Receivable."
		}
	}
	return nil, ""
}

// 3. Customer Advance: Cash received today for future work (unearned).
type CustomerAdvanceRule struct{}

func (CustomerAdvanceRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "unearned_revenue", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (CustomerAdvanceRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Customer paid %s cash in advance for future service.", amt.FormatDollars()),
		DebitRationale:  "Cash is an Asset that increased; Assets increase by Debit.",
		CreditRationale: "Unearned Revenue is a Liability representing the obligation to perform work; Liabilities increase by Credit.",
		EquationSummary: "Assets increase (+Cash); Liabilities increase (+Unearned Revenue); Equity is unaffected.",
	}
}

func (CustomerAdvanceRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "service_revenue" {
			return []string{TagRevenueRecognizedPrematurely}, "Revenue cannot be recognized yet because the service has not been performed. The advance creates a liability (Unearned Revenue)."
		}
	}
	return nil, ""
}

// 4. Collect Receivable: Cash collected on a previously recorded receivable.
type CollectReceivableRule struct{}

func (CollectReceivableRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "accounts_receivable", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (CollectReceivableRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Collected %s cash on an invoice previously billed and recognized.", amt.FormatDollars()),
		DebitRationale:  "Cash is an Asset that increased; Assets increase by Debit.",
		CreditRationale: "Accounts Receivable is an Asset that decreased upon collection; Assets decrease by Credit.",
		EquationSummary: "Asset exchange: Cash increases (+), Accounts Receivable decreases (-); Total Assets and Equity are unchanged.",
	}
}

func (CollectReceivableRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "service_revenue" {
			return []string{TagDuplicateRevenueOnCollection}, "Revenue was already recorded when the service was billed last month. Recording revenue again now causes duplicate revenue recognition."
		}
	}
	return nil, ""
}

// 5. Earn Advance: Services completed today for work previously paid in advance.
type EarnAdvanceRule struct{}

func (EarnAdvanceRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "unearned_revenue", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (EarnAdvanceRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Completed %s of services previously collected in advance.", amt.FormatDollars()),
		DebitRationale:  "Unearned Revenue is a Liability that decreased as work was fulfilled; Liabilities decrease by Debit.",
		CreditRationale: "Service Revenue is recognized now that performance is complete; Revenue increases by Credit.",
		EquationSummary: "Liabilities decrease (-Unearned Revenue); Equity increases (+Service Revenue); Cash is unaffected today.",
	}
}

func (EarnAdvanceRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "cash" {
			return []string{TagCashRecordedOnEarningAdvance}, "No cash changed hands today; cash was received in the prior period. You are fulfilling the liability (Unearned Revenue)."
		}
	}
	return nil, ""
}

// 6. Cash Operating Expense (e.g. Rent): Current period expense paid in cash.
type CashExpenseRule struct {
	ExpenseAccount domain.AccountID
}

func (r CashExpenseRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	acc := r.ExpenseAccount
	if acc == "" {
		acc = "rent_expense"
	}
	return domain.NewEntry(
		domain.Posting{AccountID: acc, Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (r CashExpenseRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Paid %s cash for current period operating expenses.", amt.FormatDollars()),
		DebitRationale:  "Expense accounts reduce Equity and have a normal Debit balance; Expenses increase by Debit.",
		CreditRationale: "Cash is an Asset that decreased; Assets decrease by Credit.",
		EquationSummary: "Assets decrease (-Cash); Equity decreases (-Rent Expense).",
	}
}

func (CashExpenseRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	return nil, ""
}

// 7. Prepaid Purchase: Cash paid in advance for future benefit.
type PrepaidPurchaseRule struct{}

func (PrepaidPurchaseRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "prepaid_insurance", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (PrepaidPurchaseRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Paid %s cash in advance for insurance policy covering future periods.", amt.FormatDollars()),
		DebitRationale:  "Prepaid Insurance is an Asset (future economic benefit) that increased; Assets increase by Debit.",
		CreditRationale: "Cash is an Asset that decreased; Assets decrease by Credit.",
		EquationSummary: "Asset exchange: Prepaid Insurance increases (+), Cash decreases (-); Total Assets and Equity are unchanged.",
	}
}

func (PrepaidPurchaseRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "insurance_expense" {
			return []string{TagExpenseRecordedOnPrepaidPurchase}, "The policy covers future periods; it cannot be expensed immediately. It must be recorded as an asset (Prepaid Insurance)."
		}
	}
	return nil, ""
}

// 8. Prepaid Consumption: Using up prepaid benefits over time.
type PrepaidConsumptionRule struct{}

func (PrepaidConsumptionRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "insurance_expense", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "prepaid_insurance", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (PrepaidConsumptionRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Recognized %s of expired insurance coverage from prepaids.", amt.FormatDollars()),
		DebitRationale:  "Insurance Expense reflects the used benefit and reduces Equity; Expenses increase by Debit.",
		CreditRationale: "Prepaid Insurance is an Asset that decreased as coverage expired; Assets decrease by Credit.",
		EquationSummary: "Assets decrease (-Prepaid Insurance); Equity decreases (-Insurance Expense); Cash is unaffected.",
	}
}

func (PrepaidConsumptionRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "cash" {
			return []string{TagCashRecordedOnPrepaidExpiration}, "No cash is paid when insurance expires; the asset Prepaid Insurance was already purchased earlier."
		}
	}
	return nil, ""
}

// 9. Equipment Purchase Cash: Capital expenditure for cash.
type EquipmentPurchaseCashRule struct{}

func (EquipmentPurchaseCashRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "equipment", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (EquipmentPurchaseCashRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Purchased %s of equipment for cash.", amt.FormatDollars()),
		DebitRationale:  "Equipment is a long-term productive Asset that increased; Assets increase by Debit.",
		CreditRationale: "Cash is an Asset that decreased; Assets decrease by Credit.",
		EquationSummary: "Asset exchange: Equipment increases (+), Cash decreases (-); Total Assets and Equity are unchanged.",
	}
}

func (EquipmentPurchaseCashRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "rent_expense" || p.AccountID == "insurance_expense" {
			return []string{TagExpenseRecordedOnEquipmentPurchase}, "Purchasing equipment is a capital expenditure (acquiring an asset), not an immediate operating expense."
		}
	}
	return nil, ""
}

// 10. Borrow Cash on Note: Financing cash receipt via promissory note.
type BorrowCashRule struct{}

func (BorrowCashRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "notes_payable", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (BorrowCashRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Borrowed %s cash and signed a promissory note.", amt.FormatDollars()),
		DebitRationale:  "Cash is an Asset that increased; Assets increase by Debit.",
		CreditRationale: "Notes Payable is a Liability representing the obligation to repay; Liabilities increase by Credit.",
		EquationSummary: "Assets increase (+Cash); Liabilities increase (+Notes Payable); Equity is unaffected (borrowing is NOT revenue).",
	}
}

func (BorrowCashRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "service_revenue" {
			return []string{TagRevenueRecordedOnBorrowing}, "Borrowing money creates a liability (debt to be repaid), not earned revenue."
		}
	}
	return nil, ""
}

// 11. Repay Note Principal: Cash paid to reduce loan principal.
type RepayNotePrincipalRule struct{}

func (RepayNotePrincipalRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "notes_payable", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (RepayNotePrincipalRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Repaid %s cash of loan principal on notes payable.", amt.FormatDollars()),
		DebitRationale:  "Notes Payable is a Liability that decreased upon payment; Liabilities decrease by Debit.",
		CreditRationale: "Cash is an Asset that decreased; Assets decrease by Credit.",
		EquationSummary: "Assets decrease (-Cash); Liabilities decrease (-Notes Payable); Equity is unaffected (principal repayment is NOT an expense).",
	}
}

func (RepayNotePrincipalRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "rent_expense" || p.AccountID == "insurance_expense" {
			return []string{TagExpenseRecordedOnLoanRepayment}, "Repaying loan principal settles a liability; principal reduction is not an expense."
		}
	}
	return nil, ""
}

// 12. Issue Shares: Owners invest cash in exchange for common stock.
type IssueSharesRule struct{}

func (IssueSharesRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "common_stock", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (IssueSharesRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Issued common stock for %s cash investment.", amt.FormatDollars()),
		DebitRationale:  "Cash is an Asset that increased; Assets increase by Debit.",
		CreditRationale: "Common Stock represents contributed owner capital; Equity increases by Credit.",
		EquationSummary: "Assets increase (+Cash); Equity increases (+Common Stock); Contributed capital is NOT revenue.",
	}
}

func (IssueSharesRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "service_revenue" {
			return []string{TagRevenueRecordedOnShareIssue}, "Owner investment provides capital, not earned revenue from goods or services."
		}
	}
	return nil, ""
}

// 13. Cash Dividends: Cash distribution to stockholders.
type DividendCashRule struct{}

func (DividendCashRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "dividends", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (DividendCashRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Paid %s cash dividends to shareholders.", amt.FormatDollars()),
		DebitRationale:  "Dividends represent a direct distribution of earnings and reduce Equity; Dividends increase by Debit.",
		CreditRationale: "Cash is an Asset that decreased; Assets decrease by Credit.",
		EquationSummary: "Assets decrease (-Cash); Equity decreases (-Dividends); Dividends are NOT an expense.",
	}
}

func (DividendCashRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "rent_expense" || p.AccountID == "insurance_expense" {
			return []string{TagExpenseRecordedOnDividend}, "Dividends are distributions of profit to owners, not an operating expense of doing business."
		}
	}
	return nil, ""
}

// 14. Declare Dividend: Board declares dividend payable in future.
// Incurs liability (Dividends Payable) and reduces equity (Dividends). No cash moves.
type DeclareDividendRule struct{}

func (DeclareDividendRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "dividends", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "dividends_payable", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (DeclareDividendRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Declared %s dividend payable to shareholders at future date.", amt.FormatDollars()),
		DebitRationale:  "Dividends account is Debited, reducing Stockholders' Equity.",
		CreditRationale: "Dividends Payable is a Liability that increased; Liabilities increase by Credit.",
		EquationSummary: "Liabilities increase (+Dividends Payable); Equity decreases (-Dividends); No cash flow upon declaration.",
	}
}

func (DeclareDividendRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "cash" {
			return []string{TagCashRecordedOnDividendDeclaration}, "Declaration of dividends creates a legal liability (Dividends Payable); cash is not distributed until the payment date."
		}
		if p.AccountID == "rent_expense" || p.AccountID == "insurance_expense" {
			return []string{TagExpenseRecordedOnDividend}, "Dividends are distributions of profit to shareholders, not an operating expense."
		}
	}
	return nil, ""
}

// 15. Pay Dividend Payable: Distribute cash to settle previously declared dividend liability.
// Settles liability (Dividends Payable decreases) with cash (Cash decreases). Equity is unaffected.
type PayDividendPayableRule struct{}

func (PayDividendPayableRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "dividends_payable", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (PayDividendPayableRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Paid %s cash to satisfy previously declared dividend obligation.", amt.FormatDollars()),
		DebitRationale:  "Dividends Payable is a Liability that decreased; Liabilities decrease by Debit.",
		CreditRationale: "Cash is an Asset that decreased; Assets decrease by Credit.",
		EquationSummary: "Assets decrease (-Cash); Liabilities decrease (-Dividends Payable); Stockholders' Equity was already reduced on declaration date.",
	}
}

func (PayDividendPayableRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "dividends" && p.Side == domain.SideDebit {
			return []string{TagDividendsDebitedOnPayment}, "The reduction to equity was already recorded when the dividend was declared. Paying the dividend satisfies the Dividends Payable liability."
		}
		if p.AccountID == "rent_expense" || p.AccountID == "insurance_expense" {
			return []string{TagExpenseRecordedOnDividend}, "Settling a dividend liability is not an expense."
		}
	}
	return nil, ""
}

// 16. Bad Debt Allowance: Record estimated uncollectible accounts under the allowance method.
// Debits Bad Debt Expense (reducing equity), Credits Allowance for Doubtful Accounts (contra-asset).
type BadDebtAllowanceRule struct{}

func (BadDebtAllowanceRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "bad_debt_expense", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "allowance_for_doubtful_accounts", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (BadDebtAllowanceRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Recorded %s estimated uncollectible accounts under the allowance method.", amt.FormatDollars()),
		DebitRationale:  "Bad Debt Expense is an Expense that reduces Stockholders' Equity (-E) in the period sales occur under matching.",
		CreditRationale: "Allowance for Doubtful Accounts is a Contra-Asset that reduces net Accounts Receivable (-A); Contra-Assets increase by Credit.",
		EquationSummary: "Assets decrease (-Allowance for Doubtful Accounts); Equity decreases (-Bad Debt Expense); Liabilities are unchanged.",
	}
}

func (BadDebtAllowanceRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "accounts_receivable" && p.Side == domain.SideCredit {
			return []string{TagDirectWriteOffExpense}, "The direct write-off method violates GAAP matching. Under the allowance method, credit Allowance for Doubtful Accounts (a contra-asset), not Accounts Receivable directly."
		}
		if p.AccountID == "cash" {
			return []string{TagCashRecordedWhenUncollected}, "Estimating uncollectible accounts is an adjusting entry; no cash is exchanged."
		}
	}
	return nil, ""
}

// 17. Write Off Receivable: Specific uncollectible account written off against the established allowance.
// Debits Allowance for Doubtful Accounts (reducing contra-asset), Credits Accounts Receivable.
// Net realizable value of A/R and Net Income are completely unchanged.
type WriteOffReceivableRule struct{}

func (WriteOffReceivableRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "allowance_for_doubtful_accounts", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "accounts_receivable", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (WriteOffReceivableRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Wrote off %s specific uncollectible customer balance against the allowance.", amt.FormatDollars()),
		DebitRationale:  "Allowance for Doubtful Accounts is Debited, reducing the previously established contra-asset reserve.",
		CreditRationale: "Accounts Receivable is Credited, removing the specific uncollectible customer account.",
		EquationSummary: "Total Assets change is $0: Gross A/R decreases and contra-asset decreases by the exact same amount. Net Income and Equity are unaffected.",
	}
}

func (WriteOffReceivableRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "bad_debt_expense" && p.Side == domain.SideDebit {
			return []string{TagWriteOffExpensedAgain}, "Bad Debt Expense was already recognized when estimating uncollectibles. Writing off a specific uncollectible account reduces the allowance, not net income again."
		}
	}
	return nil, ""
}

// 18. Issue Bonds at Par: Issued long-term debt securities at face value for cash.
// Debits Cash (Asset increases), Credits Bonds Payable (Liability increases).
type IssueBondsParRule struct{}

func (IssueBondsParRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "bonds_payable", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (IssueBondsParRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Issued %s of long-term bonds at face value (par) for cash.", amt.FormatDollars()),
		DebitRationale:  "Cash is an Asset that increased; Assets increase by Debit.",
		CreditRationale: "Bonds Payable is a long-term Liability that increased; Liabilities increase by Credit.",
		EquationSummary: "Assets increase (+Cash); Liabilities increase (+Bonds Payable); Equity is unaffected.",
	}
}

func (IssueBondsParRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "common_stock" {
			return []string{TagStockIssuedInsteadOfBond}, "Bonds are debt obligations that must be repaid, not equity ownership stock."
		}
		if p.AccountID == "service_revenue" {
			return []string{TagRevenueRecordedOnBorrowing}, "Issuing bonds is borrowing from investors, not earned revenue."
		}
	}
	return nil, ""
}

// 19. Acquire Goodwill: Business acquisition recognizing goodwill for excess purchase price.
// Debits Goodwill (Intangible Asset), Credits Cash (Asset decreases).
type AcquireGoodwillRule struct{}

func (AcquireGoodwillRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "goodwill", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (AcquireGoodwillRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Acquired business unit recognizing %s of Goodwill for excess purchase price.", amt.FormatDollars()),
		DebitRationale:  "Goodwill is an indefinite-lived intangible Asset representing excess purchase price; Assets increase by Debit.",
		CreditRationale: "Cash is an Asset that decreased; Assets decrease by Credit.",
		EquationSummary: "Asset swap: Goodwill increases (+A), Cash decreases (-A); Net Assets, Liabilities, and Equity remain unchanged.",
	}
}

func (AcquireGoodwillRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "rent_expense" || p.AccountID == "insurance_expense" || p.AccountID == "bad_debt_expense" {
			return []string{TagGoodwillExpensed}, "Goodwill acquired in a business combination has future economic value and must be capitalized as an asset, not expensed at acquisition."
		}
	}
	return nil, ""
}

// 20. Credit Sales Cash Flow: Service on credit with Statement of Cash Flows operating analysis.
// Debits Accounts Receivable, Credits Service Revenue.
// In the indirect SCF, the increase in Accounts Receivable is deducted from Net Income.
type CreditSalesCashFlowRule struct{}

func (CreditSalesCashFlowRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "accounts_receivable", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (CreditSalesCashFlowRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Completed %s service on credit; under indirect Statement of Cash Flows, increase in A/R is deducted from Net Income.", amt.FormatDollars()),
		DebitRationale:  "Accounts Receivable is an Asset that increased; Assets increase by Debit.",
		CreditRationale: "Service Revenue increases Equity; Revenue increases by Credit.",
		EquationSummary: "Assets increase (+Accounts Receivable); Equity increases (+Service Revenue). On the Statement of Cash Flows (indirect method), this +A/R is deducted from Net Income because revenues were recognized without cash collection.",
	}
}

func (CreditSalesCashFlowRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "cash" {
			return []string{TagCashRecordedWhenUncollected}, "No cash was received today; services were billed on credit, creating Accounts Receivable."
		}
	}
	return nil, ""
}

// 21. Issue Bonds at a Discount: Issued long-term bonds when market rate > coupon rate.
// Debits Cash (Asset increases by cash proceeds),
// Debits Discount on Bonds Payable (Contra-Liability increases, debit normal balance),
// Credits Bonds Payable (Liability increases by par/face value).
type IssueBondsDiscountRule struct{}

func (IssueBondsDiscountRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	discountVal, ok := params["discount_minor_units"]
	if !ok || discountVal <= 0 {
		discountVal = amt.Cents() * 5 / 100
	}
	discount := domain.NewMoney(discountVal)
	cashVal, ok := params["cash_minor_units"]
	if !ok || cashVal <= 0 {
		cashVal = amt.Cents() - discountVal
	}
	cash := domain.NewMoney(cashVal)

	return domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: cash},
		domain.Posting{AccountID: "discount_on_bonds_payable", Side: domain.SideDebit, Amount: discount},
		domain.Posting{AccountID: "bonds_payable", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (IssueBondsDiscountRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	discountVal := params["discount_minor_units"]
	if discountVal <= 0 {
		discountVal = amt.Cents() * 5 / 100
	}
	discount := domain.NewMoney(discountVal)
	cash := domain.NewMoney(amt.Cents() - discountVal)

	return EventExplanation{
		Summary:         fmt.Sprintf("Issued %s face value bonds at a discount, receiving %s cash and recording %s Discount on Bonds Payable.", amt.FormatDollars(), cash.FormatDollars(), discount.FormatDollars()),
		DebitRationale:  fmt.Sprintf("Cash increases by %s (+A); Discount on Bonds Payable is a Contra-Liability that increases by Debit (%s) to offset bonds to net carrying value.", cash.FormatDollars(), discount.FormatDollars()),
		CreditRationale: fmt.Sprintf("Bonds Payable is a Liability that increases by %s (+L) for the face maturity obligation.", amt.FormatDollars()),
		EquationSummary: fmt.Sprintf("Assets increase by %s (+Cash); Net Liabilities increase by %s (+Bonds Payable %s, -Discount %s); Equity is unaffected.", cash.FormatDollars(), cash.FormatDollars(), amt.FormatDollars(), discount.FormatDollars()),
	}
}

func (IssueBondsDiscountRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "interest_expense" || p.AccountID == "rent_expense" {
			return []string{TagDiscountRecordedAsLoss}, "The bond discount represents additional borrowing cost amortized over the bond term as interest expense, not recognized as an immediate loss or expense upon issuance."
		}
		if p.AccountID == "common_stock" {
			return []string{TagStockIssuedInsteadOfBond}, "Bonds are debt obligations that must be repaid, not equity ownership stock."
		}
	}
	return nil, ""
}

// 22. Capitalize Software Feasibility: Software coding costs incurred after technological feasibility.
// Debits Capitalized Software Costs (Intangible Asset increases), Credits Cash (Asset decreases).
type CapitalizeSoftwareFeasibilityRule struct{}

func (CapitalizeSoftwareFeasibilityRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "capitalized_software_costs", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (CapitalizeSoftwareFeasibilityRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Capitalized %s of software development costs incurred after technological feasibility.", amt.FormatDollars()),
		DebitRationale:  "Capitalized Software Costs is an intangible Asset that increased; Assets increase by Debit under US GAAP ASC 985-20.",
		CreditRationale: "Cash is an Asset that decreased; Assets decrease by Credit.",
		EquationSummary: "Asset swap: Capitalized Software Costs increases (+A), Cash decreases (-A); Total Assets, Liabilities, and Equity remain unchanged.",
	}
}

func (CapitalizeSoftwareFeasibilityRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "rent_expense" || p.AccountID == "insurance_expense" || p.AccountID == "bad_debt_expense" {
			return []string{TagExpensingFeasibleSoftware}, "Under US GAAP (ASC 985-20), software development costs incurred after technological feasibility is established must be capitalized as an intangible asset, not expensed as R&D."
		}
	}
	return nil, ""
}

// 23. Current Ratio Analysis: Settle short-term accounts payable with cash and evaluate retail liquidity impact.
// Debits Accounts Payable (Current Liability decreases), Credits Cash (Current Asset decreases).
type CurrentRatioAnalysisRule struct{}

func (CurrentRatioAnalysisRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "accounts_payable", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (CurrentRatioAnalysisRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Paid %s cash to settle Accounts Payable; in retail liquidity analysis, paying current liabilities decreases both Current Assets and Current Liabilities equally.", amt.FormatDollars()),
		DebitRationale:  "Accounts Payable is a Current Liability that decreased; Liabilities decrease by Debit.",
		CreditRationale: "Cash is a Current Asset that decreased; Assets decrease by Credit.",
		EquationSummary: "Assets decrease (-Cash); Liabilities decrease (-Accounts Payable); Equity is unchanged. Working Capital (CA - CL) is unchanged, but for a firm with Current Ratio > 1.0, the Current Ratio increases.",
	}
}

func (CurrentRatioAnalysisRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "rent_expense" || p.AccountID == "insurance_expense" {
			return []string{TagExpenseRecordedOnLoanRepayment}, "Paying off previously invoiced Accounts Payable settles an existing liability; it is not an additional operating expense."
		}
	}
	return nil, ""
}

// 24. Repurchase Treasury Stock: Buy back company common stock from the market to hold in treasury.
// Debits Treasury Stock (Contra-Equity increases, reducing total stockholders' equity),
// Credits Cash (Asset decreases).
type RepurchaseTreasuryStockRule struct{}

func (RepurchaseTreasuryStockRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "treasury_stock", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (RepurchaseTreasuryStockRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Repurchased company shares for %s cash, holding them in treasury.", amt.FormatDollars()),
		DebitRationale:  "Treasury Stock is a contra-equity account that increased; contra-equity accounts increase by Debit and reduce total Stockholders' Equity.",
		CreditRationale: "Cash is an Asset that decreased; Assets decrease by Credit.",
		EquationSummary: "Assets decrease (-Cash); Equity decreases (-Treasury Stock); Liabilities are unchanged. Total Stockholders' Equity is reduced by the cost of repurchased shares.",
	}
}

func (RepurchaseTreasuryStockRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "goodwill" || p.AccountID == "patents" || p.AccountID == "capitalized_software_costs" || p.AccountID == "equipment" {
			return []string{TagTreasuryStockRecordedAsAsset}, "Treasury stock represents repurchased company equity (contra-equity), not an investment asset or operating equipment."
		}
		if p.AccountID == "cash" && p.Side == domain.SideDebit {
			return []string{TagReversedSides}, "Cash was paid out to repurchase shares from stockholders, not received."
		}
	}
	return nil, ""
}

// 25. Pay Bond Interest with Discount Amortization: Pay semi-annual cash coupon and amortize bond discount.
// Debits Interest Expense (Expense increases by total borrowing cost: coupon + discount amortization),
// Credits Cash (Asset decreases by contractual cash coupon paid),
// Credits Discount on Bonds Payable (Contra-Liability decreases / is amortized toward par).
type PayBondInterestRule struct{}

func (PayBondInterestRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	discountVal, ok := params["discount_minor_units"]
	if !ok || discountVal <= 0 {
		discountVal = amt.Cents() * 20 / 100
	}
	discount := domain.NewMoney(discountVal)
	cashVal, ok := params["cash_minor_units"]
	if !ok || cashVal <= 0 {
		cashVal = amt.Cents() - discountVal
	}
	cash := domain.NewMoney(cashVal)

	return domain.NewEntry(
		domain.Posting{AccountID: "interest_expense", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: cash},
		domain.Posting{AccountID: "discount_on_bonds_payable", Side: domain.SideCredit, Amount: discount},
	), nil
}

func (PayBondInterestRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	discountVal := params["discount_minor_units"]
	if discountVal <= 0 {
		discountVal = amt.Cents() * 20 / 100
	}
	discount := domain.NewMoney(discountVal)
	cash := domain.NewMoney(amt.Cents() - discountVal)

	return EventExplanation{
		Summary:         fmt.Sprintf("Paid %s cash coupon interest and amortized %s bond discount, recording %s total Interest Expense.", cash.FormatDollars(), discount.FormatDollars(), amt.FormatDollars()),
		DebitRationale:  "Interest Expense is an Expense that increases by Debit, reflecting total periodic borrowing cost (coupon payment plus discount amortization).",
		CreditRationale: "Cash is an Asset that decreased by Credit; Discount on Bonds Payable is a contra-liability that is reduced (credited) as it is amortized toward par.",
		EquationSummary: "Assets decrease (-Cash); Net Liabilities increase (+carrying value as contra-liability discount shrinks); Equity decreases (-Interest Expense). Total debits equal total credits.",
	}
}

func (PayBondInterestRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	hasDiscount := false
	for _, p := range candidate.Postings {
		if p.AccountID == "discount_on_bonds_payable" {
			hasDiscount = true
		}
	}
	if !hasDiscount {
		return []string{TagIgnoringBondAmortization}, "Under US GAAP, bonds issued at a discount accrue higher effective interest than the cash coupon paid; the discount must be amortized into Interest Expense over the bond term."
	}
	return nil, ""
}

// 26. Repo 105 Analysis: Short-term repurchase agreement accounting under economic substance.
// Debits Cash (Asset increases by borrowing proceeds),
// Credits Notes Payable (Liability increases, recognizing collateralized borrowing obligation).
type Repo105AnalysisRule struct{}

func (Repo105AnalysisRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "notes_payable", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (Repo105AnalysisRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Received %s cash under a short-term repurchase agreement, recognizing collateralized financing (Notes Payable).", amt.FormatDollars()),
		DebitRationale:  "Cash is an Asset that increased; Assets increase by Debit.",
		CreditRationale: "Notes Payable is a Liability that increased; Liabilities increase by Credit reflecting the obligation to repurchase collateral.",
		EquationSummary: "Assets increase (+Cash); Liabilities increase (+Notes Payable); Equity is unchanged. Under economic substance, transparent collateralized borrowing expands the balance sheet and preserves leverage transparency, rather than masking it as an asset sale.",
	}
}

func (Repo105AnalysisRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "service_revenue" || p.AccountID == "equipment" || p.AccountID == "goodwill" || p.AccountID == "patents" {
			return []string{TagRepoTreatedAsSale}, "Under US GAAP economic substance, short-term repurchase agreements where collateral is repurchased shortly thereafter are collateralized borrowings (Notes Payable), not true asset sales."
		}
	}
	return nil, ""
}

// 27. Apple Multi-Element Revenue: Bundled hardware and deferred multi-period software updates.
// Debits Cash (Asset increases by full proceeds),
// Credits Service Revenue (Revenue increases by 80% allocated to delivered hardware),
// Credits Unearned Revenue (Liability increases by 20% allocated to undelivered software update obligations).
type AppleMultiElementRevenueRule struct{}

func (AppleMultiElementRevenueRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	deliveredVal, ok := params["delivered_revenue_minor_units"]
	if !ok || deliveredVal <= 0 {
		deliveredVal = amt.Cents() * 80 / 100
	}
	unearnedVal, ok := params["unearned_revenue_minor_units"]
	if !ok || unearnedVal <= 0 {
		unearnedVal = amt.Cents() - deliveredVal
	}
	delivered := domain.NewMoney(deliveredVal)
	unearned := domain.NewMoney(unearnedVal)

	return domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: delivered},
		domain.Posting{AccountID: "unearned_revenue", Side: domain.SideCredit, Amount: unearned},
	), nil
}

func (AppleMultiElementRevenueRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	deliveredVal := params["delivered_revenue_minor_units"]
	if deliveredVal <= 0 {
		deliveredVal = amt.Cents() * 80 / 100
	}
	delivered := domain.NewMoney(deliveredVal)
	unearned := domain.NewMoney(amt.Cents() - deliveredVal)

	return EventExplanation{
		Summary:         fmt.Sprintf("Received %s cash for bundled hardware and software updates, recognizing %s point-of-sale revenue for delivered hardware and deferring %s in Unearned Revenue for future software updates.", amt.FormatDollars(), delivered.FormatDollars(), unearned.FormatDollars()),
		DebitRationale:  "Cash is an Asset that increased by total proceeds; Assets increase by Debit.",
		CreditRationale: "Service Revenue is credited for delivered hardware performance obligations; Unearned Revenue is a Liability credited for undelivered software obligations.",
		EquationSummary: "Assets increase (+Cash); Liabilities increase (+Unearned Revenue for deferred updates); Equity increases (+Service Revenue for delivered device). Performance obligations are properly allocated under ASC 606.",
	}
}

func (AppleMultiElementRevenueRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	hasUnearned := false
	for _, p := range candidate.Postings {
		if p.AccountID == "unearned_revenue" {
			hasUnearned = true
		}
	}
	if !hasUnearned {
		return []string{TagPrematureRevenueRecognition}, "Bundled contracts with future software deliverables require allocating consideration; recognizing 100% as immediate revenue violates GAAP performance obligation allocation."
	}
	return nil, ""
}

// 28. Goodwill Impairment: Record goodwill impairment loss when reporting unit carrying value exceeds fair value.
// Debits Loss on Impairment (Expense/Loss increases, reducing Stockholders' Equity),
// Credits Goodwill (Intangible Asset decreases).
type GoodwillImpairmentRule struct{}

func (GoodwillImpairmentRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "loss_on_impairment", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "goodwill", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (GoodwillImpairmentRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Recorded %s goodwill impairment loss under US GAAP (ASC 350).", amt.FormatDollars()),
		DebitRationale:  "Loss on Impairment is an Expense/Loss account that reduces Stockholders' Equity (-E) in the period fair value declines.",
		CreditRationale: "Goodwill is an Intangible Asset that decreased (-A); Assets decrease by Credit.",
		EquationSummary: "Assets decrease (-Goodwill); Equity decreases (-Loss on Impairment); Liabilities are unchanged. Under US GAAP, goodwill impairment losses cannot be reversed in subsequent periods.",
	}
}

func (GoodwillImpairmentRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "cash" {
			return []string{TagCashRecordedWhenUncollected}, "Goodwill impairment is a non-cash valuation write-down; no cash is exchanged."
		}
		if p.AccountID == "rent_expense" || p.AccountID == "insurance_expense" {
			return []string{TagGoodwillAmortized}, "Under US GAAP (ASC 350), goodwill has an indefinite useful life and is tested for impairment, not amortized periodically."
		}
	}
	return nil, ""
}

// 29. Digital Asset Accounting: Purchase cryptocurrency (Tesla Bitcoin) as an indefinite-lived intangible asset.
// Debits Digital Assets (Intangible Asset increases),
// Credits Cash (Asset decreases).
type DigitalAssetAccountingRule struct{}

func (DigitalAssetAccountingRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "digital_assets", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (DigitalAssetAccountingRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Purchased %s of digital assets (cryptocurrency) for cash, classified as an indefinite-lived intangible asset.", amt.FormatDollars()),
		DebitRationale:  "Digital Assets is an Intangible Asset that increased; Assets increase by Debit.",
		CreditRationale: "Cash is an Asset that decreased; Assets decrease by Credit.",
		EquationSummary: "Asset swap: Digital Assets increases (+A), Cash decreases (-A); Total Assets, Liabilities, and Equity remain unchanged. Under US GAAP (ASC 350), digital assets are tested for impairment when market value drops below carrying value, but unrealized price increases cannot be recognized as gains until sold.",
	}
}

func (DigitalAssetAccountingRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "rent_expense" || p.AccountID == "insurance_expense" || p.AccountID == "bad_debt_expense" {
			return []string{TagCryptoRecordedAsCashEquivalent}, "Under US GAAP, cryptocurrencies purchased for investment are classified as indefinite-lived intangible assets, not expensed immediately."
		}
		if p.AccountID == "service_revenue" {
			return []string{TagCryptoRecordedAsCashEquivalent}, "Purchasing digital assets for cash is an asset acquisition, not an earnings transaction."
		}
	}
	return nil, ""
}

// 30. Reporting Policy Choice: Conservative revenue timing under customer commitments with pending obligations.
// Debits Cash (Asset increases by cash collected),
// Credits Unearned Revenue (Liability increases, deferring revenue until performance obligations are satisfied).
type ReportingPolicyChoiceRule struct{}

func (ReportingPolicyChoiceRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "unearned_revenue", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (ReportingPolicyChoiceRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Received %s cash under customer commitments, adopting a conservative reporting policy by deferring revenue to Unearned Revenue until performance obligations are satisfied.", amt.FormatDollars()),
		DebitRationale:  "Cash is an Asset that increased; Assets increase by Debit.",
		CreditRationale: "Unearned Revenue is a Liability that increased; Liabilities increase by Credit reflecting the obligation to perform.",
		EquationSummary: "Assets increase (+Cash); Liabilities increase (+Unearned Revenue); Equity is unchanged. A conservative policy defers revenue recognition to protect earnings quality, whereas an aggressive policy would prematurely book revenue.",
	}
}

func (ReportingPolicyChoiceRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "service_revenue" && p.Side == domain.SideCredit {
			return []string{TagAggressiveRevenueAcceleration}, "Recognizing immediate revenue before satisfying performance obligations reflects an aggressive reporting policy that inflates current-period earnings at the expense of earnings quality."
		}
	}
	return nil, ""
}

// 31. Financial Forecasting: Pro-forma projected credit sales derived from disaggregated unit-level drivers.
// Debits Accounts Receivable (Asset increases by projected credit billings),
// Credits Service Revenue (Revenue increases Stockholders' Equity).
type FinancialForecastingRule struct{}

func (FinancialForecastingRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "accounts_receivable", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (FinancialForecastingRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Projected %s in pro-forma credit sales using disaggregated unit-level drivers (unit volume times price per unit).", amt.FormatDollars()),
		DebitRationale:  "Accounts Receivable is an Asset that increases; Assets increase by Debit.",
		CreditRationale: "Service Revenue is a Revenue account that increases Stockholders' Equity; Revenue increases by Credit.",
		EquationSummary: "Assets increase (+Accounts Receivable); Equity increases (+Service Revenue); Liabilities are unchanged. Pro-forma modeling from disaggregated drivers provides superior forecasting fidelity compared to crude aggregated growth extrapolations.",
	}
}

func (FinancialForecastingRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "cash" && p.Side == domain.SideDebit {
			return []string{TagCashRecordedWhenUncollected}, "In modeling credit sales, Accounts Receivable tracks uncollected claims; cash collections are modeled separately via collections schedules."
		}
	}
	return nil, ""
}
