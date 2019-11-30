# gocfdi [![Actions Status](https://github.com/inkatze/gocfdi/workflows/test/badge.svg)](https://github.com/{owner}/{repo}/actions) [![GoDoc](https://godoc.org/github.com/inkatze/gocfdi?status.png)](https://godoc.org/github.com/inkatze/gocfdi) [![Go Report Card](https://goreportcard.com/badge/github.com/inkatze/gocfdi)](https://goreportcard.com/report/github.com/inkatze/gocfdi) [![codecov](https://codecov.io/gh/inkatze/gocfdi/branch/master/graph/badge.svg)](https://codecov.io/gh/inkatze/gocfdi)

Package to help validating CFDI invoices with the mexican Tax Administration System.

# Install

```
go get github.com/inkatze/gocfdi
```

# Example

```go
package main

import (
	"fmt"

	"github.com/inkatze/gocfdi"
)

func main() {
	document := cfdi.DocumentHeaders{
		IssuerRFC:    "LSO1306189R5",
		AddresseeRFC: "GACJ940911ASA",
		TotalAmount:  "4999.99",
		InvoiceUUID:  "e7df3047-f8de-425d-b469-37abe5b4dabb",
	}
	result, _ := document.Validate()
	fmt.Printf(result.RawResponse)
}
```
