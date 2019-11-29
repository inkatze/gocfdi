# gocfdi [![GoDoc](https://godoc.org/github.com/inkatze/gocfdi?status.png)](https://godoc.org/github.com/inkatze/gocfdi)

Package to help validating CFDI invoices with the mexican Tax Administration System.

# Install

```
go get github.com/inkatze/gocfdi
```

# Example

```go
package main

import cfdi "github.com/inkatze/gocfdi"

func main() {
	invoice := cfdi.InvoiceHeaders{
		IssuerRFC:    "LSO1306189R5",
		AddresseeRFC: "GACJ940911ASA",
		TotalAmount:  "4999.99",
		InvoiceUUID:  "e7df3047-f8de-425d-b469-37abe5b4dabb",
	}
	cfdi.ValidateDocument(invoice)
}
