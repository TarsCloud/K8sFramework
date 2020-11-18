// Code generated by go-swagger; DO NOT EDIT.

package applications

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"tafadmin/openapi/models"
)

// SelectAppOKCode is the HTTP code returned for type SelectAppOK
const SelectAppOKCode int = 200

/*SelectAppOK OK

swagger:response selectAppOK
*/
type SelectAppOK struct {

	/*
	  In: Body
	*/
	Payload *models.SelectResult `json:"body,omitempty"`
}

// NewSelectAppOK creates SelectAppOK with default headers values
func NewSelectAppOK() *SelectAppOK {

	return &SelectAppOK{}
}

// WithPayload adds the payload to the select app o k response
func (o *SelectAppOK) WithPayload(payload *models.SelectResult) *SelectAppOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the select app o k response
func (o *SelectAppOK) SetPayload(payload *models.SelectResult) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SelectAppOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SelectAppInternalServerErrorCode is the HTTP code returned for type SelectAppInternalServerError
const SelectAppInternalServerErrorCode int = 500

/*SelectAppInternalServerError 内部错误

swagger:response selectAppInternalServerError
*/
type SelectAppInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSelectAppInternalServerError creates SelectAppInternalServerError with default headers values
func NewSelectAppInternalServerError() *SelectAppInternalServerError {

	return &SelectAppInternalServerError{}
}

// WithPayload adds the payload to the select app internal server error response
func (o *SelectAppInternalServerError) WithPayload(payload *models.Error) *SelectAppInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the select app internal server error response
func (o *SelectAppInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SelectAppInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}