// Code generated by go-swagger; DO NOT EDIT.

package notify

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"tarsadmin/openapi/models"
)

// SelectNotifyOKCode is the HTTP code returned for type SelectNotifyOK
const SelectNotifyOKCode int = 200

/*SelectNotifyOK OK

swagger:response selectNotifyOK
*/
type SelectNotifyOK struct {

	/*
	  In: Body
	*/
	Payload *models.SelectResult `json:"body,omitempty"`
}

// NewSelectNotifyOK creates SelectNotifyOK with default headers values
func NewSelectNotifyOK() *SelectNotifyOK {

	return &SelectNotifyOK{}
}

// WithPayload adds the payload to the select notify o k response
func (o *SelectNotifyOK) WithPayload(payload *models.SelectResult) *SelectNotifyOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the select notify o k response
func (o *SelectNotifyOK) SetPayload(payload *models.SelectResult) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SelectNotifyOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SelectNotifyInternalServerErrorCode is the HTTP code returned for type SelectNotifyInternalServerError
const SelectNotifyInternalServerErrorCode int = 500

/*SelectNotifyInternalServerError 内部错误

swagger:response selectNotifyInternalServerError
*/
type SelectNotifyInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSelectNotifyInternalServerError creates SelectNotifyInternalServerError with default headers values
func NewSelectNotifyInternalServerError() *SelectNotifyInternalServerError {

	return &SelectNotifyInternalServerError{}
}

// WithPayload adds the payload to the select notify internal server error response
func (o *SelectNotifyInternalServerError) WithPayload(payload *models.Error) *SelectNotifyInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the select notify internal server error response
func (o *SelectNotifyInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SelectNotifyInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
