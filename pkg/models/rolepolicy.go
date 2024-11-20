// Copyright (C) 2024 IOTech Ltd

package models

import (
	edgexModels "github.com/edgexfoundry/go-mod-core-contracts/v4/models"
)

type RolePolicy struct {
	edgexModels.DBTimestamp
	Id             string
	Role           string
	Description    string
	AccessPolicies []AccessPolicy
}

type AccessPolicy struct {
	Path        string
	HttpMethods []string
	Effect      string
}
