// Copyright (C) 2020-2024 IOTech Ltd

package v2models

// Device and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.x#/Device
// Model fields are same as the DTOs documented by this swagger. Exceptions, if any, are noted below.
type Device struct {
	DBTimestamp
	Id             string
	Name           string
	Description    string
	AdminState     AdminState
	OperatingState OperatingState
	ProtocolName   string
	Protocols      map[string]ProtocolProperties
	LastConnected  int64 // Deprecated: will be replaced by Metrics in v3
	LastReported   int64 // Deprecated: will be replaced by Metrics in v3
	Labels         []string
	Location       interface{}
	Tags           map[string]interface{}
	ServiceName    string
	ProfileName    string
	AutoEvents     []AutoEvent
	Notify         bool
	// Properties are required for device discovery, the feedback from JCI was that when discovering a device they
	// want as much info as we can find about the device (so for example in a big system they have a better idea of what
	// actual device is being provisioned). So we added a properties field to carry this information. IMHO this is a valid
	// generic requirement to support discovery.
	Properties map[string]any
}

// ProtocolProperties contains the device connection information in key/value pair
type ProtocolProperties map[string]string

// AdminState controls the range of values which constitute valid administrative states for a device
type AdminState string

// OperatingState is an indication of the operations of the device.
type OperatingState string
