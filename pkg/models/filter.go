// Copyright (C) 2025 IOTech Ltd

package models

import edgexModels "github.com/edgexfoundry/go-mod-core-contracts/v4/models"

// Filter defining the device, event source, and resource name patterns to filterIn or filterOut events
type Filter struct {
	edgexModels.DBTimestamp
	Id                string
	Type              string
	DeviceName        string
	EventSourceName   string
	ResourceName      string
	OnChange          bool
	OnChangeThreshold float64
}
