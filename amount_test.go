// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package currency_test

import (
	"testing"

	"github.com/bojanz/currency"
)

func TestNewAmount(t *testing.T) {
	_, err := currency.NewAmount("INVALID", "USD")
	if e, ok := err.(currency.InvalidNumberError); ok {
		if e.Op != "NewAmount" {
			t.Errorf("got %v, want NewAmount", e.Op)
		}
		if e.Number != "INVALID" {
			t.Errorf("got %v, want INVALID", e.Number)
		}
	} else {
		t.Errorf("got %T, want currency.InvalidNumberError", err)
	}

	_, err = currency.NewAmount("10.99", "usd")
	if e, ok := err.(currency.InvalidCurrencyCodeError); ok {
		if e.Op != "NewAmount" {
			t.Errorf("got %v, want NewAmount", e.Op)
		}
		if e.CurrencyCode != "usd" {
			t.Errorf("got %v, want usd", e.CurrencyCode)
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

func TestAmount_ToMinorUnits(t *testing.T) {
	tests := []struct {
		number string
		want   int64
	}{
		{"20.99", 2099},
		// Number with additional decimals.
		{"12.3564", 1236},
		// Number with no decimals.
		{"50", 5000},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			a, _ := currency.NewAmount(tt.number, "USD")
			got := a.ToMinorUnits()
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
		if e.Op != "Amount.Convert" {
			t.Errorf("got %v, want Amount.Convert", e.Op)
		}
		if e.CurrencyCode != "eur" {
			t.Errorf("got %v, want eur", e.CurrencyCode)
		}
	} else {
		t.Errorf("got %T, want currency.InvalidCurrencyCodeError", err)
	}

	_, err = a.Convert("EUR", "INVALID")
	if e, ok := err.(currency.InvalidNumberError); ok {
		if e.Op != "Amount.Convert" {
			t.Errorf("got %v, want Amount.Convert", e.Op)
		}
		if e.Number != "INVALID" {
			t.Errorf("got %v, want INVALID", e.Number)
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
}

func TestAmount_Add(t *testing.T) {
	a, _ := currency.NewAmount("20.99", "USD")
	b, _ := currency.NewAmount("3.50", "USD")
	x, _ := currency.NewAmount("99.99", "EUR")

	_, err := a.Add(x)
	if e, ok := err.(currency.MismatchError); ok {
		if e.Op != "Amount.Add" {
			t.Errorf("got %v, want Amount.Add", e.Op)
		}
		if e.A != a {
			t.Errorf("got %v, want %v", e.A, a)
		}
		if e.B != x {
			t.Errorf("got %v, want %v", e.B, x)
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
}

func TestAmount_Sub(t *testing.T) {
	a, _ := currency.NewAmount("20.99", "USD")
	b, _ := currency.NewAmount("3.50", "USD")
	x, _ := currency.NewAmount("99.99", "EUR")

	_, err := a.Sub(x)
	if e, ok := err.(currency.MismatchError); ok {
		if e.Op != "Amount.Sub" {
			t.Errorf("got %v, want Amount.Sub", e.Op)
		}
		if e.A != a {
			t.Errorf("got %v, want %v", e.A, a)
		}
		if e.B != x {
			t.Errorf("got %v, want %v", e.B, x)
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
}

func TestAmount_Mul(t *testing.T) {
	a, _ := currency.NewAmount("20.99", "USD")

	_, err := a.Mul("INVALID")
	if e, ok := err.(currency.InvalidNumberError); ok {
		if e.Op != "Amount.Mul" {
			t.Errorf("got %v, want Amount.Mul", e.Op)
		}
		if e.Number != "INVALID" {
			t.Errorf("got %v, want INVALID", e.Number)
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
}

func TestAmount_Div(t *testing.T) {
	a, _ := currency.NewAmount("99.99", "USD")

	for _, n := range []string{"INVALID", "0"} {
		_, err := a.Div(n)
		if e, ok := err.(currency.InvalidNumberError); ok {
			if e.Op != "Amount.Div" {
				t.Errorf("got %v, want Amount.Div", e.Op)
			}
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
		digits byte
		want   string
	}{
		{"12.345", 4, "12.3450"},
		{"12.345", 3, "12.345"},
		{"12.345", 2, "12.35"},
		{"12.345", 1, "12.3"},
		{"12.345", 0, "12"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			a, _ := currency.NewAmount(tt.number, "USD")
			b := a.RoundTo(tt.digits)
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

func TestAmount_Cmp(t *testing.T) {
	a, _ := currency.NewAmount("3.33", "USD")
	b, _ := currency.NewAmount("3.33", "EUR")
	_, err := a.Cmp(b)
	if e, ok := err.(currency.MismatchError); ok {
		if e.Op != "Amount.Cmp" {
			t.Errorf("got %v, want Amount.Cmp", e.Op)
		}
		if e.A != a {
			t.Errorf("got %v, want %v", e.A, a)
		}
		if e.B != b {
			t.Errorf("got %v, want %v", e.B, b)
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
