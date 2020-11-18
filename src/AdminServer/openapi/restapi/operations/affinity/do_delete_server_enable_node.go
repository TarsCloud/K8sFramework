// Code generated by go-swagger; DO NOT EDIT.

package affinity

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// DoDeleteServerEnableNodeHandlerFunc turns a function with the right signature into a do delete server enable node handler
type DoDeleteServerEnableNodeHandlerFunc func(DoDeleteServerEnableNodeParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DoDeleteServerEnableNodeHandlerFunc) Handle(params DoDeleteServerEnableNodeParams) middleware.Responder {
	return fn(params)
}

// DoDeleteServerEnableNodeHandler interface for that can handle valid do delete server enable node params
type DoDeleteServerEnableNodeHandler interface {
	Handle(DoDeleteServerEnableNodeParams) middleware.Responder
}

// NewDoDeleteServerEnableNode creates a new http.Handler for the do delete server enable node operation
func NewDoDeleteServerEnableNode(ctx *middleware.Context, handler DoDeleteServerEnableNodeHandler) *DoDeleteServerEnableNode {
	return &DoDeleteServerEnableNode{Context: ctx, Handler: handler}
}

/*DoDeleteServerEnableNode swagger:route DELETE /affinities/servers affinity k8s doDeleteServerEnableNode

服务禁止可部署的Node，web=/affinity_del_node

*/
type DoDeleteServerEnableNode struct {
	Context *middleware.Context
	Handler DoDeleteServerEnableNodeHandler
}

func (o *DoDeleteServerEnableNode) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDoDeleteServerEnableNodeParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// DoDeleteServerEnableNodeBody do delete server enable node body
//
// swagger:model DoDeleteServerEnableNodeBody
type DoDeleteServerEnableNodeBody struct {

	// metadata
	Metadata *DoDeleteServerEnableNodeParamsBodyMetadata `json:"metadata,omitempty"`
}

// Validate validates this do delete server enable node body
func (o *DoDeleteServerEnableNodeBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMetadata(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *DoDeleteServerEnableNodeBody) validateMetadata(formats strfmt.Registry) error {

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
func (o *DoDeleteServerEnableNodeBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *DoDeleteServerEnableNodeBody) UnmarshalBinary(b []byte) error {
	var res DoDeleteServerEnableNodeBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// DoDeleteServerEnableNodeOKBody do delete server enable node o k body
//
// swagger:model DoDeleteServerEnableNodeOKBody
type DoDeleteServerEnableNodeOKBody struct {

	// result
	Result int32 `json:"result,omitempty"`
}

// Validate validates this do delete server enable node o k body
func (o *DoDeleteServerEnableNodeOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *DoDeleteServerEnableNodeOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *DoDeleteServerEnableNodeOKBody) UnmarshalBinary(b []byte) error {
	var res DoDeleteServerEnableNodeOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// DoDeleteServerEnableNodeParamsBodyMetadata do delete server enable node params body metadata
//
// swagger:model DoDeleteServerEnableNodeParamsBodyMetadata
type DoDeleteServerEnableNodeParamsBodyMetadata struct {

	// node name
	NodeName []string `json:"NodeName"`

	// server app
	ServerApp string `json:"ServerApp,omitempty"`
}

// Validate validates this do delete server enable node params body metadata
func (o *DoDeleteServerEnableNodeParamsBodyMetadata) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *DoDeleteServerEnableNodeParamsBodyMetadata) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *DoDeleteServerEnableNodeParamsBodyMetadata) UnmarshalBinary(b []byte) error {
	var res DoDeleteServerEnableNodeParamsBodyMetadata
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}