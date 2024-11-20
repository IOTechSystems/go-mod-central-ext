// Copyright (C) 2022-2024 IOTech Ltd

package dtos

import "encoding/json"

type DisconnectionNotificationContent struct {
	Servers  []string `json:"servers"`
	ClientId string   `json:"clientId"`
	Message  string   `json:"message"`
	DateTime string   `json:"dateTime"`
}

// NewDisconnectionNotificationContent creates a new DisconnectionNotificationContent object for the specified data
func NewDisconnectionNotificationContent(servers []string, clientId, message, dateTime string) DisconnectionNotificationContent {
	return DisconnectionNotificationContent{
		Servers:  servers,
		ClientId: clientId,
		Message:  message,
		DateTime: dateTime,
	}
}

// JsonString returns a JSON encoded string representation of the DisconnectionNotificationContent object
func (d DisconnectionNotificationContent) JsonString() (string, error) {
	if b, err := json.Marshal(d); err == nil {
		return string(b), nil
	} else {
		return "", err
	}
}
