// Copyright (C) 2024 IOTech Ltd

package interfaces

import (
	"context"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/requests"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/responses"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

// AuthClient defines the interface for interactions with the auth API endpoint on the IOTech security-proxy-auth service.
type AuthClient interface {
	// Auth authenticates and authorizes a user by request headers
	Auth(ctx context.Context, headers map[string]string) (errors.EdgeX, any)
	// AuthRoutes check user permissions on a sets of path-method pairs with the authorization header
	AuthRoutes(ctx context.Context, headers map[string]string, reqs []requests.AuthRouteRequest) (responses.AuthRouteResponse, errors.EdgeX)
}
