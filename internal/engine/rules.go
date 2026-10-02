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
