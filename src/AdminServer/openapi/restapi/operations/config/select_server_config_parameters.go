// Code generated by go-swagger; DO NOT EDIT.

package config

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

// NewSelectServerConfigParams creates a new SelectServerConfigParams object
// no default values defined in spec.
func NewSelectServerConfigParams() SelectServerConfigParams {

	return SelectServerConfigParams{}
}

// SelectServerConfigParams contains all the bound params for the select server config operation
// typically these are obtained from a http.Request
//
// swagger:parameters selectServerConfig
type SelectServerConfigParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*string=json.encode(SelectRequestFilter)
	  In: query
	*/
	Filter *string
	/*string=json.encode(SelectRequestLimiter)
	  In: query
	*/
	Limiter *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewSelectServerConfigParams() beforehand.
func (o *SelectServerConfigParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qFilter, qhkFilter, _ := qs.GetOK("Filter")
	if err := o.bindFilter(qFilter, qhkFilter, route.Formats); err != nil {
		res = append(res, err)
	}

	qLimiter, qhkLimiter, _ := qs.GetOK("Limiter")
	if err := o.bindLimiter(qLimiter, qhkLimiter, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindFilter binds and validates parameter Filter from query.
func (o *SelectServerConfigParams) bindFilter(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Filter = &raw

	return nil
}

// bindLimiter binds and validates parameter Limiter from query.
func (o *SelectServerConfigParams) bindLimiter(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Limiter = &raw

	return nil
}
