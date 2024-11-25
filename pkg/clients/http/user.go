// Copyright (C) 2024 IOTech Ltd

package http

import (
	"context"
	"net/url"
	"strconv"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/clients/interfaces"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/requests"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/responses"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients/http/utils"
	clientsInterfaces "github.com/edgexfoundry/go-mod-core-contracts/v4/clients/interfaces"
	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

type UserClient struct {
	baseUrl               string
	authInjector          clientsInterfaces.AuthenticationInjector
	enableNameFieldEscape bool
}

// NewUserClient creates an instance of UserClient
func NewUserClient(baseUrl string, authInjector clientsInterfaces.AuthenticationInjector, enableNameFieldEscape bool) interfaces.UserClient {
	return &UserClient{
		baseUrl:               baseUrl,
		authInjector:          authInjector,
		enableNameFieldEscape: enableNameFieldEscape,
	}
}

func (uc UserClient) Add(ctx context.Context, reqs []requests.AddUserRequest) (res []dtoCommon.BaseWithIdResponse, err errors.EdgeX) {
	err = utils.PostRequestWithRawData(ctx, &res, uc.baseUrl, common.ApiUserRoute, nil, reqs, uc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (uc UserClient) Update(ctx context.Context, reqs []requests.UpdateUserRequest) (res []dtoCommon.BaseResponse, err errors.EdgeX) {
	err = utils.PatchRequest(ctx, &res, uc.baseUrl, common.ApiUserRoute, nil, reqs, uc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (uc UserClient) AllUsers(ctx context.Context, offset int, limit int) (res responses.MultiUsersResponse, err errors.EdgeX) {
	requestParams := url.Values{}

	requestParams.Set(edgexCommon.Offset, strconv.Itoa(offset))
	requestParams.Set(edgexCommon.Limit, strconv.Itoa(limit))
	err = utils.GetRequest(ctx, &res, uc.baseUrl, common.ApiAllUserRoute, requestParams, uc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (uc UserClient) UserByName(ctx context.Context, name string) (res responses.UserResponse, err errors.EdgeX) {
	path := edgexCommon.NewPathBuilder().EnableNameFieldEscape(uc.enableNameFieldEscape).
		SetPath(common.ApiUserRoute).SetPath(edgexCommon.Name).SetNameFieldPath(name).BuildPath()
	err = utils.GetRequest(ctx, &res, uc.baseUrl, path, nil, uc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (uc UserClient) DeleteUserByName(ctx context.Context, name string) (res dtoCommon.BaseResponse, err errors.EdgeX) {
	path := edgexCommon.NewPathBuilder().EnableNameFieldEscape(uc.enableNameFieldEscape).
		SetPath(common.ApiUserRoute).SetPath(edgexCommon.Name).SetNameFieldPath(name).BuildPath()
	err = utils.DeleteRequest(ctx, &res, uc.baseUrl, path, uc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (uc UserClient) Login(ctx context.Context, reqs requests.LoginRequest) (res responses.TokenResponse, err errors.EdgeX) {
	err = utils.PostRequestWithRawData(ctx, &res, uc.baseUrl, common.ApiLoginRoute, nil, reqs, uc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (uc UserClient) Authenticate(ctx context.Context, headers map[string]string) (res any, err errors.EdgeX) {
	err = utils.PostRequestWithRawDataAndHeaders(ctx, &res, uc.baseUrl, common.ApiAuthRoute, nil, nil, uc.authInjector, headers)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (uc UserClient) AuthRoutes(ctx context.Context, headers map[string]string, routeReqs []requests.AuthRouteRequest) (res responses.AuthRouteResponse, err errors.EdgeX) {
	err = utils.PostRequestWithRawDataAndHeaders(ctx, &res, uc.baseUrl, common.ApiAuthRoutesRoute, nil, routeReqs, uc.authInjector, headers)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}
