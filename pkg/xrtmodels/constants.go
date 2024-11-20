// Copyright (C) 2024 IOTech Ltd

package xrtmodels

// Constants related to defined message type and version
const (
	MessageTypeRequest         = "xrt.request:1.0"
	MessageTypeReply           = "xrt.reply:1.0"
	MessageTypeTelemetry       = "xrt.telemetry:1.0"
	MessageTypeDeviceDiscovery = "xrt.device.discovery:1.0"
	MessageTypeDiscovery       = "xrt.discovery:1.0"
	MessageTypeEvent           = "xrt.event:1.0"
)

// Constant related to define the event type of notification
const (
	EventTypeDeviceAdded   = "device:added"
	EventTypeDeviceUpdated = "device:updated"
	EventTypeDeviceDeleted = "device:deleted"
)
