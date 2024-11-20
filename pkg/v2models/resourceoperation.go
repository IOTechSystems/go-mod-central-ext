// Copyright (C) 2020-2024 IOTech Ltd

package v2models

// ResourceOperation and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.x#/ResourceOperation
// Model fields are same as the DTOs documented by this swagger. Exceptions, if any, are noted below.
type ResourceOperation struct {
	DeviceResource string
	DefaultValue   string
	Mappings       map[string]string
}
