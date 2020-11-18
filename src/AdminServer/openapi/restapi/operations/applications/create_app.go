// Code generated by go-swagger; DO NOT EDIT.

package applications

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

// CreateAppHandlerFunc turns a function with the right signature into a create app handler
type CreateAppHandlerFunc func(CreateAppParams) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateAppHandlerFunc) Handle(params CreateAppParams) middleware.Responder {
	return fn(params)
}

// CreateAppHandler interface for that can handle valid create app params
type CreateAppHandler interface {
	Handle(CreateAppParams) middleware.Responder
}

// NewCreateApp creates a new http.Handler for the create app operation
func NewCreateApp(ctx *middleware.Context, handler CreateAppHandler) *CreateApp {
	return &CreateApp{Context: ctx, Handler: handler}
}

/*CreateApp swagger:route POST /applications applications k8s createApp

创建应用列表，web=/application_create

*/
type CreateApp struct {
	Context *middleware.Context
	Handler CreateAppHandler
}

func (o *CreateApp) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewCreateAppParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// CreateAppBody create app body
//
// swagger:model CreateAppBody
type CreateAppBody struct {

	// metadata
	Metadata *CreateAppParamsBodyMetadata `json:"metadata,omitempty"`
}

// Validate validates this create app body
func (o *CreateAppBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMetadata(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *CreateAppBody) validateMetadata(formats strfmt.Registry) error {

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
func (o *CreateAppBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *CreateAppBody) UnmarshalBinary(b []byte) error {
	var res CreateAppBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// CreateAppOKBody create app o k body
//
// swagger:model CreateAppOKBody
type CreateAppOKBody struct {

	// result
	Result int32 `json:"result,omitempty"`
}

// Validate validates this create app o k body
func (o *CreateAppOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *CreateAppOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *CreateAppOKBody) UnmarshalBinary(b []byte) error {
	var res CreateAppOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// CreateAppParamsBodyMetadata create app params body metadata
//
// swagger:model CreateAppParamsBodyMetadata
type CreateAppParamsBodyMetadata struct {

	// app mark
	AppMark string `json:"AppMark,omitempty"`

	// app name
	// Required: true
	// Pattern: ^[a-zA-Z0-9]{1,24}\z
	AppName *string `json:"AppName"`

	// business name
	// Pattern: ^[a-zA-Z0-9.:_-]{1,128}\z
	BusinessName string `json:"BusinessName,omitempty"`

	// create person
	CreatePerson string `json:"CreatePerson,omitempty"`
}

// Validate validates this create app params body metadata
func (o *CreateAppParamsBodyMetadata) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateAppName(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateBusinessName(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *CreateAppParamsBodyMetadata) validateAppName(formats strfmt.Registry) error {

	if err := validate.Required("Params"+"."+"metadata"+"."+"AppName", "body", o.AppName); err != nil {
		return err
	}

	if err := validate.Pattern("Params"+"."+"metadata"+"."+"AppName", "body", string(*o.AppName), `^[a-zA-Z0-9]{1,24}\z`); err != nil {
		return err
	}

	return nil
}

func (o *CreateAppParamsBodyMetadata) validateBusinessName(formats strfmt.Registry) error {

	if swag.IsZero(o.BusinessName) { // not required
		return nil
	}

	if err := validate.Pattern("Params"+"."+"metadata"+"."+"BusinessName", "body", string(o.BusinessName), `^[a-zA-Z0-9.:_-]{1,128}\z`); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *CreateAppParamsBodyMetadata) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *CreateAppParamsBodyMetadata) UnmarshalBinary(b []byte) error {
	var res CreateAppParamsBodyMetadata
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
