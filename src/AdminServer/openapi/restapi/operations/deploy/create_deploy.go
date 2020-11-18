// Code generated by go-swagger; DO NOT EDIT.

package deploy

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"tafadmin/openapi/models"
)

// CreateDeployHandlerFunc turns a function with the right signature into a create deploy handler
type CreateDeployHandlerFunc func(CreateDeployParams) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateDeployHandlerFunc) Handle(params CreateDeployParams) middleware.Responder {
	return fn(params)
}

// CreateDeployHandler interface for that can handle valid create deploy params
type CreateDeployHandler interface {
	Handle(CreateDeployParams) middleware.Responder
}

// NewCreateDeploy creates a new http.Handler for the create deploy operation
func NewCreateDeploy(ctx *middleware.Context, handler CreateDeployHandler) *CreateDeploy {
	return &CreateDeploy{Context: ctx, Handler: handler}
}

/*CreateDeploy swagger:route POST /deploys deploy k8s createDeploy

创建服务部署，web=/deploy_create

*/
type CreateDeploy struct {
	Context *middleware.Context
	Handler CreateDeployHandler
}

func (o *CreateDeploy) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewCreateDeployParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// CreateDeployBody create deploy body
//
// swagger:model CreateDeployBody
type CreateDeployBody struct {

	// metadata
	Metadata *models.DeployMeta `json:"metadata,omitempty"`
}

// Validate validates this create deploy body
func (o *CreateDeployBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMetadata(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *CreateDeployBody) validateMetadata(formats strfmt.Registry) error {

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
func (o *CreateDeployBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *CreateDeployBody) UnmarshalBinary(b []byte) error {
	var res CreateDeployBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// CreateDeployOKBody create deploy o k body
//
// swagger:model CreateDeployOKBody
type CreateDeployOKBody struct {

	// result
	Result int32 `json:"result,omitempty"`
}

// Validate validates this create deploy o k body
func (o *CreateDeployOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *CreateDeployOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *CreateDeployOKBody) UnmarshalBinary(b []byte) error {
	var res CreateDeployOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}