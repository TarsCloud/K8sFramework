// Code generated by go-swagger; DO NOT EDIT.

package server_option

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

// NewUpdateServerOptionParams creates a new UpdateServerOptionParams object
// no default values defined in spec.
func NewUpdateServerOptionParams() UpdateServerOptionParams {

	return UpdateServerOptionParams{}
}

// UpdateServerOptionParams contains all the bound params for the update server option operation
// typically these are obtained from a http.Request
//
// swagger:parameters updateServerOption
type UpdateServerOptionParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  In: body
	*/
	Params UpdateServerOptionBody
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewUpdateServerOptionParams() beforehand.
func (o *UpdateServerOptionParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body UpdateServerOptionBody
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			res = append(res, errors.NewParseError("params", "body", "", err))
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Params = body
			}
		}
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
