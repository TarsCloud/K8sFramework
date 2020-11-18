// Code generated by go-swagger; DO NOT EDIT.

package server_pod

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// SelectPodPerishedHandlerFunc turns a function with the right signature into a select pod perished handler
type SelectPodPerishedHandlerFunc func(SelectPodPerishedParams) middleware.Responder

// Handle executing the request and returning a response
func (fn SelectPodPerishedHandlerFunc) Handle(params SelectPodPerishedParams) middleware.Responder {
	return fn(params)
}

// SelectPodPerishedHandler interface for that can handle valid select pod perished params
type SelectPodPerishedHandler interface {
	Handle(SelectPodPerishedParams) middleware.Responder
}

// NewSelectPodPerished creates a new http.Handler for the select pod perished operation
func NewSelectPodPerished(ctx *middleware.Context, handler SelectPodPerishedHandler) *SelectPodPerished {
	return &SelectPodPerished{Context: ctx, Handler: handler}
}

/*SelectPodPerished swagger:route GET /servers/perishedPods server-pod k8s selectPodPerished

获取当前Pod，web=/pod_list_history

not be implemented

*/
type SelectPodPerished struct {
	Context *middleware.Context
	Handler SelectPodPerishedHandler
}

func (o *SelectPodPerished) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewSelectPodPerishedParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}