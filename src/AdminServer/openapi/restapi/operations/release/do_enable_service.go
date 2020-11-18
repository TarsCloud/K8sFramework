// Code generated by go-swagger; DO NOT EDIT.

package release

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// DoEnableServiceHandlerFunc turns a function with the right signature into a do enable service handler
type DoEnableServiceHandlerFunc func(DoEnableServiceParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DoEnableServiceHandlerFunc) Handle(params DoEnableServiceParams) middleware.Responder {
	return fn(params)
}

// DoEnableServiceHandler interface for that can handle valid do enable service params
type DoEnableServiceHandler interface {
	Handle(DoEnableServiceParams) middleware.Responder
}

// NewDoEnableService creates a new http.Handler for the do enable service operation
func NewDoEnableService(ctx *middleware.Context, handler DoEnableServiceHandler) *DoEnableService {
	return &DoEnableService{Context: ctx, Handler: handler}
}

/*DoEnableService swagger:route PUT /releases release k8s doEnableService

发布服务版本，web=/patch_publish

*/
type DoEnableService struct {
	Context *middleware.Context
	Handler DoEnableServiceHandler
}

func (o *DoEnableService) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDoEnableServiceParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// DoEnableServiceBody do enable service body
//
// swagger:model DoEnableServiceBody
type DoEnableServiceBody struct {

	// metadata
	Metadata *DoEnableServiceParamsBodyMetadata `json:"metadata,omitempty"`
}

// Validate validates this do enable service body
func (o *DoEnableServiceBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMetadata(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *DoEnableServiceBody) validateMetadata(formats strfmt.Registry) error {

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

// MarshalBinary interface implementation
func (o *DoEnableServiceBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *DoEnableServiceBody) UnmarshalBinary(b []byte) error {
	var res DoEnableServiceBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// DoEnableServiceOKBody do enable service o k body
//
// swagger:model DoEnableServiceOKBody
type DoEnableServiceOKBody struct {

	// result
	Result int32 `json:"result,omitempty"`
}

// Validate validates this do enable service o k body
func (o *DoEnableServiceOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *DoEnableServiceOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *DoEnableServiceOKBody) UnmarshalBinary(b []byte) error {
	var res DoEnableServiceOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// DoEnableServiceParamsBodyMetadata do enable service params body metadata
//
// swagger:model DoEnableServiceParamsBodyMetadata
type DoEnableServiceParamsBodyMetadata struct {

	// enable mark
	EnableMark string `json:"EnableMark,omitempty"`

	// replicas
	// Required: true
	Replicas *int32 `json:"Replicas"`

	// server Id
	// Required: true
	ServerID *string `json:"ServerId"`

	// service Id
	// Required: true
	ServiceID *string `json:"ServiceId"`
}

// Validate validates this do enable service params body metadata
func (o *DoEnableServiceParamsBodyMetadata) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateReplicas(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateServerID(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateServiceID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *DoEnableServiceParamsBodyMetadata) validateReplicas(formats strfmt.Registry) error {

	if err := validate.Required("Params"+"."+"metadata"+"."+"Replicas", "body", o.Replicas); err != nil {
		return err
	}

	return nil
}

func (o *DoEnableServiceParamsBodyMetadata) validateServerID(formats strfmt.Registry) error {

	if err := validate.Required("Params"+"."+"metadata"+"."+"ServerId", "body", o.ServerID); err != nil {
		return err
	}

	return nil
}

func (o *DoEnableServiceParamsBodyMetadata) validateServiceID(formats strfmt.Registry) error {

	if err := validate.Required("Params"+"."+"metadata"+"."+"ServiceId", "body", o.ServiceID); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *DoEnableServiceParamsBodyMetadata) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *DoEnableServiceParamsBodyMetadata) UnmarshalBinary(b []byte) error {
	var res DoEnableServiceParamsBodyMetadata
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
