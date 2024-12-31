// Copyright (C) 2024 IOTech Ltd

package dtos

import (
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/models"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
)

type User struct {
	edgexDtos.DBTimestamp
	Id          string   `json:"id"`
	Name        string   `json:"name" validate:"required,edgex-dto-none-empty-string,edgex-dto-username"`
	DisplayName string   `json:"displayName"`
	Password    string   `json:"password,omitempty" validate:"required,edgex-dto-none-empty-string,edgex-dto-password"`
	Description string   `json:"description"`
	Roles       []string `json:"roles,omitempty" validate:"unique"`
}

type UpdateUser struct {
	edgexDtos.DBTimestamp
	Id          *string  `json:"id" validate:"required_without=Name,edgex-dto-uuid"`
	Name        *string  `json:"name" validate:"required_without=Id,edgex-dto-none-empty-string"`
	DisplayName *string  `json:"displayName"`
	Password    *string  `json:"password" validate:"omitempty,edgex-dto-password"`
	Description *string  `json:"description" validate:"omitempty"`
	Roles       []string `json:"roles" validate:"unique"`
}

// ToUserModel transforms the User DTO to the User Model
func ToUserModel(user User) models.User {
	return models.User{
		Id:          user.Id,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Password:    user.Password,
		Description: user.Description,
		Roles:       user.Roles,
	}
}

// FromUserModelToDTO transforms the User Model to the User DTO
func FromUserModelToDTO(d models.User) User {
	return User{
		DBTimestamp: edgexDtos.DBTimestamp(d.DBTimestamp),
		Id:          d.Id,
		Name:        d.Name,
		DisplayName: d.DisplayName,
		Description: d.Description,
		Roles:       d.Roles,
	}
}

// UpdateUserReqToUserModel transforms the UpdateUserRequest DTO to the User model
func UpdateUserReqToUserModel(userModel *models.User, updateUser UpdateUser) {
	if updateUser.Password != nil {
		userModel.Password = *updateUser.Password
	}

	if updateUser.DisplayName != nil {
		userModel.DisplayName = *updateUser.DisplayName
	}

	if updateUser.Description != nil {
		userModel.Description = *updateUser.Description
	}

	if updateUser.Roles != nil {
		userModel.Roles = updateUser.Roles
	}
}
