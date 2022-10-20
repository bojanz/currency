# currency [![Build](https://github.com/bojanz/currency/actions/workflows/build.yml/badge.svg)](https://github.com/bojanz/currency/actions/workflows/build.yml) [![Coverage Status](https://coveralls.io/repos/github/bojanz/currency/badge.svg?branch=master)](https://coveralls.io/github/bojanz/currency?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/bojanz/currency)](https://goreportcard.com/report/github.com/bojanz/currency) [![PkgGoDev](https://pkg.go.dev/badge/github.com/bojanz/currency)](https://pkg.go.dev/github.com/bojanz/currency)

Handles currency amounts, provides currency information and formatting.

Powered by CLDR v42, in just ~30kb of data.

Backstory: https://bojanz.github.io/price-currency-handling-go/

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

### Real decimal implementation under the hood.

Currency amounts can't be floats. Storing integer minor units (2.99 => 299)
becomes problematic once there are multiple currencies (difficult to sort in the
DB), or there is a need for sub-minor-unit precision (due to merchant or tax
requirements, etc). A real arbitrary-precision decimal type is required. Since
Go doesn't have one natively, a userspace implementation is used, provided by
the [cockroachdb/apd](https://github.com/cockroachdb/apd) package. The Amount struct provides an easy to use
abstraction on top of it, allowing the underlying implementation to be replaced
in the future without a backwards compatibility break.

### Smart filtering of CLDR data.

The "modern" subset of CLDR locales is used, reducing the list from ~560 to ~370 locales.

Once gathered, locales are filtered to remove all data not used by this package,
and then deduplicated by parent (e.g. don't keep `fr-CH` if `fr` has the
same data).

Currency symbols are grouped together to avoid repetition. For example:

    "ARS": {
        {"ARS", []string{"en", "fr-CA"}},
        {"$", []string{"es-AR"}},
        {"$AR", []string{"fr"}},
    }

Currency names are not included because they are rarely shown, but need
significant space. Instead, they can be fetched on the frontend via [Intl.DisplayNames](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Intl/DisplayNames).

### Easy to compare.

Amount structs can be compared via [google/go-cmp](https://github.com/google/go-cmp) thanks to the built-in Equal() method.

### Usable with a PostgreSQL composite type.

Thanks to the driver.Valuer and sql.Scanner interfaces, applications using the [pgx](https://github.com/jackc/pgx) driver can store amounts in a composite type.

Example schema:
```
CREATE TYPE price AS (
   number NUMERIC,
   currency_code TEXT
);

CREATE TABLE products (
   id CHAR(26) PRIMARY KEY,
   name TEXT NOT NULL,
   price price NOT NULL,
   created_at TIMESTAMPTZ NOT NULL,
   updated_at TIMESTAMPTZ
);
```
Note that the number and currency_code columns can have any name, only their ordering matters.

Example struct:
```go
type Product struct {
	ID          string
	Name        string
	Price       currency.Amount
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
```

Example scan:
```go
p := Product{}
row := tx.QueryRow(ctx, `SELECT id, name, price, created_at, updated_at FROM products WHERE id = $1`, id)
err := row.Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt, &p.UpdatedAt)
```

See our [database integration notes](https://github.com/bojanz/currency/wiki/Database-integration-notes) for other examples (MySQL/MariaDB, SQLite).
