// Copyright (C) 2024-2025 IOTech Ltd

package requests

import (
	"encoding/json"
	"fmt"

	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/models"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
	"github.com/stretchr/testify/require"
)

var (
	mockRole        = "mockAdmin"
	mockDescription = "A test admin role"
	mockPath        = "/core-command"
	mockMethods     = []string{"GET", "PUT"}
	mockEffect      = "allow"

	mockAccessPolicyDTO = dtos.AccessPolicy{
		Path:    mockPath,
		Methods: mockMethods,
		Effect:  mockEffect,
	}
	mockAccessPolicyModel = models.AccessPolicy{
		Path:    mockPath,
		Methods: mockMethods,
		Effect:  mockEffect,
	}
	mockRolePolicyDTO = dtos.RolePolicy{
		Role:           mockRole,
		AccessPolicies: []dtos.AccessPolicy{mockAccessPolicyDTO},
		RolePolicyBasicInfo: dtos.RolePolicyBasicInfo{
			Description: mockDescription,
		},
	}
	testAddRolePolicyRequest = AddRolePolicyRequest{
		BaseRequest: dtoCommon.BaseRequest{
			RequestId:   common.ExampleUUID,
			Versionable: dtoCommon.NewVersionable(),
		},
		RolePolicy: mockRolePolicyDTO,
	}
	mockAuthRoute = dtos.AuthRoute{
		Path:   "/core-metadata/.*/device/name/abc",
		Method: "DELETE",
	}
	testAuthRouteRequest = AuthRouteRequest{
		BaseRequest: dtoCommon.BaseRequest{
			RequestId:   common.ExampleUUID,
			Versionable: dtoCommon.NewVersionable(),
		},
		AuthRoute: mockAuthRoute,
	}
	mockAuthGraphQL = dtos.AuthGraphQL{
		Path:   "^/alarms-service/graphql/[^/]*Alarm$",
		Method: "MUTATION",
	}
	testAuthGraphQLRequest = AuthGraphQLRequest{
		BaseRequest: dtoCommon.BaseRequest{
			RequestId:   common.ExampleUUID,
			Versionable: dtoCommon.NewVersionable(),
		},
		AuthGraphQL: mockAuthGraphQL,
	}
)

func TestAddRolePolicyRequest_Validate(t *testing.T) {
	valid := testAddRolePolicyRequest

	emptyRole := valid
	emptyRole.RolePolicy.Role = ""

	emptyAccessPolicy := valid
	emptyAccessPolicy.RolePolicy.AccessPolicies = nil

	emptyPathAccessPolicy := mockAccessPolicyDTO
	emptyPathAccessPolicy.Path = ""
	invalidAccessPolicy := valid
	invalidAccessPolicy.RolePolicy.AccessPolicies = []dtos.AccessPolicy{emptyPathAccessPolicy}

	emptyMethodAccessPolicy := mockAccessPolicyDTO
	emptyMethodAccessPolicy.Methods = nil
	emptyMethodPolicy := valid
	emptyMethodPolicy.RolePolicy.AccessPolicies = []dtos.AccessPolicy{emptyMethodAccessPolicy}

	invalidMethodAccessPolicy := mockAccessPolicyDTO
	invalidMethodAccessPolicy.Methods = []string{"invalid"}
	invalidMethodPolicy := valid
	invalidMethodPolicy.RolePolicy.AccessPolicies = []dtos.AccessPolicy{invalidMethodAccessPolicy}

	dupMethodAccessPolicy := mockAccessPolicyDTO
	dupMethodAccessPolicy.Methods = []string{"GET", "GET"}
	dupMethodPolicy := valid
	dupMethodPolicy.RolePolicy.AccessPolicies = []dtos.AccessPolicy{dupMethodAccessPolicy}
	tests := []struct {
		name        string
		RolePolicy  AddRolePolicyRequest
		expectError bool
	}{
		{"valid AddRolePolicyRequest", valid, false},
		{"invalid AddRolePolicyRequest, empty role", emptyRole, true},
		{"invalid AddRolePolicyRequest, empty AccessPolicy", emptyAccessPolicy, true},
		{"invalid AddRolePolicyRequest, invalid AccessPolicy", invalidAccessPolicy, true},
		{"invalid AddRolePolicyRequest, empty HttpMethod", emptyMethodPolicy, true},
		{"invalid AddRolePolicyRequest, invalid HttpMethod", invalidMethodPolicy, true},
		{"invalid AddRolePolicyRequest, duplicate HttpMethod", dupMethodPolicy, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.RolePolicy.Validate()
			if tt.expectError {
				require.Error(t, err, fmt.Sprintf("expect error but not : %s", tt.name))
			} else {
				require.NoError(t, err, fmt.Sprintf("unexpected error occurs : %s", tt.name))
			}
		})
	}
}

func TestAddRolePolicy_UnmarshalJSON(t *testing.T) {
	valid := testAddRolePolicyRequest
	resultTestBytes, _ := json.Marshal(testAddRolePolicyRequest)
	type args struct {
		data []byte
	}

	tests := []struct {
		name          string
		addRolePolicy AddRolePolicyRequest
		args          args
		wantErr       bool
	}{
		{"unmarshal AddRolePolicyRequest with success", valid, args{resultTestBytes}, false},
		{"unmarshal invalid AddRolePolicyRequest, empty data", AddRolePolicyRequest{}, args{[]byte{}}, true},
		{"unmarshal invalid AddRolePolicyRequest, string data", AddRolePolicyRequest{}, args{[]byte("Invalid AddUserRequest")}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expected = tt.addRolePolicy
			err := tt.addRolePolicy.UnmarshalJSON(tt.args.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, expected, tt.addRolePolicy, "Unmarshal did not result in expected AddRolePolicyRequest")
			}
		})
	}
}

func Test_AddRolePolicyReqToRolePolicyModels(t *testing.T) {
	requests := []AddRolePolicyRequest{testAddRolePolicyRequest}
	expectedRoleModels := []models.RolePolicy{
		{
			Role:           mockRole,
			AccessPolicies: []models.AccessPolicy{mockAccessPolicyModel},
			RolePolicyBasicInfo: models.RolePolicyBasicInfo{
				Description: mockDescription,
			},
		},
	}
	resultModels := AddRolePolicyReqToRolePolicyModels(requests)
	require.Equal(t, expectedRoleModels, resultModels, "AddRolePolicyReqToRolePolicyModels did not result in expected RolePolicy model")
}

func TestAuthRouteRequest_Validate(t *testing.T) {
	valid := testAuthRouteRequest

	emptyPath := valid
	emptyPath.AuthRoute.Path = ""

	emptyMethod := valid
	emptyMethod.AuthRoute.Method = ""

	invalidMethod := valid
	invalidMethod.AuthRoute.Method = "invalid"
	tests := []struct {
		name         string
		authRouteReq AuthRouteRequest
		expectError  bool
	}{
		{"valid AuthRouteRequest", valid, false},
		{"invalid AuthRouteRequest, empty path", emptyPath, true},
		{"invalid AuthRouteRequest, empty HttpMethod", emptyMethod, true},
		{"invalid AuthRouteRequest, invalid HttpMethod", invalidMethod, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.authRouteReq.Validate()
			if tt.expectError {
				require.Error(t, err, fmt.Sprintf("expect error but not : %s", tt.name))
			} else {
				require.NoError(t, err, fmt.Sprintf("unexpected error occurs : %s", tt.name))
			}
		})
	}
}

func TestAuthRouteRequest_UnmarshalJSON(t *testing.T) {
	valid := testAuthRouteRequest
	resultTestBytes, _ := json.Marshal(testAuthRouteRequest)
	type args struct {
		data []byte
	}

	tests := []struct {
		name         string
		authRouteReq AuthRouteRequest
		args         args
		wantErr      bool
	}{
		{"unmarshal AddRolePolicyRequest with success", valid, args{resultTestBytes}, false},
		{"unmarshal invalid AddRolePolicyRequest, empty data", AuthRouteRequest{}, args{[]byte{}}, true},
		{"unmarshal invalid AddRolePolicyRequest, string data", AuthRouteRequest{}, args{[]byte("Invalid AddUserRequest")}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expected = tt.authRouteReq
			err := tt.authRouteReq.UnmarshalJSON(tt.args.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, expected, tt.authRouteReq, "Unmarshal did not result in expected AuthRouteRequest")
			}
		})
	}
}

func TestAuthGraphQLRequest_Validate(t *testing.T) {
	valid := testAuthGraphQLRequest

	emptyPath := valid
	emptyPath.Path = ""

	emptyMethod := valid
	emptyMethod.Method = ""

	invalidMethod := valid
	invalidMethod.Method = "invalid"
	tests := []struct {
		name           string
		authGraphQLReq AuthGraphQLRequest
		expectError    bool
	}{
		{"valid AuthGraphQLRequest", valid, false},
		{"invalid AuthGraphQLRequest, empty path", emptyPath, true},
		{"invalid AuthGraphQLRequest, empty Method", emptyMethod, true},
		{"invalid AuthGraphQLRequest, invalid Method", invalidMethod, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.authGraphQLReq.Validate()
			if tt.expectError {
				require.Error(t, err, fmt.Sprintf("expect error but not : %s", tt.name))
			} else {
				require.NoError(t, err, fmt.Sprintf("unexpected error occurs : %s", tt.name))
			}
		})
	}
}

func TestAuthGraphQLRequest_UnmarshalJSON(t *testing.T) {
	valid := testAuthGraphQLRequest
	resultTestBytes, _ := json.Marshal(testAuthGraphQLRequest)
	type args struct {
		data []byte
	}

	tests := []struct {
		name           string
		authGraphQLReq AuthGraphQLRequest
		args           args
		wantErr        bool
	}{
		{"unmarshal AuthGraphQLRequest with success", valid, args{resultTestBytes}, false},
		{"unmarshal invalid AuthGraphQLRequest, empty data", AuthGraphQLRequest{}, args{[]byte{}}, true},
		{"unmarshal invalid AuthGraphQLRequest, string data", AuthGraphQLRequest{}, args{[]byte("Invalid AuthGraphQLRequest")}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expected = tt.authGraphQLReq
			err := tt.authGraphQLReq.UnmarshalJSON(tt.args.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, expected, tt.authGraphQLReq, "Unmarshal did not result in expected AuthGraphQLRequest")
			}
		})
	}
}
