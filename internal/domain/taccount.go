package domain

import (
	"fmt"
	"sort"
	"strings"
)

// GenericVisualNotice explicitly identifies that classroom visuals follow a generic baseline
// until verified against official professor slide decks or syllabus materials.
const GenericVisualNotice = "Generic visual baseline (standard classroom T-account format; unverified against professor-specific slide templates)"

// TAccount models an individual account's T-account ledger card.
// In classroom accounting:
// - Left column is ALWAYS Debit
// - Right column is ALWAYS Credit
// - Normal side dictates whether Debit or Credit increases the balance:
//   - Asset, Expense, Dividend: Debit increases (+), Credit decreases (-)
//   - Liability, Equity, Revenue: Credit increases (+), Debit decreases (-)
type TAccount struct {
	AccountID    AccountID `json:"account_id"`
	AccountName  string    `json:"account_name"`
	Category     Category  `json:"category"`
	NormalSide   Side      `json:"normal_side"`
	Debits       []Money   `json:"debits"`
	Credits      []Money   `json:"credits"`
	TotalDebits  Money     `json:"total_debits"`
	TotalCredits Money     `json:"total_credits"`
	NetBalance   Money     `json:"net_balance"`
	BalanceSide  Side      `json:"balance_side"`
}

// NewTAccount creates an initialized TAccount for a given account definition.
func NewTAccount(acc Account) TAccount {
	return TAccount{
		AccountID:    acc.ID,
		AccountName:  acc.Name,
		Category:     acc.Category,
		NormalSide:   acc.NormalSide,
		Debits:       make([]Money, 0),
		Credits:      make([]Money, 0),
		TotalDebits:  ZeroMoney(),
		TotalCredits: ZeroMoney(),
		NetBalance:   ZeroMoney(),
		BalanceSide:  acc.NormalSide,
	}
}

// Post records a debit or credit posting onto this T-account and recomputes balances.
func (t *TAccount) Post(side Side, amount Money) {
	if side == SideDebit {
		t.Debits = append(t.Debits, amount)
		t.TotalDebits = t.TotalDebits.Add(amount)
	} else if side == SideCredit {
		t.Credits = append(t.Credits, amount)
		t.TotalCredits = t.TotalCredits.Add(amount)
	}
	t.recomputeBalance()
}

func (t *TAccount) recomputeBalance() {
	if t.TotalDebits.Cents() >= t.TotalCredits.Cents() {
		t.NetBalance = t.TotalDebits.Sub(t.TotalCredits)
		t.BalanceSide = SideDebit
	} else {
		t.NetBalance = t.TotalCredits.Sub(t.TotalDebits)
		t.BalanceSide = SideCredit
	}
}

// DebitHeading returns the pedagogical label for the left (Debit) column.
func (t *TAccount) DebitHeading() string {
	switch t.Category {
	case CategoryAsset, CategoryExpense, CategoryDividends:
		return "Debit (+)"
	default:
		return "Debit (-)"
	}
}

// CreditHeading returns the pedagogical label for the right (Credit) column.
func (t *TAccount) CreditHeading() string {
	switch t.Category {
	case CategoryLiability, CategoryEquity, CategoryRevenue:
		return "Credit (+)"
	default:
		return "Credit (-)"
	}
}

// Render produces a clean standard classroom T-account ASCII visualization.
func (t *TAccount) Render() string {
	var b strings.Builder

	catTitle := string(t.Category)
	if len(catTitle) > 0 {
		catTitle = strings.ToUpper(catTitle[:1]) + strings.ToLower(catTitle[1:])
	}
	title := fmt.Sprintf("%s (%s)", t.AccountName, catTitle)
	leftHead := t.DebitHeading()
	rightHead := t.CreditHeading()

	b.WriteString(fmt.Sprintf("       %-28s\n", title))
	b.WriteString("─────────────────────────────┼─────────────────────────────\n")
	b.WriteString(fmt.Sprintf("  %-25s  │  %-25s  \n", leftHead, rightHead))
	b.WriteString("─────────────────────────────┼─────────────────────────────\n")

	maxLines := len(t.Debits)
	if len(t.Credits) > maxLines {
		maxLines = len(t.Credits)
	}
	if maxLines == 0 {
		maxLines = 1
	}

	for i := 0; i < maxLines; i++ {
		drStr := ""
		crStr := ""
		if i < len(t.Debits) {
			drStr = t.Debits[i].FormatExact()
		}
		if i < len(t.Credits) {
			crStr = t.Credits[i].FormatExact()
		}
		b.WriteString(fmt.Sprintf("  %-25s  │  %-25s  \n", drStr, crStr))
	}

	b.WriteString("─────────────────────────────┼─────────────────────────────\n")

	balDr := ""
	balCr := ""
	balStr := fmt.Sprintf("Bal: %s", t.NetBalance.FormatExact())
	if t.BalanceSide == SideDebit {
		balDr = balStr
	} else {
		balCr = balStr
	}
	b.WriteString(fmt.Sprintf("  %-25s  │  %-25s  \n", balDr, balCr))

	return b.String()
}

// TAccountLedger collects and organizes T-accounts across a catalog for a transaction or session.
type TAccountLedger struct {
	catalog  *AccountCatalog
	accounts map[AccountID]*TAccount
}

// NewTAccountLedger initializes a ledger against an account catalog.
func NewTAccountLedger(catalog *AccountCatalog) *TAccountLedger {
	return &TAccountLedger{
		catalog:  catalog,
		accounts: make(map[AccountID]*TAccount),
	}
}

// GetOrCreate returns the TAccount for the accountID, initializing it if absent.
func (l *TAccountLedger) GetOrCreate(id AccountID) *TAccount {
	if ta, ok := l.accounts[id]; ok {
		return ta
	}
	name := string(id)
	cat := CategoryAsset
	normal := SideDebit
	if l.catalog != nil {
		if acc, ok := l.catalog.Get(id); ok {
			name = acc.Name
			cat = acc.Category
			normal = acc.NormalSide
		}
	}
	ta := &TAccount{
		AccountID:    id,
		AccountName:  name,
		Category:     cat,
		NormalSide:   normal,
		Debits:       make([]Money, 0),
		Credits:      make([]Money, 0),
		TotalDebits:  ZeroMoney(),
		TotalCredits: ZeroMoney(),
		NetBalance:   ZeroMoney(),
		BalanceSide:  normal,
	}
	l.accounts[id] = ta
	return ta
}

// PostEntry records all postings from an entry into their respective T-accounts.
func (l *TAccountLedger) PostEntry(entry Entry) {
	for _, p := range entry.Postings {
		ta := l.GetOrCreate(p.AccountID)
		ta.Post(p.Side, p.Amount)
	}
}

// TotalDebits calculates the grand total of all debit postings across all accounts in the ledger.
func (l *TAccountLedger) TotalDebits() Money {
	var total Money
	for _, ta := range l.accounts {
		total = total.Add(ta.TotalDebits)
	}
	return total
}

// TotalCredits calculates the grand total of all credit postings across all accounts in the ledger.
func (l *TAccountLedger) TotalCredits() Money {
	var total Money
	for _, ta := range l.accounts {
		total = total.Add(ta.TotalCredits)
	}
	return total
}

// AccountsInOrder returns the T-accounts sorted deterministically:
// First by Category (Asset, Liability, Equity, Revenue, Expense, Dividends), then by AccountID.
func (l *TAccountLedger) AccountsInOrder() []*TAccount {
	catOrder := map[Category]int{
		CategoryAsset:     1,
		CategoryLiability: 2,
		CategoryEquity:    3,
		CategoryRevenue:   4,
		CategoryExpense:   5,
		CategoryDividends: 6,
	}

	result := make([]*TAccount, 0, len(l.accounts))
	for _, ta := range l.accounts {
		result = append(result, ta)
	}

	sort.Slice(result, func(i, j int) bool {
		cI := catOrder[result[i].Category]
		cJ := catOrder[result[j].Category]
		if cI != cJ {
			return cI < cJ
		}
		return result[i].AccountID < result[j].AccountID
	})

	return result
}
