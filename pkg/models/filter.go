// Copyright (C) 2025 IOTech Ltd

package models

// Filter defining the device, event source, and resource name patterns to filterIn or filterOut events
type Filter struct {
	Id                string
	Type              string
	DeviceName        string
	EventSourceName   string
	ResourceName      string
	OnChange          bool
	OnChangeThreshold float64
}
