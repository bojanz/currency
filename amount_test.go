// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package currency_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/bojanz/currency"
)

func TestNewAmount(t *testing.T) {
	_, err := currency.NewAmount("INVALID", "USD")
	if e, ok := err.(currency.InvalidNumberError); ok {
		if e.Number != "INVALID" {
			t.Errorf("got %v, want INVALID", e.Number)
		}
		wantError := `invalid number "INVALID"`
		if e.Error() != wantError {
			t.Errorf("got %v, want %v", e.Error(), wantError)
		}
	} else {
		t.Errorf("got %T, want currency.InvalidNumberError", err)
	}

	_, err = currency.NewAmount("10.99", "usd")
	if e, ok := err.(currency.InvalidCurrencyCodeError); ok {
		if e.CurrencyCode != "usd" {
			t.Errorf("got %v, want usd", e.CurrencyCode)
		}
		wantError := `invalid currency code "usd"`
		if e.Error() != wantError {
			t.Errorf("got %v, want %v", e.Error(), wantError)
		}
	} else {
		t.Errorf("got %T, want currency.InvalidCurrencyCodeError", err)
	}

	a, err := currency.NewAmount("10.99", "USD")
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if a.Number() != "10.99" {
		t.Errorf("got %v, want 10.99", a.Number())
	}
	if a.CurrencyCode() != "USD" {
		t.Errorf("got %v, want USD", a.CurrencyCode())
	}
	if a.String() != "10.99 USD" {
		t.Errorf("got %v, want 10.99 USD", a.String())
	}
}

func TestNewAmountFromBigInt(t *testing.T) {
	_, err := currency.NewAmountFromBigInt(nil, "USD")
	if e, ok := err.(currency.InvalidNumberError); ok {
		if e.Number != "nil" {
			t.Errorf("got %v, want nil", e.Number)
		}
		wantError := `invalid number "nil"`
		if e.Error() != wantError {
			t.Errorf("got %v, want %v", e.Error(), wantError)
		}
	} else {
		t.Errorf("got %T, want currency.InvalidNumberError", err)
	}

	_, err = currency.NewAmountFromBigInt(big.NewInt(1099), "usd")
	if e, ok := err.(currency.InvalidCurrencyCodeError); ok {
		if e.CurrencyCode != "usd" {
			t.Errorf("got %v, want usd", e.CurrencyCode)
		}
		wantError := `invalid currency code "usd"`
		if e.Error() != wantError {
			t.Errorf("got %v, want %v", e.Error(), wantError)
		}
	} else {
		t.Errorf("got %T, want currency.InvalidCurrencyCodeError", err)
	}

	// An integer larger than math.MaxInt64.
	hugeInt, _ := big.NewInt(0).SetString("922337203685477598799", 10)
	tests := []struct {
		n            *big.Int
		currencyCode string
		wantNumber   string
	}{
		{big.NewInt(2099), "USD", "20.99"},
		{big.NewInt(5000), "USD", "50.00"},
		{big.NewInt(50), "JPY", "50"},
		{hugeInt, "USD", "9223372036854775987.99"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			a, err := currency.NewAmountFromBigInt(tt.n, tt.currencyCode)
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}
			if a.Number() != tt.wantNumber {
				t.Errorf("got %v, want %v", a.Number(), tt.wantNumber)
			}
			if a.CurrencyCode() != tt.currencyCode {
				t.Errorf("got %v, want %v", a.CurrencyCode(), tt.currencyCode)
			}
		})
	}
}

func TestNewAmountFromInt64(t *testing.T) {
	_, err := currency.NewAmountFromInt64(1099, "usd")
	if e, ok := err.(currency.InvalidCurrencyCodeError); ok {
		if e.CurrencyCode != "usd" {
			t.Errorf("got %v, want usd", e.CurrencyCode)
		}
		wantError := `invalid currency code "usd"`
		if e.Error() != wantError {
			t.Errorf("got %v, want %v", e.Error(), wantError)
		}
	} else {
		t.Errorf("got %T, want currency.InvalidCurrencyCodeError", err)
	}

	tests := []struct {
		n            int64
		currencyCode string
		wantNumber   string
	}{
		{2099, "USD", "20.99"},
		{5000, "USD", "50.00"},
		{50, "JPY", "50"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			a, err := currency.NewAmountFromInt64(tt.n, tt.currencyCode)
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}
			if a.Number() != tt.wantNumber {
				t.Errorf("got %v, want %v", a.Number(), tt.wantNumber)
			}
			if a.CurrencyCode() != tt.currencyCode {
				t.Errorf("got %v, want %v", a.CurrencyCode(), tt.currencyCode)
			}
		})
	}
}

func TestAmount_BigInt(t *testing.T) {
	tests := []struct {
		number       string
		currencyCode string
		want         *big.Int
	}{
		{"20.99", "USD", big.NewInt(2099)},
		// Number with additional decimals.
		{"12.3564", "USD", big.NewInt(1236)},
		// Number with no decimals.
		{"50", "USD", big.NewInt(5000)},
		{"50", "JPY", big.NewInt(50)},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			a, _ := currency.NewAmount(tt.number, tt.currencyCode)
			got := a.BigInt()
			if got.Cmp(tt.want) != 0 {
				t.Errorf("got %v, want %v", got, tt.want)
			}
			// Confirm that a is unchanged.
			if a.Number() != tt.number {
				t.Errorf("got %v, want %v", a.Number(), tt.number)
			}
		})
	}
}

func TestAmount_Int64(t *testing.T) {
	// Number that can't be represented as an int64.
	a, _ := currency.NewAmount("922337203685477598799", "USD")
	n, err := a.Int64()
	if n != 0 {
		t.Error("expected a.Int64() to be 0")
	}
	if err == nil {
		t.Error("expected a.Int64() to return an error")
	}

	tests := []struct {
		number       string
		currencyCode string
		want         int64
	}{
		{"20.99", "USD", 2099},
		// Number with additional decimals.
		{"12.3564", "USD", 1236},
		// Number with no decimals.
		{"50", "USD", 5000},
		{"50", "JPY", 50},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			a, _ := currency.NewAmount(tt.number, tt.currencyCode)
			got, _ := a.Int64()
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
			// Confirm that a is unchanged.
			if a.Number() != tt.number {
				t.Errorf("got %v, want %v", a.Number(), tt.number)
			}
		})
	}
}

func TestAmount_Convert(t *testing.T) {
	a, _ := currency.NewAmount("20.99", "USD")

	_, err := a.Convert("eur", "0.91")
	if e, ok := err.(currency.InvalidCurrencyCodeError); ok {
		if e.CurrencyCode != "eur" {
			t.Errorf("got %v, want eur", e.CurrencyCode)
		}
		wantError := `invalid currency code "eur"`
		if e.Error() != wantError {
			t.Errorf("got %v, want %v", e.Error(), wantError)
		}
	} else {
		t.Errorf("got %T, want currency.InvalidCurrencyCodeError", err)
	}

	_, err = a.Convert("EUR", "INVALID")
	if e, ok := err.(currency.InvalidNumberError); ok {
		if e.Number != "INVALID" {
			t.Errorf("got %v, want INVALID", e.Number)
		}
		wantError := `invalid number "INVALID"`
		if e.Error() != wantError {
			t.Errorf("got %v, want %v", e.Error(), wantError)
		}
	} else {
		t.Errorf("got %T, want currency.InvalidNumberError", err)
	}

	b, err := a.Convert("EUR", "0.91")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if b.String() != "19.1009 EUR" {
		t.Errorf("got %v, want 19.1009 EUR", b.String())
	}
	// Confirm that a is unchanged.
	if a.String() != "20.99 USD" {
		t.Errorf("got %v, want 20.99 USD", a.String())
	}

	// An amount larger than math.MaxInt64.
	c, _ := currency.NewAmount("922337203685477598799", "USD")
	d, err := c.Convert("RSD", "100")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if d.String() != "92233720368547759879900 RSD" {
		t.Errorf("got %v, want 92233720368547759879900 RSD", d.String())
	}
}

func TestAmount_Add(t *testing.T) {
	a, _ := currency.NewAmount("20.99", "USD")
	b, _ := currency.NewAmount("3.50", "USD")
	x, _ := currency.NewAmount("99.99", "EUR")
	var z currency.Amount

	_, err := a.Add(x)
	if e, ok := err.(currency.MismatchError); ok {
		if e.A != a {
			t.Errorf("got %v, want %v", e.A, a)
		}
		if e.B != x {
			t.Errorf("got %v, want %v", e.B, x)
		}
		wantError := `amounts "20.99 USD" and "99.99 EUR" have mismatched currency codes`
		if e.Error() != wantError {
			t.Errorf("got %v, want %v", e.Error(), wantError)
		}
	} else {
		t.Errorf("got %T, want currency.MismatchError", err)
	}

	c, err := a.Add(b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.String() != "24.49 USD" {
		t.Errorf("got %v, want 24.49 USD", c.String())
	}
	// Confirm that a and b are unchanged.
	if a.String() != "20.99 USD" {
		t.Errorf("got %v, want 20.99 USD", a.String())
	}
	if b.String() != "3.50 USD" {
		t.Errorf("got %v, want 3.50 USD", b.String())
	}

	// An amount equal to math.MaxInt64.
	d, _ := currency.NewAmount("9223372036854775807", "USD")
	e, err := d.Add(a)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if e.String() != "9223372036854775827.99 USD" {
		t.Errorf("got %v, want 9223372036854775827.99 USD", e.String())
	}

	// Test that addition with the zero value works and yields the other operand.
	f, err := a.Add(z)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !f.Equal(a) {
		t.Errorf("%v + zero = %v, want %v", a, f, a)
	}

	g, err := z.Add(a)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !g.Equal(a) {
		t.Errorf("%v + zero = %v, want %v", a, g, a)
	}
}

func TestAmount_Sub(t *testing.T) {
	a, _ := currency.NewAmount("20.99", "USD")
	b, _ := currency.NewAmount("3.50", "USD")
	x, _ := currency.NewAmount("99.99", "EUR")
	var z currency.Amount

	_, err := a.Sub(x)
	if e, ok := err.(currency.MismatchError); ok {
		if e.A != a {
			t.Errorf("got %v, want %v", e.A, a)
		}
		if e.B != x {
			t.Errorf("got %v, want %v", e.B, x)
		}
		wantError := `amounts "20.99 USD" and "99.99 EUR" have mismatched currency codes`
		if e.Error() != wantError {
			t.Errorf("got %v, want %v", e.Error(), wantError)
		}
	} else {
		t.Errorf("got %T, want currency.MismatchError", err)
	}

	c, err := a.Sub(b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.String() != "17.49 USD" {
		t.Errorf("got %v, want 17.49 USD", c.String())
	}
	// Confirm that a and b are unchanged.
	if a.String() != "20.99 USD" {
		t.Errorf("got %v, want 20.99 USD", a.String())
	}
	if b.String() != "3.50 USD" {
		t.Errorf("got %v, want 3.50 USD", b.String())
	}

	// An amount larger than math.MaxInt64.
	d, _ := currency.NewAmount("922337203685477598799", "USD")
	e, err := d.Sub(a)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if e.String() != "922337203685477598778.01 USD" {
		t.Errorf("got %v, want 922337203685477598778.01 USD", e.String())
	}

	// Test that subtraction with the zero value works.
	f, err := a.Sub(z)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !f.Equal(a) {
		t.Errorf("%v - zero = %v, want %v", a, f, a)
	}

	g, err := z.Sub(a)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	negA, err := a.Mul("-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !g.Equal(negA) {
		t.Errorf("zero - %v = %v, want %v", a, g, negA)
	}
}

func TestAmount_Mul(t *testing.T) {
	a, _ := currency.NewAmount("20.99", "USD")

	_, err := a.Mul("INVALID")
	if e, ok := err.(currency.InvalidNumberError); ok {
		if e.Number != "INVALID" {
			t.Errorf("got %v, want INVALID", e.Number)
		}
		wantError := `invalid number "INVALID"`
		if e.Error() != wantError {
			t.Errorf("got %v, want %v", e.Error(), wantError)
		}
	} else {
		t.Errorf("got %T, want currency.InvalidNumberError", err)
	}

	b, err := a.Mul("0.20")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if b.String() != "4.1980 USD" {
		t.Errorf("got %v, want 4.1980 USD", b.String())
	}
	// Confirm that a is unchanged.
	if a.String() != "20.99 USD" {
		t.Errorf("got %v, want 20.99 USD", a.String())
	}

	// An amount equal to math.MaxInt64.
	d, _ := currency.NewAmount("9223372036854775807", "USD")
	e, err := d.Mul("10")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if e.String() != "92233720368547758070 USD" {
		t.Errorf("got %v, want 92233720368547758070 USD", e.String())
	}
}

func TestAmount_Div(t *testing.T) {
	a, _ := currency.NewAmount("99.99", "USD")

	for _, n := range []string{"INVALID", "0"} {
		_, err := a.Div(n)
		if e, ok := err.(currency.InvalidNumberError); ok {
			if e.Number != n {
				t.Errorf("got %v, want %v", e.Number, n)
			}
		} else {
			t.Errorf("got %T, want currency.InvalidNumberError", err)
		}
	}

	b, err := a.Div("3")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if b.String() != "33.33 USD" {
		t.Errorf("got %v, want 33.33 USD", b.String())
	}
	// Confirm that a is unchanged.
	if a.String() != "99.99 USD" {
		t.Errorf("got %v, want 99.99 USD", a.String())
	}

	// An amount equal to math.MaxInt64.
	d, _ := currency.NewAmount("9223372036854775807", "USD")
	e, err := d.Div("0.5")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if e.String() != "18446744073709551614 USD" {
		t.Errorf("got %v, want 18446744073709551614 USD", e.String())
	}
}

func TestAmount_Round(t *testing.T) {
	tests := []struct {
		number       string
		currencyCode string
		want         string
	}{
		{"12.345", "USD", "12.35"},
		{"12.345", "JPY", "12"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			a, _ := currency.NewAmount(tt.number, tt.currencyCode)
			b := a.Round()
			if b.Number() != tt.want {
				t.Errorf("got %v, want %v", b.Number(), tt.want)
			}
			// Confirm that a is unchanged.
			if a.Number() != tt.number {
				t.Errorf("got %v, want %v", a.Number(), tt.number)
			}
		})
	}
}

func TestAmount_RoundTo(t *testing.T) {
	tests := []struct {
		number string
		digits uint8
		mode   currency.RoundingMode
		want   string
	}{
		{"12.343", 2, currency.RoundHalfUp, "12.34"},
		{"12.345", 2, currency.RoundHalfUp, "12.35"},
		{"12.347", 2, currency.RoundHalfUp, "12.35"},

		{"12.343", 2, currency.RoundHalfDown, "12.34"},
		{"12.345", 2, currency.RoundHalfDown, "12.34"},
		{"12.347", 2, currency.RoundHalfDown, "12.35"},

		{"12.343", 2, currency.RoundUp, "12.35"},
		{"12.345", 2, currency.RoundUp, "12.35"},
		{"12.347", 2, currency.RoundUp, "12.35"},

		{"12.343", 2, currency.RoundDown, "12.34"},
		{"12.345", 2, currency.RoundDown, "12.34"},
		{"12.347", 2, currency.RoundDown, "12.34"},

		{"12.344", 2, currency.RoundHalfEven, "12.34"},
		{"12.345", 2, currency.RoundHalfEven, "12.34"},
		{"12.346", 2, currency.RoundHalfEven, "12.35"},

		{"12.334", 2, currency.RoundHalfEven, "12.33"},
		{"12.335", 2, currency.RoundHalfEven, "12.34"},
		{"12.336", 2, currency.RoundHalfEven, "12.34"},

		// Negative amounts.
		{"-12.345", 2, currency.RoundHalfUp, "-12.35"},
		{"-12.345", 2, currency.RoundHalfDown, "-12.34"},
		{"-12.345", 2, currency.RoundUp, "-12.35"},
		{"-12.345", 2, currency.RoundDown, "-12.34"},
		{"-12.345", 2, currency.RoundHalfEven, "-12.34"},
		{"-12.335", 2, currency.RoundHalfEven, "-12.34"},

		// More digits that the amount has.
		{"12.345", 4, currency.RoundHalfUp, "12.3450"},
		{"12.345", 4, currency.RoundHalfDown, "12.3450"},

		// Same number of digits that the amount has.
		{"12.345", 3, currency.RoundHalfUp, "12.345"},
		{"12.345", 3, currency.RoundHalfDown, "12.345"},
		{"12.345", 3, currency.RoundUp, "12.345"},
		{"12.345", 3, currency.RoundDown, "12.345"},

		// 0 digits.
		{"12.345", 0, currency.RoundHalfUp, "12"},
		{"12.345", 0, currency.RoundHalfDown, "12"},
		{"12.345", 0, currency.RoundUp, "13"},
		{"12.345", 0, currency.RoundDown, "12"},

		// Amounts larger than math.MaxInt64.
		{"12345678901234567890.0345", 3, currency.RoundHalfUp, "12345678901234567890.035"},
		{"12345678901234567890.0345", 3, currency.RoundHalfDown, "12345678901234567890.034"},
		{"12345678901234567890.0345", 3, currency.RoundUp, "12345678901234567890.035"},
		{"12345678901234567890.0345", 3, currency.RoundDown, "12345678901234567890.034"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			a, _ := currency.NewAmount(tt.number, "USD")
			b := a.RoundTo(tt.digits, tt.mode)
			if b.Number() != tt.want {
				t.Errorf("got %v, want %v", b.Number(), tt.want)
			}
			// Confirm that a is unchanged.
			if a.Number() != tt.number {
				t.Errorf("got %v, want %v", a.Number(), tt.number)
			}
		})
	}
}

func TestAmount_RoundToWithConcurrency(t *testing.T) {
	n := 2
	roundingModes := []currency.RoundingMode{
		currency.RoundHalfUp,
		currency.RoundHalfDown,
		currency.RoundUp,
		currency.RoundDown,
	}

	for _, roundingMode := range roundingModes {
		t.Run(fmt.Sprintf("rounding_mode_%d", roundingMode), func(t *testing.T) {
			t.Parallel()

			var allDone sync.WaitGroup
			allDone.Add(n)

			for i := 0; i < n; i++ {
				go func() {
					defer allDone.Done()
					amount, _ := currency.NewAmount("10.99", "EUR")
					amount.RoundTo(1, roundingMode)
				}()
			}

			allDone.Wait()
		})
	}
}

func TestAmount_Cmp(t *testing.T) {
	a, _ := currency.NewAmount("3.33", "USD")
	b, _ := currency.NewAmount("3.33", "EUR")
	_, err := a.Cmp(b)
	if e, ok := err.(currency.MismatchError); ok {
		if e.A != a {
			t.Errorf("got %v, want %v", e.A, a)
		}
		if e.B != b {
			t.Errorf("got %v, want %v", e.B, b)
		}
		wantError := `amounts "3.33 USD" and "3.33 EUR" have mismatched currency codes`
		if e.Error() != wantError {
			t.Errorf("got %v, want %v", e.Error(), wantError)
		}
	} else {
		t.Errorf("got %T, want currency.MismatchError", err)
	}

	tests := []struct {
		aNumber string
		bNumber string
		want    int
	}{
		{"3.33", "6.66", -1},
		{"3.33", "3.33", 0},
		{"6.66", "3.33", 1},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			a, _ := currency.NewAmount(tt.aNumber, "USD")
			b, _ := currency.NewAmount(tt.bNumber, "USD")
			got, err := a.Cmp(b)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAmount_Equal(t *testing.T) {
	tests := []struct {
		aNumber       string
		aCurrencyCode string
		bNumber       string
		bCurrencyCode string
		want          bool
	}{
		{"3.33", "USD", "6.66", "EUR", false},
		{"3.33", "USD", "3.33", "EUR", false},
		{"3.33", "USD", "3.33", "USD", true},
		{"3.33", "USD", "6.66", "USD", false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			a, _ := currency.NewAmount(tt.aNumber, tt.aCurrencyCode)
			b, _ := currency.NewAmount(tt.bNumber, tt.bCurrencyCode)
			got := a.Equal(b)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAmount_Checks(t *testing.T) {
	tests := []struct {
		number       string
		wantPositive bool
		wantNegative bool
		wantZero     bool
	}{
		{"9.99", true, false, false},
		{"-9.99", false, true, false},
		{"0", false, false, true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			a, _ := currency.NewAmount(tt.number, "USD")
			gotPositive := a.IsPositive()
			gotNegative := a.IsNegative()
			gotZero := a.IsZero()
			if gotPositive != tt.wantPositive {
				t.Errorf("positive: got %v, want %v", gotPositive, tt.wantPositive)
			}
			if gotNegative != tt.wantNegative {
				t.Errorf("negative: got %v, want %v", gotNegative, tt.wantNegative)
			}
			if gotZero != tt.wantZero {
				t.Errorf("zero: got %v, want %v", gotZero, tt.wantZero)
			}
		})
	}
}

func TestAmount_MarshalBinary(t *testing.T) {
	a, _ := currency.NewAmount("3.45", "USD")
	d, err := a.MarshalBinary()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	got := string(d)
	want := "USD3.45"
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAmount_UnmarshalBinary(t *testing.T) {
	d := []byte("US")
	a := &currency.Amount{}
	err := a.UnmarshalBinary(d)
	if e, ok := err.(currency.InvalidCurrencyCodeError); ok {
		if e.CurrencyCode != "US" {
			t.Errorf("got %v, want US", e.CurrencyCode)
		}
		wantError := `invalid currency code "US"`
		if e.Error() != wantError {
			t.Errorf("got %v, want %v", e.Error(), wantError)
		}
	} else {
		t.Errorf("got %T, want currency.InvalidCurrencyCodeError", err)
	}

	d = []byte("USD3,60")
	err = a.UnmarshalBinary(d)
	if e, ok := err.(currency.InvalidNumberError); ok {
		if e.Number != "3,60" {
			t.Errorf("got %v, want 3,60", e.Number)
		}
		wantError := `invalid number "3,60"`
		if e.Error() != wantError {
			t.Errorf("got %v, want %v", e.Error(), wantError)
		}
	} else {
		t.Errorf("got %T, want currency.InvalidNumberError", err)
	}

	d = []byte("XXX2.60")
	err = a.UnmarshalBinary(d)
	if e, ok := err.(currency.InvalidCurrencyCodeError); ok {
		if e.CurrencyCode != "XXX" {
			t.Errorf("got %v, want XXX", e.CurrencyCode)
		}
		wantError := `invalid currency code "XXX"`
		if e.Error() != wantError {
			t.Errorf("got %v, want %v", e.Error(), wantError)
		}
	} else {
		t.Errorf("got %T, want currency.InvalidCurrencyCodeError", err)
	}

	d = []byte("USD3.45")
	err = a.UnmarshalBinary(d)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if a.Number() != "3.45" {
		t.Errorf("got %v, want 3.45", a.Number())
	}
	if a.CurrencyCode() != "USD" {
		t.Errorf("got %v, want USD", a.CurrencyCode())
	}
}

func TestAmount_MarshalJSON(t *testing.T) {
	a, _ := currency.NewAmount("3.45", "USD")
	d, err := json.Marshal(a)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	got := string(d)
	want := `{"number":"3.45","currency":"USD"}`
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAmount_UnmarshalJSON(t *testing.T) {
	d := []byte(`{"number":"INVALID","currency":"USD"}`)
	unmarshalled := &currency.Amount{}
	err := json.Unmarshal(d, unmarshalled)
	if e, ok := err.(currency.InvalidNumberError); ok {
		if e.Number != "INVALID" {
			t.Errorf("got %v, want INVALID", e.Number)
		}
		wantError := `invalid number "INVALID"`
		if e.Error() != wantError {
			t.Errorf("got %v, want %v", e.Error(), wantError)
		}
	} else {
		t.Errorf("got %T, want currency.InvalidNumberError", err)
	}

	d = []byte(`{"number":"3.45","currency":"usd"}`)
	err = json.Unmarshal(d, unmarshalled)
	if e, ok := err.(currency.InvalidCurrencyCodeError); ok {
		if e.CurrencyCode != "usd" {
			t.Errorf("got %v, want usd", e.CurrencyCode)
		}
		wantError := `invalid currency code "usd"`
		if e.Error() != wantError {
			t.Errorf("got %v, want %v", e.Error(), wantError)
		}
	} else {
		t.Errorf("got %T, want currency.InvalidCurrencyCodeError", err)
	}

	d = []byte(`{"number":"3.45","currency":"USD"}`)
	err = json.Unmarshal(d, unmarshalled)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if unmarshalled.Number() != "3.45" {
		t.Errorf("got %v, want 3.45", unmarshalled.Number())
	}
	if unmarshalled.CurrencyCode() != "USD" {
		t.Errorf("got %v, want USD", unmarshalled.CurrencyCode())
	}
}

func TestAmount_Value(t *testing.T) {
	a, _ := currency.NewAmount("3.45", "USD")
	got, _ := a.Value()
	want := "(3.45,USD)"
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	var b currency.Amount
	got, _ = b.Value()
	want = "(0,)"
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAmount_Scan(t *testing.T) {
	tests := []struct {
		src              string
		wantNumber       string
		wantCurrencyCode string
		wantError        string
	}{
		{"", "0", "", ""},
		{"(3.45,USD)", "3.45", "USD", ""},
		{"(3.45,)", "0", "", `invalid currency code ""`},
		{"(,USD)", "0", "", `invalid number ""`},
		{"(0,)", "0", "", ""},
		{"(0,   )", "0", "", ""},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			var a currency.Amount
			err := a.Scan(tt.src)
			if a.Number() != tt.wantNumber {
				t.Errorf("number: got %v, want %v", a.Number(), tt.wantNumber)
			}
			if a.CurrencyCode() != tt.wantCurrencyCode {
				t.Errorf("currency code: got %v, want %v", a.CurrencyCode(), tt.wantCurrencyCode)
			}
			errStr := ""
			if err != nil {
				errStr = err.Error()
			}
			if errStr != tt.wantError {
				t.Errorf("error: got %v, want %v", errStr, tt.wantError)
			}
		})
	}
}
