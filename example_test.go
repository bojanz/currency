// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package currency_test

import (
	"fmt"

	"github.com/bojanz/currency"
)

func ExampleNewAmount() {
	amount, _ := currency.NewAmount("24.49", "USD")
	fmt.Println(amount)
	fmt.Println(amount.Number())
	fmt.Println(amount.CurrencyCode())
	// Output: 24.49 USD
	// 24.49
	// USD
}

func ExampleAmount_ToMinorUnits() {
	firstAmount, _ := currency.NewAmount("20.99", "USD")
	secondAmount, _ := currency.NewAmount("50", "USD")
	fmt.Println(firstAmount.ToMinorUnits())
	fmt.Println(secondAmount.ToMinorUnits())
	// Output: 2099
	// 5000
}

func ExampleAmount_Convert() {
	amount, _ := currency.NewAmount("20.99", "USD")
	amount, _ = amount.Convert("EUR", "0.91")
	fmt.Println(amount)
	fmt.Println(amount.Round())
	// Output: 19.1009 EUR
	// 19.10 EUR
}

func ExampleAmount_Add() {
	firstAmount, _ := currency.NewAmount("20.99", "USD")
	secondAmount, _ := currency.NewAmount("3.50", "USD")
	totalAmount, _ := firstAmount.Add(secondAmount)
	fmt.Println(totalAmount)
	// Output: 24.49 USD
}

func ExampleAmount_Sub() {
	baseAmount, _ := currency.NewAmount("20.99", "USD")
	discountAmount, _ := currency.NewAmount("5.00", "USD")
	amount, _ := baseAmount.Sub(discountAmount)
	fmt.Println(amount)
	// Output: 15.99 USD
}

func ExampleAmount_Mul() {
	amount, _ := currency.NewAmount("20.99", "USD")
	taxAmount, _ := amount.Mul("0.20")
	fmt.Println(taxAmount)
	fmt.Println(taxAmount.Round())
	// Output: 4.1980 USD
	// 4.20 USD
}

func ExampleAmount_Div() {
	totalAmount, _ := currency.NewAmount("99.99", "USD")
	amount, _ := totalAmount.Div("3")
	fmt.Println(amount)
	// Output: 33.33 USD
}

func ExampleAmount_Round() {
	firstAmount, _ := currency.NewAmount("12.345", "USD")
	secondAmount, _ := currency.NewAmount("12.345", "JPY")
	fmt.Println(firstAmount.Round())
	fmt.Println(secondAmount.Round())
	// Output: 12.35 USD
	// 12 JPY
}

func ExampleAmount_RoundTo() {
	amount, _ := currency.NewAmount("12.345", "USD")
	for _, digits := range []byte{4, 3, 2, 1, 0} {
		fmt.Println(amount.RoundTo(digits))
	}
	// Output: 12.3450 USD
	// 12.345 USD
	// 12.35 USD
	// 12.3 USD
	// 12 USD
}

func ExampleGetNumericCode() {
	numericCode, ok := currency.GetNumericCode("USD")
	fmt.Println(numericCode, ok)

	// Non-existent currency code.
	numericCode, ok = currency.GetNumericCode("XXX")
	fmt.Println(numericCode, ok)
	// Output: 840 true
	// 000 false
}

func ExampleGetDigits() {
	digits, ok := currency.GetDigits("USD")
	fmt.Println(digits, ok)

	// Non-existent currency code.
	digits, ok = currency.GetDigits("XXX")
	fmt.Println(digits, ok)
	// Output: 2 true
	// 0 false
}
