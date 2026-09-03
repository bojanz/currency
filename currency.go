// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

// Package currency handles currency amounts, provides currency information and formatting.
package currency

import "strings"

// DefaultDigits is a placeholder for each currency's number of fraction digits.
const DefaultDigits uint8 = 255

// ForCountryCode returns the currency code for a country code.
func ForCountryCode(countryCode string) (currencyCode string, ok bool) {
	currencyCode, ok = countryCurrencies[countryCode]

	return currencyCode, ok
}

// GetCurrencyCodes returns all known currency codes.
func GetCurrencyCodes() []string {
	return currencyCodes
}

// IsValid checks whether a currency code is valid.
//
// An empty currency code is considered valid.
func IsValid(currencyCode string) bool {
	if currencyCode == "" {
		return true
	}
	_, ok := currencies[currencyCode]

	return ok
}

// GetNumericCode returns the numeric code for a currency code.
func GetNumericCode(currencyCode string) (numericCode string, ok bool) {
	if currencyCode == "" || !IsValid(currencyCode) {
		return "000", false
	}
	return currencies[currencyCode].numericCode, true
}

// GetDigits returns the number of fraction digits for a currency code.
func GetDigits(currencyCode string) (digits uint8, ok bool) {
	if currencyCode == "" || !IsValid(currencyCode) {
		return 0, false
	}
	return currencies[currencyCode].digits, true
}

// GetSymbol returns the symbol for a currency code.
func GetSymbol(currencyCode string, locale Locale) (symbol string, ok bool) {
	if currencyCode == "" || !IsValid(currencyCode) {
		return currencyCode, false
	}
	symbols, ok := currencySymbols[currencyCode]
	if !ok {
		return currencyCode, true
	}
	enLocale := Locale{Language: "en"}
	enUSLocale := Locale{Language: "en", Territory: "US"}
	if locale == enLocale || locale == enUSLocale || locale.IsEmpty() {
		// The "en"/"en-US" symbol is always first.
		return symbols[0].symbol, true
	}

	for {
		localeID := locale.String()
		for _, s := range symbols {
			if containsLocale(s.locales, localeID) {
				symbol = s.symbol
				break
			}
		}
		if symbol != "" {
			break
		}
		locale = locale.GetParent()
		if locale.IsEmpty() {
			break
		}
	}

	return symbol, true
}

// getFormat returns the format for a locale.
func getFormat(locale Locale) currencyFormat {
	// CLDR considers "en" and "en-US" to be equivalent.
	// Fall back immediately for better performance
	enUSLocale := Locale{Language: "en", Territory: "US"}
	if locale == enUSLocale || locale.IsEmpty() {
		return currencyFormats["en"]
	}

	var format currencyFormat
	for {
		localeID := locale.String()
		if cf, ok := currencyFormats[localeID]; ok {
			format = cf
			break
		}
		locale = locale.GetParent()
		if locale.IsEmpty() {
			break
		}
	}

	return format
}

// containsLocale returns whether the space-separated list contains localeID.
//
// An empty localeID is never contained.
func containsLocale(list, localeID string) bool {
	if localeID == "" {
		return false
	}
	for i := 0; ; {
		j := strings.Index(list[i:], localeID)
		if j < 0 {
			return false
		}
		// Confirm the match spans a whole entry, so that "en" doesn't match "en-AU".
		start, end := i+j, i+j+len(localeID)
		if (start == 0 || list[start-1] == ' ') && (end == len(list) || list[end] == ' ') {
			return true
		}
		i = start + 1
	}
}

// Definition contains information for registering a currency.
type Definition struct {
	// NumericCode is a three-digit code such as "999".
	NumericCode string

	// Digits is the number of fraction digits.
	Digits uint8

	// DefaultSymbol is the default symbol, used for all locales.
	//
	// When overriding an existing currency, keep DefaultSymbol
	// empty to retain the built-in locale-specific symbols.
	DefaultSymbol string
}

// Register adds a currency to the internal registry.
//
// This can be a non-ISO currency such as BTC, or an existing
// currency for which we want to override the predefined data.
func Register(currencyCode string, d Definition) {
	if currencyCode == "" {
		return
	}

	if _, ok := currencies[currencyCode]; !ok {
		// Overriding an existing currency must not duplicate its code.
		currencyCodes = append(currencyCodes, currencyCode)
	}
	currencies[currencyCode] = currencyInfo{
		numericCode: d.NumericCode,
		digits:      d.Digits,
	}

	if d.DefaultSymbol != "" {
		currencySymbols[currencyCode] = []symbolInfo{
			{symbol: d.DefaultSymbol, locales: "en"},
		}
	}
}
