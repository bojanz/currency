// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package currency_test

import (
	"testing"

	"github.com/bojanz/currency"
)

func TestForCountryCode(t *testing.T) {
	tests := []struct {
		countryCode      string
		wantCurrencyCode string
		wantOK           bool
	}{
		{"FR", "EUR", true},
		{"RS", "RSD", true},
		{"XX", "", false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			gotCurrencyCode, gotOK := currency.ForCountryCode(tt.countryCode)
			if gotOK != tt.wantOK {
				t.Errorf("got %v, want %v", gotOK, tt.wantOK)
			}
			if gotCurrencyCode != tt.wantCurrencyCode {
				t.Errorf("got %q, want %q", gotCurrencyCode, tt.wantCurrencyCode)
			}
		})
	}
}

func TestGetCurrencyCodes(t *testing.T) {
	currencyCodes := currency.GetCurrencyCodes()
	var got [10]string
	copy(got[:], currencyCodes[0:10])
	want := [10]string{"AUD", "CAD", "CHF", "EUR", "GBP", "JPY", "NOK", "NZD", "SEK", "USD"}
	// Confirm that the first 10 currency codes are the "G10" ones.
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		currencyCode string
		want         bool
	}{
		{"", true},
		{"INVALID", false},
		{"XXX", false},
		{"usd", false},
		{"USD", true},
		{"EUR", true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := currency.IsValid(tt.currencyCode)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNumericCode(t *testing.T) {
	numericCode, ok := currency.GetNumericCode("USD")
	if !ok {
		t.Errorf("got %v, want true", ok)
	}
	if numericCode != "840" {
		t.Errorf("got %v, want 840", numericCode)
	}

	// Non-existent currency code.
	numericCode, ok = currency.GetNumericCode("XXX")
	if ok {
		t.Errorf("got %v, want false", ok)
	}
	if numericCode != "000" {
		t.Errorf("got %v, want 000", numericCode)
	}
}

func TestGetDigits(t *testing.T) {
	digits, ok := currency.GetDigits("USD")
	if !ok {
		t.Errorf("got %v, want true", ok)
	}
	if digits != 2 {
		t.Errorf("got %v, want 2", digits)
	}

	// Non-existent currency code.
	digits, ok = currency.GetDigits("XXX")
	if ok {
		t.Errorf("got %v, want false", ok)
	}
	if digits != 0 {
		t.Errorf("got %v, want 0", digits)
	}
}

func TestGetSymbol(t *testing.T) {
	tests := []struct {
		currencyCode string
		locale       currency.Locale
		wantSymbol   string
		wantOk       bool
	}{
		{"XXX", currency.NewLocale("en"), "XXX", false},
		{"usd", currency.NewLocale("en"), "usd", false},
		{"CHF", currency.NewLocale("en"), "CHF", true},
		{"USD", currency.NewLocale("en"), "$", true},
		{"USD", currency.NewLocale("en-US"), "$", true},
		{"USD", currency.NewLocale("en-AU"), "US$", true},
		{"USD", currency.NewLocale("es"), "US$", true},
		{"USD", currency.NewLocale("es-ES"), "US$", true},
		// An empty locale should use "en" data.
		{"USD", currency.NewLocale(""), "$", true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			gotSymbol, gotOk := currency.GetSymbol(tt.currencyCode, tt.locale)
			if gotSymbol != tt.wantSymbol {
				t.Errorf("got %v, want %v", gotSymbol, tt.wantSymbol)
			}
			if gotOk != tt.wantOk {
				t.Errorf("got %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_Register_NoCurrencyCode(t *testing.T) {
	currency.Register("", currency.Definition{
		NumericCode: "123",
	})

	numericCode, ok := currency.GetNumericCode("")
	if ok {
		t.Errorf("got %v, want false", ok)
	}
	if numericCode != "000" {
		t.Errorf("got %v, want 000", numericCode)
	}
}

func Test_Register_New(t *testing.T) {
	currency.Register("BTC", currency.Definition{
		NumericCode:   "1000",
		Digits:        8,
		DefaultSymbol: "₿",
	})

	if !currency.IsValid("BTC") {
		t.Errorf("got false, want true")
	}
	numericCode, ok := currency.GetNumericCode("BTC")
	if !ok {
		t.Errorf("got %v, want true", ok)
	}
	if numericCode != "1000" {
		t.Errorf("got %v, want 1000", numericCode)
	}

	d, ok := currency.GetDigits("BTC")
	if !ok {
		t.Errorf("got %v, want true", ok)
	}
	if d != 8 {
		t.Errorf("got %v, want 8", d)
	}

	symbol, _ := currency.GetSymbol("BTC", currency.NewLocale("en"))
	if symbol != "₿" {
		t.Errorf("got %v, want ₿", symbol)
	}
	symbol, _ = currency.GetSymbol("BTC", currency.NewLocale("es-ES"))
	if symbol != "₿" {
		t.Errorf("got %v, want ₿", symbol)
	}
}

func Test_Register_OverrideExisting(t *testing.T) {
	currency.Register("CAD", currency.Definition{
		NumericCode: "125",
		Digits:      3,
	})

	if !currency.IsValid("CAD") {
		t.Errorf("got false, want true")
	}
	numericCode, ok := currency.GetNumericCode("CAD")
	if !ok {
		t.Errorf("got %v, want true", ok)
	}
	if numericCode != "125" {
		t.Errorf("got %v, want 125", numericCode)
	}

	d, ok := currency.GetDigits("CAD")
	if !ok {
		t.Errorf("got %v, want true", ok)
	}
	if d != 3 {
		t.Errorf("got %v, want 3", d)
	}

	symbol, _ := currency.GetSymbol("CAD", currency.NewLocale("en"))
	if symbol != "CA$" {
		t.Errorf("got %v, want CA$", symbol)
	}
	symbol, _ = currency.GetSymbol("CAD", currency.NewLocale("fr"))
	if symbol != "$CA" {
		t.Errorf("got %v, want $CA", symbol)
	}

	// Override the symbols.
	currency.Register("CAD", currency.Definition{
		NumericCode:   "125",
		Digits:        3,
		DefaultSymbol: "$$",
	})

	symbol, _ = currency.GetSymbol("CAD", currency.NewLocale("en"))
	if symbol != "$$" {
		t.Errorf("got %v, want $$", symbol)
	}
	symbol, _ = currency.GetSymbol("CAD", currency.NewLocale("fr"))
	if symbol != "$$" {
		t.Errorf("got %v, want $$", symbol)
	}
}
