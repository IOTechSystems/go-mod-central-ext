// Copyright (C) 2020-2024 IOTech Ltd

package v2dtos

import (
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/v2models"
)

// ResourceProperties and its properties care defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.1.0#/ResourceProperties
type ResourceProperties struct {
	ValueType    string `json:"valueType" yaml:"valueType" validate:"required,edgex-dto-value-type"`
	ReadWrite    string `json:"readWrite" yaml:"readWrite" validate:"required,oneof='R' 'W' 'RW' 'WR'"`
	Units        string `json:"units,omitempty" yaml:"units,omitempty"`
	Minimum      string `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	Maximum      string `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	DefaultValue string `json:"defaultValue,omitempty" yaml:"defaultValue,omitempty"`
	Mask         string `json:"mask,omitempty" yaml:"mask,omitempty"`
	Shift        string `json:"shift,omitempty" yaml:"shift,omitempty"`
	Scale        string `json:"scale,omitempty" yaml:"scale,omitempty"`
	Offset       string `json:"offset,omitempty" yaml:"offset,omitempty"`
	Base         string `json:"base,omitempty" yaml:"base,omitempty"`
	Assertion    string `json:"assertion,omitempty" yaml:"assertion,omitempty"`
	MediaType    string `json:"mediaType,omitempty" yaml:"mediaType,omitempty"`
}

// ToResourcePropertiesModel transforms the ResourceProperties DTO to the ResourceProperties model
func ToResourcePropertiesModel(p ResourceProperties) v2models.ResourceProperties {
	return v2models.ResourceProperties{
		ValueType:    p.ValueType,
		ReadWrite:    p.ReadWrite,
		Units:        p.Units,
		Minimum:      p.Minimum,
		Maximum:      p.Maximum,
		DefaultValue: p.DefaultValue,
		Mask:         p.Mask,
		Shift:        p.Shift,
		Scale:        p.Scale,
		Offset:       p.Offset,
		Base:         p.Base,
		Assertion:    p.Assertion,
		MediaType:    p.MediaType,
	}
}

// FromResourcePropertiesModelToDTO transforms the ResourceProperties Model to the ResourceProperties DTO
func FromResourcePropertiesModelToDTO(p v2models.ResourceProperties) ResourceProperties {
	return ResourceProperties{
		ValueType:    p.ValueType,
		ReadWrite:    p.ReadWrite,
		Units:        p.Units,
		Minimum:      p.Minimum,
		Maximum:      p.Maximum,
		DefaultValue: p.DefaultValue,
		Mask:         p.Mask,
		Shift:        p.Shift,
		Scale:        p.Scale,
		Offset:       p.Offset,
		Base:         p.Base,
		Assertion:    p.Assertion,
		MediaType:    p.MediaType,
	}
}
