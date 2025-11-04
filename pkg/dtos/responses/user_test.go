// Copyright (C) 2024-2025 IOTech Ltd

package responses

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"
)

var (
	testUsername  = "bob"
	testUsername2 = "alice"
)

func TestNewUserResponse(t *testing.T) {
	expectedUser := dtos.User{Name: testUsername}
	actual := NewUserResponse(common.ExpectedRequestId, common.ExpectedMessage, common.ExpectedStatusCode, expectedUser)

	require.Equal(t, common.ExpectedRequestId, actual.RequestId)
	require.Equal(t, common.ExpectedStatusCode, actual.StatusCode)
	require.Equal(t, common.ExpectedMessage, actual.Message)
	require.Equal(t, expectedUser, actual.User)
}

func TestNewMultiUsersResponse(t *testing.T) {
	expectedUsers := []dtos.User{
		{Name: testUsername},
		{Name: testUsername2},
	}
	expectedTotalCount := int64(2)
	actual := NewMultiUsersResponse(common.ExpectedRequestId, common.ExpectedMessage, common.ExpectedStatusCode, expectedTotalCount, expectedUsers)

	require.Equal(t, common.ExpectedRequestId, actual.RequestId)
	require.Equal(t, common.ExpectedStatusCode, actual.StatusCode)
	require.Equal(t, common.ExpectedMessage, actual.Message)
	require.Equal(t, expectedTotalCount, actual.TotalCount)
	require.Equal(t, expectedUsers, actual.Users)
}
