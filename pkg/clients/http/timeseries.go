// Copyright (C) 2026 IOTech Ltd

package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients/http/utils"
	clientsInterfaces "github.com/edgexfoundry/go-mod-core-contracts/v4/clients/interfaces"
	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/clients/interfaces"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/responses"
)

type TimeSeriesClient struct {
	baseUrlFunc           clients.ClientBaseUrlFunc
	authInjector          clientsInterfaces.AuthenticationInjector
	enableNameFieldEscape bool
}

// NewTimeSeriesClient creates an instance of TimeSeriesClient
func NewTimeSeriesClient(baseUrl string, authInjector clientsInterfaces.AuthenticationInjector, enableNameFieldEscape bool) interfaces.TimeSeriesClient {
	return &TimeSeriesClient{
		baseUrlFunc:           clients.GetDefaultClientBaseUrlFunc(baseUrl),
		authInjector:          authInjector,
		enableNameFieldEscape: enableNameFieldEscape,
	}
}

func (tsc *TimeSeriesClient) TimeSeriesByDeviceNameAndResourceNameAndTimeRange(ctx context.Context, deviceName, resourceName string, start, end int64) (responses.TimeSeriesResponse, errors.EdgeX) {
	requestPath := edgexCommon.NewPathBuilder().EnableNameFieldEscape(tsc.enableNameFieldEscape).
		SetPath(common.ApiTimeSeriesRoute).SetPath(edgexCommon.Device).SetPath(edgexCommon.Name).SetNameFieldPath(deviceName).SetPath(edgexCommon.ResourceName).SetNameFieldPath(resourceName).
		SetPath(edgexCommon.Start).SetPath(strconv.FormatInt(start, 10)).SetPath(edgexCommon.End).SetPath(strconv.FormatInt(end, 10)).BuildPath()
	res := responses.TimeSeriesResponse{}
	baseUrl, err := clients.GetBaseUrl(tsc.baseUrlFunc)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	requestParams := url.Values{}
	err = utils.GetRequest(ctx, &res, baseUrl, requestPath, requestParams, tsc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (tsc *TimeSeriesClient) TimeSeriesByDeviceNameAndMultiResourceNamesAndTimeRange(ctx context.Context, deviceName string, resourceNames []string, start, end int64) (responses.TimeSeriesResponse, errors.EdgeX) {
	requestPath := edgexCommon.NewPathBuilder().EnableNameFieldEscape(tsc.enableNameFieldEscape).
		SetPath(common.ApiTimeSeriesRoute).SetPath(edgexCommon.Device).SetPath(edgexCommon.Name).SetNameFieldPath(deviceName).
		SetPath(edgexCommon.Start).SetPath(strconv.FormatInt(start, 10)).SetPath(edgexCommon.End).SetPath(strconv.FormatInt(end, 10)).BuildPath()
	res := responses.TimeSeriesResponse{}
	var base64Payload string
	if len(resourceNames) > 0 {
		jsonPayload := map[string]any{
			edgexCommon.ResourceNames: resourceNames,
		}
		jsonBytes, err := json.Marshal(jsonPayload)
		if err != nil {
			return res, errors.NewCommonEdgeXWrapper(err)
		}
		base64Payload = base64.StdEncoding.EncodeToString(jsonBytes)
	}

	baseUrl, err := clients.GetBaseUrl(tsc.baseUrlFunc)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	requestParams := url.Values{}
	if base64Payload != "" {
		requestParams.Set(common.Payload, base64Payload)
	}
	err = utils.GetRequest(ctx, &res, baseUrl, requestPath, requestParams, tsc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}
