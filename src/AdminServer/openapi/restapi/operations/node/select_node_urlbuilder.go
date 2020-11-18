// Code generated by go-swagger; DO NOT EDIT.

package node

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"
)

// SelectNodeURL generates an URL for the select node operation
type SelectNodeURL struct {
	Filter  *string
	Limiter *string

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *SelectNodeURL) WithBasePath(bp string) *SelectNodeURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *SelectNodeURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *SelectNodeURL) Build() (*url.URL, error) {
	var _result url.URL

	var _path = "/nodes"

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/admin/v1alpha1/"
	}
	_result.Path = golangswaggerpaths.Join(_basePath, _path)

	qs := make(url.Values)

	var filterQ string
	if o.Filter != nil {
		filterQ = *o.Filter
	}
	if filterQ != "" {
		qs.Set("Filter", filterQ)
	}

	var limiterQ string
	if o.Limiter != nil {
		limiterQ = *o.Limiter
	}
	if limiterQ != "" {
		qs.Set("Limiter", limiterQ)
	}

	_result.RawQuery = qs.Encode()

	return &_result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *SelectNodeURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *SelectNodeURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *SelectNodeURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on SelectNodeURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on SelectNodeURL")
	}

	base, err := o.Build()
	if err != nil {
		return nil, err
	}

	base.Scheme = scheme
	base.Host = host
	return base, nil
}

// StringFull returns the string representation of a complete url
func (o *SelectNodeURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
