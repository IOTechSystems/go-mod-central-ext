// Copyright (C) 2024 IOTech Ltd

package requests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/edgexfoundry/go-mod-central-ext/v4/pkg/common"
	"github.com/edgexfoundry/go-mod-central-ext/v4/pkg/dtos"
	"github.com/edgexfoundry/go-mod-central-ext/v4/pkg/models"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"

	"github.com/stretchr/testify/require"
)

var (
	testUserId          = "431f0134-ae07-45ac-8577-1c1ac2b74fb3"
	testUsername        = "bob"
	testPassword        = "MyPassword@123"
	testUserDescription = "A test user"
	testUserDTO         = dtos.User{
		Id:          testUserId,
		Name:        testUsername,
		Password:    testPassword,
		Description: testUserDescription,
	}
	testAddUser = AddUserRequest{
		BaseRequest: dtoCommon.BaseRequest{
			RequestId:   common.ExampleUUID,
			Versionable: dtoCommon.NewVersionable(),
		},
		User: testUserDTO,
	}
	testUpdateUser = UpdateUserRequest{
		BaseRequest: dtoCommon.BaseRequest{
			RequestId:   common.ExampleUUID,
			Versionable: dtoCommon.NewVersionable(),
		},
		User: dtos.UpdateUser{
			Name:        &testUsername,
			Description: &testUserDescription,
		},
	}
	testLoginRequest = LoginRequest{
		Username: "bob",
		Password: "MySecret@123",
	}
)

func TestAddUserRequest_Validate(t *testing.T) {
	valid := testAddUser
	invalidUsername := testAddUser
	invalidUsername.User.Name = "_testInvalid"
	invalidUserPass := testAddUser
	invalidUserPass.User.Password = "invalidpass"
	invalidUserPassTooShort := testAddUser
	invalidUserPassTooShort.User.Password = "Inv3&"
	tests := []struct {
		name        string
		User        AddUserRequest
		expectError bool
	}{
		{"valid AddUserRequest", valid, false},
		{"invalid AddUserRequest, username not starts with letter or digit", invalidUsername, true},
		{"invalid AddUserRequest, password not contains digit and special character", invalidUserPass, true},
		{"invalid AddUserRequest, password is too short", invalidUserPassTooShort, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.User.Validate()
			if tt.expectError {
				require.Error(t, err, fmt.Sprintf("expect error but not : %s", tt.name))
			} else {
				require.NoError(t, err, fmt.Sprintf("unexpected error occurs : %s", tt.name))
			}
		})
	}
}

func TestAddUser_UnmarshalJSON(t *testing.T) {
	valid := testAddUser
	resultTestBytes, _ := json.Marshal(testAddUser)
	type args struct {
		data []byte
	}

	tests := []struct {
		name    string
		addUser AddUserRequest
		args    args
		wantErr bool
	}{
		{"unmarshal AddUserRequest with success", valid, args{resultTestBytes}, false},
		{"unmarshal invalid AddUserRequest, empty data", AddUserRequest{}, args{[]byte{}}, true},
		{"unmarshal invalid AddUserRequest, string data", AddUserRequest{}, args{[]byte("Invalid AddUserRequest")}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expected = tt.addUser
			err := tt.addUser.UnmarshalJSON(tt.args.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, expected, tt.addUser, "Unmarshal did not result in expected AddUserRequest")
			}
		})
	}
}

func Test_AddUserReqToUserModels(t *testing.T) {
	requests := []AddUserRequest{testAddUser}
	expectedUserModel := []models.User{
		{
			Id:          testUserId,
			Name:        testUsername,
			Password:    testPassword,
			Description: testUserDescription,
		},
	}
	resultModels := AddUserReqToUserModels(requests)
	require.Equal(t, expectedUserModel, resultModels, "AddUserReqToUserModels did not result in expected User model")
}

func TestUpdateUserRequest_Validate(t *testing.T) {
	valid := testUpdateUser
	invalidNoName := testUpdateUser
	emptyName := ""
	invalidNoName.User.Name = &emptyName
	invalidUserPass := testUpdateUser
	invalidPass := "invalidpass"
	invalidUserPass.User.Password = &invalidPass
	invalidPassTooShort := "Inv3&"
	invalidUserPassTooShort := testUpdateUser
	invalidUserPassTooShort.User.Password = &invalidPassTooShort

	tests := []struct {
		name        string
		User        UpdateUserRequest
		expectError bool
	}{
		{"valid UpdateUserRequest", valid, false},
		{"invalid UpdateUserRequest, username empty", invalidNoName, true},
		{"invalid UpdateUserRequest, password not contains digit and special character", invalidUserPass, true},
		{"invalid UpdateUserRequest, password is too short", invalidUserPassTooShort, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.User.Validate()
			if tt.expectError {
				require.Error(t, err, fmt.Sprintf("expect error but not : %s", tt.name))
			} else {
				require.NoError(t, err, fmt.Sprintf("unexpected error occurs : %s", tt.name))
			}
		})
	}
}

func TestUpdateUser_UnmarshalJSON(t *testing.T) {
	valid := testUpdateUser
	resultTestBytes, _ := json.Marshal(testUpdateUser)
	type args struct {
		data []byte
	}

	tests := []struct {
		name       string
		updateUser UpdateUserRequest
		args       args
		wantErr    bool
	}{
		{"unmarshal UpdateUserRequest with success", valid, args{resultTestBytes}, false},
		{"unmarshal invalid UpdateUserRequest, empty data", UpdateUserRequest{}, args{[]byte{}}, true},
		{"unmarshal invalid UpdateUserRequest, string data", UpdateUserRequest{}, args{[]byte("Invalid UpdateUserRequest")}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expected = tt.updateUser
			err := tt.updateUser.UnmarshalJSON(tt.args.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, expected, tt.updateUser, "Unmarshal did not result in expected UpdateUserRequest")
			}
		})
	}
}

func TestLoginRequest_Validate(t *testing.T) {
	valid := testLoginRequest
	invalidNoName := valid
	invalidNoName.Username = ""
	invalidNoPass := valid
	invalidNoPass.Password = ""

	tests := []struct {
		name        string
		login       LoginRequest
		expectError bool
	}{
		{"valid LoginRequest", valid, false},
		{"invalid LoginRequest, username empty", invalidNoName, true},
		{"invalid LoginRequest, password empty", invalidNoPass, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.login.Validate()
			if tt.expectError {
				require.Error(t, err, fmt.Sprintf("expect error but not : %s", tt.name))
			} else {
				require.NoError(t, err, fmt.Sprintf("unexpected error occurs : %s", tt.name))
			}
		})
	}
}

func TestLoginRequest_UnmarshalJSON(t *testing.T) {
	valid := testLoginRequest
	resultTestBytes, _ := json.Marshal(valid)
	type args struct {
		data []byte
	}

	tests := []struct {
		name         string
		loginRequest LoginRequest
		args         args
		wantErr      bool
	}{
		{"unmarshal LoginRequest with success", valid, args{resultTestBytes}, false},
		{"unmarshal invalid LoginRequest, empty data", LoginRequest{}, args{[]byte{}}, true},
		{"unmarshal invalid LoginRequest, string data", LoginRequest{}, args{[]byte("Invalid LoginRequest")}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expected = tt.loginRequest
			err := tt.loginRequest.UnmarshalJSON(tt.args.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, expected, tt.loginRequest, "Unmarshal did not result in expected LoginRequest")
			}
		})
	}
}
