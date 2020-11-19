// Code generated by go-swagger; DO NOT EDIT.

package tree

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"tarsadmin/openapi/models"
)

// SelectServerTreeOKCode is the HTTP code returned for type SelectServerTreeOK
const SelectServerTreeOKCode int = 200

/*SelectServerTreeOK OK

swagger:response selectServerTreeOK
*/
type SelectServerTreeOK struct {

	/*
	  In: Body
	*/
	Payload *SelectServerTreeOKBody `json:"body,omitempty"`
}

// NewSelectServerTreeOK creates SelectServerTreeOK with default headers values
func NewSelectServerTreeOK() *SelectServerTreeOK {

	return &SelectServerTreeOK{}
}

// WithPayload adds the payload to the select server tree o k response
func (o *SelectServerTreeOK) WithPayload(payload *SelectServerTreeOKBody) *SelectServerTreeOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the select server tree o k response
func (o *SelectServerTreeOK) SetPayload(payload *SelectServerTreeOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SelectServerTreeOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SelectServerTreeInternalServerErrorCode is the HTTP code returned for type SelectServerTreeInternalServerError
const SelectServerTreeInternalServerErrorCode int = 500

/*SelectServerTreeInternalServerError 内部错误

swagger:response selectServerTreeInternalServerError
*/
type SelectServerTreeInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSelectServerTreeInternalServerError creates SelectServerTreeInternalServerError with default headers values
func NewSelectServerTreeInternalServerError() *SelectServerTreeInternalServerError {

	return &SelectServerTreeInternalServerError{}
}

// WithPayload adds the payload to the select server tree internal server error response
func (o *SelectServerTreeInternalServerError) WithPayload(payload *models.Error) *SelectServerTreeInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the select server tree internal server error response
func (o *SelectServerTreeInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SelectServerTreeInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
