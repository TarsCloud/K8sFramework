// Code generated by go-swagger; DO NOT EDIT.

package node

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"tafadmin/openapi/models"
)

// DoSetPublicNodeOKCode is the HTTP code returned for type DoSetPublicNodeOK
const DoSetPublicNodeOKCode int = 200

/*DoSetPublicNodeOK OK

swagger:response doSetPublicNodeOK
*/
type DoSetPublicNodeOK struct {

	/*
	  In: Body
	*/
	Payload *DoSetPublicNodeOKBody `json:"body,omitempty"`
}

// NewDoSetPublicNodeOK creates DoSetPublicNodeOK with default headers values
func NewDoSetPublicNodeOK() *DoSetPublicNodeOK {

	return &DoSetPublicNodeOK{}
}

// WithPayload adds the payload to the do set public node o k response
func (o *DoSetPublicNodeOK) WithPayload(payload *DoSetPublicNodeOKBody) *DoSetPublicNodeOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the do set public node o k response
func (o *DoSetPublicNodeOK) SetPayload(payload *DoSetPublicNodeOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DoSetPublicNodeOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DoSetPublicNodeInternalServerErrorCode is the HTTP code returned for type DoSetPublicNodeInternalServerError
const DoSetPublicNodeInternalServerErrorCode int = 500

/*DoSetPublicNodeInternalServerError 内部错误

swagger:response doSetPublicNodeInternalServerError
*/
type DoSetPublicNodeInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDoSetPublicNodeInternalServerError creates DoSetPublicNodeInternalServerError with default headers values
func NewDoSetPublicNodeInternalServerError() *DoSetPublicNodeInternalServerError {

	return &DoSetPublicNodeInternalServerError{}
}

// WithPayload adds the payload to the do set public node internal server error response
func (o *DoSetPublicNodeInternalServerError) WithPayload(payload *models.Error) *DoSetPublicNodeInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the do set public node internal server error response
func (o *DoSetPublicNodeInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DoSetPublicNodeInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
