// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ConfigMeta config meta
//
// swagger:model ConfigMeta
type ConfigMeta struct {

	// config content
	// Required: true
	ConfigContent *string `json:"ConfigContent"`

	// config Id
	// Read Only: true
	ConfigID string `json:"ConfigId,omitempty"`

	// config mark
	ConfigMark string `json:"ConfigMark,omitempty"`

	// config name
	// Required: true
	ConfigName *string `json:"ConfigName"`

	// config version
	// Read Only: true
	ConfigVersion int32 `json:"ConfigVersion,omitempty"`

	// create person
	CreatePerson string `json:"CreatePerson,omitempty"`

	// pod seq
	PodSeq string `json:"PodSeq,omitempty"`

	// server app
	ServerApp string `json:"ServerApp,omitempty"`

	// server name
	ServerName string `json:"ServerName,omitempty"`
}

// Validate validates this config meta
func (m *ConfigMeta) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateConfigContent(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateConfigName(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ConfigMeta) validateConfigContent(formats strfmt.Registry) error {

	if err := validate.Required("ConfigContent", "body", m.ConfigContent); err != nil {
		return err
	}

	return nil
}

func (m *ConfigMeta) validateConfigName(formats strfmt.Registry) error {

	if err := validate.Required("ConfigName", "body", m.ConfigName); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ConfigMeta) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ConfigMeta) UnmarshalBinary(b []byte) error {
	var res ConfigMeta
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}