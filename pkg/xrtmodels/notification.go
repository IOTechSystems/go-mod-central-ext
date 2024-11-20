// Copyright (C) 2024 IOTech Ltd

package xrtmodels

// Notification is used to send the status of a device changes
// https://docs.iotechsys.com/edge-xrt22/mqtt-management/mqtt-management.html#notification-format
type Notification struct {
	DeviceServiceName string `json:"device_service"`
	Event             any    `json:"event"`
	EventType         string `json:"event_type"`
	Type              string `json:"type"`
}
