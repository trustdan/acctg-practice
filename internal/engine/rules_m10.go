package engine

import (
	"fmt"

	"github.com/trustdan/acctg-practice/internal/domain"
)

// 1. M10ForecastingOperatingExpenseAPRule: Pro-forma operating expense incurred on trade terms (Accounts Payable).
// Debits Rent Expense (Expense increases, reducing projected Equity),
// Credits Accounts Payable (Liability increases).
type M10ForecastingOperatingExpenseAPRule struct{}

func (M10ForecastingOperatingExpenseAPRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "rent_expense", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "accounts_payable", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M10ForecastingOperatingExpenseAPRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Modeled projected operating expense of %s incurred on trade terms in pro-forma financial statements.", amt.FormatDollars()),
		DebitRationale:  "Rent Expense is an operating expense that increased; expenses increase by Debit, reducing projected Equity (-E).",
		CreditRationale: "Accounts Payable is a current operating liability that increased by Credit (+L).",
		EquationSummary: "Liabilities increase (+Accounts Payable); Equity decreases (-Rent Expense); Assets are unchanged. In integrated financial forecasting, linking projected expenses to working capital drivers (such as days-payable-outstanding) ensures that the pro-forma income statement seamlessly feeds projected balance sheet trade liabilities and indirect operating cash flows.",
	}
}

func (M10ForecastingOperatingExpenseAPRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "cash" && p.Side == domain.SideCredit {
			return []string{TagM10ForecastingAPOmitted}, "The scenario models operating expenses incurred on trade credit terms (Accounts Payable) rather than immediate cash settlement."
		}
		if p.AccountID == "service_revenue" {
			return []string{TagWrongAccount}, "Operating expenses reduce income, not revenues."
		}
	}
	return nil, ""
}

// 2. M10ForecastingCapExRule: Pro-forma capital expenditures (CapEx) purchasing productive equipment.
// Debits Equipment (Asset increases),
// Credits Cash (Asset decreases).
type M10ForecastingCapExRule struct{}

func (M10ForecastingCapExRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "equipment", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M10ForecastingCapExRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Modeled projected capital expenditure (CapEx) of %s for equipment capacity expansion in pro-forma statements.", amt.FormatDollars()),
		DebitRationale:  "Equipment is a long-term productive asset that increased by Debit (+A).",
		CreditRationale: "Cash is a current asset that decreased by Credit (-A).",
		EquationSummary: "Asset swap: Long-term assets increase (+Equipment) and liquid assets decrease (-Cash); Liabilities and Equity are unchanged. Projecting capital expenditures ensures productive capacity keeps pace with forecasted sales growth and flows directly into investing cash outflows on the pro-forma cash flow statement.",
	}
}

func (M10ForecastingCapExRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "rent_expense" || p.AccountID == "depreciation_expense" {
			return []string{TagM10CapExExpensedImmediately}, "Capital expenditures for productive assets that benefit multiple future periods must be capitalized on the balance sheet as Equipment, not immediately expensed in full."
		}
		if p.AccountID == "common_stock" {
			return []string{TagWrongAccount}, "Purchasing physical equipment with cash affects assets, not contributed common stock."
		}
	}
	return nil, ""
}

// 3. M10ForecastingDebtPlugRule: External Financing Needed (EFN) modeled through long-term debt issuance.
// Debits Cash (Asset increases),
// Credits Notes Payable (Liability increases).
type M10ForecastingDebtPlugRule struct{}

func (M10ForecastingDebtPlugRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "notes_payable", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M10ForecastingDebtPlugRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Modeled projected External Financing Needed (EFN) debt plug of %s via long-term notes payable in pro-forma statements.", amt.FormatDollars()),
		DebitRationale:  "Cash is a liquid asset that increased by Debit (+A) to fund projected asset expansion.",
		CreditRationale: "Notes Payable is a debt liability that increased by Credit (+L) to balance projected assets with liabilities and equity.",
		EquationSummary: "Assets increase (+Cash); Liabilities increase (+Notes Payable); Equity is unchanged. In pro-forma financial statement forecasting, the External Financing Needed (EFN) debt plug reconciles projected asset requirements with projected funding, linking the projected balance sheet to financing cash inflows.",
	}
}

func (M10ForecastingDebtPlugRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "common_stock" {
			return []string{TagM10EFNTreatedAsEquityWithoutShares}, "When financing pro-forma asset expansion through debt borrowing, the obligation must be credited to Notes Payable, not Common Stock."
		}
		if p.AccountID == "service_revenue" {
			return []string{TagWrongAccount}, "Borrowing debt proceeds creates a liability obligation, not operating earned revenue."
		}
	}
	return nil, ""
}
