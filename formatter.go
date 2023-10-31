// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package currency

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Display represents the currency display type.
type Display uint8

const (
	// DisplaySymbol shows the currency symbol.
	DisplaySymbol Display = iota
	// DisplayCode shows the currency code.
	DisplayCode
	// DisplayNone shows nothing, hiding the currency.
	DisplayNone
)

var localDigits = map[numberingSystem]string{
	numArab:    "٠١٢٣٤٥٦٧٨٩",
	numArabExt: "۰۱۲۳۴۵۶۷۸۹",
	numBeng:    "০১২৩৪৫৬৭৮৯",
	numDeva:    "०१२३४५६७८९",
	numMymr:    "၀၁၂၃၄၅၆၇၈၉",
}

// Formatter formats and parses currency amounts.
type Formatter struct {
	locale Locale
	format currencyFormat
	// AccountingStyle formats the amount using the accounting style.
	// For example, "-3.00 USD" in the "en" locale is formatted as "($3.00)" instead of "-$3.00".
	// Defaults to false.
	AccountingStyle bool
	// AddPlusSign inserts the plus sign in front of positive amounts.
	// Defaults to false.
	AddPlusSign bool
	// NoGrouping turns off grouping of major digits.
	// Defaults to false.
	NoGrouping bool
	// MinDigits specifies the minimum number of fraction digits.
	// All zeroes past the minimum will be removed (0 => no trailing zeroes).
	// Defaults to currency.DefaultDigits (e.g. 2 for USD, 0 for RSD).
	MinDigits uint8
	// MaxDigits specifies the maximum number of fraction digits.
	// Formatted amounts will be rounded to this number of digits.
	// Defaults to 6, so that most amounts are shown as-is (without rounding).
	MaxDigits uint8
	// RoundingMode specifies how the formatted amount will be rounded.
	// Defaults to currency.RoundHalfUp.
	RoundingMode RoundingMode
	// CurrencyDisplay specifies how the currency will be displayed (symbol/code/none).
	// Defaults to currency.DisplaySymbol.
	CurrencyDisplay Display
	// SymbolMap specifies custom symbols for individual currency codes.
	// For example, "USD": "$" means that the $ symbol will be used even if
	// the current locale's symbol is different ("US$", "$US", etc).
	SymbolMap map[string]string
}

// NewFormatter creates a new formatter for the given locale.
func NewFormatter(locale Locale) *Formatter {
	f := &Formatter{
		locale:          locale,
		format:          getFormat(locale),
		MinDigits:       DefaultDigits,
		MaxDigits:       6,
		RoundingMode:    RoundHalfUp,
		CurrencyDisplay: DisplaySymbol,
		SymbolMap:       make(map[string]string),
	}
	return f
}

// Locale returns the locale.
func (f *Formatter) Locale() Locale {
	return f.locale
}

// Format formats a currency amount.
func (f *Formatter) Format(amount Amount) string {
	pattern := f.getPattern(amount)
	if amount.IsNegative() {
		// The minus sign will be provided by the pattern.
		amount, _ = amount.Mul("-1")
	}
	formattedNumber := f.formatNumber(amount)
	formattedCurrency := f.formatCurrency(amount.CurrencyCode())
	if formattedCurrency != "" {
		// CLDR requires having a space between the letters
		// in a currency symbol and adjacent numbers.
		if strings.Contains(pattern, "0¤") {
			r, _ := utf8.DecodeRuneInString(formattedCurrency)
			if unicode.IsLetter(r) {
				formattedCurrency = "\u00a0" + formattedCurrency
			}
		} else if strings.Contains(pattern, "¤0") {
			r, _ := utf8.DecodeLastRuneInString(formattedCurrency)
			if unicode.IsLetter(r) {
				formattedCurrency = formattedCurrency + "\u00a0"
			}
		}
	}

	replacements := []string{
		"0.00", formattedNumber,
		"+", f.format.plusSign,
		"-", f.format.minusSign,
	}
	if formattedCurrency == "" {
		// Many patterns have a non-breaking space between
		// the number and currency, not needed in this case.
		replacements = append(replacements, "\u00a0¤", "", "¤\u00a0", "", "¤", "")
	} else {
		replacements = append(replacements, "¤", formattedCurrency)
	}
	r := strings.NewReplacer(replacements...)

	return r.Replace(pattern)
}

// Parse parses a formatted amount.
func (f *Formatter) Parse(s, currencyCode string) (Amount, error) {
	symbol, _ := GetSymbol(currencyCode, f.locale)
	replacements := []string{
		f.format.decimalSeparator, ".",
		f.format.groupingSeparator, "",
		f.format.plusSign, "+",
		f.format.minusSign, "-",
		symbol, "",
		currencyCode, "",
		"\u200e", "",
		"\u200f", "",
		"\u00a0", "",
		" ", "",
	}
	if f.format.numberingSystem != numLatn {
		digits := localDigits[f.format.numberingSystem]
		for i, v := range strings.Split(digits, "") {
			replacements = append(replacements, v, strconv.Itoa(i))
		}
	}
	if f.AccountingStyle {
		replacements = append(replacements, "(", "-", ")", "")
	}
	r := strings.NewReplacer(replacements...)
	n := r.Replace(s)

	return NewAmount(n, currencyCode)
}

// getPattern returns a positive or negative pattern for a currency amount.
func (f *Formatter) getPattern(amount Amount) string {
	var patterns []string
	if f.usesAccountingPattern() {
		patterns = strings.Split(f.format.accountingPattern, ";")
	} else {
		patterns = strings.Split(f.format.standardPattern, ";")
	}

	switch {
	case amount.IsNegative():
		if len(patterns) == 1 {
			return "-" + patterns[0]
		}
		return patterns[1]
	case f.AddPlusSign:
		if len(patterns) == 1 || f.usesAccountingPattern() {
			return "+" + patterns[0]
		}
		return strings.Replace(patterns[1], "-", "+", 1)
	default:
		return patterns[0]
	}
}

// usesAccountingPattern returns whether the formatter needs to use the accounting pattern.
func (f *Formatter) usesAccountingPattern() bool {
	return f.AccountingStyle && f.format.accountingPattern != ""
}

// formatNumber formats the number for display.
func (f *Formatter) formatNumber(amount Amount) string {
	minDigits := f.MinDigits
	if minDigits == DefaultDigits {
		minDigits, _ = GetDigits(amount.CurrencyCode())
	}
	maxDigits := f.MaxDigits
	if maxDigits == DefaultDigits {
		maxDigits, _ = GetDigits(amount.CurrencyCode())
	}
	amount = amount.RoundTo(maxDigits, f.RoundingMode)
	numberParts := strings.Split(amount.Number(), ".")
	majorDigits := f.groupMajorDigits(numberParts[0])
	minorDigits := ""
	if len(numberParts) == 2 {
		minorDigits = numberParts[1]
	}
	if minDigits < maxDigits {
		// Strip any trailing zeroes.
		minorDigits = strings.TrimRight(minorDigits, "0")
		if len(minorDigits) < int(minDigits) {
			// Now there are too few digits, re-add trailing zeroes
			// until minDigits is reached.
			minorDigits += strings.Repeat("0", int(minDigits)-len(minorDigits))
		}
	}
	b := strings.Builder{}
	b.WriteString(majorDigits)
	if minorDigits != "" {
		b.WriteString(f.format.decimalSeparator)
		b.WriteString(minorDigits)
	}
	formatted := f.localizeDigits(b.String())

	return formatted
}

// formatCurrency formats the currency for display.
func (f *Formatter) formatCurrency(currencyCode string) string {
	var formatted string
	switch f.CurrencyDisplay {
	case DisplaySymbol:
		if symbol, ok := f.SymbolMap[currencyCode]; ok {
			formatted = symbol
		} else {
			formatted, _ = GetSymbol(currencyCode, f.locale)
		}
	case DisplayCode:
		formatted = currencyCode
	default:
		formatted = ""
	}

	return formatted
}

// groupMajorDigits groups major digits according to the currency format.
func (f *Formatter) groupMajorDigits(majorDigits string) string {
	if f.NoGrouping || f.format.primaryGroupingSize == 0 {
		return majorDigits
	}
	numDigits := len(majorDigits)
	minDigits := int(f.format.minGroupingDigits)
	primarySize := int(f.format.primaryGroupingSize)
	secondarySize := int(f.format.secondaryGroupingSize)
	if numDigits < (minDigits + primarySize) {
		return majorDigits
	}

	// Digits are grouped from right to left.
	// First the primary group, then the secondary groups.
	var groups []string
	groups = append(groups, majorDigits[numDigits-primarySize:numDigits])
	for i := numDigits - primarySize; i > 0; i = i - secondarySize {
		low := i - secondarySize
		if low < 0 {
			low = 0
		}
		groups = append(groups, majorDigits[low:i])
	}
	// Reverse the groups and reconstruct the digits.
	for i, j := 0, len(groups)-1; i < j; i, j = i+1, j-1 {
		groups[i], groups[j] = groups[j], groups[i]
	}
	majorDigits = strings.Join(groups, f.format.groupingSeparator)

	return majorDigits
}

// localizeDigits replaces digits with their localized equivalents.
func (f *Formatter) localizeDigits(number string) string {
	if f.format.numberingSystem == numLatn {
		return number
	}
	digits := localDigits[f.format.numberingSystem]
	replacements := make([]string, 0, 20)
	for i, v := range strings.Split(digits, "") {
		replacements = append(replacements, strconv.Itoa(i), v)
	}
	r := strings.NewReplacer(replacements...)
	number = r.Replace(number)

	return number
}
