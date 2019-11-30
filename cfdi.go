package cfdi

import (
	"encoding/xml"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/tiaguinho/gosoap"
)

// ValidationURL for the SOAP service to validate CFDI documents
const ValidationURL = "https://consultaqr.facturaelectronica.sat.gob.mx/ConsultaCFDIService.svc?wsdl"

const successMessage = "S - Comprobante obtenido satisfactoriamente."
const invalidInvoice = "N - 601: La expresión impresa proporcionada no es válida."
const invalidNotFound = "N - 602: Comprobante no encontrado"

// DocumentHeaders values used to uniquely identify CFDI documents
type DocumentHeaders struct {
	IssuerRFC, AddresseeRFC, TotalAmount, InvoiceUUID string
	Client                                            *gosoap.Client
}

// ValidationResult contains the processed results from the validation response
type ValidationResult struct {
	RawResponse string
}

type consultaResponse struct {
	XMLName            xml.Name `xml:"ConsultaResponse"`
	ResponseStatus     string   `xml:"ConsultaResult>CodigoEstatus"`
	CFDIStatus         string   `xml:"ConsultaResult>Estado"`
	CancellationStatus string   `xml:"ConsultaResult>EstatusCancelacion"`
	IsCancellable      string   `xml:"ConsultaResult>EsCancelable"`
}

var (
	c consultaResponse
	r ValidationResult
)

// SoapClient creates a new SOAP client for the given url.
func (d *DocumentHeaders) SoapClient(url string) error {
	client, err := gosoap.SoapClient(url)
	if err != nil {
		log.Errorf("Error while creating SOAP client for url: %s", url)
		return fmt.Errorf("Error while creating a SOAP client for %s: %w", url, err)
	}
	log.Debugf("SOAP client successfully created")
	d.Client = client
	return nil
}

func (d *DocumentHeaders) call() (*gosoap.Response, error) {
	if d.Client == nil {
		if err := d.SoapClient(ValidationURL); err != nil {
			return nil, err
		}
	}
	parsedValues := fmt.Sprintf(
		"re=%s&rr=%s&tt=%s&id=%s",
		d.IssuerRFC, d.AddresseeRFC, d.TotalAmount, d.InvoiceUUID,
	)
	params := gosoap.Params{"expresionImpresa": parsedValues}

	// Work around a bug that uses the wrong namespace
	res, err := d.Client.Call("Consulta", params)
	callError := "Error while calling SOAP action: %w"
	if err != nil {
		return res, fmt.Errorf(callError, err)
	}
	d.Client.Definitions.Types[0].XsdSchema[0].TargetNamespace = "http://tempuri.org/"

	// Actual call
	res, err = d.Client.Call("Consulta", params)
	if err != nil {
		return res, fmt.Errorf(callError, err)
	}
	return res, nil
}

// Validate checks if the document with the given parameters exists in SAT's system and is valid.
// BUG(inkatze): The namespace is not being parsed correctly in gosoap.
func (d *DocumentHeaders) Validate() (ValidationResult, error) {
	res, err := d.call()
	if err != nil {
		// Error
	}
	err = res.Unmarshal(&c)
	if err != nil {
		return r, fmt.Errorf("Error while unmarshalling response: %w", err)
	}

	switch c.ResponseStatus {
	case invalidInvoice:
		return r, fmt.Errorf("The given parameters for the CFDI document are invalid: %s", c.ResponseStatus)
	case invalidNotFound:
		return r, fmt.Errorf("Couldn't find a CFDI document with the given parameters")
	default:
		if c.ResponseStatus != successMessage {
			return r, fmt.Errorf("Unrecognized status response %s", c.ResponseStatus)

		}
	}
	r.RawResponse = string(res.Body)
	return r, nil
}
