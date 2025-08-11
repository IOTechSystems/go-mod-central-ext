// Copyright (C) 2025 IOTech Ltd

package models

import edgexModels "github.com/edgexfoundry/go-mod-core-contracts/v4/models"

// RetentionPolicy specifies the patterns for device and source name associated with the saved event, as well as the data retention duration for the policy
type RetentionPolicy struct {
	edgexModels.DBTimestamp
	Id         string
	DeviceName string
	SourceName string
	Duration   string
}
