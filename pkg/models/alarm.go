// Copyright (C) 2026 IOTech Ltd

package models

// AlarmSetting is the general data to config alarm service
type AlarmSetting struct {
	Id   string
	Name string
}

type AlarmMetadataResponse struct {
	Count int
}

type AlarmMultiResponse struct {
	Templates  []AlarmSetting
	Conditions []AlarmSetting
	Actions    []AlarmSetting
	Metadata   AlarmMetadataResponse
}

type AlarmMultiAssociationResponse struct {
	Associations []any
	Metadata     AlarmMetadataResponse
}

// AlarmAssociation represents a single association definition from the provision JSON file
type AlarmAssociation struct {
	SourceType           string `json:"sourceType"`
	ConfigName           string `json:"configName"`
	DeviceName           string `json:"deviceName,omitempty"`
	ProfileName          string `json:"profileName,omitempty"`
	ResourceName         string `json:"resourceName,omitempty"`
	MessageBusSourceName string `json:"messageBusSourceName,omitempty"`
	SparkplugNodeId      string `json:"sparkplugNodeId,omitempty"`
	SparkplugDeviceName  string `json:"sparkplugDeviceName,omitempty"`
	SparkplugMetricName  string `json:"sparkplugMetricName,omitempty"`
}
