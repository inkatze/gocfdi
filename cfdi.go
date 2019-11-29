package cfdi

import (
	"encoding/xml"
	"fmt"

	"github.com/tiaguinho/gosoap"
)

// ValidationURL for the SOAP service to validate CFDI documents
const ValidationURL = "https://consultaqr.facturaelectronica.sat.gob.mx/ConsultaCFDIService.svc?wsdl"

// InvoiceHeaders values used to uniquely identify CFDI documents
type InvoiceHeaders struct {
	IssuerRFC, AddresseeRFC, TotalAmount, InvoiceUUID string
	Client                                            *gosoap.Client
}

// ConsultaResponse has the request's unmarshalled values
type ConsultaResponse struct {
	XMLName            xml.Name `xml:"ConsultaResponse"`
	ResponseStatus     string   `xml:"ConsultaResult>CodigoEstatus"`
	CFDIStatus         string   `xml:"ConsultaResult>Estado"`
	CancellationStatus string   `xml:"ConsultaResult>EstatusCancelacion"`
	IsCancellable      string   `xml:"ConsultaResult>EsCancelable"`
	RawResponse        string
}

var r ConsultaResponse

// SoapClient creates a new SOAP client for the given url.
func (i *InvoiceHeaders) SoapClient(url string) error {
	client, err := gosoap.SoapClient(url)
	if err != nil {
		return fmt.Errorf("Error while creating a SOAP client for %s: %w", url, err)
	}
	i.Client = client
	return nil
}

// ValidateDocument checks if the document with the given parameters exists in SAT's system and is valid.
// BUG(inkatze): The namespace is not being parsed correctly in gosoap.
func ValidateDocument(i InvoiceHeaders) (ConsultaResponse, error) {
	if i.Client == nil {
		if err := i.SoapClient(ValidationURL); err != nil {
			return r, err
		}
	}
	parsedValues := fmt.Sprintf(
		"re=%s&rr=%s&tt=%s&id=%s",
		i.IssuerRFC, i.AddresseeRFC, i.TotalAmount, i.InvoiceUUID,
	)
	params := gosoap.Params{"expresionImpresa": parsedValues}

	// Work around a bug that uses the wrong namespace
	res, err := i.Client.Call("Consulta", params)
	callError := "Error while calling SOAP action: %w"
	if err != nil {
		return r, fmt.Errorf(callError, err)
	}
	i.Client.Definitions.Types[0].XsdSchema[0].TargetNamespace = "http://tempuri.org/"

	// Actual call
	res, err = i.Client.Call("Consulta", params)
	if err != nil {
		return r, fmt.Errorf(callError, err)
	}
	r.RawResponse = string(res.Body)
	err = res.Unmarshal(&r)
	if err != nil {
		return r, fmt.Errorf("Error while unmarshalling response: %w", err)
	}

	return r, nil
}
