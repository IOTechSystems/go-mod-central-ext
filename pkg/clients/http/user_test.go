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

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"

	"github.com/stretchr/testify/require"
)

const mockUserName = "testUser"

func TestAddUser(t *testing.T) {
	ts := newTestServer(http.MethodPost, common.ApiUserRoute, []dtoCommon.BaseWithIdResponse{})
	defer ts.Close()

	client := NewUserClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.Add(context.Background(), []requests.AddUserRequest{})
	require.NoError(t, err)
	require.IsType(t, []dtoCommon.BaseWithIdResponse{}, res)
}

func TestUpdateUser(t *testing.T) {
	ts := newTestServer(http.MethodPatch, common.ApiUserRoute, []dtoCommon.BaseResponse{})
	defer ts.Close()

	mockName := "test"
	client := NewUserClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.Update(context.Background(), []requests.UpdateUserRequest{{User: dtos.UpdateUser{Name: &mockName}}})
	require.NoError(t, err)
	require.IsType(t, []dtoCommon.BaseResponse{}, res)
}

func TestQueryAllUsers(t *testing.T) {
	ts := newTestServer(http.MethodGet, common.ApiAllUserRoute, responses.MultiUsersResponse{})
	defer ts.Close()

	client := NewUserClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.AllUsers(context.Background(), 1, 10)
	require.NoError(t, err)
	require.IsType(t, responses.MultiUsersResponse{}, res)
}

func TestQueryUserByName(t *testing.T) {
	urlPath := path.Join(common.ApiUserRoute, edgexCommon.Name, mockUserName)
	ts := newTestServer(http.MethodGet, urlPath, responses.UserResponse{})
	defer ts.Close()

	client := NewUserClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.UserByName(context.Background(), mockUserName)
	require.NoError(t, err)
	require.IsType(t, responses.UserResponse{}, res)
}

func TestDeleteUserByName(t *testing.T) {
	urlPath := path.Join(common.ApiUserRoute, edgexCommon.Name, mockUserName)
	ts := newTestServer(http.MethodDelete, urlPath, dtoCommon.BaseResponse{})
	defer ts.Close()

	client := NewUserClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.DeleteUserByName(context.Background(), mockUserName)
	require.NoError(t, err)
	require.IsType(t, dtoCommon.BaseResponse{}, res)
}

func TestAuthenticate(t *testing.T) {
	ts := newTestServer(http.MethodPost, common.ApiAuthRoute, nil)
	defer ts.Close()

	client := NewUserClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.Authenticate(context.Background(), map[string]string{"mock": "mockHeader"})
	require.NoError(t, err)
	require.IsType(t, nil, res)
}

func TestAuthRoutes(t *testing.T) {
	ts := newTestServer(http.MethodPost, common.ApiAuthRoutesRoute, responses.AuthRouteResponse{})
	defer ts.Close()

	client := NewUserClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.AuthRoutes(context.Background(), map[string]string{"mock": "mockHeader"}, []requests.AuthRouteRequest{})
	require.NoError(t, err)
	require.IsType(t, responses.AuthRouteResponse{}, res)
}
