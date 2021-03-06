# gocfdi ![](https://github.com/inkatze/gocfdi/workflows/test/badge.svg?branch=master) [![GoDoc](https://godoc.org/github.com/inkatze/gocfdi?status.png)](https://godoc.org/github.com/inkatze/gocfdi) [![Go Report Card](https://goreportcard.com/badge/github.com/inkatze/gocfdi)](https://goreportcard.com/report/github.com/inkatze/gocfdi)

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

	"github.com/inkatze/gocfdi/validation"
)

func main() {
	document := validation.InvoiceHeader{
		IssuerRFC:    "LSO1306189R5",
		AddresseeRFC: "GACJ940911ASA",
		TotalAmount:  "4999.99",
		UUID:         "e7df3047-f8de-425d-b469-37abe5b4dabb",
	}
	result, err := document.Validate()
	if err != nil {
		fmt.Printf("Error while trying to validate the document")
		return
	}
	fmt.Printf("Validation results: %+v\n", *result)
}
```
