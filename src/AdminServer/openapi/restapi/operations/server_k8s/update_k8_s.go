// Code generated by go-swagger; DO NOT EDIT.

package server_k8s

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"tarsadmin/openapi/models"
)

// UpdateK8SHandlerFunc turns a function with the right signature into a update k8 s handler
type UpdateK8SHandlerFunc func(UpdateK8SParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateK8SHandlerFunc) Handle(params UpdateK8SParams) middleware.Responder {
	return fn(params)
}

// UpdateK8SHandler interface for that can handle valid update k8 s params
type UpdateK8SHandler interface {
	Handle(UpdateK8SParams) middleware.Responder
}

// NewUpdateK8S creates a new http.Handler for the update k8 s operation
func NewUpdateK8S(ctx *middleware.Context, handler UpdateK8SHandler) *UpdateK8S {
	return &UpdateK8S{Context: ctx, Handler: handler}
}

/*UpdateK8S swagger:route PATCH /servers/k8s server-k8s k8s updateK8S

更新k8s属性，web=/server_k8s_update

*/
type UpdateK8S struct {
	Context *middleware.Context
	Handler UpdateK8SHandler
}

func (o *UpdateK8S) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUpdateK8SParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// UpdateK8SBody update k8 s body
//
// swagger:model UpdateK8SBody
type UpdateK8SBody struct {

	// metadata
	Metadata *UpdateK8SParamsBodyMetadata `json:"metadata,omitempty"`

	// target
	Target *models.ServerK8S `json:"target,omitempty"`
}

// Validate validates this update k8 s body
func (o *UpdateK8SBody) Validate(formats strfmt.Registry) error {
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

func (o *UpdateK8SBody) validateMetadata(formats strfmt.Registry) error {

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

func (o *UpdateK8SBody) validateTarget(formats strfmt.Registry) error {

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
func (o *UpdateK8SBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *UpdateK8SBody) UnmarshalBinary(b []byte) error {
	var res UpdateK8SBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// UpdateK8SOKBody update k8 s o k body
//
// swagger:model UpdateK8SOKBody
type UpdateK8SOKBody struct {

	// result
	Result int32 `json:"result,omitempty"`
}

// Validate validates this update k8 s o k body
func (o *UpdateK8SOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *UpdateK8SOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *UpdateK8SOKBody) UnmarshalBinary(b []byte) error {
	var res UpdateK8SOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// UpdateK8SParamsBodyMetadata update k8 s params body metadata
//
// swagger:model UpdateK8SParamsBodyMetadata
type UpdateK8SParamsBodyMetadata struct {

	// server Id
	// Required: true
	ServerID *string `json:"ServerId"`
}

// Validate validates this update k8 s params body metadata
func (o *UpdateK8SParamsBodyMetadata) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateServerID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *UpdateK8SParamsBodyMetadata) validateServerID(formats strfmt.Registry) error {

	if err := validate.Required("Params"+"."+"metadata"+"."+"ServerId", "body", o.ServerID); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *UpdateK8SParamsBodyMetadata) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *UpdateK8SParamsBodyMetadata) UnmarshalBinary(b []byte) error {
	var res UpdateK8SParamsBodyMetadata
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
