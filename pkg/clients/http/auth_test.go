// Copyright (C) 2024 IOTech Ltd

package http

import (
	"context"
	"net/http"
	"path"
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/requests"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/responses"

	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	ts := newTestServer(http.MethodPost, common.ApiAuthRoute, "abc")
	defer ts.Close()

	client := NewAuthClient(ts.URL, NewNullAuthenticationInjector())
	err, jwt := client.Auth(context.Background(), map[string]string{"mock": "mockHeader"})
	require.NoError(t, err)
	require.IsType(t, "abc", jwt)
}

func TestAuthRoutes(t *testing.T) {
	ts := newTestServer(http.MethodPost, common.ApiAuthRoutesRoute, responses.AuthRouteResponse{})
	defer ts.Close()

	client := NewAuthClient(ts.URL, NewNullAuthenticationInjector())
	res, err := client.AuthRoutes(context.Background(), map[string]string{"mock": "mockHeader"}, []requests.AuthRouteRequest{})
	require.NoError(t, err)
	require.IsType(t, responses.AuthRouteResponse{}, res)
}

func TestVerificationKeyByIssuer(t *testing.T) {
	mockIssuer := "mockIssuer"
	urlPath := path.Join(common.ApiKeyRoute, common.VerificationKeyType, common.Issuer, mockIssuer)
	ts := newTestServer(http.MethodGet, urlPath, responses.KeyDataResponse{})
	defer ts.Close()

	client := NewAuthClient(ts.URL, NewNullAuthenticationInjector())
	res, err := client.VerificationKeyByIssuer(context.Background(), mockIssuer)
	require.NoError(t, err)
	require.IsType(t, responses.KeyDataResponse{}, res)
}
