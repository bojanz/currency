// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package currency

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestMinor(t *testing.T) {
	tests := []struct {
		amt    string
		cur    string
		err    error
		number string
		strVal string
	}{
		{
			amt:    "10.99",
			cur:    "EUR",
			err:    InvalidNumberError{"NewMinor", "10.99"},
			number: "0",
			strVal: "0 ",
		},
		{
			amt:    "NaN",
			cur:    "USD",
			err:    InvalidNumberError{"NewMinor", "NaN"},
			number: "0",
			strVal: "0 ",
		},
		{
			amt:    "10050",
			cur:    "USD",
			err:    nil,
			number: "10050",
			strVal: "100.50 USD",
		},
	}
	for i, tt := range tests {
		m, err := NewMinor(tt.amt, tt.cur)
		if err != nil {
			if !errors.Is(err, tt.err) {
				t.Errorf("%d: got %T, want %T", i, err, tt.err)
			}
			if err.Error() != tt.err.Error() {
				t.Errorf("%d: got %v, want %v", i, err.Error(), tt.err.Error())
			}
		}
		if m.Number() != tt.number {
			t.Errorf("%d number: got %v, want %v", i, m.Number(), tt.number)
		}
		if m.String() != tt.strVal {
			t.Errorf("%d string: got %v, want %v", i, m.String(), tt.strVal)
		}
	}
}

func TestMinorMarshal(t *testing.T) {
	tests := []struct {
		amt    string
		cur    string
		strVal string
	}{
		{
			amt:    "0",
			cur:    "USD",
			strVal: "0.00 USD",
		},
		{
			amt:    "100000",
			cur:    "USD",
			strVal: "1000.00 USD",
		},
		{
			amt:    "3000001",
			cur:    "USD",
			strVal: "30000.01 USD",
		},
		{
			amt:    "1",
			cur:    "USD",
			strVal: "0.01 USD",
		},
	}
	for i, tt := range tests {
		m, err := NewMinor(tt.amt, tt.cur)
		if err != nil {
			t.Fatalf("%d: %v", i, err)
		}
		jBuf, err := json.Marshal(m)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}
		var amtJSON *Minor
		if err := json.Unmarshal(jBuf, &amtJSON); err != nil {
			t.Errorf("%d: %v", i, err)
		}
		if !amtJSON.Equal(m) {
			t.Errorf("%d: %v != %v", i, amtJSON, m)
		}

		bBuf, err := m.MarshalBinary()
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}
		amtBinary := &Minor{}
		if err := amtBinary.UnmarshalBinary(bBuf); err != nil {
			t.Errorf("%d: %v", i, err)
		}
		if !amtBinary.Equal(m) {
			t.Errorf("%d: %v != %v", i, amtBinary, m)
		}
	}
}
