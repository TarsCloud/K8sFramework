// Code generated by go-swagger; DO NOT EDIT.

package business

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"tafadmin/openapi/models"
)

// DoAddBusinessAppOKCode is the HTTP code returned for type DoAddBusinessAppOK
const DoAddBusinessAppOKCode int = 200

/*DoAddBusinessAppOK OK

swagger:response doAddBusinessAppOK
*/
type DoAddBusinessAppOK struct {

	/*
	  In: Body
	*/
	Payload *DoAddBusinessAppOKBody `json:"body,omitempty"`
}

// NewDoAddBusinessAppOK creates DoAddBusinessAppOK with default headers values
func NewDoAddBusinessAppOK() *DoAddBusinessAppOK {

	return &DoAddBusinessAppOK{}
}

// WithPayload adds the payload to the do add business app o k response
func (o *DoAddBusinessAppOK) WithPayload(payload *DoAddBusinessAppOKBody) *DoAddBusinessAppOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the do add business app o k response
func (o *DoAddBusinessAppOK) SetPayload(payload *DoAddBusinessAppOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DoAddBusinessAppOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DoAddBusinessAppInternalServerErrorCode is the HTTP code returned for type DoAddBusinessAppInternalServerError
const DoAddBusinessAppInternalServerErrorCode int = 500

/*DoAddBusinessAppInternalServerError 内部错误

swagger:response doAddBusinessAppInternalServerError
*/
type DoAddBusinessAppInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDoAddBusinessAppInternalServerError creates DoAddBusinessAppInternalServerError with default headers values
func NewDoAddBusinessAppInternalServerError() *DoAddBusinessAppInternalServerError {

	return &DoAddBusinessAppInternalServerError{}
}

// WithPayload adds the payload to the do add business app internal server error response
func (o *DoAddBusinessAppInternalServerError) WithPayload(payload *models.Error) *DoAddBusinessAppInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the do add business app internal server error response
func (o *DoAddBusinessAppInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DoAddBusinessAppInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}