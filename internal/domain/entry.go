package domain

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Posting represents a single debit or credit line to an account.
type Posting struct {
	AccountID AccountID `json:"account_id"`
	Side      Side      `json:"side"`
	Amount    Money     `json:"amount"`
}

func (p Posting) Validate(catalog *AccountCatalog) error {
	if p.AccountID == "" {
		return errors.New("posting account ID cannot be empty")
	}
	if !p.Side.IsValid() {
		return fmt.Errorf("posting side is invalid: %q", p.Side)
	}
	if !p.Amount.IsPositive() {
		return fmt.Errorf("posting amount must be positive, got %d", p.Amount.Cents())
	}
	if catalog != nil && !catalog.Has(p.AccountID) {
		return fmt.Errorf("posting references unknown account: %s", p.AccountID)
	}
	return nil
}

// Entry represents a complete journal transaction consisting of multiple postings.
type Entry struct {
	Postings []Posting `json:"postings"`
}

func NewEntry(postings ...Posting) Entry {
	return Entry{Postings: postings}
}

// TotalDebit computes the total debit sum of all postings in the entry.
func (e Entry) TotalDebit() Money {
	var total Money
	for _, p := range e.Postings {
		if p.Side == SideDebit {
			total = total.Add(p.Amount)
		}
	}
	return total
}

// TotalCredit computes the total credit sum of all postings in the entry.
func (e Entry) TotalCredit() Money {
	var total Money
	for _, p := range e.Postings {
		if p.Side == SideCredit {
			total = total.Add(p.Amount)
		}
	}
	return total
}

// IsBalanced checks whether total debits equal total credits.
func (e Entry) IsBalanced() bool {
	if len(e.Postings) < 2 {
		return false
	}
	dr := e.TotalDebit()
	cr := e.TotalCredit()
	return dr.IsPositive() && dr == cr
}

// Validate ensures every posting is valid and the overall entry balances.
func (e Entry) Validate(catalog *AccountCatalog) error {
	if len(e.Postings) < 2 {
		return fmt.Errorf("entry must have at least 2 postings, got %d", len(e.Postings))
	}
	for i, p := range e.Postings {
		if err := p.Validate(catalog); err != nil {
			return fmt.Errorf("posting[%d] invalid: %w", i, err)
		}
	}
	if !e.IsBalanced() {
		return fmt.Errorf("entry is not balanced: total debits (%s) != total credits (%s)",
			e.TotalDebit().FormatDollars(), e.TotalCredit().FormatDollars())
	}
	return nil
}

// EquationDelta represents the net change in the balance sheet equation (A = L + E).
type EquationDelta struct {
	DeltaAssets      Money `json:"delta_assets"`
	DeltaLiabilities Money `json:"delta_liabilities"`
	DeltaEquity      Money `json:"delta_equity"`
}

// EquationEffects calculates the net impact on Assets, Liabilities, and Equity.
// It verifies that DeltaAssets == DeltaLiabilities + DeltaEquity.
func (e Entry) EquationEffects(catalog *AccountCatalog) (EquationDelta, error) {
	if err := e.Validate(catalog); err != nil {
		return EquationDelta{}, err
	}

	var dAssets, dLiab, dEquity Money

	for _, p := range e.Postings {
		acc, ok := catalog.Get(p.AccountID)
		if !ok {
			return EquationDelta{}, fmt.Errorf("unknown account in catalog: %s", p.AccountID)
		}

		switch acc.Category {
		case CategoryAsset:
			if acc.ContraOf != "" {
				// Contra-asset: credit increases contra, debit decreases
				if p.Side == SideCredit {
					dAssets = dAssets.Sub(p.Amount)
				} else {
					dAssets = dAssets.Add(p.Amount)
				}
			} else {
				if p.Side == SideDebit {
					dAssets = dAssets.Add(p.Amount)
				} else {
					dAssets = dAssets.Sub(p.Amount)
				}
			}

		case CategoryLiability:
			if p.Side == SideCredit {
				dLiab = dLiab.Add(p.Amount)
			} else {
				dLiab = dLiab.Sub(p.Amount)
			}

		case CategoryEquity:
			if p.Side == SideCredit {
				dEquity = dEquity.Add(p.Amount)
			} else {
				dEquity = dEquity.Sub(p.Amount)
			}

		case CategoryRevenue:
			// Revenue increases Equity
			if p.Side == SideCredit {
				dEquity = dEquity.Add(p.Amount)
			} else {
				dEquity = dEquity.Sub(p.Amount)
			}

		case CategoryExpense:
			// Expense reduces Equity
			if p.Side == SideDebit {
				dEquity = dEquity.Sub(p.Amount)
			} else {
				dEquity = dEquity.Add(p.Amount)
			}

		case CategoryDividends:
			// Dividends reduce Equity
			if p.Side == SideDebit {
				dEquity = dEquity.Sub(p.Amount)
			} else {
				dEquity = dEquity.Add(p.Amount)
			}
		}
	}

	expectedAssets := dLiab.Add(dEquity)
	if dAssets != expectedAssets {
		return EquationDelta{}, fmt.Errorf("equation does not balance: delta assets (%s) != delta liabilities (%s) + delta equity (%s)",
			dAssets.FormatDollars(), dLiab.FormatDollars(), dEquity.FormatDollars())
	}

	return EquationDelta{
		DeltaAssets:      dAssets,
		DeltaLiabilities: dLiab,
		DeltaEquity:      dEquity,
	}, nil
}

// AggregatePostings aggregates duplicate lines having the same AccountID and Side,
// and orders debits first (sorted alphabetically by AccountID), followed by credits (sorted by AccountID).
func (e Entry) AggregatePostings() Entry {
	type key struct {
		accountID AccountID
		side      Side
	}

	totals := make(map[key]Money)
	var debitAccounts []AccountID
	var creditAccounts []AccountID
	seenDebit := make(map[AccountID]bool)
	seenCredit := make(map[AccountID]bool)

	for _, p := range e.Postings {
		k := key{accountID: p.AccountID, side: p.Side}
		totals[k] = totals[k].Add(p.Amount)
		if p.Side == SideDebit {
			if !seenDebit[p.AccountID] {
				seenDebit[p.AccountID] = true
				debitAccounts = append(debitAccounts, p.AccountID)
			}
		} else if p.Side == SideCredit {
			if !seenCredit[p.AccountID] {
				seenCredit[p.AccountID] = true
				creditAccounts = append(creditAccounts, p.AccountID)
			}
		}
	}

	sort.Slice(debitAccounts, func(i, j int) bool {
		return debitAccounts[i] < debitAccounts[j]
	})
	sort.Slice(creditAccounts, func(i, j int) bool {
		return creditAccounts[i] < creditAccounts[j]
	})

	var result []Posting
	for _, acc := range debitAccounts {
		amt := totals[key{accountID: acc, side: SideDebit}]
		if amt.IsPositive() {
			result = append(result, Posting{AccountID: acc, Side: SideDebit, Amount: amt})
		}
	}
	for _, acc := range creditAccounts {
		amt := totals[key{accountID: acc, side: SideCredit}]
		if amt.IsPositive() {
			result = append(result, Posting{AccountID: acc, Side: SideCredit, Amount: amt})
		}
	}

	return Entry{Postings: result}
}

// EqualPostings checks whether two entries have identical postings (ignoring line order),
// without aggregating duplicate lines.
func (e Entry) EqualPostings(other Entry) bool {
	if len(e.Postings) != len(other.Postings) {
		return false
	}
	matched := make([]bool, len(other.Postings))
	for _, pA := range e.Postings {
		found := false
		for j, pB := range other.Postings {
			if !matched[j] && pA.AccountID == pB.AccountID && pA.Side == pB.Side && pA.Amount == pB.Amount {
				matched[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// EqualNormalized checks whether two entries have identical postings regardless of line order
// AND after aggregating duplicate lines for the same account and side.
func (e Entry) EqualNormalized(other Entry) bool {
	normA := e.AggregatePostings()
	normB := other.AggregatePostings()
	return normA.EqualPostings(normB)
}

// PostingEquationEffect determines the (+A), (-A), (+L), (-L), (+E), or (-E) tag
// based on the account's category, contra status, and debit/credit side.
func PostingEquationEffect(p Posting, catalog *AccountCatalog) string {
	if catalog == nil {
		return ""
	}
	acc, ok := catalog.Get(p.AccountID)
	if !ok {
		return ""
	}

	switch acc.Category {
	case CategoryAsset:
		if acc.ContraOf != "" {
			if p.Side == SideCredit {
				return "(-A)"
			}
			return "(+A)"
		}
		if p.Side == SideDebit {
			return "(+A)"
		}
		return "(-A)"

	case CategoryLiability:
		if p.Side == SideCredit {
			return "(+L)"
		}
		return "(-L)"

	case CategoryEquity:
		if p.Side == SideCredit {
			return "(+E)"
		}
		return "(-E)"

	case CategoryRevenue:
		// Revenue increases equity
		if p.Side == SideCredit {
			return "(+E)"
		}
		return "(-E)"

	case CategoryExpense:
		// Expense reduces equity
		if p.Side == SideDebit {
			return "(-E)"
		}
		return "(+E)"

	case CategoryDividends:
		// Dividends reduce equity
		if p.Side == SideDebit {
			return "(-E)"
		}
		return "(+E)"

	default:
		return ""
	}
}

// RenderClassroomGrid renders postings into the 4-column bordered grid format
// directly matching the instructor slides:
// Column 1: Dr. / Cr.
// Column 2: Account Name with Equation Effect (e.g. Cash (+A), Loan (+L))
// Column 3: Debit amount
// Column 4: Credit amount
func RenderClassroomGrid(postings []Posting, catalog *AccountCatalog, totalWidth int, includeTotals bool) string {
	if totalWidth < 60 {
		totalWidth = 60
	}
	if totalWidth > 86 {
		totalWidth = 86
	}

	sideWidth := 6
	amtWidth := 13
	// Borders: │ Side │ Account │ Dr │ Cr │ -> 5 vertical bars
	accWidth := totalWidth - sideWidth - (amtWidth * 2) - 5
	if accWidth < 24 {
		accWidth = 24
	}

	var b strings.Builder

	// Top border
	b.WriteString("┌" + strings.Repeat("─", sideWidth) +
		"┬" + strings.Repeat("─", accWidth) +
		"┬" + strings.Repeat("─", amtWidth) +
		"┬" + strings.Repeat("─", amtWidth) + "┐\n")

	var totalDr, totalCr Money

	for i, p := range postings {
		if i > 0 {
			b.WriteString("├" + strings.Repeat("─", sideWidth) +
				"┼" + strings.Repeat("─", accWidth) +
				"┼" + strings.Repeat("─", amtWidth) +
				"┼" + strings.Repeat("─", amtWidth) + "┤\n")
		}

		accName := string(p.AccountID)
		if catalog != nil && catalog.Has(p.AccountID) {
			acc, _ := catalog.Get(p.AccountID)
			accName = acc.Name
		}
		effect := PostingEquationEffect(p, catalog)

		var sideStr, accLabel, drStr, crStr string
		if p.Side == SideDebit {
			sideStr = " Dr.  "
			accLabel = fmt.Sprintf(" %s %s", accName, effect)
			drStr = fmt.Sprintf("%11s  ", p.Amount.FormatCommas())
			crStr = strings.Repeat(" ", amtWidth)
			totalDr = totalDr.Add(p.Amount)
		} else {
			sideStr = " Cr.  "
			accLabel = fmt.Sprintf("   %s %s", accName, effect)
			drStr = strings.Repeat(" ", amtWidth)
			crStr = fmt.Sprintf("%11s  ", p.Amount.FormatCommas())
			totalCr = totalCr.Add(p.Amount)
		}

		if len(accLabel) > accWidth {
			accLabel = accLabel[:accWidth]
		} else {
			accLabel = accLabel + strings.Repeat(" ", accWidth-len(accLabel))
		}

		b.WriteString(fmt.Sprintf("│%s│%s│%s│%s│\n", sideStr, accLabel, drStr, crStr))
	}

	if includeTotals && len(postings) > 0 {
		b.WriteString("╞" + strings.Repeat("═", sideWidth) +
			"╪" + strings.Repeat("═", accWidth) +
			"╪" + strings.Repeat("═", amtWidth) +
			"╪" + strings.Repeat("═", amtWidth) + "╡\n")

		totalsLabel := " Totals"
		if len(totalsLabel) < accWidth {
			totalsLabel = totalsLabel + strings.Repeat(" ", accWidth-len(totalsLabel))
		}
		totalDrStr := fmt.Sprintf("%11s  ", totalDr.FormatCommas())
		totalCrStr := fmt.Sprintf("%11s  ", totalCr.FormatCommas())

		b.WriteString(fmt.Sprintf("│%s│%s│%s│%s│\n",
			strings.Repeat(" ", sideWidth),
			totalsLabel,
			totalDrStr,
			totalCrStr,
		))
	}

	// Bottom border
	b.WriteString("└" + strings.Repeat("─", sideWidth) +
		"┴" + strings.Repeat("─", accWidth) +
		"┴" + strings.Repeat("─", amtWidth) +
		"┴" + strings.Repeat("─", amtWidth) + "┘")

	return b.String()
}
