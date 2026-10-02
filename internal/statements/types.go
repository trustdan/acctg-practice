package statements

import (
	"fmt"
	"strings"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// CashFlowActivity categorizes cash receipts and payments.
type CashFlowActivity string

const (
	ActivityOperating CashFlowActivity = "Operating"
	ActivityInvesting CashFlowActivity = "Investing"
	ActivityFinancing CashFlowActivity = "Financing"
)

// LineItem represents a financial statement item line.
type LineItem struct {
	AccountID   domain.AccountID `json:"account_id"`
	AccountName string           `json:"account_name"`
	Amount      domain.Money     `json:"amount"`
}

// CashFlowLine represents a specific cash flow activity item.
type CashFlowLine struct {
	Description string           `json:"description"`
	Amount      domain.Money     `json:"amount"` // positive = cash inflow, negative = cash outflow
	Activity    CashFlowActivity `json:"activity"`
}

// TrialBalanceAccount holds debit/credit balance for an account in trial balance.
type TrialBalanceAccount struct {
	AccountID   domain.AccountID `json:"account_id"`
	AccountName string           `json:"account_name"`
	Category    domain.Category  `json:"category"`
	Debit       domain.Money     `json:"debit"`
	Credit      domain.Money     `json:"credit"`
}

// TrialBalance holds an unadjusted or adjusted trial balance.
type TrialBalance struct {
	Title       string                `json:"title"`
	Accounts    []TrialBalanceAccount `json:"accounts"`
	TotalDebit  domain.Money          `json:"total_debit"`
	TotalCredit domain.Money          `json:"total_credit"`
	IsBalanced  bool                  `json:"is_balanced"`
}

// FormatText formats the trial balance as a clean ASCII financial schedule.
func (tb TrialBalance) FormatText() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s\n", tb.Title))
	b.WriteString("────────────────────────────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("%-36s %16s %16s\n", "Account Title", "Debit", "Credit"))
	b.WriteString("────────────────────────────────────────────────────────────────────────\n")
	for _, acc := range tb.Accounts {
		drStr := ""
		if acc.Debit.Cents() > 0 {
			drStr = acc.Debit.FormatDollars()
		}
		crStr := ""
		if acc.Credit.Cents() > 0 {
			crStr = acc.Credit.FormatDollars()
		}
		b.WriteString(fmt.Sprintf("%-36s %16s %16s\n", acc.AccountName, drStr, crStr))
	}
	b.WriteString("────────────────────────────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("%-36s %16s %16s\n", "Totals", tb.TotalDebit.FormatDollars(), tb.TotalCredit.FormatDollars()))
	b.WriteString("════════════════════════════════════════════════════════════════════════\n")
	return b.String()
}

// IncomeStatement reports revenues, expenses, and net income over a period.
type IncomeStatement struct {
	CompanyName   string       `json:"company_name"`
	Period        string       `json:"period"`
	Revenues      []LineItem   `json:"revenues"`
	TotalRevenue  domain.Money `json:"total_revenue"`
	Expenses      []LineItem   `json:"expenses"`
	TotalExpenses domain.Money `json:"total_expenses"`
	NetIncome     domain.Money `json:"net_income"`
}

// FormatText renders the Income Statement in standard GAAP single-step/multi-step format.
func (is IncomeStatement) FormatText() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s\n", strings.ToUpper(is.CompanyName)))
	b.WriteString("INCOME STATEMENT\n")
	b.WriteString(fmt.Sprintf("%s\n", is.Period))
	b.WriteString("────────────────────────────────────────────────────────────────────────\n")
	b.WriteString("Revenues:\n")
	for _, rev := range is.Revenues {
		b.WriteString(fmt.Sprintf("  %-38s %28s\n", rev.AccountName, rev.Amount.FormatDollars()))
	}
	b.WriteString(fmt.Sprintf("    Total Revenues %49s\n", is.TotalRevenue.FormatDollars()))
	b.WriteString("\nExpenses:\n")
	for _, exp := range is.Expenses {
		b.WriteString(fmt.Sprintf("  %-38s %28s\n", exp.AccountName, exp.Amount.FormatDollars()))
	}
	b.WriteString(fmt.Sprintf("    Total Expenses %49s\n", is.TotalExpenses.FormatDollars()))
	b.WriteString("────────────────────────────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("NET INCOME %57s\n", is.NetIncome.FormatDollars()))
	b.WriteString("════════════════════════════════════════════════════════════════════════\n")
	return b.String()
}

// RetainedEarningsStatement reports the roll-forward of retained earnings.
// Beginning Retained Earnings + Net Income - Dividends Declared = Ending Retained Earnings.
type RetainedEarningsStatement struct {
	CompanyName               string       `json:"company_name"`
	Period                    string       `json:"period"`
	BeginningRetainedEarnings domain.Money `json:"beginning_retained_earnings"`
	NetIncome                 domain.Money `json:"net_income"`
	LessDividendsDeclared     domain.Money `json:"less_dividends_declared"`
	EndingRetainedEarnings    domain.Money `json:"ending_retained_earnings"`
}

// FormatText renders the Statement of Retained Earnings.
func (re RetainedEarningsStatement) FormatText() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s\n", strings.ToUpper(re.CompanyName)))
	b.WriteString("STATEMENT OF RETAINED EARNINGS\n")
	b.WriteString(fmt.Sprintf("%s\n", re.Period))
	b.WriteString("────────────────────────────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("Retained Earnings, Beginning %40s\n", re.BeginningRetainedEarnings.FormatDollars()))
	b.WriteString(fmt.Sprintf("  Add: Net Income %51s\n", re.NetIncome.FormatDollars()))
	subtotal := re.BeginningRetainedEarnings.Add(re.NetIncome)
	b.WriteString(fmt.Sprintf("  Subtotal %58s\n", subtotal.FormatDollars()))
	b.WriteString(fmt.Sprintf("  Less: Dividends Declared %42s\n", fmt.Sprintf("(%s)", re.LessDividendsDeclared.FormatDollars())))
	b.WriteString("────────────────────────────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("Retained Earnings, Ending %43s\n", re.EndingRetainedEarnings.FormatDollars()))
	b.WriteString("════════════════════════════════════════════════════════════════════════\n")
	return b.String()
}

// BalanceSheet reports assets, liabilities, and stockholders' equity as of a specific date.
type BalanceSheet struct {
	CompanyName               string       `json:"company_name"`
	AsOfDate                  string       `json:"as_of_date"`
	CurrentAssets             []LineItem   `json:"current_assets"`
	TotalCurrentAssets        domain.Money `json:"total_current_assets"`
	NonCurrentAssets          []LineItem   `json:"non_current_assets"`
	TotalNonCurrentAssets     domain.Money `json:"total_non_current_assets"`
	TotalAssets               domain.Money `json:"total_assets"`
	CurrentLiabilities        []LineItem   `json:"current_liabilities"`
	TotalCurrentLiabilities   domain.Money `json:"total_current_liabilities"`
	LongTermLiabilities       []LineItem   `json:"long_term_liabilities"`
	TotalLongTermLiabilities  domain.Money `json:"total_long_term_liabilities"`
	TotalLiabilities          domain.Money `json:"total_liabilities"`
	StockholdersEquity        []LineItem   `json:"stockholders_equity"`
	TotalEquity               domain.Money `json:"total_equity"`
	TotalLiabilitiesAndEquity domain.Money `json:"total_liabilities_and_equity"`
	IsBalanced                bool         `json:"is_balanced"`
}

// FormatText renders the classified Balance Sheet.
func (bs BalanceSheet) FormatText() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s\n", strings.ToUpper(bs.CompanyName)))
	b.WriteString("BALANCE SHEET\n")
	b.WriteString(fmt.Sprintf("As of %s\n", bs.AsOfDate))
	b.WriteString("────────────────────────────────────────────────────────────────────────\n")
	b.WriteString("ASSETS\n")
	b.WriteString("Current Assets:\n")
	for _, it := range bs.CurrentAssets {
		b.WriteString(fmt.Sprintf("  %-38s %28s\n", it.AccountName, it.Amount.FormatDollars()))
	}
	b.WriteString(fmt.Sprintf("    Total Current Assets %43s\n", bs.TotalCurrentAssets.FormatDollars()))
	if len(bs.NonCurrentAssets) > 0 {
		b.WriteString("\nProperty, Plant & Equipment:\n")
		for _, it := range bs.NonCurrentAssets {
			b.WriteString(fmt.Sprintf("  %-38s %28s\n", it.AccountName, it.Amount.FormatDollars()))
		}
		b.WriteString(fmt.Sprintf("    Total Property, Plant & Equipment %30s\n", bs.TotalNonCurrentAssets.FormatDollars()))
	}
	b.WriteString("────────────────────────────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("TOTAL ASSETS %55s\n", bs.TotalAssets.FormatDollars()))
	b.WriteString("════════════════════════════════════════════════════════════════════════\n\n")

	b.WriteString("LIABILITIES & STOCKHOLDERS' EQUITY\n")
	b.WriteString("Current Liabilities:\n")
	for _, it := range bs.CurrentLiabilities {
		b.WriteString(fmt.Sprintf("  %-38s %28s\n", it.AccountName, it.Amount.FormatDollars()))
	}
	b.WriteString(fmt.Sprintf("    Total Current Liabilities %38s\n", bs.TotalCurrentLiabilities.FormatDollars()))
	if len(bs.LongTermLiabilities) > 0 {
		b.WriteString("\nLong-Term Liabilities:\n")
		for _, it := range bs.LongTermLiabilities {
			b.WriteString(fmt.Sprintf("  %-38s %28s\n", it.AccountName, it.Amount.FormatDollars()))
		}
		b.WriteString(fmt.Sprintf("    Total Long-Term Liabilities %35s\n", bs.TotalLongTermLiabilities.FormatDollars()))
	}
	b.WriteString(fmt.Sprintf("  TOTAL LIABILITIES %50s\n", bs.TotalLiabilities.FormatDollars()))

	b.WriteString("\nStockholders' Equity:\n")
	for _, it := range bs.StockholdersEquity {
		b.WriteString(fmt.Sprintf("  %-38s %28s\n", it.AccountName, it.Amount.FormatDollars()))
	}
	b.WriteString(fmt.Sprintf("  TOTAL STOCKHOLDERS' EQUITY %41s\n", bs.TotalEquity.FormatDollars()))
	b.WriteString("────────────────────────────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("TOTAL LIABILITIES & STOCKHOLDERS' EQUITY %27s\n", bs.TotalLiabilitiesAndEquity.FormatDollars()))
	b.WriteString("════════════════════════════════════════════════════════════════════════\n")
	return b.String()
}

// CashFlowStatement reports cash flows classified into Operating, Investing, and Financing activities.
type CashFlowStatement struct {
	CompanyName             string         `json:"company_name"`
	Period                  string         `json:"period"`
	OperatingActivities     []CashFlowLine `json:"operating_activities"`
	NetOperatingCash        domain.Money   `json:"net_operating_cash"`
	InvestingActivities     []CashFlowLine `json:"investing_activities"`
	NetInvestingCash        domain.Money   `json:"net_investing_cash"`
	FinancingActivities     []CashFlowLine `json:"financing_activities"`
	NetFinancingCash        domain.Money   `json:"net_financing_cash"`
	NetChangeInCash         domain.Money   `json:"net_change_in_cash"`
	BeginningCash           domain.Money   `json:"beginning_cash"`
	EndingCash              domain.Money   `json:"ending_cash"`
	NoncashDisclosures      []string       `json:"noncash_disclosures"`
	OperatingReconciliation string         `json:"operating_reconciliation"`
}

// FormatText renders the Statement of Cash Flows.
func (cfs CashFlowStatement) FormatText() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s\n", strings.ToUpper(cfs.CompanyName)))
	b.WriteString("STATEMENT OF CASH FLOWS (DIRECT METHOD)\n")
	b.WriteString(fmt.Sprintf("%s\n", cfs.Period))
	b.WriteString("────────────────────────────────────────────────────────────────────────\n")

	b.WriteString("Cash flows from operating activities:\n")
	for _, l := range cfs.OperatingActivities {
		amtStr := l.Amount.FormatDollars()
		if l.Amount.Cents() < 0 {
			amtStr = fmt.Sprintf("(%s)", domain.NewMoney(-l.Amount.Cents()).FormatDollars())
		}
		b.WriteString(fmt.Sprintf("  %-48s %18s\n", l.Description, amtStr))
	}
	netOpStr := cfs.NetOperatingCash.FormatDollars()
	if cfs.NetOperatingCash.Cents() < 0 {
		netOpStr = fmt.Sprintf("(%s)", domain.NewMoney(-cfs.NetOperatingCash.Cents()).FormatDollars())
	}
	b.WriteString(fmt.Sprintf("    Net cash provided by (used in) operating activities %14s\n\n", netOpStr))

	b.WriteString("Cash flows from investing activities:\n")
	for _, l := range cfs.InvestingActivities {
		amtStr := l.Amount.FormatDollars()
		if l.Amount.Cents() < 0 {
			amtStr = fmt.Sprintf("(%s)", domain.NewMoney(-l.Amount.Cents()).FormatDollars())
		}
		b.WriteString(fmt.Sprintf("  %-48s %18s\n", l.Description, amtStr))
	}
	netInvStr := cfs.NetInvestingCash.FormatDollars()
	if cfs.NetInvestingCash.Cents() < 0 {
		netInvStr = fmt.Sprintf("(%s)", domain.NewMoney(-cfs.NetInvestingCash.Cents()).FormatDollars())
	}
	b.WriteString(fmt.Sprintf("    Net cash provided by (used in) investing activities %14s\n\n", netInvStr))

	b.WriteString("Cash flows from financing activities:\n")
	for _, l := range cfs.FinancingActivities {
		amtStr := l.Amount.FormatDollars()
		if l.Amount.Cents() < 0 {
			amtStr = fmt.Sprintf("(%s)", domain.NewMoney(-l.Amount.Cents()).FormatDollars())
		}
		b.WriteString(fmt.Sprintf("  %-48s %18s\n", l.Description, amtStr))
	}
	netFinStr := cfs.NetFinancingCash.FormatDollars()
	if cfs.NetFinancingCash.Cents() < 0 {
		netFinStr = fmt.Sprintf("(%s)", domain.NewMoney(-cfs.NetFinancingCash.Cents()).FormatDollars())
	}
	b.WriteString(fmt.Sprintf("    Net cash provided by (used in) financing activities %14s\n", netFinStr))

	b.WriteString("────────────────────────────────────────────────────────────────────────\n")
	netChangeStr := cfs.NetChangeInCash.FormatDollars()
	if cfs.NetChangeInCash.Cents() < 0 {
		netChangeStr = fmt.Sprintf("(%s)", domain.NewMoney(-cfs.NetChangeInCash.Cents()).FormatDollars())
	}
	b.WriteString(fmt.Sprintf("Net increase (decrease) in cash %38s\n", netChangeStr))
	b.WriteString(fmt.Sprintf("Cash balance at beginning of period %34s\n", cfs.BeginningCash.FormatDollars()))
	b.WriteString(fmt.Sprintf("Cash balance at end of period %40s\n", cfs.EndingCash.FormatDollars()))
	b.WriteString("════════════════════════════════════════════════════════════════════════\n")

	if len(cfs.NoncashDisclosures) > 0 {
		b.WriteString("\nNoncash Investing and Financing Activities:\n")
		for _, disc := range cfs.NoncashDisclosures {
			b.WriteString(fmt.Sprintf("  • %s\n", disc))
		}
	}

	if cfs.OperatingReconciliation != "" {
		b.WriteString("\nIndirect Operating Reconciliation:\n")
		b.WriteString(cfs.OperatingReconciliation)
	}

	return b.String()
}

// CaseTransaction models a single chronological event in a multi-event case study.
type CaseTransaction struct {
	Date               string                  `json:"date"`
	Description        string                  `json:"description"`
	Event              engine.TransactionEvent `json:"event"`
	IsAdjustment       bool                    `json:"is_adjustment"`
	CashClassification *CashFlowActivity       `json:"cash_classification,omitempty"`
	CashDescription    string                  `json:"cash_description,omitempty"`
	NoncashDisclosure  string                  `json:"noncash_disclosure,omitempty"`
}

// ComprehensiveCase models a small comprehensive business cycle with opening balances and transaction batches.
type ComprehensiveCase struct {
	CaseID          string                            `json:"case_id"`
	Title           string                            `json:"title"`
	Description     string                            `json:"description"`
	CompanyName     string                            `json:"company_name"`
	Period          string                            `json:"period"`
	AsOfDate        string                            `json:"as_of_date"`
	OpeningBalances map[domain.AccountID]domain.Money `json:"opening_balances"`
	Transactions    []CaseTransaction                 `json:"transactions"`
	Adjustments     []CaseTransaction                 `json:"adjustments"`
}

// AccountingCycleReport contains the complete multi-statement output and audit tie-outs.
type AccountingCycleReport struct {
	CaseID                    string                    `json:"case_id"`
	Title                     string                    `json:"title"`
	OpeningTrialBalance       TrialBalance              `json:"opening_trial_balance"`
	JournalEntries            []domain.Entry            `json:"journal_entries"`
	UnadjustedTrialBalance    TrialBalance              `json:"unadjusted_trial_balance"`
	AdjustingEntries          []domain.Entry            `json:"adjusting_entries"`
	AdjustedTrialBalance      TrialBalance              `json:"adjusted_trial_balance"`
	IncomeStatement           IncomeStatement           `json:"income_statement"`
	RetainedEarningsStatement RetainedEarningsStatement `json:"retained_earnings_statement"`
	BalanceSheet              BalanceSheet              `json:"balance_sheet"`
	CashFlowStatement         CashFlowStatement         `json:"cash_flow_statement"`
	Ledger                    *domain.TAccountLedger    `json:"ledger"`
	AuditReconciliationNotes  []string                  `json:"audit_reconciliation_notes"`
	AllGatesPassed            bool                      `json:"all_gates_passed"`
}
