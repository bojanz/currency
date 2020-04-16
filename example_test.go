// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package currency_test

import (
	"fmt"

	"github.com/bojanz/currency"
)

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
