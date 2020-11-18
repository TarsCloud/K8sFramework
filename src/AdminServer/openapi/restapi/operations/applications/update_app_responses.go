// Code generated by go-swagger; DO NOT EDIT.

package applications

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"tafadmin/openapi/models"
)

// UpdateAppOKCode is the HTTP code returned for type UpdateAppOK
const UpdateAppOKCode int = 200

/*UpdateAppOK OK

swagger:response updateAppOK
*/
type UpdateAppOK struct {

	/*
	  In: Body
	*/
	Payload *UpdateAppOKBody `json:"body,omitempty"`
}

// NewUpdateAppOK creates UpdateAppOK with default headers values
func NewUpdateAppOK() *UpdateAppOK {

	return &UpdateAppOK{}
}

// WithPayload adds the payload to the update app o k response
func (o *UpdateAppOK) WithPayload(payload *UpdateAppOKBody) *UpdateAppOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update app o k response
func (o *UpdateAppOK) SetPayload(payload *UpdateAppOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateAppOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateAppInternalServerErrorCode is the HTTP code returned for type UpdateAppInternalServerError
const UpdateAppInternalServerErrorCode int = 500

/*UpdateAppInternalServerError 内部错误

swagger:response updateAppInternalServerError
*/
type UpdateAppInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewUpdateAppInternalServerError creates UpdateAppInternalServerError with default headers values
func NewUpdateAppInternalServerError() *UpdateAppInternalServerError {

	return &UpdateAppInternalServerError{}
}

// WithPayload adds the payload to the update app internal server error response
func (o *UpdateAppInternalServerError) WithPayload(payload *models.Error) *UpdateAppInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update app internal server error response
func (o *UpdateAppInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateAppInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}