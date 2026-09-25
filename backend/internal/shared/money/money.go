package money

import (
	"github.com/shopspring/decimal"
)

// Zero returns a zero decimal
func Zero() decimal.Decimal {
	return decimal.Zero
}

// FromFloat creates a Decimal from float64 (only for initial parsing, never for calculations)
func FromFloat(f float64) decimal.Decimal {
	return decimal.NewFromFloat(f)
}

// FromString parses string into Decimal
func FromString(s string) (decimal.Decimal, error) {
	return decimal.NewFromString(s)
}

// RoundBanker rounds to standard currency decimal places (2 decimal places) using Banker's Rounding
func RoundBanker(d decimal.Decimal) decimal.Decimal {
	return d.RoundBanker(2)
}

// RoundTax rounds tax calculations to 4 decimal places
func RoundTax(d decimal.Decimal) decimal.Decimal {
	return d.RoundBanker(4)
}

// CalculateTax computes total tax amount: base * (taxRate / 100)
func CalculateTax(base decimal.Decimal, taxRate decimal.Decimal) decimal.Decimal {
	return base.Mul(taxRate).Div(decimal.NewFromInt(100)).RoundBanker(2)
}
