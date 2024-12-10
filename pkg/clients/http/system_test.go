// Copyright (C) 2021-2024 IOTech Ltd

package http

import (
	"context"
	"encoding/json"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients/interfaces"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/requests"
	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
)

func newTestServer(httpMethod string, apiRoute string, expectedResponse interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethod {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.EscapedPath() != apiRoute {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		if expectedResponse != nil {
			b, _ := json.Marshal(expectedResponse)
			_, _ = w.Write(b)
		}
	}))
}

type emptyAuthenticationInjector struct {
}

// NewNullAuthenticationInjector creates an instance of AuthenticationInjector
func NewNullAuthenticationInjector() interfaces.AuthenticationInjector {
	return &emptyAuthenticationInjector{}
}

func (_ *emptyAuthenticationInjector) AddAuthenticationData(_ *http.Request) error {
	// Do nothing to the request; used for unit tests
	return nil
}

func (_ *emptyAuthenticationInjector) RoundTripper() http.RoundTripper {
	// Do nothing to the request; used for unit tests
	return nil
}

func TestSystemManagementClient_GetHealth(t *testing.T) {
	ts := newTestServer(http.MethodGet, edgexCommon.ApiHealthRoute, []dtoCommon.BaseWithServiceNameResponse{})
	defer ts.Close()

	client := NewSystemManagementClient(ts.URL, NewNullAuthenticationInjector())
	res, err := client.GetHealth(context.Background(), []string{"core-data"})
	require.NoError(t, err)
	require.IsType(t, []dtoCommon.BaseWithServiceNameResponse{}, res)
}

func TestSystemManagementClient_GetConfig(t *testing.T) {
	ts := newTestServer(http.MethodGet, edgexCommon.ApiMultiConfigRoute, []dtoCommon.BaseWithConfigResponse{})
	defer ts.Close()

	client := NewSystemManagementClient(ts.URL, NewNullAuthenticationInjector())
	res, err := client.GetConfig(context.Background(), []string{"core-data"})
	require.NoError(t, err)
	require.IsType(t, []dtoCommon.BaseWithConfigResponse{}, res)
}

func TestSystemManagementClient_DoOperation(t *testing.T) {
	ts := newTestServer(http.MethodPost, edgexCommon.ApiOperationRoute, []dtoCommon.BaseResponse{})
	defer ts.Close()

	client := NewSystemManagementClient(ts.URL, NewNullAuthenticationInjector())
	res, err := client.DoOperation(context.Background(), []requests.OperationRequest{})
	require.NoError(t, err)
	require.IsType(t, []dtoCommon.BaseResponse{}, res)
}
