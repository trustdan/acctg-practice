package engine

import (
	"fmt"

	"github.com/trustdan/acctg-practice/internal/domain"
)

// 1. M05InventoryFIFORule: Records Cost of Goods Sold and reduces Inventory under perpetual FIFO.
// Debits Cost of Goods Sold (Expense increases, reducing Equity),
// Credits Inventory (Asset decreases).
type M05InventoryFIFORule struct{}

func (M05InventoryFIFORule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "cost_of_goods_sold", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "inventory", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M05InventoryFIFORule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Recognized %s of Cost of Goods Sold and reduced Inventory under perpetual FIFO.", amt.FormatDollars()),
		DebitRationale:  "Cost of Goods Sold is an Expense that increased; Expenses increase by Debit, reducing Stockholders' Equity (-E).",
		CreditRationale: "Inventory is an Asset that decreased; Assets decrease by Credit (-A).",
		EquationSummary: "Assets decrease (-Inventory); Equity decreases (-Cost of Goods Sold); Liabilities are unchanged. Under FIFO during inflation, older and lower cost batches are expensed first, producing lower COGS and higher ending inventory.",
	}
}

func (M05InventoryFIFORule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "service_revenue" {
			return []string{TagM05InventoryRecordedAtSelling}, "Cost of Goods Sold reflects inventory cost, not the retail selling price or service revenue."
		}
		if p.AccountID == "rent_expense" || p.AccountID == "insurance_expense" {
			return []string{TagM05COGSConfusedWithOperating}, "Product costs must be expensed as Cost of Goods Sold, not periodic operating overhead."
		}
	}
	return nil, ""
}

// 2. M05InventoryLIFORule: Records Cost of Goods Sold and reduces Inventory under perpetual LIFO.
// Debits Cost of Goods Sold (Expense increases, reducing Equity),
// Credits Inventory (Asset decreases).
type M05InventoryLIFORule struct{}

func (M05InventoryLIFORule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "cost_of_goods_sold", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "inventory", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M05InventoryLIFORule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Recognized %s of Cost of Goods Sold and reduced Inventory under perpetual LIFO.", amt.FormatDollars()),
		DebitRationale:  "Cost of Goods Sold is an Expense that increased; Expenses increase by Debit, reducing Stockholders' Equity (-E).",
		CreditRationale: "Inventory is an Asset that decreased; Assets decrease by Credit (-A).",
		EquationSummary: "Assets decrease (-Inventory); Equity decreases (-Cost of Goods Sold); Liabilities are unchanged. Under LIFO during inflation, the newest and highest acquisition costs are matched against revenue, creating tax deferral benefits. Note: LIFO is permitted under US GAAP but prohibited under IFRS.",
	}
}

func (M05InventoryLIFORule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "cash" {
			return []string{TagCashRecordedWhenUncollected}, "The inventory cost outflow is a non-cash adjustment of inventory to COGS; cash was recorded in the separate sales transaction."
		}
	}
	return nil, ""
}

// 3. M05InventoryNRVWriteDownRule: Writes down damaged or obsolete inventory to Net Realizable Value (ASC 330).
// Debits Cost of Goods Sold / Loss on Inventory Write-Down (Expense increases),
// Credits Inventory (Asset decreases).
type M05InventoryNRVWriteDownRule struct{}

func (M05InventoryNRVWriteDownRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "cost_of_goods_sold", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "inventory", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M05InventoryNRVWriteDownRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Recorded %s inventory write-down to Net Realizable Value under US GAAP (ASC 330).", amt.FormatDollars()),
		DebitRationale:  "Cost of Goods Sold absorbs the valuation impairment as an Expense, reducing Net Income and Stockholders' Equity (-E).",
		CreditRationale: "Inventory is an Asset that is directly written down to net realizable value (-A); Assets decrease by Credit.",
		EquationSummary: "Assets decrease (-Inventory); Equity decreases (-Cost of Goods Sold); Liabilities are unchanged. Under US GAAP (ASC 330), once inventory is written down to NRV, the write-down establishes a new cost basis and subsequent market value recoveries cannot be reversed.",
	}
}

func (M05InventoryNRVWriteDownRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "accounts_payable" {
			return []string{TagPayableRecordedForCashPayment}, "An inventory write-down reflects asset obsolescence or market price drop, not an adjustment of trade debt to suppliers."
		}
	}
	return nil, ""
}

// 4. M05FreightInCapitalizeRule: Capitalizes inbound freight charges into Inventory historical cost (FOB Shipping Point).
// Debits Inventory (Asset increases),
// Credits Cash (Asset decreases).
type M05FreightInCapitalizeRule struct{}

func (M05FreightInCapitalizeRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "inventory", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M05FreightInCapitalizeRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Capitalized %s cash freight-in transportation costs into Inventory under FOB shipping point terms.", amt.FormatDollars()),
		DebitRationale:  "Inventory is an Asset that increased; inbound freight is a necessary cost of acquisition capitalized by Debit (+A).",
		CreditRationale: "Cash is an Asset that decreased; Assets decrease by Credit (-A).",
		EquationSummary: "Asset swap: Inventory increases (+A) and Cash decreases (-A); Total Assets, Liabilities, and Equity remain unchanged. Under US GAAP, freight-in is capitalized into inventory, whereas freight-out to customers is an operating selling expense.",
	}
}

func (M05FreightInCapitalizeRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "rent_expense" || p.AccountID == "insurance_expense" {
			return []string{TagM05FreightInExpensed}, "Inbound freight costs (freight-in) on purchases are product costs capitalized into Inventory, not period operating expenses."
		}
	}
	return nil, ""
}
