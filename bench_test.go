package currency_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/bojanz/currency"
)

var result currency.Amount
var cmpResult int

func BenchmarkNewAmount(b *testing.B) {
	var z currency.Amount
	for n := 0; n < b.N; n++ {
		z, _ = currency.NewAmount("99.99", "EUR")
	}
	result = z
}

func BenchmarkNewAmountFromBigInt(b *testing.B) {
	x := big.NewInt(9999)

	var z currency.Amount
	for n := 0; n < b.N; n++ {
		z, _ = currency.NewAmountFromBigInt(x, "EUR")
	}
	result = z
}

func BenchmarkNewAmountFromInt64(b *testing.B) {
	var z currency.Amount
	for n := 0; n < b.N; n++ {
		z, _ = currency.NewAmountFromInt64(9999, "EUR")
	}
	result = z
}

func BenchmarkAmount_Add(b *testing.B) {
	x, _ := currency.NewAmount("34.99", "USD")
	y, _ := currency.NewAmount("12.99", "USD")

	var z currency.Amount
	for n := 0; n < b.N; n++ {
		z, _ = x.Add(y)
	}
	result = z
}

func BenchmarkAmount_Sub(b *testing.B) {
	x, _ := currency.NewAmount("34.99", "USD")
	y, _ := currency.NewAmount("12.99", "USD")

	var z currency.Amount
	for n := 0; n < b.N; n++ {
		z, _ = x.Sub(y)
	}
	result = z
}

func BenchmarkAmount_Mul(b *testing.B) {
	x, _ := currency.NewAmount("34.99", "USD")

	var z currency.Amount
	for n := 0; n < b.N; n++ {
		z, _ = x.Mul("2")
	}
	result = z
}

func BenchmarkAmount_MulDec(b *testing.B) {
	x, _ := currency.NewAmount("34.99", "USD")

	var z currency.Amount
	for n := 0; n < b.N; n++ {
		z, _ = x.Mul("2.5")
	}
	result = z
}

func BenchmarkAmount_Div(b *testing.B) {
	x, _ := currency.NewAmount("34.99", "USD")

	var z currency.Amount
	for n := 0; n < b.N; n++ {
		z, _ = x.Div("2")
	}
	result = z
}

func BenchmarkAmount_DivDec(b *testing.B) {
	x, _ := currency.NewAmount("34.99", "USD")

	var z currency.Amount
	for n := 0; n < b.N; n++ {
		z, _ = x.Div("2.5")
	}
	result = z
}

func BenchmarkAmount_Round(b *testing.B) {
	x, _ := currency.NewAmount("34.9876", "USD")

	var z currency.Amount
	for n := 0; n < b.N; n++ {
		z = x.Round()
	}
	result = z
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
			var z currency.Amount
			for n := 0; n < b.N; n++ {
				z = x.RoundTo(2, roundingMode)
			}
			result = z
		})
	}
}

func BenchmarkAmount_Cmp(b *testing.B) {
	x, _ := currency.NewAmount("34.99", "USD")
	y, _ := currency.NewAmount("12.99", "USD")

	var z int
	for n := 0; n < b.N; n++ {
		z, _ = x.Cmp(y)
	}
	cmpResult = z
}
