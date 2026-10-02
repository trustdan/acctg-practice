package domain

import (
	"errors"
	"fmt"
	"strings"
)

// Money represents a monetary amount stored in integer minor units (cents).
// This guarantees zero floating-point rounding issues in all accounting operations.
type Money int64

// Common Money constructors and helpers.
func NewMoney(minorUnits int64) Money {
	return Money(minorUnits)
}

// ZeroMoney returns a zero monetary amount ($0.00).
func ZeroMoney() Money {
	return Money(0)
}

func FromDollarsAndCents(dollars int64, cents int64) (Money, error) {
	if cents < 0 || cents >= 100 {
		return 0, errors.New("cents must be between 0 and 99")
	}
	if dollars < 0 {
		return Money(dollars*100 - cents), nil
	}
	return Money(dollars*100 + cents), nil
}

// Cents returns the raw integer minor units.
func (m Money) Cents() int64 {
	return int64(m)
}

// Dollars returns a float64 representation for display/calculation purposes only.
func (m Money) Dollars() float64 {
	return float64(m) / 100.0
}

// FormatDollars returns a formatted dollar string, e.g. "$50" for whole dollars or "$50.25" for cents.
func (m Money) FormatDollars() string {
	cents := int64(m)
	if cents < 0 {
		abs := -cents
		if abs%100 == 0 {
			return fmt.Sprintf("-$%d", abs/100)
		}
		return fmt.Sprintf("-$%d.%02d", abs/100, abs%100)
	}
	if cents%100 == 0 {
		return fmt.Sprintf("$%d", cents/100)
	}
	return fmt.Sprintf("$%d.%02d", cents/100, cents%100)
}

// FormatExact returns an exact 2-decimal dollar string, e.g. "$50.00".
func (m Money) FormatExact() string {
	cents := int64(m)
	if cents < 0 {
		abs := -cents
		return fmt.Sprintf("-$%d.%02d", abs/100, abs%100)
	}
	return fmt.Sprintf("$%d.%02d", cents/100, cents%100)
}

// FormatCommas returns a formatted dollar string with thousands comma separators, e.g. "$80,000" or "$1,250.50".
func (m Money) FormatCommas() string {
	cents := int64(m)
	sign := ""
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	dollars := cents / 100
	remainder := cents % 100

	dStr := fmt.Sprintf("%d", dollars)
	var formattedDollars strings.Builder
	l := len(dStr)
	for i, c := range dStr {
		if i > 0 && (l-i)%3 == 0 {
			formattedDollars.WriteRune(',')
		}
		formattedDollars.WriteRune(c)
	}

	if remainder == 0 {
		return fmt.Sprintf("%s$%s", sign, formattedDollars.String())
	}
	return fmt.Sprintf("%s$%s.%02d", sign, formattedDollars.String(), remainder)
}

func (m Money) Add(other Money) Money {
	return m + other
}

func (m Money) Sub(other Money) Money {
	return m - other
}

func (m Money) Abs() Money {
	if m < 0 {
		return -m
	}
	return m
}

func (m Money) IsZero() bool {
	return m == 0
}

func (m Money) IsPositive() bool {
	return m > 0
}

func (m Money) IsNegative() bool {
	return m < 0
}
