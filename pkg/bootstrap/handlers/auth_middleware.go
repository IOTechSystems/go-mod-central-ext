// Copyright (C) 2024 IOTech Ltd

package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/edgexfoundry/go-mod-bootstrap/v4/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v4/bootstrap/secret"
	"github.com/edgexfoundry/go-mod-bootstrap/v4/di"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// openBaoIssuer defines the issuer if JWT was issued from OpenBao
const openBaoIssuer = "/v1/identity/oidc"

// SecretStoreAuthenticationHandlerFunc prefixes an existing HandlerFunc
// with a OpenBao-based JWT authentication check or a security-proxy-auth issued JWT check.  Usage:
//
//	 authenticationHook := handlers.NilAuthenticationHandlerFunc()
//	 if secret.IsSecurityEnabled() {
//			lc := container.LoggingClientFrom(dic.Get)
//	     secretProvider := container.SecretProviderFrom(dic.Get)
//	     authenticationHook = handlers.SecretStoreAuthenticationHandlerFunc(secretProvider, lc)
//	 }
//	 For optionally-authenticated requests
//	 r.HandleFunc("path", authenticationHook(handlerFunc)).Methods(http.MethodGet)
//
//	 For unauthenticated requests
//	 r.HandleFunc("path", handlerFunc).Methods(http.MethodGet)
//
// For typical usage, it is preferred to use AutoConfigAuthenticationFunc which
// will automatically select between a real and a fake JWT validation handler.
// func SecretStoreAuthenticationHandlerFunc(secretProvider interfaces.SecretProviderExt, lc logger.LoggingClient, serviceConfig interfaces.Configuration, authInjector contractInterfaces.AuthenticationInjector) echo.MiddlewareFunc {
func SecretStoreAuthenticationHandlerFunc(dic *di.Container) echo.MiddlewareFunc {
	return func(inner echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()
			w := c.Response()
			lc := container.LoggingClientFrom(dic.Get)
			secretProvider := container.SecretProviderExtFrom(dic.Get)
			serviceConfig := container.ConfigurationFrom(dic.Get)
			authInjector := secret.NewJWTSecretProvider(secretProvider)
			authHeader := r.Header.Get("Authorization")
			lc.Debugf("Authorizing incoming call to '%s' via JWT (Authorization len=%d), %v", r.URL.Path, len(authHeader), secretProvider.IsZeroTrustEnabled())

			if secretProvider.IsZeroTrustEnabled() {
				// this implementation will be pick up in the build when build tag no_openziti is specified, where
				// OpenZiti packages are not included and the Zero Trust feature is not available.
				lc.Info("zero trust was enabled, but service is built with no_openziti flag. falling back to token-based auth")
			}

			authParts := strings.Split(authHeader, " ")
			if len(authParts) >= 2 && strings.EqualFold(authParts[0], "Bearer") {
				token := authParts[1]

				parser := jwt.NewParser()
				parsedToken, _, jwtErr := parser.ParseUnverified(token, &jwt.MapClaims{})
				if jwtErr != nil {
					w.Committed = false
					return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				}
				issuer, jwtErr := parsedToken.Claims.GetIssuer()
				if jwtErr != nil {
					w.Committed = false
					return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				}

				if issuer == openBaoIssuer {
					// Verify the JWT by the secret provider IsJWTValid method
					validToken, err := secretProvider.IsJWTValid(token)
					if err != nil {
						lc.Errorf("Error checking JWT validity: %v", err)
						// set Response.Committed to true in order to rewrite the status code
						w.Committed = false
						return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
					} else if !validToken {
						lc.Warnf("Request to '%s' UNAUTHORIZED", r.URL.Path)
						// set Response.Committed to true in order to rewrite the status code
						w.Committed = false
						return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
					}

				} else {
					// Verify the JWT by invoking security-proxy-auth http client
					err := verifyJWT(token, issuer, serviceConfig, authInjector, lc, r.Context())
					if err != nil {
						errResp := dtoCommon.NewBaseResponse("", err.Error(), err.Code())
						return c.JSON(err.Code(), errResp)
					}
				}
				lc.Debugf("Request to '%s' authorized", r.URL.Path)
				return inner(c)
			}
			err := fmt.Errorf("unable to parse JWT for call to '%s'; unauthorized", r.URL.Path)
			lc.Errorf("%v", err)
			// set Response.Committed to true in order to rewrite the status code
			w.Committed = false
			return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		}
	}
}
