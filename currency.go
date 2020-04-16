// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package currency

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
