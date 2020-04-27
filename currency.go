// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package currency

import "sort"

// GetCurrencyCodes returns all known currency codes.
func GetCurrencyCodes() []string {
	return currencyCodes
}

// IsValid checks whether a currencyCode is valid.
func IsValid(currencyCode string) bool {
	if len(currencyCode) != 3 {
		return false
	}
	_, ok := currencies[currencyCode]

	return ok
}

// GetNumericCode returns the numeric code for a currencyCode.
func GetNumericCode(currencyCode string) (numericCode string, ok bool) {
	if !IsValid(currencyCode) {
		return "000", false
	}
	return currencies[currencyCode].numericCode, true
}

// GetDigits returns the number of fraction digits for a currencyCode.
func GetDigits(currencyCode string) (digits byte, ok bool) {
	if !IsValid(currencyCode) {
		return 0, false
	}
	return currencies[currencyCode].digits, true
}

// GetSymbol returns the symbol for a currencyCode.
func GetSymbol(currencyCode string, locale Locale) (symbol string, ok bool) {
	if !IsValid(currencyCode) {
		return currencyCode, false
	}
	symbols, ok := currencySymbols[currencyCode]
	if !ok {
		return currencyCode, true
	}
	enLocale := Locale{Language: "en"}
	enUSLocale := Locale{Language: "en", Region: "US"}
	if locale == enLocale || locale == enUSLocale {
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
		if i := sort.SearchStrings(a, x); i < n {
			return true
		}
	}
	return false
}
