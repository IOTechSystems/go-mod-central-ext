// Copyright (C) 2020-2024 IOTech Ltd

package v2dtos

import (
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/v2models"
)

// DeviceCommand and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.1.0#/DeviceCommand
type DeviceCommand struct {
	Name               string                 `json:"name" yaml:"name" validate:"required,edgex-dto-none-empty-string,edgex-dto-no-reserved-chars"`
	IsHidden           bool                   `json:"isHidden" yaml:"isHidden"`
	ReadWrite          string                 `json:"readWrite" yaml:"readWrite" validate:"required,oneof='R' 'W' 'RW' 'WR'"`
	ResourceOperations []ResourceOperation    `json:"resourceOperations" yaml:"resourceOperations" validate:"gt=0,dive"`
	Tags               map[string]interface{} `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// UpdateDeviceComman and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.2.0#/DeviceCommand
type UpdateDeviceCommand struct {
	Name     *string `json:"name" validate:"required,edgex-dto-none-empty-string,edgex-dto-no-reserved-chars"`
	IsHidden *bool   `json:"isHidden"`
}

// ToDeviceCommandModel transforms the DeviceCommand DTO to the DeviceCommand model
func ToDeviceCommandModel(dto DeviceCommand) v2models.DeviceCommand {
	resourceOperations := ToResourceOperationModels(dto.ResourceOperations)

	return v2models.DeviceCommand{
		Name:               dto.Name,
		IsHidden:           dto.IsHidden,
		ReadWrite:          dto.ReadWrite,
		ResourceOperations: resourceOperations,
		Tags:               dto.Tags,
	}
}

// ToDeviceCommandModels transforms the DeviceCommand DTOs to the DeviceCommand models
func ToDeviceCommandModels(deviceCommandDTOs []DeviceCommand) []v2models.DeviceCommand {
	deviceCommandModels := make([]v2models.DeviceCommand, len(deviceCommandDTOs))
	for i, p := range deviceCommandDTOs {
		deviceCommandModels[i] = ToDeviceCommandModel(p)
	}
	return deviceCommandModels
}

// FromDeviceCommandModelToDTO transforms the DeviceCommand model to the DeviceCommand DTO
func FromDeviceCommandModelToDTO(d v2models.DeviceCommand) DeviceCommand {
	resourceOperations := make([]ResourceOperation, len(d.ResourceOperations))
	for i, ro := range d.ResourceOperations {
		resourceOperations[i] = FromResourceOperationModelToDTO(ro)
	}
	return DeviceCommand{
		Name:               d.Name,
		IsHidden:           d.IsHidden,
		ReadWrite:          d.ReadWrite,
		ResourceOperations: resourceOperations,
		Tags:               d.Tags,
	}
}

// FromDeviceCommandModelsToDTOs transforms the DeviceCommand models to the DeviceCommand DTOs
func FromDeviceCommandModelsToDTOs(deviceCommandModels []v2models.DeviceCommand) []DeviceCommand {
	deviceCommandDTOs := make([]DeviceCommand, len(deviceCommandModels))
	for i, p := range deviceCommandModels {
		deviceCommandDTOs[i] = FromDeviceCommandModelToDTO(p)
	}
	return deviceCommandDTOs
}
