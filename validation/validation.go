// Package validation contains a collection of tools to read and manage CFDI documents.
package validation

import (
	"encoding/xml"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tiaguinho/gosoap"
)

const validationURL = "https://consultaqr.facturaelectronica.sat.gob.mx/ConsultaCFDIService.svc?wsdl"
const successMessage = "S - Comprobante obtenido satisfactoriamente."
const invalidInvoice = "N - 601: La expresión impresa proporcionada no es válida."
const invalidNotFound = "N - 602: Comprobante no encontrado"
const validCFDIStatus = "Vigente"
const notCancellable = "No cancelable"

// Result contains the processed results from the validation response
type Result struct {
	RawResponse                             string
	IsDocumentFound, IsValid, IsCancellable bool
	Timestamp                               int64
}

type serviceResponse struct {
	XMLName            xml.Name `xml:"ConsultaResponse"`
	ResponseStatus     string   `xml:"ConsultaResult>CodigoEstatus"`
	CFDIStatus         string   `xml:"ConsultaResult>Estado"`
	CancellationStatus string   `xml:"ConsultaResult>EstatusCancelacion"`
	Cancellable        string   `xml:"ConsultaResult>EsCancelable"`
}

// InvoiceHeader contains data used to validate a CFDI document.
type InvoiceHeader struct {
	IssuerRFC, AddresseeRFC, TotalAmount, UUID string
}

// Validate checks if the document with the given parameters exists in SAT's system and is valid.
// BUG(inkatze): The namespace is not being parsed correctly in gosoap, so we have to make two calls to fix the element's namespace
func (i *InvoiceHeader) Validate() (*Result, error) {
	log.Debugf("Preparing validation service call")
	soapResponse, err := i.callService()
	if err != nil {
		log.Errorf("There was an error during the validation service call")
		return nil, err
	}

	log.Debugf("Parsing values from service response")
	var unmarshaller *serviceResponse
	err = soapResponse.Unmarshal(&unmarshaller)
	if err := soapResponse.Unmarshal(&unmarshaller); err != nil {
		return nil, fmt.Errorf("Error while unmarshalling response: %w", err)
	}

	r := &Result{
		RawResponse:     string(soapResponse.Body),
		IsDocumentFound: validateResponseStatus(unmarshaller.ResponseStatus),
		IsValid:         validateCFDIStatus(unmarshaller.CFDIStatus),
		IsCancellable:   validateCancelationStatus(unmarshaller.Cancellable),
		Timestamp:       time.Now().Unix(),
	}

	return r, nil
}

func (i *InvoiceHeader) validationRequestBody() string {
	return fmt.Sprintf(
		"re=%s&rr=%s&tt=%s&id=%s",
		i.IssuerRFC, i.AddresseeRFC, i.TotalAmount, i.UUID,
	)
}

func (i *InvoiceHeader) callService() (*gosoap.Response, error) {
	client, err := soapClient(validationURL)
	if err != nil {
		return nil, err
	}
	parsedValues := i.validationRequestBody()
	params := gosoap.Params{"expresionImpresa": parsedValues}

	// Work around a bug that uses the wrong namespace
	log.Debugf("Fetching results from first query")
	res, err := client.Call("Consulta", params)
	callError := "Error while calling SOAP action: %w"
	if err != nil {
		return res, fmt.Errorf(callError, err)
	}
	log.Debugf("Fixing target namespace")
	client.Definitions.Types[0].XsdSchema[0].TargetNamespace = "http://tempuri.org/"

	// Actual call
	log.Debugf("Running query with fixed target namespace")
	res, err = client.Call("Consulta", params)
	if err != nil {
		return res, fmt.Errorf(callError, err)
	}
	return res, nil
}

func soapClient(url string) (*gosoap.Client, error) {
	log.Debugf("Creating SOAP client")
	client, err := gosoap.SoapClient(url)
	if err != nil {
		log.Errorf("Error while creating SOAP client for url: %s", url)
		return nil, fmt.Errorf("Error while creating a SOAP client for %s: %w", url, err)
	}
	log.Debugf("SOAP client successfully created")
	return client, nil
}

func validateResponseStatus(responseStatus string) bool {
	switch responseStatus {
	case successMessage:
		log.Debugf("Response seems to be OK")
		return true
	case invalidInvoice:
		log.Debugf("The given parameters for the CFDI document are invalid: %s", responseStatus)
	case invalidNotFound:
		log.Debugf("Couldn't find a CFDI document with the given parameters")
	default:
		log.Errorf("Unrecognized status response %s", responseStatus)
	}
	return false
}

func validateCFDIStatus(cfdiStatus string) bool {
	log.Debugf("Validating CFDI status")
	return cfdiStatus == validCFDIStatus
}

func validateCancelationStatus(cancellable string) bool {
	log.Debugf("Validating cancellation status")
	return cancellable != notCancellable
}
