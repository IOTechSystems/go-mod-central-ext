// Copyright (C) 2023-2024 IOTech Ltd

package xlsx

// constants relates to the stylesheet names
const (
	devicesSheetName        = "Devices"
	mappingTableSheetName   = "MappingTable"
	autoEventsSheetName     = "AutoEvents"
	deviceInfoSheetName     = "DeviceInfo"
	deviceResourceSheetName = "DeviceResource"
	deviceCommandSheetName  = "DeviceCommand"
)

// constants relates to the header names
const (
	objectCol       = "object"
	pathCol         = "path"
	defaultValueCol = "default value"
	refDeviceName   = "Reference Device Name"
)

// constants relates to the Device DTO field names
const (
	protocols    = "Protocols"
	protocolName = "ProtocolName"
	autoEvents   = "AutoEvents"
	tags         = "Tags"
)

// constants relates to the Device/DeviceProfile/DeviceResource/DeviceCommand DTO field names
const (
	attributes         = "Attributes"
	properties         = "Properties"
	resourceOperations = "ResourceOperations"
	apiVersion         = "ApiVersion"
	description        = "Description"
	isHidden           = "IsHidden"
	readWrite          = "ReadWrite"
	units              = "Units"
	minimum            = "Minimum"
	maximum            = "Maximum"
	defaultValue       = "DefaultValue"
	mask               = "Mask"
	shift              = "Shift"
	scale              = "Scale"
	base               = "Base"
	assertion          = "Assertion"
	mediaType          = "MediaType"
	resourceOperation  = "ResourceOperation"
	adminState         = "AdminState"
	operatingState     = "OperatingState"
	onChange           = "OnChange"
)

// constants relates to the Device protocols
const (
	modbusRTU = "modbus-rtu"
)

const mappingPathSeparator = "."
