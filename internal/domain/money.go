package domain

import (
	"errors"
	"fmt"
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
