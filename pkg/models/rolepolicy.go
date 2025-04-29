// Copyright (C) 2024-2025 IOTech Ltd

package models

import (
	edgexModels "github.com/edgexfoundry/go-mod-core-contracts/v4/models"
)

type RolePolicy struct {
	Id             string
	Role           string
	AccessPolicies []AccessPolicy
	RolePolicyBasicInfo
}

// RolePolicyBasicInfo includes additional fields that provide descriptive information about a role policy. When stored in the Casbin database tables, an instance of RolePolicyBasicInfo is encoded as a base64 string.
type RolePolicyBasicInfo struct {
	edgexModels.DBTimestamp
	Description string
}

type AccessPolicy struct {
	Path    string
	Methods []string
	Effect  string
}
