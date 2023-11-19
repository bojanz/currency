// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

// Package currency handles currency amounts, provides currency information and formatting.
package currency

import "sort"

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
			if contains(s.locales, localeID) {
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

// contains returns whether the sorted slice a contains x.
// The slice must be sorted in ascending order.
func contains(a []string, x string) bool {
	n := len(a)
	if n < 10 {
		// Linear search is faster with a small number of elements.
		for _, v := range a {
			if v == x {
				return true
			}
		}
	} else {
		// Binary search is faster with a large number of elements.
		i := sort.SearchStrings(a, x)
		if i < n && a[i] == x {
			return true
		}
	}
	return false
}
