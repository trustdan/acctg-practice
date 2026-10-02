package domain

import (
	"fmt"
	"strings"
)

// Category represents the fundamental accounting classification.
type Category string

const (
	CategoryAsset     Category = "asset"
	CategoryLiability Category = "liability"
	CategoryEquity    Category = "equity"
	CategoryRevenue   Category = "revenue"
	CategoryExpense   Category = "expense"
	CategoryDividends Category = "dividends"
)

func (c Category) IsValid() bool {
	switch c {
	case CategoryAsset, CategoryLiability, CategoryEquity, CategoryRevenue, CategoryExpense, CategoryDividends:
		return true
	default:
		return false
	}
}

// NormalSide returns the standard normal balance side for a category.
func (c Category) NormalSide() Side {
	switch c {
	case CategoryAsset, CategoryExpense, CategoryDividends:
		return SideDebit
	case CategoryLiability, CategoryEquity, CategoryRevenue:
		return SideCredit
	default:
		return SideUnknown
	}
}

// Side represents Debit or Credit.
type Side string

const (
	SideDebit   Side = "debit"
	SideCredit  Side = "credit"
	SideUnknown Side = "unknown"
)

func (s Side) IsValid() bool {
	return s == SideDebit || s == SideCredit
}

func (s Side) Opposite() Side {
	switch s {
	case SideDebit:
		return SideCredit
	case SideCredit:
		return SideDebit
	default:
		return SideUnknown
	}
}

// Direction represents whether an account balance increases or decreases.
type Direction string

const (
	DirectionIncrease Direction = "increase"
	DirectionDecrease Direction = "decrease"
	DirectionUnknown  Direction = "unknown"
)

func (d Direction) IsValid() bool {
	return d == DirectionIncrease || d == DirectionDecrease
}

// SideForDirection computes the posting side (debit or credit) given the account's normal side and change direction.
func SideForDirection(normal Side, dir Direction) Side {
	if !normal.IsValid() || !dir.IsValid() {
		return SideUnknown
	}
	if dir == DirectionIncrease {
		return normal
	}
	return normal.Opposite()
}

// DirectionForSide computes whether a debit or credit increases or decreases an account with the given normal side.
func DirectionForSide(normal Side, side Side) Direction {
	if !normal.IsValid() || !side.IsValid() {
		return DirectionUnknown
	}
	if side == normal {
		return DirectionIncrease
	}
	return DirectionDecrease
}

// AccountID uniquely identifies an account in the system (e.g., "cash", "service_revenue").
type AccountID string

// Account represents a single ledger account with its category and normal side.
type Account struct {
	ID         AccountID `json:"id"`
	Name       string    `json:"name"`
	Category   Category  `json:"category"`
	NormalSide Side      `json:"normal_side"`
	ContraOf   AccountID `json:"contra_of,omitempty"`
}

func (a Account) Validate() error {
	if strings.TrimSpace(string(a.ID)) == "" {
		return fmt.Errorf("account ID cannot be empty")
	}
	if strings.TrimSpace(a.Name) == "" {
		return fmt.Errorf("account %s name cannot be empty", a.ID)
	}
	if !a.Category.IsValid() {
		return fmt.Errorf("account %s has invalid category: %q", a.ID, a.Category)
	}
	if !a.NormalSide.IsValid() {
		return fmt.Errorf("account %s has invalid normal side: %q", a.ID, a.NormalSide)
	}

	expectedNormal := a.Category.NormalSide()
	if a.ContraOf == "" {
		if a.NormalSide != expectedNormal {
			return fmt.Errorf("account %s (%s) has normal side %s, expected %s",
				a.ID, a.Category, a.NormalSide, expectedNormal)
		}
	} else {
		// Contra accounts have the opposite normal side of the base account category.
		if a.NormalSide != expectedNormal.Opposite() {
			return fmt.Errorf("contra account %s (%s) has normal side %s, expected %s",
				a.ID, a.Category, a.NormalSide, expectedNormal.Opposite())
		}
	}
	return nil
}

// AccountCatalog is an in-memory collection of validated accounts.
type AccountCatalog struct {
	accounts map[AccountID]Account
	order    []AccountID
}

func NewAccountCatalog() *AccountCatalog {
	return &AccountCatalog{
		accounts: make(map[AccountID]Account),
		order:    make([]AccountID, 0),
	}
}

func (c *AccountCatalog) Add(acc Account) error {
	if err := acc.Validate(); err != nil {
		return err
	}
	if _, exists := c.accounts[acc.ID]; exists {
		return fmt.Errorf("duplicate account ID: %s", acc.ID)
	}
	c.accounts[acc.ID] = acc
	c.order = append(c.order, acc.ID)
	return nil
}

func (c *AccountCatalog) Get(id AccountID) (Account, bool) {
	acc, ok := c.accounts[id]
	return acc, ok
}

func (c *AccountCatalog) Has(id AccountID) bool {
	_, ok := c.accounts[id]
	return ok
}

func (c *AccountCatalog) All() []Account {
	res := make([]Account, len(c.order))
	for i, id := range c.order {
		res[i] = c.accounts[id]
	}
	return res
}

func (c *AccountCatalog) Count() int {
	return len(c.accounts)
}
