// Code generated by go-swagger; DO NOT EDIT.

package default_operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"tarsadmin/openapi/models"
)

// SelectDefaultValueOKCode is the HTTP code returned for type SelectDefaultValueOK
const SelectDefaultValueOKCode int = 200

/*SelectDefaultValueOK OK?

swagger:response selectDefaultValueOK
*/
type SelectDefaultValueOK struct {

	/*
	  In: Body
	*/
	Payload *SelectDefaultValueOKBody `json:"body,omitempty"`
}

// NewSelectDefaultValueOK creates SelectDefaultValueOK with default headers values
func NewSelectDefaultValueOK() *SelectDefaultValueOK {

	return &SelectDefaultValueOK{}
}

// WithPayload adds the payload to the select default value o k response
func (o *SelectDefaultValueOK) WithPayload(payload *SelectDefaultValueOKBody) *SelectDefaultValueOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the select default value o k response
func (o *SelectDefaultValueOK) SetPayload(payload *SelectDefaultValueOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SelectDefaultValueOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SelectDefaultValueInternalServerErrorCode is the HTTP code returned for type SelectDefaultValueInternalServerError
const SelectDefaultValueInternalServerErrorCode int = 500

/*SelectDefaultValueInternalServerError 内部错误

swagger:response selectDefaultValueInternalServerError
*/
type SelectDefaultValueInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSelectDefaultValueInternalServerError creates SelectDefaultValueInternalServerError with default headers values
func NewSelectDefaultValueInternalServerError() *SelectDefaultValueInternalServerError {

	return &SelectDefaultValueInternalServerError{}
}

// WithPayload adds the payload to the select default value internal server error response
func (o *SelectDefaultValueInternalServerError) WithPayload(payload *models.Error) *SelectDefaultValueInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the select default value internal server error response
func (o *SelectDefaultValueInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SelectDefaultValueInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
