package currency

import (
	"fmt"
)

// EmptyCurrencyCodeError indicates that the currency code was empty.
type EmptyCurrencyCodeError struct{}

func (e EmptyCurrencyCodeError) Error() string {
	return "register currency error: empty currency code"
}

// CurrencyAlreadyExistsError indicates that the currency code already exists in the ISO list.
type CurrencyAlreadyExistsError struct {
	Code string
}

func (e CurrencyAlreadyExistsError) Error() string {
	return fmt.Sprintf("register currency error: code %q already exists in ISO list", e.Code)
}

// RegisterCurrencyOptions defines parameters for registering a new or custom currency.
type RegisterCurrencyOptions struct {
	// NumericCode is usually a three-digit code, for example "999".
	NumericCode string

	// Digits is the number of decimal fraction digits.
	Digits uint8

	// SymbolData is a list of possible symbols and the locales
	// in which each symbol is used.
	//
	// Example:
	//    []SymbolData{
	//       {Symbol: "₿", Locales: []string{"en"}},
	//       {Symbol: "BTC", Locales: []string{"uk"}},
	//    }
	SymbolData []SymbolData
}

// SymbolData describes one symbol and the set of locales
// for which that symbol applies.
type SymbolData struct {
	Symbol  string
	Locales []string
}

// RegisterCurrency adds a non-ISO currency to the global structures:
//   - currencies
//   - currencyCodes
//   - currencySymbols
//
// It returns an error if the code already exists in the ISO list, or if the code is empty.
func RegisterCurrency(code string, opts RegisterCurrencyOptions) error {
	if code == "" {
		return EmptyCurrencyCodeError{}
	}
	if _, isoExists := currencies[code]; isoExists {
		return CurrencyAlreadyExistsError{Code: code}
	}

	// Insert into the global `currencies` map.
	currencies[code] = currencyInfo{
		numericCode: opts.NumericCode,
		digits:      opts.Digits,
	}

	// Also append to currencyCodes, so that GetCurrencyCodes() is aware of it.
	currencyCodes = append(currencyCodes, code)

	// If SymbolData is provided, insert symbols into currencySymbols.
	if len(opts.SymbolData) > 0 {
		// Ensure there's a slice for 'code' in currencySymbols.
		if _, ok := currencySymbols[code]; !ok {
			currencySymbols[code] = []symbolInfo{}
		}
		// Add each entry from opts.SymbolData.
		for _, s := range opts.SymbolData {
			currencySymbols[code] = append(
				currencySymbols[code],
				symbolInfo{
					symbol:  s.Symbol,
					locales: s.Locales,
				},
			)
		}
	}

	return nil
}
