// Code generated by go-swagger; DO NOT EDIT.

package node

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"tarsadmin/openapi/models"
)

// SelectNodeOKCode is the HTTP code returned for type SelectNodeOK
const SelectNodeOKCode int = 200

/*SelectNodeOK OK

swagger:response selectNodeOK
*/
type SelectNodeOK struct {

	/*
	  In: Body
	*/
	Payload *models.SelectResult `json:"body,omitempty"`
}

// NewSelectNodeOK creates SelectNodeOK with default headers values
func NewSelectNodeOK() *SelectNodeOK {

	return &SelectNodeOK{}
}

// WithPayload adds the payload to the select node o k response
func (o *SelectNodeOK) WithPayload(payload *models.SelectResult) *SelectNodeOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the select node o k response
func (o *SelectNodeOK) SetPayload(payload *models.SelectResult) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SelectNodeOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SelectNodeInternalServerErrorCode is the HTTP code returned for type SelectNodeInternalServerError
const SelectNodeInternalServerErrorCode int = 500

/*SelectNodeInternalServerError 内部错误

swagger:response selectNodeInternalServerError
*/
type SelectNodeInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSelectNodeInternalServerError creates SelectNodeInternalServerError with default headers values
func NewSelectNodeInternalServerError() *SelectNodeInternalServerError {

	return &SelectNodeInternalServerError{}
}

// WithPayload adds the payload to the select node internal server error response
func (o *SelectNodeInternalServerError) WithPayload(payload *models.Error) *SelectNodeInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the select node internal server error response
func (o *SelectNodeInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SelectNodeInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
