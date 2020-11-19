// Code generated by go-swagger; DO NOT EDIT.

package template

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"tarsadmin/openapi/models"
)

// CreateTemplateOKCode is the HTTP code returned for type CreateTemplateOK
const CreateTemplateOKCode int = 200

/*CreateTemplateOK OK

swagger:response createTemplateOK
*/
type CreateTemplateOK struct {

	/*
	  In: Body
	*/
	Payload *CreateTemplateOKBody `json:"body,omitempty"`
}

// NewCreateTemplateOK creates CreateTemplateOK with default headers values
func NewCreateTemplateOK() *CreateTemplateOK {

	return &CreateTemplateOK{}
}

// WithPayload adds the payload to the create template o k response
func (o *CreateTemplateOK) WithPayload(payload *CreateTemplateOKBody) *CreateTemplateOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create template o k response
func (o *CreateTemplateOK) SetPayload(payload *CreateTemplateOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTemplateOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateTemplateInternalServerErrorCode is the HTTP code returned for type CreateTemplateInternalServerError
const CreateTemplateInternalServerErrorCode int = 500

/*CreateTemplateInternalServerError 内部错误

swagger:response createTemplateInternalServerError
*/
type CreateTemplateInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateTemplateInternalServerError creates CreateTemplateInternalServerError with default headers values
func NewCreateTemplateInternalServerError() *CreateTemplateInternalServerError {

	return &CreateTemplateInternalServerError{}
}

// WithPayload adds the payload to the create template internal server error response
func (o *CreateTemplateInternalServerError) WithPayload(payload *models.Error) *CreateTemplateInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create template internal server error response
func (o *CreateTemplateInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTemplateInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
