package cfdi

import (
	"testing"
)

func TestValidation(t *testing.T) {
	invoice := InvoiceHeaders{
		IssuerRFC:    "LSO1306189R5",
		AddresseeRFC: "GACJ940911ASA",
		TotalAmount:  "4999.99",
		InvoiceUUID:  "e7df3047-f8de-425d-b469-37abe5b4dabb",
	}
	response, err := ValidateDocument(invoice)
	if err != nil {
		t.Errorf("Failed to validate CFDI document: %w", err)
	}
	t.Logf("Got valid response object %v", response)
}
