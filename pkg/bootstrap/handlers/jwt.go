// Copyright (C) 2024 IOTech Ltd

package handlers

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/pem"
	stdErrs "errors"
	"fmt"

	httpClients "github.com/IOTechSystems/go-mod-central-ext/v4/pkg/clients/http"

	"github.com/edgexfoundry/go-mod-bootstrap/v4/bootstrap/interfaces"
	contractInterfaces "github.com/edgexfoundry/go-mod-core-contracts/v4/clients/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"

	"github.com/golang-jwt/jwt/v5"
)

// verifyJWT verifies if the JWT is valid using the verification key from proxy-auth
func verifyJWT(token string,
	issuer string,
	serviceConfig interfaces.Configuration,
	authInjector contractInterfaces.AuthenticationInjector,
	lc logger.LoggingClient,
	ctx context.Context) errors.EdgeX {
	bootstrapClients := *serviceConfig.GetBootstrap().Clients

	proxyAuthConfig, ok := bootstrapClients[common.SecurityProxyAuthServiceKey]
	if !ok {
		return errors.NewCommonEdgeX(errors.KindServerError, "security-proxy-auth client not defined in the service config", nil)
	}

	proxyAuthURL := proxyAuthConfig.Url()
	authClient := httpClients.NewAuthClient(proxyAuthURL, authInjector)
	keyResponse, edgexErr := authClient.VerificationKeyByIssuer(ctx, issuer)
	if edgexErr != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to obtaining the verification key from proxy-auth for issuer %s", issuer), edgexErr)
	}

	_, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
		return processVerificationKey(keyResponse.KeyData.Key, token.Method.Alg(), lc, ctx)
	})
	if err != nil {
		if stdErrs.Is(err, jwt.ErrTokenExpired) {
			// Skip the JWT expired error
			lc.Debug("JWT is valid but expired")
			return nil
		} else {
			if stdErrs.Is(err, jwt.ErrTokenMalformed) ||
				stdErrs.Is(err, jwt.ErrTokenUnverifiable) ||
				stdErrs.Is(err, jwt.ErrTokenSignatureInvalid) ||
				stdErrs.Is(err, jwt.ErrTokenRequiredClaimMissing) {
				lc.Errorf("Invalid jwt : %v\n", err)
				return errors.NewCommonEdgeX(errors.KindForbidden, "invalid jwt", err)
			}
			lc.Errorf("Error occurred while validating JWT: %v", err)
			return errors.NewCommonEdgeX(errors.KindServerError, "failed to parse jwt", err)
		}
	}
	return nil
}

// processVerificationKey processes the verification key obtained proxy-auth and return the public key with the corresponding format based on the JWT signing algorithm
func processVerificationKey(keyString string, alg string, lc logger.LoggingClient, ctx context.Context) (any, errors.EdgeX) {
	keyBytes := []byte(keyString)

	switch alg {
	case jwt.SigningMethodHS256.Alg(), jwt.SigningMethodHS384.Alg(), jwt.SigningMethodHS512.Alg():
		binaryKey, err := base64.StdEncoding.DecodeString(string(keyBytes))
		if err != nil {
			lc.Debugf("the key is not a valid base64, err: '%v', using the key '%s' without base64 encoding.", err, keyBytes)
			return keyBytes, nil
		}

		return binaryKey, nil
	case jwt.SigningMethodEdDSA.Alg():
		block, _ := pem.Decode(keyBytes)
		if block == nil || block.Type != "PUBLIC KEY" {
			return nil, errors.NewCommonEdgeX(errors.KindServerError, "failed to decode the verification key PEM block", nil)
		}

		lc.Debugf("public key: %s", string(keyBytes))
		edPublicKey := ed25519.PublicKey(block.Bytes)
		return edPublicKey, nil
	case jwt.SigningMethodRS256.Alg(), jwt.SigningMethodRS384.Alg(), jwt.SigningMethodRS512.Alg(),
		jwt.SigningMethodPS256.Alg(), jwt.SigningMethodPS384.Alg(), jwt.SigningMethodPS512.Alg():
		rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to parse '%s' rsa verification key", alg), err)
		}

		return rsaPublicKey, nil
	case jwt.SigningMethodES256.Alg(), jwt.SigningMethodES384.Alg(), jwt.SigningMethodES512.Alg():
		ecdsaPublicKey, err := jwt.ParseECPublicKeyFromPEM(keyBytes)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to parse '%s' es verification key", alg), err)
		}

		return ecdsaPublicKey, nil
	default:
		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("unsupported signing algorithm '%s'", alg), nil)
	}
}
