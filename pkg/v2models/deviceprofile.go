// Copyright (C) 2020-2024 IOTech Ltd

package v2models

// DeviceProfile and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.x#/DeviceProfile
// Model fields are same as the DTOs documented by this swagger. Exceptions, if any, are noted below.
type DeviceProfile struct {
	DBTimestamp
	ApiVersion      string
	Description     string
	Id              string
	Name            string
	Manufacturer    string
	Model           string
	Labels          []string
	DeviceResources []DeviceResource
	DeviceCommands  []DeviceCommand
}
