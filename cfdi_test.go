package cfdi

import (
	"testing"
)

func TestValidation(t *testing.T) {
	document := DocumentHeaders{
		IssuerRFC:    "LSO1306189R5",
		AddresseeRFC: "GACJ940911ASA",
		TotalAmount:  "4999.99",
		InvoiceUUID:  "e7df3047-f8de-425d-b469-37abe5b4dabb",
	}
	response, err := document.Validate()
	if err != nil {
		t.Errorf("Failed to validate CFDI document: %w", err)
	}
	t.Logf("Got valid response object %v", response)
}
