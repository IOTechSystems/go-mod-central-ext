// Copyright (C) 2022-2025 IOTech Ltd

package common

// Constants related for provisionWatcher discoveredDevice or properties
const (
	IOTechPrefix = "IOTech_"

	ProtocolName       = IOTechPrefix + "ProtocolName"
	DeviceNamePattern  = IOTechPrefix + "DeviceNamePattern"
	DeviceDescription  = IOTechPrefix + "DeviceDescription"
	DeviceLabels       = IOTechPrefix + "DeviceLabels"
	ProfileNamePattern = IOTechPrefix + "ProfileNamePattern"
	ProfileDescription = IOTechPrefix + "ProfileDescription"
	ProfileLabels      = IOTechPrefix + "ProfileLabels"
	ProfileScanOptions = IOTechPrefix + "ProfileScanOptions"

	TransformScript = IOTechPrefix + "TransformScript"
)

// constants relate to the remote edge node
const (
	BrokerName                            = "brokerName"
	EdgeInstName                          = "edgeinstName"
	DeviceServiceName                     = "deviceServiceName"
	TopicPatternFieldGroupName            = "GROUP_NAME"
	TopicPatternFieldInstName             = "INSTANCE_NAME"
	TopicPatternFieldCentralGroupName     = "EDGECENTRAL_GROUP_NAME"
	TopicPatternFieldCentralInstName      = "EDGECENTRAL_INSTANCE_NAME"
	TopicPatternFieldDeviceServiceName    = "DEVICE_SERVICE_NAME"
	TopicPatternFieldKeyCentralGroupName  = "${" + TopicPatternFieldCentralGroupName + "}"
	TopicPatternFieldKeyCentralInstName   = "${" + TopicPatternFieldCentralInstName + "}"
	TopicPatternFieldKeyDeviceServiceName = "${" + TopicPatternFieldDeviceServiceName + "}"
	CentralNodeRequestTopicKey            = "RequestTopic"
	CentralNodeReplyTopicKey              = "ReplyTopic"
)

// Constants relate to the service status error from sys-mgmt inspect operation
const (
	ServiceIsNotRunningButShouldBe = "service is not running but should be"
	ServiceIsRunningButShouldNotBe = "service is running but shouldn't be"
)

// Constants related to how services identify themselves in the Service Registry
const (
	SupportProvisionServiceKey          = "support-provision"
	SupportSparkplugServiceKey          = "support-sparkplug"
	SupportSparkplugHistorianServiceKey = "support-sparkplug-historian"
	SupportRulesEngineServiceKey        = "support-rulesengine"
	EdgeHistorianServiceKey             = "edge-historian"
)

// Constants related for Notification Category
const (
	DisconnectAlert      = "Disconnection"
	DeviceOperatingState = "DeviceOperatingState"
)

// Constants related for DeviceChangedNotification
const (
	DeviceCreateAction = "Device creation"
	DeviceUpdateAction = "Device update"
	DeviceRemoveAction = "Device removal"

	DeviceChangedNotificationCategory = "DEVICE_CHANGED"
)

const (
	SystemManagementServiceKey = "sys-mgmt"

	BacnetAddress        = "Address"
	BacnetCOV            = "COV"
	BacnetCOVConfirmed   = "Confirmed"
	BacnetCOVLifetime    = "Lifetime"
	BacnetCOVPropName    = "BACnet-COVs"
	BacnetDeviceInstance = "DeviceInstance"
	BacnetIP             = "BACnet-IP"
	BacnetMSTP           = "BACnet-MSTP"
	BacnetPort           = "Port"

	Gps                   = "GPS"
	GpsGpsdPort           = "GpsdPort"
	GpsGpsdRetries        = "GpsdRetries"
	GpsGpsdConnTimeout    = "GpsdConnTimeout"
	GpsGpsdRequestTimeout = "GpsdRequestTimeout"

	ModbusTcp                       = "modbus-tcp"
	ModbusRtu                       = "modbus-rtu"
	ModbusAddress                   = "Address"
	ModbusBaudRate                  = "BaudRate"
	ModbusDataBits                  = "DataBits"
	ModbusParity                    = "Parity"
	ModbusPort                      = "Port"
	ModbusReadMaxBitsCoils          = "ReadMaxBitsCoils"
	ModbusReadMaxBitsDiscreteInputs = "ReadMaxBitsDiscreteInputs"
	ModbusReadMaxHoldingRegisters   = "ReadMaxHoldingRegisters"
	ModbusReadMaxInputRegisters     = "ReadMaxInputRegisters"
	ModbusStopBits                  = "StopBits"
	ModbusUnitID                    = "UnitID"
	ModbusWriteMaxBitsCoils         = "WriteMaxBitsCoils"
	ModbusWriteMaxHoldingRegisters  = "WriteMaxHoldingRegisters"

	Opcua                           = "OPC-UA"
	OpcuaBrowseDepth                = "BrowseDepth"
	OpcuaBrowsePublishInterval      = "BrowsePublishInterval"
	OpcuaConnectionReadingPostDelay = "ConnectionReadingPostDelay"
	OpcuaIDType                     = "IDType"
	OpcuaNodesPerBrowse             = "NodesPerBrowse"
	OpcuaReadBatchSize              = "ReadBatchSize"
	OpcuaRequestedSessionTimeout    = "RequestedSessionTimeout"
	OpcuaSessionKeepAliveInterval   = "SessionKeepAliveInterval"
	OpcuaWriteBatchSize             = "WriteBatchSize"

	S7     = "S7"
	S7Rack = "Rack"
	S7Slot = "Slot"

	EtherNetIP                  = "ethernet-ip"
	EtherNetIPXRT               = "EtherNet-IP" // XRT only accept EtherNet-IP as protocol name
	EtherNetIPO2T               = "O2T"
	EtherNetIPT2O               = "T2O"
	EtherNetIPExplicitConnected = "ExplicitConnected"
	EtherNetIPDeviceResource    = "DeviceResource"
	EtherNetIPSaveValue         = "SaveValue"
	EtherNetIPConnectionType    = "ConnectionType"
	EtherNetIPRPI               = "RPI"
	EtherNetIPPriority          = "Priority"
	EtherNetIPOwnership         = "Ownership"
	EtherNetIPKey               = "Key"
	EtherNetIPMethod            = "Method"
	EtherNetIPVendorID          = "VendorID"
	EtherNetIPDeviceType        = "DeviceType"
	EtherNetIPProductCode       = "ProductCode"
	EtherNetIPMajorRevision     = "MajorRevision"
	EtherNetIPMinorRevision     = "MinorRevision"
	EtherNetIPAddress           = "Address"
)

// Constants related for proxy auth
const (
	VerificationKeyType = "verification"
	SigningKeyType      = "signing"
	Issuer              = "issuer"

	Query        = "QUERY"
	Mutation     = "MUTATION"
	Subscription = "SUBSCRIPTION"
)

// Constants related for filter
const (
	OUT = "OUT"
	IN  = "IN"
)
