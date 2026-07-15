// Copyright (C) 2026 IOTech Ltd

package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path"
	"strconv"
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/requests"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/responses"

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"

	"github.com/stretchr/testify/require"
)

const mockFilterId = "0a8555c7-8d70-4266-8db2-2f8fb8bd0ce8"

// newPagingTestServer behaves like newTestServer but also asserts the offset and
// limit query params match the expected values, returning 400 on mismatch so a
// wrong param key/value surfaces as a request error in the test.
func newPagingTestServer(httpMethod, apiRoute string, offset, limit int, expectedResponse interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if r.Method != httpMethod || r.URL.EscapedPath() != apiRoute ||
			q.Get(edgexCommon.Offset) != strconv.Itoa(offset) || q.Get(edgexCommon.Limit) != strconv.Itoa(limit) {
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

func TestQueryAllFilters(t *testing.T) {
	ts := newPagingTestServer(http.MethodGet, common.ApiAllFilterRoute, 0, 10, responses.MultiFiltersResponse{})
	defer ts.Close()

	client := NewFilterClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.AllFilters(context.Background(), 0, 10)
	require.NoError(t, err)
	require.IsType(t, responses.MultiFiltersResponse{}, res)
}

func TestQueryFilterById(t *testing.T) {
	urlPath := path.Join(common.ApiFilterRoute, edgexCommon.Id, mockFilterId)
	ts := newTestServer(http.MethodGet, urlPath, responses.FilterResponse{})
	defer ts.Close()

	client := NewFilterClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.FilterById(context.Background(), mockFilterId)
	require.NoError(t, err)
	require.IsType(t, responses.FilterResponse{}, res)
}

func TestQueryFiltersByDeviceName(t *testing.T) {
	deviceName := "device"
	urlPath := path.Join(common.ApiFilterRoute, edgexCommon.DeviceName, deviceName)
	ts := newPagingTestServer(http.MethodGet, urlPath, 0, 10, responses.MultiFiltersResponse{})
	defer ts.Close()

	client := NewFilterClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.FiltersByDeviceName(context.Background(), deviceName, 0, 10)
	require.NoError(t, err)
	require.IsType(t, responses.MultiFiltersResponse{}, res)
}

func TestAddFilter(t *testing.T) {
	ts := newTestServer(http.MethodPost, common.ApiFilterRoute, []dtoCommon.BaseWithIdResponse{})
	defer ts.Close()

	client := NewFilterClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.Add(context.Background(), []requests.FilterRequest{})
	require.NoError(t, err)
	require.IsType(t, []dtoCommon.BaseWithIdResponse{}, res)
}

func TestPutFilters(t *testing.T) {
	ts := newTestServer(http.MethodPut, common.ApiAllFilterRoute, []dtoCommon.BaseWithIdResponse{})
	defer ts.Close()

	client := NewFilterClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.PutFilters(context.Background(), []requests.FilterRequest{})
	require.NoError(t, err)
	require.IsType(t, []dtoCommon.BaseWithIdResponse{}, res)
}

func TestDeleteFilters(t *testing.T) {
	ts := newTestServer(http.MethodDelete, common.ApiAllFilterRoute, dtoCommon.BaseResponse{})
	defer ts.Close()

	client := NewFilterClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.DeleteFilters(context.Background())
	require.NoError(t, err)
	require.IsType(t, dtoCommon.BaseResponse{}, res)
}

func TestUpdateFilter(t *testing.T) {
	urlPath := path.Join(common.ApiFilterRoute, edgexCommon.Id, mockFilterId)
	ts := newTestServer(http.MethodPut, urlPath, dtoCommon.BaseResponse{})
	defer ts.Close()

	client := NewFilterClient(ts.URL, NewNullAuthenticationInjector(), false)
	req := requests.FilterRequest{}
	req.Filter.Id = mockFilterId
	res, err := client.Update(context.Background(), req)
	require.NoError(t, err)
	require.IsType(t, dtoCommon.BaseResponse{}, res)
}

func TestDeleteFilterById(t *testing.T) {
	urlPath := path.Join(common.ApiFilterRoute, edgexCommon.Id, mockFilterId)
	ts := newTestServer(http.MethodDelete, urlPath, dtoCommon.BaseResponse{})
	defer ts.Close()

	client := NewFilterClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.DeleteFilterById(context.Background(), mockFilterId)
	require.NoError(t, err)
	require.IsType(t, dtoCommon.BaseResponse{}, res)
}

func TestDeleteFiltersByDeviceName(t *testing.T) {
	deviceName := "device"
	urlPath := path.Join(common.ApiFilterRoute, edgexCommon.DeviceName, deviceName)
	ts := newTestServer(http.MethodDelete, urlPath, dtoCommon.BaseResponse{})
	defer ts.Close()

	client := NewFilterClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.DeleteFiltersByDeviceName(context.Background(), deviceName)
	require.NoError(t, err)
	require.IsType(t, dtoCommon.BaseResponse{}, res)
}
