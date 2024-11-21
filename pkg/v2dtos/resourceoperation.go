// Copyright (C) 2020-2024 IOTech Ltd

package v2dtos

import (
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/v2models"
)

// ResourceOperation and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.1.0#/ResourceOperation
type ResourceOperation struct {
	DeviceResource string            `json:"deviceResource" yaml:"deviceResource" validate:"required"` // The replacement of Object field
	DefaultValue   string            `json:"defaultValue,omitempty" yaml:"defaultValue,omitempty"`
	Mappings       map[string]string `json:"mappings,omitempty" yaml:"mappings,omitempty"`
}

// ToResourceOperationModel transforms the ResourceOperation DTO to the ResourceOperation model
func ToResourceOperationModel(ro ResourceOperation) v2models.ResourceOperation {
	return v2models.ResourceOperation{
		DeviceResource: ro.DeviceResource,
		DefaultValue:   ro.DefaultValue,
		Mappings:       ro.Mappings,
	}
}

// FromResourceOperationModelToDTO transforms the ResourceOperation model to the ResourceOperation DTO
func FromResourceOperationModelToDTO(ro v2models.ResourceOperation) ResourceOperation {
	return ResourceOperation{
		DeviceResource: ro.DeviceResource,
		DefaultValue:   ro.DefaultValue,
		Mappings:       ro.Mappings,
	}
}

func ToResourceOperationModels(dtos []ResourceOperation) []v2models.ResourceOperation {
	resourceOperations := make([]v2models.ResourceOperation, len(dtos))
	for i, ro := range dtos {
		resourceOperations[i] = ToResourceOperationModel(ro)
	}
	return resourceOperations
}
