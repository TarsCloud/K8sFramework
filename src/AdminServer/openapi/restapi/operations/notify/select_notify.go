// Code generated by go-swagger; DO NOT EDIT.

package notify

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// SelectNotifyHandlerFunc turns a function with the right signature into a select notify handler
type SelectNotifyHandlerFunc func(SelectNotifyParams) middleware.Responder

// Handle executing the request and returning a response
func (fn SelectNotifyHandlerFunc) Handle(params SelectNotifyParams) middleware.Responder {
	return fn(params)
}

// SelectNotifyHandler interface for that can handle valid select notify params
type SelectNotifyHandler interface {
	Handle(SelectNotifyParams) middleware.Responder
}

// NewSelectNotify creates a new http.Handler for the select notify operation
func NewSelectNotify(ctx *middleware.Context, handler SelectNotifyHandler) *SelectNotify {
	return &SelectNotify{Context: ctx, Handler: handler}
}

/*SelectNotify swagger:route GET /notifies notify mysql selectNotify

拉取notify，web=/server_notify_list

columns key=['NotifyId', 'AppServer', 'PodName', 'NotifyLevel', 'NotifyMessage', 'NotifyTime', 'NotifyThread', 'NotifySource']

*/
type SelectNotify struct {
	Context *middleware.Context
	Handler SelectNotifyHandler
}

func (o *SelectNotify) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewSelectNotifyParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
