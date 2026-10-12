package engine

import (
	"fmt"

	"github.com/trustdan/acctg-practice/internal/domain"
)

// RegisterM03Rules registers the Stream-M03 rules into the Engine.
func RegisterM03Rules(e *Engine) {
	e.rules[FamilyM03OCFAccountsPayableIncrease] = M03OCFAccountsPayableIncreaseRule{}
	e.rules[FamilyM03OCFPrepaidIncrease] = M03OCFPrepaidIncreaseRule{}
	e.rules[FamilyM03OCFUnearnedRevenueIncrease] = M03OCFUnearnedRevenueIncreaseRule{}
	e.rules[FamilyM03OCFQuickRatio] = M03OCFQuickRatioRule{}
}

// 1. M03OCFAccountsPayableIncreaseRule: Incurring operating rent expense on account (A/P increase).
// Under the indirect Statement of Cash Flows, an increase in Accounts Payable is added back to Net Income.
type M03OCFAccountsPayableIncreaseRule struct{}

func (M03OCFAccountsPayableIncreaseRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "rent_expense", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "accounts_payable", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M03OCFAccountsPayableIncreaseRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Incurred %s of operating rent on account, increasing Accounts Payable.", amt.FormatDollars()),
		DebitRationale:  "Rent Expense is an operating expense that decreases Equity; expenses increase by Debit.",
		CreditRationale: "Accounts Payable is a Current Liability that increased; liabilities increase by Credit.",
		EquationSummary: fmt.Sprintf("Liabilities increase (+Accounts Payable %s); Equity decreases (-Rent Expense %s). On the indirect Statement of Cash Flows, the increase in Accounts Payable is ADDED back to Net Income because the expense was recognized without a cash outflow.", amt.FormatDollars(), amt.FormatDollars()),
	}
}

func (M03OCFAccountsPayableIncreaseRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "cash" {
			return []string{TagCashRecordedWhenUncollected}, "No cash was paid today; the operating expense was billed on account, creating Accounts Payable."
		}
	}
	return nil, ""
}

// 2. M03OCFPrepaidIncreaseRule: Cash paid in advance for insurance coverage (Prepaid Insurance increase).
// Under the indirect Statement of Cash Flows, an increase in Prepaid Insurance is deducted from Net Income.
type M03OCFPrepaidIncreaseRule struct{}

func (M03OCFPrepaidIncreaseRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "prepaid_insurance", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M03OCFPrepaidIncreaseRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Paid %s cash in advance for insurance coverage, increasing Prepaid Insurance.", amt.FormatDollars()),
		DebitRationale:  "Prepaid Insurance is an Asset (future economic resource) that increased; Assets increase by Debit.",
		CreditRationale: "Cash is an Asset that was disbursed; Assets decrease by Credit.",
		EquationSummary: fmt.Sprintf("Total Assets are unchanged (+Prepaid Insurance %s, -Cash %s). On the indirect Statement of Cash Flows, the increase in Prepaid Insurance is DEDUCTED from Net Income because cash was paid out without an income statement expense.", amt.FormatDollars(), amt.FormatDollars()),
	}
}

func (M03OCFPrepaidIncreaseRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "insurance_expense" {
			return []string{TagExpenseRecordedOnPrepaidPurchase}, "The insurance coverage covers future operating periods; capitalize as Prepaid Insurance rather than expensing immediately."
		}
	}
	return nil, ""
}

// 3. M03OCFUnearnedRevenueIncreaseRule: Advance cash received for future operating performance (Unearned Revenue increase).
// Under the indirect Statement of Cash Flows, an increase in Unearned Revenue is added to Net Income.
type M03OCFUnearnedRevenueIncreaseRule struct{}

func (M03OCFUnearnedRevenueIncreaseRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "unearned_revenue", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M03OCFUnearnedRevenueIncreaseRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Received %s cash in advance for operating services to be delivered next quarter.", amt.FormatDollars()),
		DebitRationale:  "Cash is an Asset received; Assets increase by Debit.",
		CreditRationale: "Unearned Revenue is a Liability representing the obligation to perform; Liabilities increase by Credit.",
		EquationSummary: fmt.Sprintf("Assets increase (+Cash %s); Liabilities increase (+Unearned Revenue %s). On the indirect Statement of Cash Flows, the increase in Unearned Revenue is ADDED to Net Income because cash was collected before revenue was earned.", amt.FormatDollars(), amt.FormatDollars()),
	}
}

func (M03OCFUnearnedRevenueIncreaseRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "service_revenue" {
			return []string{TagRevenueRecognizedPrematurely}, "Services have not been performed yet; the performance obligation is unearned and must be recorded as Unearned Revenue."
		}
	}
	return nil, ""
}

// 4. M03OCFQuickRatioRule: Cash payment to settle Accounts Payable (Quick Ratio analysis).
// When Quick Ratio > 1.0, subtracting equal amounts from quick assets and current liabilities increases the Quick Ratio.
type M03OCFQuickRatioRule struct{}

func (M03OCFQuickRatioRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "accounts_payable", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M03OCFQuickRatioRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Paid %s cash to settle supplier balances in Accounts Payable.", amt.FormatDollars()),
		DebitRationale:  "Accounts Payable is a Current Liability satisfied; Liabilities decrease by Debit.",
		CreditRationale: "Cash is an Asset disbursed; Assets decrease by Credit.",
		EquationSummary: fmt.Sprintf("Assets decrease (-Cash %s) and Liabilities decrease (-Accounts Payable %s). Total Working Capital (CA - CL) is unchanged. Because initial Quick Ratio > 1.0, subtracting equal amounts from Quick Assets and Current Liabilities increases the Quick Ratio.", amt.FormatDollars(), amt.FormatDollars()),
	}
}

func (M03OCFQuickRatioRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "rent_expense" || p.AccountID == "insurance_expense" {
			return []string{TagExpenseRecordedOnLoanRepayment}, "The expense was already recognized when supplier goods/services were billed; paying the invoice settles Accounts Payable."
		}
	}
	return nil, ""
}
