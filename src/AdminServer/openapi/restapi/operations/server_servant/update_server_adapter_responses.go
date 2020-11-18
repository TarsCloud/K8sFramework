// Code generated by go-swagger; DO NOT EDIT.

package server_servant

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"tafadmin/openapi/models"
)

// UpdateServerAdapterOKCode is the HTTP code returned for type UpdateServerAdapterOK
const UpdateServerAdapterOKCode int = 200

/*UpdateServerAdapterOK OK

swagger:response updateServerAdapterOK
*/
type UpdateServerAdapterOK struct {

	/*
	  In: Body
	*/
	Payload *UpdateServerAdapterOKBody `json:"body,omitempty"`
}

// NewUpdateServerAdapterOK creates UpdateServerAdapterOK with default headers values
func NewUpdateServerAdapterOK() *UpdateServerAdapterOK {

	return &UpdateServerAdapterOK{}
}

// WithPayload adds the payload to the update server adapter o k response
func (o *UpdateServerAdapterOK) WithPayload(payload *UpdateServerAdapterOKBody) *UpdateServerAdapterOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update server adapter o k response
func (o *UpdateServerAdapterOK) SetPayload(payload *UpdateServerAdapterOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateServerAdapterOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateServerAdapterInternalServerErrorCode is the HTTP code returned for type UpdateServerAdapterInternalServerError
const UpdateServerAdapterInternalServerErrorCode int = 500

/*UpdateServerAdapterInternalServerError 内部错误

swagger:response updateServerAdapterInternalServerError
*/
type UpdateServerAdapterInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewUpdateServerAdapterInternalServerError creates UpdateServerAdapterInternalServerError with default headers values
func NewUpdateServerAdapterInternalServerError() *UpdateServerAdapterInternalServerError {

	return &UpdateServerAdapterInternalServerError{}
}

// WithPayload adds the payload to the update server adapter internal server error response
func (o *UpdateServerAdapterInternalServerError) WithPayload(payload *models.Error) *UpdateServerAdapterInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update server adapter internal server error response
func (o *UpdateServerAdapterInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateServerAdapterInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}