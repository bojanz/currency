package currency_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/bojanz/currency"
)

func BenchmarkNewAmount(b *testing.B) {
	for b.Loop() {
		currency.NewAmount("99.99", "EUR")
	}
}

func BenchmarkNewAmountFromBigInt(b *testing.B) {
	x := big.NewInt(9999)

	for b.Loop() {
		currency.NewAmountFromBigInt(x, "EUR")
	}
}

func BenchmarkNewAmountFromInt64(b *testing.B) {
	for b.Loop() {
		currency.NewAmountFromInt64(9999, "EUR")
	}
}

func BenchmarkAmount_Add(b *testing.B) {
	x, _ := currency.NewAmount("34.99", "USD")
	y, _ := currency.NewAmount("12.99", "USD")

	for b.Loop() {
		x.Add(y)
	}
}

func BenchmarkAmount_Sub(b *testing.B) {
	x, _ := currency.NewAmount("34.99", "USD")
	y, _ := currency.NewAmount("12.99", "USD")

	for b.Loop() {
		x.Sub(y)
	}
}

func BenchmarkAmount_Mul(b *testing.B) {
	x, _ := currency.NewAmount("34.99", "USD")

	for b.Loop() {
		x.Mul("2")
	}
}

func BenchmarkAmount_MulDec(b *testing.B) {
	x, _ := currency.NewAmount("34.99", "USD")

	for b.Loop() {
		x.Mul("2.5")
	}
}

func BenchmarkAmount_Div(b *testing.B) {
	x, _ := currency.NewAmount("34.99", "USD")

	for b.Loop() {
		x.Div("2")
	}
}

func BenchmarkAmount_DivDec(b *testing.B) {
	x, _ := currency.NewAmount("34.99", "USD")

	for b.Loop() {
		x.Div("2.5")
	}
}

func BenchmarkAmount_Round(b *testing.B) {
	x, _ := currency.NewAmount("34.9876", "USD")

	for b.Loop() {
		x.Round()
	}
}

func BenchmarkAmount_RoundTo(b *testing.B) {
	x, _ := currency.NewAmount("34.9876", "USD")
	roundingModes := []currency.RoundingMode{
		currency.RoundHalfUp,
		currency.RoundHalfDown,
		currency.RoundUp,
		currency.RoundDown,
	}

	for _, roundingMode := range roundingModes {
		b.Run(fmt.Sprintf("rounding_mode_%d", roundingMode), func(b *testing.B) {
			for b.Loop() {
				x.RoundTo(2, roundingMode)
			}
		})
	}
}

func BenchmarkAmount_Cmp(b *testing.B) {
	x, _ := currency.NewAmount("34.99", "USD")
	y, _ := currency.NewAmount("12.99", "USD")

	for b.Loop() {
		x.Cmp(y)
	}
}

func BenchmarkAmount_IsPositive(b *testing.B) {
	x, _ := currency.NewAmount("34.99", "USD")

	for b.Loop() {
		x.IsPositive()
	}
}

func BenchmarkAmount_IsNegative(b *testing.B) {
	x, _ := currency.NewAmount("-12.00", "USD")

	for b.Loop() {
		x.IsNegative()
	}
}

func BenchmarkAmount_IsZero(b *testing.B) {
	x, _ := currency.NewAmount("0.00", "USD")

	for b.Loop() {
		x.IsZero()
	}
}
