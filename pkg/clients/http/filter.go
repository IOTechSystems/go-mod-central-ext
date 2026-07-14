// Copyright (C) 2026 IOTech Ltd

package http

import (
	"context"
	"net/url"
	"strconv"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients/http/utils"
	clientsInterfaces "github.com/edgexfoundry/go-mod-core-contracts/v4/clients/interfaces"
	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/clients/interfaces"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/requests"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/responses"
)

type FilterClient struct {
	baseUrlFunc           clients.ClientBaseUrlFunc
	authInjector          clientsInterfaces.AuthenticationInjector
	enableNameFieldEscape bool
}

// NewFilterClient creates an instance of FilterClient
func NewFilterClient(baseUrl string, authInjector clientsInterfaces.AuthenticationInjector, enableNameFieldEscape bool) interfaces.FilterClient {
	return &FilterClient{
		baseUrlFunc:           clients.GetDefaultClientBaseUrlFunc(baseUrl),
		authInjector:          authInjector,
		enableNameFieldEscape: enableNameFieldEscape,
	}
}

func (fc *FilterClient) AllFilters(ctx context.Context, offset int, limit int) (res responses.MultiFiltersResponse, err errors.EdgeX) {
	baseUrl, goErr := clients.GetBaseUrl(fc.baseUrlFunc)
	if goErr != nil {
		return res, errors.NewCommonEdgeXWrapper(goErr)
	}
	params := makeOffsetLimitParams(offset, limit)
	err = utils.GetRequest(ctx, &res, baseUrl, common.ApiAllFilterRoute, params, fc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (fc *FilterClient) FilterById(ctx context.Context, id string) (res responses.FilterResponse, err errors.EdgeX) {
	requestPath := edgexCommon.NewPathBuilder().EnableNameFieldEscape(fc.enableNameFieldEscape).
		SetPath(common.ApiFilterRoute).SetPath(edgexCommon.Id).SetNameFieldPath(id).BuildPath()
	baseUrl, goErr := clients.GetBaseUrl(fc.baseUrlFunc)
	if goErr != nil {
		return res, errors.NewCommonEdgeXWrapper(goErr)
	}
	err = utils.GetRequest(ctx, &res, baseUrl, requestPath, nil, fc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (fc *FilterClient) FiltersByDeviceName(ctx context.Context, deviceName string, offset int, limit int) (res responses.MultiFiltersResponse, err errors.EdgeX) {
	requestPath := edgexCommon.NewPathBuilder().EnableNameFieldEscape(fc.enableNameFieldEscape).
		SetPath(common.ApiFilterRoute).SetPath(edgexCommon.DeviceName).SetNameFieldPath(deviceName).BuildPath()
	baseUrl, goErr := clients.GetBaseUrl(fc.baseUrlFunc)
	if goErr != nil {
		return res, errors.NewCommonEdgeXWrapper(goErr)
	}
	params := makeOffsetLimitParams(offset, limit)
	err = utils.GetRequest(ctx, &res, baseUrl, requestPath, params, fc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (fc *FilterClient) Add(ctx context.Context, reqs []requests.FilterRequest) (res []dtoCommon.BaseWithIdResponse, err errors.EdgeX) {
	baseUrl, goErr := clients.GetBaseUrl(fc.baseUrlFunc)
	if goErr != nil {
		return res, errors.NewCommonEdgeXWrapper(goErr)
	}
	err = utils.PostRequestWithRawData(ctx, &res, baseUrl, common.ApiFilterRoute, nil, reqs, fc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (fc *FilterClient) PutFilters(ctx context.Context, reqs []requests.FilterRequest) (res []dtoCommon.BaseWithIdResponse, err errors.EdgeX) {
	baseUrl, goErr := clients.GetBaseUrl(fc.baseUrlFunc)
	if goErr != nil {
		return res, errors.NewCommonEdgeXWrapper(goErr)
	}
	err = utils.PutRequest(ctx, &res, baseUrl, common.ApiAllFilterRoute, nil, reqs, fc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (fc *FilterClient) Update(ctx context.Context, req requests.FilterRequest) (res dtoCommon.BaseResponse, err errors.EdgeX) {
	requestPath := edgexCommon.NewPathBuilder().EnableNameFieldEscape(fc.enableNameFieldEscape).
		SetPath(common.ApiFilterRoute).SetPath(edgexCommon.Id).SetNameFieldPath(req.Filter.Id).BuildPath()
	baseUrl, goErr := clients.GetBaseUrl(fc.baseUrlFunc)
	if goErr != nil {
		return res, errors.NewCommonEdgeXWrapper(goErr)
	}
	err = utils.PutRequest(ctx, &res, baseUrl, requestPath, nil, req, fc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (fc *FilterClient) DeleteFilterById(ctx context.Context, id string) (res dtoCommon.BaseResponse, err errors.EdgeX) {
	requestPath := edgexCommon.NewPathBuilder().EnableNameFieldEscape(fc.enableNameFieldEscape).
		SetPath(common.ApiFilterRoute).SetPath(edgexCommon.Id).SetNameFieldPath(id).BuildPath()
	baseUrl, goErr := clients.GetBaseUrl(fc.baseUrlFunc)
	if goErr != nil {
		return res, errors.NewCommonEdgeXWrapper(goErr)
	}
	err = utils.DeleteRequest(ctx, &res, baseUrl, requestPath, fc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (fc *FilterClient) DeleteFiltersByDeviceName(ctx context.Context, deviceName string) (res dtoCommon.BaseResponse, err errors.EdgeX) {
	requestPath := edgexCommon.NewPathBuilder().EnableNameFieldEscape(fc.enableNameFieldEscape).
		SetPath(common.ApiFilterRoute).SetPath(edgexCommon.DeviceName).SetNameFieldPath(deviceName).BuildPath()
	baseUrl, goErr := clients.GetBaseUrl(fc.baseUrlFunc)
	if goErr != nil {
		return res, errors.NewCommonEdgeXWrapper(goErr)
	}
	err = utils.DeleteRequest(ctx, &res, baseUrl, requestPath, fc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (fc *FilterClient) DeleteFilters(ctx context.Context) (res dtoCommon.BaseResponse, err errors.EdgeX) {
	baseUrl, goErr := clients.GetBaseUrl(fc.baseUrlFunc)
	if goErr != nil {
		return res, errors.NewCommonEdgeXWrapper(goErr)
	}
	err = utils.DeleteRequest(ctx, &res, baseUrl, common.ApiAllFilterRoute, fc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func makeOffsetLimitParams(offset int, limit int) url.Values {
	params := url.Values{}
	params.Set(edgexCommon.Offset, strconv.Itoa(offset))
	params.Set(edgexCommon.Limit, strconv.Itoa(limit))
	return params
}
