// Copyright (C) 2023-2024 IOTech Ltd

package xrtmodels

type DeviceStatus struct {
	Device      string `json:"device"`
	Operational bool   `json:"operational"`
	Type        string `json:"type"`
}
