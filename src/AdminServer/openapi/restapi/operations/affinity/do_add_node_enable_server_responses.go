// Code generated by go-swagger; DO NOT EDIT.

package affinity

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"tarsadmin/openapi/models"
)

// DoAddNodeEnableServerOKCode is the HTTP code returned for type DoAddNodeEnableServerOK
const DoAddNodeEnableServerOKCode int = 200

/*DoAddNodeEnableServerOK OK

swagger:response doAddNodeEnableServerOK
*/
type DoAddNodeEnableServerOK struct {

	/*
	  In: Body
	*/
	Payload *DoAddNodeEnableServerOKBody `json:"body,omitempty"`
}

// NewDoAddNodeEnableServerOK creates DoAddNodeEnableServerOK with default headers values
func NewDoAddNodeEnableServerOK() *DoAddNodeEnableServerOK {

	return &DoAddNodeEnableServerOK{}
}

// WithPayload adds the payload to the do add node enable server o k response
func (o *DoAddNodeEnableServerOK) WithPayload(payload *DoAddNodeEnableServerOKBody) *DoAddNodeEnableServerOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the do add node enable server o k response
func (o *DoAddNodeEnableServerOK) SetPayload(payload *DoAddNodeEnableServerOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DoAddNodeEnableServerOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DoAddNodeEnableServerInternalServerErrorCode is the HTTP code returned for type DoAddNodeEnableServerInternalServerError
const DoAddNodeEnableServerInternalServerErrorCode int = 500

/*DoAddNodeEnableServerInternalServerError 内部错误

swagger:response doAddNodeEnableServerInternalServerError
*/
type DoAddNodeEnableServerInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDoAddNodeEnableServerInternalServerError creates DoAddNodeEnableServerInternalServerError with default headers values
func NewDoAddNodeEnableServerInternalServerError() *DoAddNodeEnableServerInternalServerError {

	return &DoAddNodeEnableServerInternalServerError{}
}

// WithPayload adds the payload to the do add node enable server internal server error response
func (o *DoAddNodeEnableServerInternalServerError) WithPayload(payload *models.Error) *DoAddNodeEnableServerInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the do add node enable server internal server error response
func (o *DoAddNodeEnableServerInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DoAddNodeEnableServerInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
