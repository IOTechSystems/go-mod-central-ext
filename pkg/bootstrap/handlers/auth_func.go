// Copyright (C) 2024 IOTech Ltd

package handlers

import (
	"os"
	"strconv"

	"github.com/edgexfoundry/go-mod-bootstrap/v4/bootstrap/interfaces"
	"github.com/edgexfoundry/go-mod-bootstrap/v4/bootstrap/secret"
	"github.com/edgexfoundry/go-mod-bootstrap/v4/di"

	"github.com/labstack/echo/v4"
)

// NilAuthenticationHandlerFunc just invokes a nested handler
func NilAuthenticationHandlerFunc() echo.MiddlewareFunc {
	return func(inner echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return inner(c)
		}
	}
}

// AutoConfigAuthenticationFunc auto-selects between a HandlerFunc
// wrapper that does authentication and a HandlerFunc wrapper that does not.
// By default, JWT validation is enabled in secure mode
// (i.e. when using a real secrets provider instead of a no-op stub)
//
// Set EDGEX_DISABLE_JWT_VALIDATION to 1, t, T, TRUE, true, or True
// to disable JWT validation.  This might be wanted for an EdgeX
// adopter that wanted to only validate JWT's at the proxy layer,
// or as an escape hatch for a caller that cannot authenticate.
func AutoConfigAuthenticationFunc(dic *di.Container, serviceConfig interfaces.Configuration) echo.MiddlewareFunc {
	// Golang standard library treats an error as false
	disableJWTValidation, _ := strconv.ParseBool(os.Getenv("EDGEX_DISABLE_JWT_VALIDATION"))
	authenticationHook := NilAuthenticationHandlerFunc()

	if secret.IsSecurityEnabled() && !disableJWTValidation {
		authenticationHook = SecretStoreAuthenticationHandlerFunc(dic, serviceConfig)
	}
	return authenticationHook
}