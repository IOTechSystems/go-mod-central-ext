// Copyright (C) 2024-2025 IOTech Ltd

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

type RolePolicyClient struct {
	baseUrl               string
	authInjector          clientsInterfaces.AuthenticationInjector
	enableNameFieldEscape bool
}

// NewRolePolicyClient creates an instance of RolePolicyClient
func NewRolePolicyClient(baseUrl string, authInjector clientsInterfaces.AuthenticationInjector, enableNameFieldEscape bool) interfaces.RolePolicyClient {
	return &RolePolicyClient{
		baseUrl:               baseUrl,
		authInjector:          authInjector,
		enableNameFieldEscape: enableNameFieldEscape,
	}
}

func (uc RolePolicyClient) Add(ctx context.Context, reqs requests.AddRolePolicyRequest) (res dtoCommon.BaseResponse, err errors.EdgeX) {
	err = utils.PostRequestWithRawData(ctx, &res, uc.baseUrl, common.ApiRolePolicyRoute, nil, reqs, uc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (uc RolePolicyClient) Update(ctx context.Context, reqs requests.AddRolePolicyRequest) (res dtoCommon.BaseResponse, err errors.EdgeX) {
	path := edgexCommon.NewPathBuilder().EnableNameFieldEscape(uc.enableNameFieldEscape).
		SetPath(common.ApiRolePolicyRoute).SetPath(common.Role).SetNameFieldPath(reqs.RolePolicy.Role).BuildPath()
	err = utils.PutRequest(ctx, &res, uc.baseUrl, path, nil, reqs, uc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (uc RolePolicyClient) AllRolePolicies(ctx context.Context, offset int, limit int) (res responses.MultiRolePolicyResponse, err errors.EdgeX) {
	requestParams := url.Values{}

	requestParams.Set(edgexCommon.Offset, strconv.Itoa(offset))
	requestParams.Set(edgexCommon.Limit, strconv.Itoa(limit))
	err = utils.GetRequest(ctx, &res, uc.baseUrl, common.ApiAllRolePolicyRoute, requestParams, uc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (uc RolePolicyClient) RolePolicyByRole(ctx context.Context, role string) (res responses.RolePolicyResponse, err errors.EdgeX) {
	path := edgexCommon.NewPathBuilder().EnableNameFieldEscape(uc.enableNameFieldEscape).
		SetPath(common.ApiRolePolicyRoute).SetPath(common.Role).SetNameFieldPath(role).BuildPath()
	err = utils.GetRequest(ctx, &res, uc.baseUrl, path, nil, uc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (uc RolePolicyClient) DeleteRolePolicyByRole(ctx context.Context, role string) (res dtoCommon.BaseResponse, err errors.EdgeX) {
	path := edgexCommon.NewPathBuilder().EnableNameFieldEscape(uc.enableNameFieldEscape).
		SetPath(common.ApiRolePolicyRoute).SetPath(common.Role).SetNameFieldPath(role).BuildPath()
	err = utils.DeleteRequest(ctx, &res, uc.baseUrl, path, uc.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}
