// Code generated by go-swagger; DO NOT EDIT.

package server_option

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"tafadmin/openapi/models"
)

// DoPreviewTemplateContentOKCode is the HTTP code returned for type DoPreviewTemplateContentOK
const DoPreviewTemplateContentOKCode int = 200

/*DoPreviewTemplateContentOK OK

swagger:response doPreviewTemplateContentOK
*/
type DoPreviewTemplateContentOK struct {

	/*
	  In: Body
	*/
	Payload *DoPreviewTemplateContentOKBody `json:"body,omitempty"`
}

// NewDoPreviewTemplateContentOK creates DoPreviewTemplateContentOK with default headers values
func NewDoPreviewTemplateContentOK() *DoPreviewTemplateContentOK {

	return &DoPreviewTemplateContentOK{}
}

// WithPayload adds the payload to the do preview template content o k response
func (o *DoPreviewTemplateContentOK) WithPayload(payload *DoPreviewTemplateContentOKBody) *DoPreviewTemplateContentOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the do preview template content o k response
func (o *DoPreviewTemplateContentOK) SetPayload(payload *DoPreviewTemplateContentOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DoPreviewTemplateContentOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DoPreviewTemplateContentInternalServerErrorCode is the HTTP code returned for type DoPreviewTemplateContentInternalServerError
const DoPreviewTemplateContentInternalServerErrorCode int = 500

/*DoPreviewTemplateContentInternalServerError 内部错误

swagger:response doPreviewTemplateContentInternalServerError
*/
type DoPreviewTemplateContentInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDoPreviewTemplateContentInternalServerError creates DoPreviewTemplateContentInternalServerError with default headers values
func NewDoPreviewTemplateContentInternalServerError() *DoPreviewTemplateContentInternalServerError {

	return &DoPreviewTemplateContentInternalServerError{}
}

// WithPayload adds the payload to the do preview template content internal server error response
func (o *DoPreviewTemplateContentInternalServerError) WithPayload(payload *models.Error) *DoPreviewTemplateContentInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the do preview template content internal server error response
func (o *DoPreviewTemplateContentInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DoPreviewTemplateContentInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}