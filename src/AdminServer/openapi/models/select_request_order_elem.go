// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// SelectRequestOrderElem select request order elem
//
// swagger:model SelectRequestOrderElem
type SelectRequestOrderElem struct {

	// column
	Column string `json:"column,omitempty"`

	// order
	// Enum: [asc desc]
	Order string `json:"order,omitempty"`
}

// Validate validates this select request order elem
func (m *SelectRequestOrderElem) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateOrder(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var selectRequestOrderElemTypeOrderPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["asc","desc"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		selectRequestOrderElemTypeOrderPropEnum = append(selectRequestOrderElemTypeOrderPropEnum, v)
	}
}

const (

	// SelectRequestOrderElemOrderAsc captures enum value "asc"
	SelectRequestOrderElemOrderAsc string = "asc"

	// SelectRequestOrderElemOrderDesc captures enum value "desc"
	SelectRequestOrderElemOrderDesc string = "desc"
)

// prop value enum
func (m *SelectRequestOrderElem) validateOrderEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, selectRequestOrderElemTypeOrderPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *SelectRequestOrderElem) validateOrder(formats strfmt.Registry) error {

	if swag.IsZero(m.Order) { // not required
		return nil
	}

	// value enum
	if err := m.validateOrderEnum("order", "body", m.Order); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SelectRequestOrderElem) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SelectRequestOrderElem) UnmarshalBinary(b []byte) error {
	var res SelectRequestOrderElem
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}