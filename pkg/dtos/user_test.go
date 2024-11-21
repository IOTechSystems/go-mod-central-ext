// Copyright (C) 2024 IOTech Ltd

package dtos

import (
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/models"

	"github.com/stretchr/testify/require"
)

var (
	testUserId          = "431f0134-ae07-45ac-8577-1c1ac2b74fb3"
	testUsername        = "bob"
	testPassword        = "MyPassword@123"
	testUserDescription = "A test user"
	testUserDTO         = User{
		Id:          testUserId,
		Name:        testUsername,
		Password:    testPassword,
		Description: testUserDescription,
	}
	testUserModel = models.User{
		Id:          testUserId,
		Name:        testUsername,
		Password:    testPassword,
		Description: testUserDescription,
	}
	createdTimestamp  = int64(1730430964)
	modifiedTimestamp = int64(1730431003)
)

func TestToUserModel(t *testing.T) {
	actualUserModel := ToUserModel(testUserDTO)

	require.Equal(t, testUserModel, actualUserModel)
}

func TestFromUserModelToDTO(t *testing.T) {
	userModel := testUserModel
	userModel.Created = createdTimestamp
	userModel.Modified = modifiedTimestamp
	expectedUserDTO := testUserDTO
	expectedUserDTO.Created = createdTimestamp
	expectedUserDTO.Modified = modifiedTimestamp
	expectedUserDTO.Password = ""

	actualUserDTO := FromUserModelToDTO(userModel)
	require.Equal(t, expectedUserDTO, actualUserDTO)
}

func TestUpdateUserReqToUserModel(t *testing.T) {
	updatePassword := "MyUpdatePassword@123"
	updateUser := UpdateUser{
		Id:       &testUserId,
		Name:     &testUsername,
		Password: &updatePassword,
	}
	expectedUserModel := testUserModel
	expectedUserModel.Password = updatePassword

	UpdateUserReqToUserModel(&testUserModel, updateUser)
	require.Equal(t, expectedUserModel, testUserModel)
}
