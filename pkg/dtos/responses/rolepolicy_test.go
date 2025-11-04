// Copyright (C) 2024-2025 IOTech Ltd

package responses

import (
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"

	"github.com/stretchr/testify/require"
)

var (
	mockRole1 = "device-admin"
	mockRole2 = "cmd-admin"
)

func TestNewRolePolicyResponse(t *testing.T) {
	expectedRolePolicy := dtos.RolePolicy{Role: mockRole1}
	actual := NewRolePolicyResponse(common.ExpectedRequestId, common.ExpectedMessage, common.ExpectedStatusCode, expectedRolePolicy)
	require.Equal(t, common.ExpectedRequestId, actual.RequestId)
	require.Equal(t, common.ExpectedStatusCode, actual.StatusCode)
	require.Equal(t, common.ExpectedMessage, actual.Message)
	require.Equal(t, expectedRolePolicy, actual.RolePolicy)
}

func TestNewMultiRolePolicyResponse(t *testing.T) {
	expectedRolePolicies := []dtos.RolePolicy{
		{Role: mockRole1},
		{Role: mockRole2},
	}
	expectedTotalCount := int64(2)
	actual := NewMultiRolePolicyResponse(common.ExpectedRequestId, common.ExpectedMessage, common.ExpectedStatusCode, expectedTotalCount, expectedRolePolicies)

	require.Equal(t, common.ExpectedRequestId, actual.RequestId)
	require.Equal(t, common.ExpectedStatusCode, actual.StatusCode)
	require.Equal(t, common.ExpectedMessage, actual.Message)
	require.Equal(t, expectedTotalCount, actual.TotalCount)
	require.Equal(t, expectedRolePolicies, actual.RolePolicies)
}

func TestNewAuthRouteResponse(t *testing.T) {
	mockAuthRouteResult1 := dtos.AuthRouteResult{
		AuthRoute: dtos.AuthRoute{
			Path:   "/core-metadata/.*/device/name/abc",
			Method: "DELETE",
		},
		AuthResult: false,
	}
	mockAuthRouteResult2 := dtos.AuthRouteResult{
		AuthRoute: dtos.AuthRoute{
			Path:   "/core-data/.*/event",
			Method: "GET",
		},
		AuthResult: true,
	}
	expectedResults := []dtos.AuthRouteResult{
		mockAuthRouteResult1,
		mockAuthRouteResult2,
	}

	actual := NewAuthRouteResponse(common.ExpectedRequestId, common.ExpectedMessage, common.ExpectedStatusCode, expectedResults)

	require.Equal(t, common.ExpectedRequestId, actual.RequestId)
	require.Equal(t, common.ExpectedStatusCode, actual.StatusCode)
	require.Equal(t, common.ExpectedMessage, actual.Message)
	require.Equal(t, expectedResults, actual.AuthResponses)
}
