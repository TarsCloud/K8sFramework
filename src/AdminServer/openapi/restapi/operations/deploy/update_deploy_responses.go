// Code generated by go-swagger; DO NOT EDIT.

package deploy

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"tafadmin/openapi/models"
)

// UpdateDeployOKCode is the HTTP code returned for type UpdateDeployOK
const UpdateDeployOKCode int = 200

/*UpdateDeployOK OK

swagger:response updateDeployOK
*/
type UpdateDeployOK struct {

	/*
	  In: Body
	*/
	Payload *UpdateDeployOKBody `json:"body,omitempty"`
}

// NewUpdateDeployOK creates UpdateDeployOK with default headers values
func NewUpdateDeployOK() *UpdateDeployOK {

	return &UpdateDeployOK{}
}

// WithPayload adds the payload to the update deploy o k response
func (o *UpdateDeployOK) WithPayload(payload *UpdateDeployOKBody) *UpdateDeployOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update deploy o k response
func (o *UpdateDeployOK) SetPayload(payload *UpdateDeployOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateDeployOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateDeployInternalServerErrorCode is the HTTP code returned for type UpdateDeployInternalServerError
const UpdateDeployInternalServerErrorCode int = 500

/*UpdateDeployInternalServerError 内部错误

swagger:response updateDeployInternalServerError
*/
type UpdateDeployInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewUpdateDeployInternalServerError creates UpdateDeployInternalServerError with default headers values
func NewUpdateDeployInternalServerError() *UpdateDeployInternalServerError {

	return &UpdateDeployInternalServerError{}
}

// WithPayload adds the payload to the update deploy internal server error response
func (o *UpdateDeployInternalServerError) WithPayload(payload *models.Error) *UpdateDeployInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update deploy internal server error response
func (o *UpdateDeployInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateDeployInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
