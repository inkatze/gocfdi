package cfdi

import (
	"testing"
)

func TestValidation(t *testing.T) {
	issuer := "LSO1306189R5"
	addressee := "GACJ940911ASA"
	total := "4999.99"
	uuid := "e7df3047-f8de-425d-b469-37abe5b4dabb"

	invoice := InvoiceHeader{
		IssuerRFC:    issuer,
		AddresseeRFC: addressee,
		TotalAmount:  total,
		UUID:         uuid,
	}
	result, err := invoice.Validate()
	if err != nil {
		t.Errorf("Failed to validate CFDI document: %w", err)
	}
	t.Logf("Got valid response object %+v", *result)
}
