// Copyright (C) 2024 IOTech Ltd

package http

import (
	"context"
	"net/http"
	"strings"

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

func (ac AuthClient) Authenticate(ctx context.Context, headers map[string]string) (errors.EdgeX, string) {
	err := utils.PostRequestWithRawDataAndHeaders(ctx, nil, ac.baseUrl, common.ApiAuthRoute, nil, nil, ac.authInjector, headers)
	if err != nil {
		if err.Code() == http.StatusForbidden {
			// Check if the returned response contains JWT
			hasJWT := strings.HasPrefix(err.Error(), "request failed, status code: 403, err: \"")
			if hasJWT {
				jwt := strings.TrimPrefix(err.Error(), "request failed, status code: 403, err: \"")
				jwt = strings.TrimSuffix(jwt, "\"\n")
				return errors.NewCommonEdgeXWrapper(err), jwt
			}
		}
		return errors.NewCommonEdgeXWrapper(err), ""
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
