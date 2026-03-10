// Copyright (C) 2026 IOTech Ltd

package http

import (
	"context"
	"net/http"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/responses"
	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
)

func TestQueryTimeSeriesByDeviceNameAndResourceNameAndTimeRange(t *testing.T) {
	deviceName := "device"
	resourceName := "resource"
	start := int64(1)
	end := int64(10)
	urlPath := path.Join(common.ApiTimeSeriesRoute, edgexCommon.Device, edgexCommon.Name, deviceName, edgexCommon.ResourceName, resourceName, edgexCommon.Start, strconv.FormatInt(start, 10), edgexCommon.End, strconv.FormatInt(end, 10))
	ts := newTestServer(http.MethodGet, urlPath, responses.TimeSeriesResponse{})
	defer ts.Close()

	client := NewTimeSeriesClient(ts.URL, NewNullAuthenticationInjector(), false)

	res, err := client.TimeSeriesByDeviceNameAndResourceNameAndTimeRange(context.Background(), deviceName, resourceName, start, end)
	require.NoError(t, err)
	assert.IsType(t, responses.TimeSeriesResponse{}, res)
}

func TestQueryTimeSeriesByDeviceNameAndMultiResourceNamesAndTimeRange(t *testing.T) {
	deviceName := "device"
	resourceNames := []string{"resource01", "resource02"}
	start := int64(1)
	end := int64(10)
	urlPath := path.Join(common.ApiTimeSeriesRoute, edgexCommon.Device, edgexCommon.Name, deviceName, edgexCommon.Start, strconv.FormatInt(start, 10), edgexCommon.End, strconv.FormatInt(end, 10))
	ts := newTestServer(http.MethodGet, urlPath, responses.TimeSeriesResponse{})
	defer ts.Close()

	client := NewTimeSeriesClient(ts.URL, NewNullAuthenticationInjector(), false)
	res, err := client.TimeSeriesByDeviceNameAndMultiResourceNamesAndTimeRange(context.Background(), deviceName, resourceNames, start, end)
	require.NoError(t, err)
	assert.IsType(t, responses.TimeSeriesResponse{}, res)
}
