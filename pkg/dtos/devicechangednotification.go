// Copyright (C) 2022-2023 IOTech Ltd

package dtos

import (
	"encoding/json"
	"time"

	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
)

type DeviceNotificationContent struct {
	DeviceName        string                                  `json:"deviceName"`
	DeviceServiceName string                                  `json:"deviceServiceName"`
	Protocols         map[string]edgexDtos.ProtocolProperties `json:"protocols"`
	ActionType        string                                  `json:"actionType"`
	OperatingState    string                                  `json:"operatingState"`
	AdminState        string                                  `json:"adminState"`
	Datetime          string                                  `json:"datetime"`
}

func NewDeviceNotificationContent(device edgexDtos.Device, action string) DeviceNotificationContent {
	return DeviceNotificationContent{
		DeviceName:        device.Name,
		DeviceServiceName: device.ServiceName,
		Protocols:         device.Protocols,
		ActionType:        action,
		OperatingState:    device.OperatingState,
		AdminState:        device.AdminState,
		Datetime:          time.Now().Format(time.RFC1123),
	}
}

func (d DeviceNotificationContent) String() (string, error) {
	if b, err := json.Marshal(d); err == nil {
		return string(b), nil
	} else {
		return "", err
	}
}
