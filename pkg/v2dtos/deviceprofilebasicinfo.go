// Copyright (C) 2022-2024 IOTech Ltd

package v2dtos

// DeviceProfileBasicInfo and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.2.0#/DeviceProfileBasicInfo
type DeviceProfileBasicInfo struct {
	Id           string   `json:"id,omitempty" validate:"omitempty,uuid" yaml:"id,omitempty"`
	Name         string   `json:"name" yaml:"name" validate:"required,edgex-dto-none-empty-string,edgex-dto-no-reserved-chars"`
	Manufacturer string   `json:"manufacturer,omitempty" yaml:"manufacturer,omitempty"`
	Description  string   `json:"description,omitempty" yaml:"description,omitempty"`
	Model        string   `json:"model,omitempty" yaml:"model,omitempty"`
	Labels       []string `json:"labels,omitempty" yaml:"labels,flow,omitempty"`
}

// UpdateDeviceProfileBasicInfo and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.2.0#/DeviceProfileBasicInfo
type UpdateDeviceProfileBasicInfo struct {
	Id           *string  `json:"id" validate:"required_without=Name,edgex-dto-uuid"`
	Name         *string  `json:"name" validate:"required_without=Id,edgex-dto-none-empty-string,edgex-dto-no-reserved-chars"`
	Manufacturer *string  `json:"manufacturer"`
	Description  *string  `json:"description"`
	Model        *string  `json:"model"`
	Labels       []string `json:"labels"`
}
