package currency

import (
	"testing"
)

func TestRegisterCurrencyBTC(t *testing.T) {
	err := RegisterCurrency("BTC", RegisterCurrencyOptions{
		NumericCode: "1000",
		Digits:      8,
		SymbolData: []SymbolData{
			{
				Symbol:  "₿",
				Locales: []string{"en"},
			},
			{
				Symbol:  "BTC",
				Locales: []string{"uk"},
			},
		},
	})
	if err != nil {
		t.Errorf("RegisterCurrency returned an error for BTC: %v", err)
	}

	if !IsValid("BTC") {
		t.Error("Expected 'BTC' to be valid after registration, but IsValid returned false.")
	}

	d, ok := GetDigits("BTC")
	if !ok {
		t.Error("Expected 'BTC' to be found, but GetDigits says not ok.")
	} else if d != 8 {
		t.Errorf("Expected 'BTC' digits=8, got %d", d)
	}

	symEN, _ := GetSymbol("BTC", NewLocale("en")) // "₿"
	if symEN != "₿" {
		t.Errorf("Expected '₿' for locale 'en', got '%s'", symEN)
	}

	symRU, _ := GetSymbol("BTC", NewLocale("uk")) // "BTC"
	if symRU != "BTC" {
		t.Errorf("Expected 'BTC' for locale 'uk', got '%s'", symRU)
	}

	err = RegisterCurrency("BTC", RegisterCurrencyOptions{})
	if err == nil {
		t.Error("Expected an error when re-registering code 'BTC', but got nil.")
	}
}
