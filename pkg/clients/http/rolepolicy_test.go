// Copyright (C) 2024 IOTech Ltd

package http

import (
	"context"
	"net/http"
	"path"
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/requests"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/responses"

	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"

	"github.com/stretchr/testify/require"
)

const mockRoleName = "testRole"

func TestAddRolePolicy(t *testing.T) {
	ts := newTestServer(http.MethodPost, common.ApiRolePolicyRoute, dtoCommon.BaseWithIdResponse{})
	defer ts.Close()

	client := NewRolePolicyClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.Add(context.Background(), requests.AddRolePolicyRequest{})
	require.NoError(t, err)
	require.IsType(t, dtoCommon.BaseWithIdResponse{}, res)
}

func TestUpdateRolePolicy(t *testing.T) {
	urlPath := path.Join(common.ApiRolePolicyRoute, common.Role, mockRoleName)
	ts := newTestServer(http.MethodPut, urlPath, dtoCommon.BaseResponse{})
	defer ts.Close()

	client := NewRolePolicyClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.Update(context.Background(), requests.AddRolePolicyRequest{
		RolePolicy: dtos.RolePolicy{Role: mockRoleName},
	})
	require.NoError(t, err)
	require.IsType(t, dtoCommon.BaseResponse{}, res)
}

func TestQueryAllRolePolicies(t *testing.T) {
	ts := newTestServer(http.MethodGet, common.ApiAllRolePolicyRoute, responses.MultiRolePolicyResponse{})
	defer ts.Close()

	client := NewRolePolicyClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.AllRolePolicies(context.Background(), 1, 10)
	require.NoError(t, err)
	require.IsType(t, responses.MultiRolePolicyResponse{}, res)
}

func TestQueryRolePolicyByName(t *testing.T) {
	path := path.Join(common.ApiRolePolicyRoute, common.Role, mockRoleName)
	ts := newTestServer(http.MethodGet, path, responses.RolePolicyResponse{})
	defer ts.Close()

	client := NewRolePolicyClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.RolePolicyByRole(context.Background(), mockRoleName)
	require.NoError(t, err)
	require.IsType(t, responses.RolePolicyResponse{}, res)
}

func TestDeleteRolePolicyByName(t *testing.T) {
	urlPath := path.Join(common.ApiRolePolicyRoute, common.Role, mockRoleName)
	ts := newTestServer(http.MethodDelete, urlPath, dtoCommon.BaseResponse{})
	defer ts.Close()

	client := NewRolePolicyClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.DeleteRolePolicyByRole(context.Background(), mockRoleName)
	require.NoError(t, err)
	require.IsType(t, dtoCommon.BaseResponse{}, res)
}
