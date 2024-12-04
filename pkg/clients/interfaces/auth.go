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
	// Auth validates and authorizes a user based on request headers and issues a new JWT if the exp claim has expired but the last chance issue claim is still valid
	Auth(ctx context.Context, headers map[string]string) (err errors.EdgeX, newJWT string)
	// AuthRoutes check user permissions on a sets of path-method pairs with the authorization header
	AuthRoutes(ctx context.Context, headers map[string]string, reqs []requests.AuthRouteRequest) (responses.AuthRouteResponse, errors.EdgeX)
	// VerificationKeyByIssuer returns the JWT verification key by the specified issuer
	VerificationKeyByIssuer(ctx context.Context, issuer string) (res responses.KeyDataResponse, err errors.EdgeX)
}
