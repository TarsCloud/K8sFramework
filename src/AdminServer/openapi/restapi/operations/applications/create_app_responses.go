// Code generated by go-swagger; DO NOT EDIT.

package applications

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"tarsadmin/openapi/models"
)

// CreateAppOKCode is the HTTP code returned for type CreateAppOK
const CreateAppOKCode int = 200

/*CreateAppOK OK

swagger:response createAppOK
*/
type CreateAppOK struct {

	/*
	  In: Body
	*/
	Payload *CreateAppOKBody `json:"body,omitempty"`
}

// NewCreateAppOK creates CreateAppOK with default headers values
func NewCreateAppOK() *CreateAppOK {

	return &CreateAppOK{}
}

// WithPayload adds the payload to the create app o k response
func (o *CreateAppOK) WithPayload(payload *CreateAppOKBody) *CreateAppOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create app o k response
func (o *CreateAppOK) SetPayload(payload *CreateAppOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateAppOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateAppInternalServerErrorCode is the HTTP code returned for type CreateAppInternalServerError
const CreateAppInternalServerErrorCode int = 500

/*CreateAppInternalServerError 内部错误

swagger:response createAppInternalServerError
*/
type CreateAppInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateAppInternalServerError creates CreateAppInternalServerError with default headers values
func NewCreateAppInternalServerError() *CreateAppInternalServerError {

	return &CreateAppInternalServerError{}
}

// WithPayload adds the payload to the create app internal server error response
func (o *CreateAppInternalServerError) WithPayload(payload *models.Error) *CreateAppInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create app internal server error response
func (o *CreateAppInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateAppInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
