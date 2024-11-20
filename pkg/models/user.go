// Copyright (C) 2022-2024 IOTech Ltd

package models

import (
	edgexModels "github.com/edgexfoundry/go-mod-core-contracts/v4/models"
)

type User struct {
	edgexModels.DBTimestamp
	Id          string
	Name        string
	DisplayName string
	Password    string
	Description string
	Roles       []string
}
