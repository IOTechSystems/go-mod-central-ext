// Copyright (C) 2021-2024 IOTech Ltd

package http

import (
	"context"
	"net/url"
	"strings"

	"github.com/edgexfoundry/go-mod-central-ext/v4/pkg/clients/interfaces"
	"github.com/edgexfoundry/go-mod-central-ext/v4/pkg/dtos/requests"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients/http/utils"
	clientsInterfaces "github.com/edgexfoundry/go-mod-core-contracts/v4/clients/interfaces"
	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

type SystemManagementClient struct {
	baseUrl      string
	authInjector clientsInterfaces.AuthenticationInjector
}

func NewSystemManagementClient(baseUrl string, authInjector clientsInterfaces.AuthenticationInjector) interfaces.SystemManagementClient {
	return &SystemManagementClient{
		baseUrl:      baseUrl,
		authInjector: authInjector,
	}
}

func (smc *SystemManagementClient) GetHealth(ctx context.Context, services []string) (res []dtoCommon.BaseWithServiceNameResponse, err errors.EdgeX) {
	requestParams := url.Values{}
	requestParams.Set(edgexCommon.Services, strings.Join(services, edgexCommon.CommaSeparator))
	err = utils.GetRequest(ctx, &res, smc.baseUrl, edgexCommon.ApiHealthRoute, requestParams, smc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}

	return
}

func (smc *SystemManagementClient) GetConfig(ctx context.Context, services []string) (res []dtoCommon.BaseWithConfigResponse, err errors.EdgeX) {
	requestParams := url.Values{}
	requestParams.Set(edgexCommon.Services, strings.Join(services, edgexCommon.CommaSeparator))
	err = utils.GetRequest(ctx, &res, smc.baseUrl, edgexCommon.ApiMultiConfigRoute, requestParams, smc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}

	return
}

func (smc *SystemManagementClient) DoOperation(ctx context.Context, reqs []requests.OperationRequest) (res []dtoCommon.BaseResponse, err errors.EdgeX) {
	err = utils.PostRequestWithRawData(ctx, &res, smc.baseUrl, edgexCommon.ApiOperationRoute, nil, reqs, smc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}

	return
}
