// Code generated by go-swagger; DO NOT EDIT.

package business

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

// CreateBusinessHandlerFunc turns a function with the right signature into a create business handler
type CreateBusinessHandlerFunc func(CreateBusinessParams) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateBusinessHandlerFunc) Handle(params CreateBusinessParams) middleware.Responder {
	return fn(params)
}

// CreateBusinessHandler interface for that can handle valid create business params
type CreateBusinessHandler interface {
	Handle(CreateBusinessParams) middleware.Responder
}

// NewCreateBusiness creates a new http.Handler for the create business operation
func NewCreateBusiness(ctx *middleware.Context, handler CreateBusinessHandler) *CreateBusiness {
	return &CreateBusiness{Context: ctx, Handler: handler}
}

/*CreateBusiness swagger:route POST /businesses business k8s createBusiness

创建业务列表，web=/business_create

*/
type CreateBusiness struct {
	Context *middleware.Context
	Handler CreateBusinessHandler
}

func (o *CreateBusiness) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewCreateBusinessParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// CreateBusinessBody create business body
//
// swagger:model CreateBusinessBody
type CreateBusinessBody struct {

	// metadata
	Metadata *CreateBusinessParamsBodyMetadata `json:"metadata,omitempty"`
}

// Validate validates this create business body
func (o *CreateBusinessBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMetadata(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *CreateBusinessBody) validateMetadata(formats strfmt.Registry) error {

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
func (o *CreateBusinessBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *CreateBusinessBody) UnmarshalBinary(b []byte) error {
	var res CreateBusinessBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// CreateBusinessOKBody create business o k body
//
// swagger:model CreateBusinessOKBody
type CreateBusinessOKBody struct {

	// result
	Result int32 `json:"result,omitempty"`
}

// Validate validates this create business o k body
func (o *CreateBusinessOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *CreateBusinessOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *CreateBusinessOKBody) UnmarshalBinary(b []byte) error {
	var res CreateBusinessOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// CreateBusinessParamsBodyMetadata create business params body metadata
//
// swagger:model CreateBusinessParamsBodyMetadata
type CreateBusinessParamsBodyMetadata struct {

	// business mark
	BusinessMark string `json:"BusinessMark,omitempty"`

	// business name
	// Required: true
	// Pattern: ^[a-zA-Z0-9.:_-]{1,128}\z
	BusinessName *string `json:"BusinessName"`

	// business order
	// Required: true
	// Maximum: 100
	// Minimum: 1
	BusinessOrder *int32 `json:"BusinessOrder"`

	// business show
	// Required: true
	BusinessShow *string `json:"BusinessShow"`

	// create person
	CreatePerson string `json:"CreatePerson,omitempty"`
}

// Validate validates this create business params body metadata
func (o *CreateBusinessParamsBodyMetadata) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateBusinessName(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateBusinessOrder(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateBusinessShow(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *CreateBusinessParamsBodyMetadata) validateBusinessName(formats strfmt.Registry) error {

	if err := validate.Required("Params"+"."+"metadata"+"."+"BusinessName", "body", o.BusinessName); err != nil {
		return err
	}

	if err := validate.Pattern("Params"+"."+"metadata"+"."+"BusinessName", "body", string(*o.BusinessName), `^[a-zA-Z0-9.:_-]{1,128}\z`); err != nil {
		return err
	}

	return nil
}

func (o *CreateBusinessParamsBodyMetadata) validateBusinessOrder(formats strfmt.Registry) error {

	if err := validate.Required("Params"+"."+"metadata"+"."+"BusinessOrder", "body", o.BusinessOrder); err != nil {
		return err
	}

	if err := validate.MinimumInt("Params"+"."+"metadata"+"."+"BusinessOrder", "body", int64(*o.BusinessOrder), 1, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("Params"+"."+"metadata"+"."+"BusinessOrder", "body", int64(*o.BusinessOrder), 100, false); err != nil {
		return err
	}

	return nil
}

func (o *CreateBusinessParamsBodyMetadata) validateBusinessShow(formats strfmt.Registry) error {

	if err := validate.Required("Params"+"."+"metadata"+"."+"BusinessShow", "body", o.BusinessShow); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *CreateBusinessParamsBodyMetadata) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *CreateBusinessParamsBodyMetadata) UnmarshalBinary(b []byte) error {
	var res CreateBusinessParamsBodyMetadata
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
