package engine

import (
	"fmt"

	"github.com/trustdan/acctg-practice/internal/domain"
)

// 1. M09CookieJarReserveRule: Cookie jar reserve buffering via discretionary allowance for doubtful accounts.
// Debits Bad Debt Expense (Expense increases, reducing Equity),
// Credits Allowance for Doubtful Accounts (Contra-Asset increases, reducing net Assets).
type M09CookieJarReserveRule struct{}

func (M09CookieJarReserveRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "bad_debt_expense", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "allowance_for_doubtful_accounts", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M09CookieJarReserveRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Recorded discretionary cookie jar reserve addition of %s to the allowance for doubtful accounts.", amt.FormatDollars()),
		DebitRationale:  "Bad Debt Expense is an operating expense that increased; expenses increase by Debit, reducing Stockholders' Equity (-E).",
		CreditRationale: "Allowance for Doubtful Accounts is a contra-asset account that increased by Credit, reducing net Accounts Receivable (-A).",
		EquationSummary: "Assets decrease (-Allowance for Doubtful Accounts); Equity decreases (-Bad Debt Expense); Liabilities are unchanged. Intentionally over-accruing uncollectible reserves in prosperous periods creates a hidden 'cookie jar' reserve that management can later bleed into earnings during lean quarters, distorting earnings quality.",
	}
}

func (M09CookieJarReserveRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "accounts_receivable" && p.Side == domain.SideCredit {
			return []string{TagDirectWriteOffExpense}, "Under US GAAP, estimating future bad debts must be credited to the contra-asset Allowance for Doubtful Accounts, not directly credited to customer accounts receivable."
		}
		if p.AccountID == "service_revenue" {
			return []string{TagWrongAccount}, "Creating an allowance reserve is an expense accrual, not a direct reduction to recorded sales revenue."
		}
	}
	return nil, ""
}

// 2. M09WarrantyAccrualRule: Product warranty expense accrual under the matching principle (ASC 460).
// Debits Warranty Expense (Expense increases, reducing Equity),
// Credits Warranty Liability (Liability increases).
type M09WarrantyAccrualRule struct{}

func (M09WarrantyAccrualRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "warranty_expense", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "warranty_liability", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M09WarrantyAccrualRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Recorded estimated product warranty expense and liability of %s under the matching principle.", amt.FormatDollars()),
		DebitRationale:  "Warranty Expense is an operating expense that increased; expenses increase by Debit, reducing Stockholders' Equity (-E).",
		CreditRationale: "Warranty Liability is an accrued claim liability that increased by Credit (+L).",
		EquationSummary: "Liabilities increase (+Warranty Liability); Equity decreases (-Warranty Expense); Assets are unchanged. Under the accrual matching principle (ASC 460), estimated future warranty claims must be recognized in the same period products are sold to fairly reflect net operating margins.",
	}
}

func (M09WarrantyAccrualRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "cash" && p.Side == domain.SideCredit {
			return []string{TagM09WarrantyExpensedWhenPaid}, "Under accrual accounting and the matching principle, warranty expense is estimated and accrued at the time of sale, not deferred until cash is disbursed for future repair claims."
		}
		if p.AccountID == "accounts_payable" {
			return []string{TagWrongAccount}, "Estimated obligations for future product warranty claims must be classified in Warranty Liability rather than trade Accounts Payable."
		}
	}
	return nil, ""
}
