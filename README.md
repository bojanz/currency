# currency [![Build Status](https://travis-ci.org/bojanz/currency.png?branch=master)](https://travis-ci.org/bojanz/currency) [![GoDoc](https://godoc.org/github.com/bojanz/currency?status.svg)](https://godoc.org/github.com/bojanz/currency)

Handles currency amounts, provides currency information and formatting.

Powered by CLDR v37, in just ~30kb of data.

## Features

1. All currency codes, their numeric codes and fraction digits.
2. Currency symbols and formats for all locales.
3. Amount struct, with value semantics (Fowler's Money pattern)
4. Formatter, for formatting amounts and parsing formatted amounts.

```go
    amount, _ := currency.NewAmount("275.98", "EUR")
    total, _ := amount.Mul("4")

    locale := currency.NewLocale("fr")
    formatter := currency.NewFormatter(locale)
    fmt.Println(formatter.Format(total)) // 1 103,92 €

    // Convert the amount to Iranian rial and show it in Farsi.
    total, _ = total.Convert("IRR", "45.538")
    total = total.Round()
    locale = currency.NewLocale("fa")
    formatter = currency.NewFormatter(locale)
    fmt.Println(formatter.Format(total)) // ‎ریال ۵۰٬۲۷۰
```

## Design goals

1. Real decimal implementation under the hood.

Currency amounts can't be floats. Storing integer minor units (2.99 => 299)
becomes problematic once there are multiple currencies (difficult to sort in the
DB), or there is a need for sub-minor-unit precision (due to merchant or tax
requirements, etc). A real arbitrary-precision decimal type is required. Since
Go doesn't have one natively, a userspace implementation is used, provided by
the cockroachdb/apd package. The Amount struct provides an easy to use
abstraction on top of it, allowing the underlying implementation to be replaced
in the future without a backwards compatibility break.

2. English as a priority.

Where possible, the ``en`` and ``en-US`` locales are given preferred treatment.
For example, looking up currency symbols is the quickest for these locales,
since they are the most commonly used.

3. Automatically generated CLDR data.

Updating to the latest CLDR release is always one ``go generate`` away.

4. Smart filtering of CLDR data.

CLDR contains 542 locales, not all of which are likely to be used. This list is
reduced to 381 using a list of ignored locales, containing constructed languages
(Esperanto, Volapuk), historical languages (Church Slavic, Prussian), secondary
non-official languages (Cornish, Gujarati, Cherokee). The ignore list is
maintained by community feedback, developers wishing to use an ignored locale
can open a pull request.

Once gathered, locales are filtered to remove all data not used by this package,
and then deduplicated by parent (e.g. don't keep ``fr-CH`` if ``fr`` has the
same data).

Currency symbols are grouped together to avoid repetition. For example:

    "ARS": {
        {"ARS", []string{"en", "fr-CA"}},
        {"$", []string{"es-AR"}},
        {"$AR", []string{"fr"}},
    }

Currency names are not included because they are rarely shown, but need
significant space.
