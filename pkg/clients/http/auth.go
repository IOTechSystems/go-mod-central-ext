// Copyright (C) 2024 IOTech Ltd

package http

import (
	"context"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/clients/interfaces"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/requests"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/responses"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients/http/utils"
	clientsInterfaces "github.com/edgexfoundry/go-mod-core-contracts/v4/clients/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

type AuthClient struct {
	baseUrl      string
	authInjector clientsInterfaces.AuthenticationInjector
}

// NewAuthClient creates an instance of AuthClient
func NewAuthClient(baseUrl string, authInjector clientsInterfaces.AuthenticationInjector) interfaces.AuthClient {
	return &AuthClient{
		baseUrl:      baseUrl,
		authInjector: authInjector,
	}
}

func (ac AuthClient) Auth(ctx context.Context, headers map[string]string) (errors.EdgeX, string) {
	var newJWT string

	err := utils.PostRequestWithRawDataAndHeaders(ctx, &newJWT, ac.baseUrl, common.ApiAuthRoute, nil, nil, ac.authInjector, headers)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err), newJWT
	}

	return nil, ""
}

func (ac AuthClient) AuthRoutes(ctx context.Context, headers map[string]string, routeReqs []requests.AuthRouteRequest) (res responses.AuthRouteResponse, err errors.EdgeX) {
	err = utils.PostRequestWithRawDataAndHeaders(ctx, &res, ac.baseUrl, common.ApiAuthRoutesRoute, nil, routeReqs, ac.authInjector, headers)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}
