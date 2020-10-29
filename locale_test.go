// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package currency_test

import (
	"testing"

	"github.com/bojanz/currency"
)

func TestNewLocale(t *testing.T) {
	tests := []struct {
		id   string
		want currency.Locale
	}{
		{"", currency.Locale{}},
		{"de", currency.Locale{Language: "de"}},
		{"de-CH", currency.Locale{Language: "de", Territory: "CH"}},
		{"es-419", currency.Locale{Language: "es", Territory: "419"}},
		{"sr-Cyrl", currency.Locale{Language: "sr", Script: "Cyrl"}},
		{"sr-Latn-RS", currency.Locale{Language: "sr", Script: "Latn", Territory: "RS"}},
		{"yue-Hans", currency.Locale{Language: "yue", Script: "Hans"}},
		// ID with the wrong case, ordering, delimeter.
		{"SR_rs_LATN", currency.Locale{Language: "sr", Script: "Latn", Territory: "RS"}},
		// ID with a variant. Variants are unsupported and ignored.
		{"ca-ES-VALENCIA", currency.Locale{Language: "ca", Territory: "ES"}},
	}
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			got := currency.NewLocale(tt.id)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocale_String(t *testing.T) {
	tests := []struct {
		locale currency.Locale
		want   string
	}{
		{currency.Locale{}, ""},
		{currency.Locale{Language: "de"}, "de"},
		{currency.Locale{Language: "de", Territory: "CH"}, "de-CH"},
		{currency.Locale{Language: "sr", Script: "Cyrl"}, "sr-Cyrl"},
		{currency.Locale{Language: "sr", Script: "Latn", Territory: "RS"}, "sr-Latn-RS"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			id := tt.locale.String()
			if id != tt.want {
				t.Errorf("got %v, want %v", id, tt.want)
			}
		})
	}
}

func TestLocale_MarshalText(t *testing.T) {
	tests := []struct {
		locale currency.Locale
		want   string
	}{
		{currency.Locale{}, ""},
		{currency.Locale{Language: "de"}, "de"},
		{currency.Locale{Language: "de", Territory: "CH"}, "de-CH"},
		{currency.Locale{Language: "sr", Script: "Cyrl"}, "sr-Cyrl"},
		{currency.Locale{Language: "sr", Script: "Latn", Territory: "RS"}, "sr-Latn-RS"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			b, _ := tt.locale.MarshalText()
			got := string(b)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocale_UnmarshalText(t *testing.T) {
	tests := []struct {
		id   string
		want currency.Locale
	}{
		{"", currency.Locale{}},
		{"de", currency.Locale{Language: "de"}},
		{"de-CH", currency.Locale{Language: "de", Territory: "CH"}},
		{"sr-Cyrl", currency.Locale{Language: "sr", Script: "Cyrl"}},
		{"sr-Latn-RS", currency.Locale{Language: "sr", Script: "Latn", Territory: "RS"}},
		// ID with the wrong case, ordering, delimeter.
		{"SR_rs_LATN", currency.Locale{Language: "sr", Script: "Latn", Territory: "RS"}},
		// ID with a variant. Variants are unsupported and ignored.
		{"ca-ES-VALENCIA", currency.Locale{Language: "ca", Territory: "ES"}},
	}
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			l := currency.Locale{}
			l.UnmarshalText([]byte(tt.id))
			if l != tt.want {
				t.Errorf("got %v, want %v", l, tt.want)
			}
		})
	}
}

func TestLocale_IsEmpty(t *testing.T) {
	tests := []struct {
		locale currency.Locale
		want   bool
	}{
		{currency.Locale{}, true},
		{currency.Locale{Language: "de"}, false},
		{currency.Locale{Language: "de", Territory: "CH"}, false},
		{currency.Locale{Language: "sr", Script: "Cyrl"}, false},
		{currency.Locale{Language: "sr", Script: "Latn", Territory: "RS"}, false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			empty := tt.locale.IsEmpty()
			if empty != tt.want {
				t.Errorf("got %v, want %v", empty, tt.want)
			}
		})
	}
}

func TestLocale_GetParent(t *testing.T) {
	tests := []struct {
		id   string
		want currency.Locale
	}{
		{"sr-Cyrl-RS", currency.Locale{Language: "sr", Script: "Cyrl"}},
		{"sr-Cyrl", currency.Locale{Language: "sr"}},
		{"sr", currency.Locale{Language: "en"}},
		{"en", currency.Locale{}},
		{"", currency.Locale{}},
		// Locales with special parents.
		{"es-AR", currency.Locale{Language: "es", Territory: "419"}},
		{"sr-Latn", currency.Locale{Language: "en"}},
	}
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			locale := currency.NewLocale(tt.id)
			parent := locale.GetParent()
			if parent != tt.want {
				t.Errorf("got %v, want %v", parent, tt.want)
			}
		})
	}
}
