// Copyright (C) 2020-2024 IOTech Ltd

package v2models

// DeviceResource and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.x#/DeviceResource
// Model fields are same as the DTOs documented by this swagger. Exceptions, if any, are noted below.
type DeviceResource struct {
	Description string
	Name        string
	IsHidden    bool
	Tag         string
	Tags        map[string]interface{}
	Properties  ResourceProperties
	Attributes  map[string]interface{}
}
