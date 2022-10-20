// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package currency_test

import (
	"testing"

	"github.com/bojanz/currency"
)

func TestFormatter_Locale(t *testing.T) {
	locale := currency.NewLocale("fr-FR")
	formatter := currency.NewFormatter(locale)
	got := formatter.Locale().String()
	if got != "fr-FR" {
		t.Errorf("got %v, want fr-FR", got)
	}
}

func TestFormatter_Format(t *testing.T) {
	tests := []struct {
		number       string
		currencyCode string
		localeID     string
		want         string
	}{
		{"1234.59", "USD", "en-US", "$1,234.59"},
		{"1234.59", "USD", "en-CA", "US$1,234.59"},
		{"1234.59", "USD", "de-CH", "$\u00a01’234.59"},
		{"1234.59", "USD", "sr", "1.234,59\u00a0US$"},

		{"-1234.59", "USD", "en-US", "-$1,234.59"},
		{"-1234.59", "USD", "en-CA", "-US$1,234.59"},
		{"-1234.59", "USD", "de-CH", "$-1’234.59"},
		{"-1234.59", "USD", "sr", "-1.234,59\u00a0US$"},

		{"1234.00", "EUR", "en", "€1,234.00"},
		{"1234.00", "EUR", "de-CH", "€\u00a01’234.00"},
		{"1234.00", "EUR", "sr", "1.234,00\u00a0€"},

		{"1234.00", "CHF", "en", "CHF\u00a01,234.00"},
		{"1234.00", "CHF", "de-CH", "CHF\u00a01’234.00"},
		{"1234.00", "CHF", "sr", "1.234,00\u00a0CHF"},

		// Arabic digits.
		{"12345678.90", "USD", "ar", "\u200f١٢٬٣٤٥٬٦٧٨٫٩٠\u00a0US$"},
		// Arabic extended (Persian) digits.
		{"12345678.90", "USD", "fa", "\u200e$۱۲٬۳۴۵٬۶۷۸٫۹۰"},
		// Bengali digits.
		{"12345678.90", "USD", "bn", "১,২৩,৪৫,৬৭৮.৯০\u00a0US$"},
		// Devanagari digits.
		{"12345678.90", "USD", "ne", "US$\u00a0१,२३,४५,६७८.९०"},
		// Myanmar (Burmese) digits.
		{"12345678.90", "USD", "my", "၁၂,၃၄၅,၆၇၈.၉၀\u00a0US$"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			amount, _ := currency.NewAmount(tt.number, tt.currencyCode)
			locale := currency.NewLocale(tt.localeID)
			formatter := currency.NewFormatter(locale)
			got := formatter.Format(amount)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatter_Grouping(t *testing.T) {
	tests := []struct {
		number       string
		currencyCode string
		localeID     string
		NoGrouping   bool
		want         string
	}{
		{"123.99", "USD", "en", false, "$123.99"},
		{"1234.99", "USD", "en", false, "$1,234.99"},
		{"1234567.99", "USD", "en", false, "$1,234,567.99"},

		{"123.99", "USD", "en", true, "$123.99"},
		{"1234.99", "USD", "en", true, "$1234.99"},
		{"1234567.99", "USD", "en", true, "$1234567.99"},

		// The "es" locale has a different minGroupingSize.
		{"123.99", "USD", "es", false, "123,99\u00a0US$"},
		{"1234.99", "USD", "es", false, "1234,99\u00a0US$"},
		{"12345.99", "USD", "es", false, "12.345,99\u00a0US$"},
		{"1234567.99", "USD", "es", false, "1.234.567,99\u00a0US$"},

		// The "hi" locale has a different secondaryGroupingSize.
		{"123.99", "USD", "hi", false, "$123.99"},
		{"1234.99", "USD", "hi", false, "$1,234.99"},
		{"1234567.99", "USD", "hi", false, "$12,34,567.99"},
		{"12345678.99", "USD", "hi", false, "$1,23,45,678.99"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			amount, _ := currency.NewAmount(tt.number, tt.currencyCode)
			locale := currency.NewLocale(tt.localeID)
			formatter := currency.NewFormatter(locale)
			formatter.NoGrouping = tt.NoGrouping
			got := formatter.Format(amount)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatter_PlusSign(t *testing.T) {
	tests := []struct {
		number       string
		currencyCode string
		localeID     string
		AddPlusSign  bool
		want         string
	}{
		{"123.99", "USD", "en", false, "$123.99"},
		{"123.99", "USD", "en", true, "+$123.99"},

		{"123.99", "USD", "de-CH", false, "$\u00a0123.99"},
		{"123.99", "USD", "de-CH", true, "$+123.99"},

		{"123.99", "USD", "fr-FR", false, "123,99\u00a0$US"},
		{"123.99", "USD", "fr-FR", true, "+123,99\u00a0$US"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			amount, _ := currency.NewAmount(tt.number, tt.currencyCode)
			locale := currency.NewLocale(tt.localeID)
			formatter := currency.NewFormatter(locale)
			formatter.AddPlusSign = tt.AddPlusSign
			got := formatter.Format(amount)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatter_Digits(t *testing.T) {
	tests := []struct {
		number       string
		currencyCode string
		localeID     string
		minDigits    uint8
		maxDigits    uint8
		want         string
	}{
		{"59", "KRW", "en", currency.DefaultDigits, 6, "₩59"},
		{"59", "USD", "en", currency.DefaultDigits, 6, "$59.00"},
		{"59", "OMR", "en", currency.DefaultDigits, 6, "OMR\u00a059.000"},

		{"59.6789", "KRW", "en", 0, currency.DefaultDigits, "₩60"},
		{"59.6789", "USD", "en", 0, currency.DefaultDigits, "$59.68"},
		{"59.6789", "OMR", "en", 0, currency.DefaultDigits, "OMR\u00a059.679"},

		// minDigits:0 strips all trailing zeroes.
		{"59", "USD", "en", 0, 6, "$59"},
		{"59.5", "USD", "en", 0, 6, "$59.5"},
		{"59.56", "USD", "en", 0, 6, "$59.56"},

		// minDigits can't override maxDigits.
		{"59.5", "USD", "en", 3, 2, "$59.50"},
		{"59.567", "USD", "en", 3, 2, "$59.57"},

		// maxDigits rounds the number.
		{"59.5", "USD", "en", 2, 3, "$59.50"},
		{"59.567", "USD", "en", 2, 3, "$59.567"},
		{"59.5678", "USD", "en", 2, 3, "$59.568"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			amount, _ := currency.NewAmount(tt.number, tt.currencyCode)
			locale := currency.NewLocale(tt.localeID)
			formatter := currency.NewFormatter(locale)
			formatter.MinDigits = tt.minDigits
			formatter.MaxDigits = tt.maxDigits
			got := formatter.Format(amount)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatter_RoundingMode(t *testing.T) {
	tests := []struct {
		number       string
		currencyCode string
		localeID     string
		roundingMode currency.RoundingMode
		want         string
	}{
		{"1234.453", "USD", "en", currency.RoundHalfUp, "$1,234.45"},
		{"1234.455", "USD", "en", currency.RoundHalfUp, "$1,234.46"},
		{"1234.456", "USD", "en", currency.RoundHalfUp, "$1,234.46"},

		{"1234.453", "USD", "en", currency.RoundHalfDown, "$1,234.45"},
		{"1234.455", "USD", "en", currency.RoundHalfDown, "$1,234.45"},
		{"1234.457", "USD", "en", currency.RoundHalfDown, "$1,234.46"},

		{"1234.453", "USD", "en", currency.RoundUp, "$1,234.46"},
		{"1234.455", "USD", "en", currency.RoundUp, "$1,234.46"},
		{"1234.457", "USD", "en", currency.RoundUp, "$1,234.46"},

		{"1234.453", "USD", "en", currency.RoundDown, "$1,234.45"},
		{"1234.455", "USD", "en", currency.RoundDown, "$1,234.45"},
		{"1234.457", "USD", "en", currency.RoundDown, "$1,234.45"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			amount, _ := currency.NewAmount(tt.number, tt.currencyCode)
			locale := currency.NewLocale(tt.localeID)
			formatter := currency.NewFormatter(locale)
			formatter.MaxDigits = currency.DefaultDigits
			formatter.RoundingMode = tt.roundingMode
			got := formatter.Format(amount)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatter_CurrencyDisplay(t *testing.T) {
	tests := []struct {
		number          string
		currencyCode    string
		localeID        string
		currencyDisplay currency.Display
		want            string
	}{
		{"1234.59", "USD", "en", currency.DisplaySymbol, "$1,234.59"},
		{"1234.59", "USD", "en", currency.DisplayCode, "USD\u00a01,234.59"},
		{"1234.59", "USD", "en", currency.DisplayNone, "1,234.59"},

		{"1234.59", "USD", "de-AT", currency.DisplaySymbol, "$\u00a01.234,59"},
		{"1234.59", "USD", "de-AT", currency.DisplayCode, "USD\u00a01.234,59"},
		{"1234.59", "USD", "de-AT", currency.DisplayNone, "1.234,59"},

		{"1234.59", "USD", "sr-Latn", currency.DisplaySymbol, "1.234,59\u00a0US$"},
		{"1234.59", "USD", "sr-Latn", currency.DisplayCode, "1.234,59\u00a0USD"},
		{"1234.59", "USD", "sr-Latn", currency.DisplayNone, "1.234,59"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			amount, _ := currency.NewAmount(tt.number, tt.currencyCode)
			locale := currency.NewLocale(tt.localeID)
			formatter := currency.NewFormatter(locale)
			formatter.CurrencyDisplay = tt.currencyDisplay
			got := formatter.Format(amount)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatter_SymbolMap(t *testing.T) {
	locale := currency.NewLocale("en")
	formatter := currency.NewFormatter(locale)
	formatter.SymbolMap["USD"] = "US$"
	formatter.SymbolMap["EUR"] = "EU"

	amount, _ := currency.NewAmount("6.99", "USD")
	got := formatter.Format(amount)
	if got != "US$6.99" {
		t.Errorf("got %v, want US$6.99", got)
	}

	amount, _ = currency.NewAmount("6.99", "EUR")
	got = formatter.Format(amount)
	if got != "EU\u00a06.99" {
		t.Errorf("got %v, want EU\u00a06.99", got)
	}
}

func TestFormatter_Parse(t *testing.T) {
	tests := []struct {
		s            string
		currencyCode string
		localeID     string
		want         string
	}{
		{"$1,234.59", "USD", "en", "1234.59"},
		{"USD\u00a01,234.59", "USD", "en", "1234.59"},
		{"1,234.59", "USD", "en", "1234.59"},
		{"1234.59", "USD", "en", "1234.59"},
		{"+1234.59", "USD", "en", "1234.59"},
		{"1234", "USD", "en", "1234"},

		{"-$1,234.59", "USD", "en", "-1234.59"},
		{"-USD\u00a01,234.59", "USD", "en", "-1234.59"},
		{"-1,234.59", "USD", "en", "-1234.59"},
		{"-1234.59", "USD", "en", "-1234.59"},

		{"€\u00a01.234,00", "EUR", "de-AT", "1234.00"},
		{"EUR\u00a01.234,00", "EUR", "de-AT", "1234.00"},
		{"1.234,00", "EUR", "de-AT", "1234.00"},
		{"1234,00", "EUR", "de-AT", "1234.00"},

		// Arabic digits.
		{"١٢٬٣٤٥٬٦٧٨٫٩٠\u00a0US$", "USD", "ar", "12345678.90"},
		// Arabic extended (Persian) digits.
		{"\u200e$۱۲٬۳۴۵٬۶۷۸٫۹۰", "USD", "fa", "12345678.90"},
		// Bengali digits.
		{"১,২৩,৪৫,৬৭৮.৯০\u00a0US$", "USD", "bn", "12345678.90"},
		// Devanagari digits.
		{"US$\u00a0१,२३,४५,६७८.९०", "USD", "ne", "12345678.90"},
		// Myanmar (Burmese) digits.
		{"၁၂,၃၄၅,၆၇၈.၉၀\u00a0US$", "USD", "my", "12345678.90"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			locale := currency.NewLocale(tt.localeID)
			formatter := currency.NewFormatter(locale)
			got, err := formatter.Parse(tt.s, tt.currencyCode)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if got.Number() != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
			if got.CurrencyCode() != tt.currencyCode {
				t.Errorf("got %v, want %v", got.CurrencyCode(), tt.currencyCode)
			}
		})
	}
}
