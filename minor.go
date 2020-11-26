// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package currency

import (
	"bytes"
	"encoding/json"
	"math/big"

	"github.com/cockroachdb/apd/v2"
)

// Minor stores a decimal number with its currency code.  All monetary amounts are in minor units.
type Minor struct {
	Amount
}

// NewMinor creates a new Amount from a numeric string and a currency code.
func NewMinor(n, currencyCode string) (Minor, error) {
	num, ok := big.NewInt(0).SetString(n, 10)
	if !ok {
		return Minor{}, InvalidNumberError{"NewMinor", n}
	}
	if currencyCode == "" || !IsValid(currencyCode) {
		return Minor{}, InvalidCurrencyCodeError{"NewAmount", currencyCode}
	}
	d, _ := GetDigits(currencyCode)
	return Minor{Amount: Amount{apd.NewWithBigInt(num, -int32(d)), currencyCode}}, nil
}

// ToMinor wraps an amount as a minor unit amount.
func ToMinor(a Amount) Minor { return Minor{Amount: a} }

// ToAmount unwraps the underlying abount, converting back to major currency units.
func (m Minor) ToAmount() Amount { return m.Amount }

// Number returns the number as a numeric string.
func (m Minor) Number() string {
	if m.number == nil {
		return "0"
	}
	return m.number.Coeff.String()
}

// MinorUnits returns a in minor units.
func (m Minor) MinorUnits() *big.Int {
	if m.number == nil {
		return nil
	}
	c := m.Round().number.Coeff
	return &c
}

// Convert converts to a divverent currency.
func (m Minor) Convert(currencyCode, rate string) (Minor, error) {
	a, err := m.Amount.Convert(currencyCode, rate)
	if err != nil {
		return Minor{}, err
	}
	return Minor{Amount: a}, nil
}

// Add adds m and b together and returns the result.
func (m Minor) Add(b Minor) (Minor, error) {
	a, err := m.Amount.Add(b.Amount)
	if err != nil {
		return Minor{}, err
	}
	return Minor{Amount: a}, nil
}

// Sub subtracts b from m and returns the result.
func (m Minor) Sub(b Minor) (Minor, error) {
	a, err := m.Amount.Sub(b.Amount)
	if err != nil {
		return Minor{}, err
	}
	return Minor{Amount: a}, nil
}

// Mul multiplies m by n and returns the result.
func (m Minor) Mul(n string) (Minor, error) {
	a, err := m.Amount.Mul(n)
	if err != nil {
		return Minor{}, err
	}
	return Minor{Amount: a}, nil
}

// Div multiplies m by n and returns the result.
func (m Minor) Div(n string) (Minor, error) {
	a, err := m.Amount.Div(n)
	if err != nil {
		return Minor{}, err
	}
	return Minor{Amount: a}, nil
}

// Round is a shortcut for RoundTo(currency.DefaultDigits, currency.RoundHalfUp).
func (a Minor) Round() Minor {
	return Minor{Amount: a.Amount.RoundTo(DefaultDigits, RoundHalfUp)}
}

// Round is a shortcut for RoundTo(currency.DefaultDigits, currency.RoundHalfUp).
func (a Minor) RoundTo(digits uint8, mode RoundingMode) Minor {
	return Minor{Amount: a.Amount.RoundTo(digits, mode)}
}

// Cmp compares a and b and returns:
//
//   -1 if a <  b
//    0 if a == b
//   +1 if a >  b
//
func (a Minor) Cmp(b Minor) (int, error) { return a.Amount.Cmp(b.Amount) }

// Equal returns whether a and b are equal.
func (a Minor) Equal(b Minor) bool { return a.Amount.Equal(b.Amount) }

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (m Minor) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteString(m.CurrencyCode())
	buf.WriteString(m.Number())
	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (m *Minor) UnmarshalBinary(data []byte) error {
	if len(data) < 3 {
		return InvalidCurrencyCodeError{"Amount.UnmarshalBinary", string(data)}
	}
	n, err := NewMinor(string(data[3:]), string(data[0:3]))
	if err != nil {
		return err
	}
	*m = n
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (m Minor) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		// Using 'amount' to ensure type mismatch is not silently converted
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	}{
		Amount:   m.Number(),
		Currency: m.CurrencyCode(),
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *Minor) UnmarshalJSON(data []byte) error {
	aux := struct {
		// Using 'amount' to ensure type mismatch is not silently converted
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	n, err := NewMinor(aux.Amount, aux.Currency)
	if err != nil {
		return err
	}
	*m = n
	return nil
}
