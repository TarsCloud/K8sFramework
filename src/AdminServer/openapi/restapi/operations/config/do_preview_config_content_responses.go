// Code generated by go-swagger; DO NOT EDIT.

package config

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"tarsadmin/openapi/models"
)

// DoPreviewConfigContentOKCode is the HTTP code returned for type DoPreviewConfigContentOK
const DoPreviewConfigContentOKCode int = 200

/*DoPreviewConfigContentOK OK

swagger:response doPreviewConfigContentOK
*/
type DoPreviewConfigContentOK struct {

	/*
	  In: Body
	*/
	Payload *DoPreviewConfigContentOKBody `json:"body,omitempty"`
}

// NewDoPreviewConfigContentOK creates DoPreviewConfigContentOK with default headers values
func NewDoPreviewConfigContentOK() *DoPreviewConfigContentOK {

	return &DoPreviewConfigContentOK{}
}

// WithPayload adds the payload to the do preview config content o k response
func (o *DoPreviewConfigContentOK) WithPayload(payload *DoPreviewConfigContentOKBody) *DoPreviewConfigContentOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the do preview config content o k response
func (o *DoPreviewConfigContentOK) SetPayload(payload *DoPreviewConfigContentOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DoPreviewConfigContentOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DoPreviewConfigContentInternalServerErrorCode is the HTTP code returned for type DoPreviewConfigContentInternalServerError
const DoPreviewConfigContentInternalServerErrorCode int = 500

/*DoPreviewConfigContentInternalServerError 内部错误

swagger:response doPreviewConfigContentInternalServerError
*/
type DoPreviewConfigContentInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDoPreviewConfigContentInternalServerError creates DoPreviewConfigContentInternalServerError with default headers values
func NewDoPreviewConfigContentInternalServerError() *DoPreviewConfigContentInternalServerError {

	return &DoPreviewConfigContentInternalServerError{}
}

// WithPayload adds the payload to the do preview config content internal server error response
func (o *DoPreviewConfigContentInternalServerError) WithPayload(payload *models.Error) *DoPreviewConfigContentInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the do preview config content internal server error response
func (o *DoPreviewConfigContentInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DoPreviewConfigContentInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
