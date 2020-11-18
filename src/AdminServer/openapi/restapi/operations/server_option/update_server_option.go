// Code generated by go-swagger; DO NOT EDIT.

package server_option

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"tafadmin/openapi/models"
)

// UpdateServerOptionHandlerFunc turns a function with the right signature into a update server option handler
type UpdateServerOptionHandlerFunc func(UpdateServerOptionParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateServerOptionHandlerFunc) Handle(params UpdateServerOptionParams) middleware.Responder {
	return fn(params)
}

// UpdateServerOptionHandler interface for that can handle valid update server option params
type UpdateServerOptionHandler interface {
	Handle(UpdateServerOptionParams) middleware.Responder
}

// NewUpdateServerOption creates a new http.Handler for the update server option operation
func NewUpdateServerOption(ctx *middleware.Context, handler UpdateServerOptionHandler) *UpdateServerOption {
	return &UpdateServerOption{Context: ctx, Handler: handler}
}

/*UpdateServerOption swagger:route PATCH /servers/options server-option k8s updateServerOption

更新私有模板，web=/server_option_update

*/
type UpdateServerOption struct {
	Context *middleware.Context
	Handler UpdateServerOptionHandler
}

func (o *UpdateServerOption) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUpdateServerOptionParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// UpdateServerOptionBody update server option body
//
// swagger:model UpdateServerOptionBody
type UpdateServerOptionBody struct {

	// confirmation
	Confirmation bool `json:"Confirmation,omitempty"`

	// metadata
	Metadata *UpdateServerOptionParamsBodyMetadata `json:"metadata,omitempty"`

	// target
	Target *models.ServerOption `json:"target,omitempty"`
}

// Validate validates this update server option body
func (o *UpdateServerOptionBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMetadata(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateTarget(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *UpdateServerOptionBody) validateMetadata(formats strfmt.Registry) error {

	if swag.IsZero(o.Metadata) { // not required
		return nil
	}

	if o.Metadata != nil {
		if err := o.Metadata.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("Params" + "." + "metadata")
			}
			return err
		}
	}

	return nil
}

func (o *UpdateServerOptionBody) validateTarget(formats strfmt.Registry) error {

	if swag.IsZero(o.Target) { // not required
		return nil
	}

	if o.Target != nil {
		if err := o.Target.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("Params" + "." + "target")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *UpdateServerOptionBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *UpdateServerOptionBody) UnmarshalBinary(b []byte) error {
	var res UpdateServerOptionBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// UpdateServerOptionOKBody update server option o k body
//
// swagger:model UpdateServerOptionOKBody
type UpdateServerOptionOKBody struct {

	// result
	Result int32 `json:"result,omitempty"`
}

// Validate validates this update server option o k body
func (o *UpdateServerOptionOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *UpdateServerOptionOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *UpdateServerOptionOKBody) UnmarshalBinary(b []byte) error {
	var res UpdateServerOptionOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// UpdateServerOptionParamsBodyMetadata update server option params body metadata
//
// swagger:model UpdateServerOptionParamsBodyMetadata
type UpdateServerOptionParamsBodyMetadata struct {

	// server Id
	// Required: true
	ServerID *string `json:"ServerId"`
}

// Validate validates this update server option params body metadata
func (o *UpdateServerOptionParamsBodyMetadata) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateServerID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *UpdateServerOptionParamsBodyMetadata) validateServerID(formats strfmt.Registry) error {

	if err := validate.Required("Params"+"."+"metadata"+"."+"ServerId", "body", o.ServerID); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *UpdateServerOptionParamsBodyMetadata) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *UpdateServerOptionParamsBodyMetadata) UnmarshalBinary(b []byte) error {
	var res UpdateServerOptionParamsBodyMetadata
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
