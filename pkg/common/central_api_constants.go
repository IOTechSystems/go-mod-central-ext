// Copyright (C) 2024 IOTech Ltd

package common

import (
	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
)

const (
	Ids           = "ids"
	User          = "user"
	Group         = "group"
	PublicKey     = "rsa_public_key"
	Ack           = "ack"
	Acknowledge   = "acknowledge"
	Unacknowledge = "unacknowledge"
	NoCallback    = "nocallback" //query string to ask core-metadata not to invoke DS callback
	Role          = "role"
)

const (
	ApiRuleRoute       = edgexCommon.ApiBase + "/rule"
	ApiAllRulesRoute   = ApiRuleRoute + "/" + edgexCommon.All
	ApiRuleByNameRoute = ApiRuleRoute + "/" + edgexCommon.Name + "/:" + edgexCommon.Name

	ApiCoreCommandsByDeviceNameRoute = edgexCommon.ApiBase + "/command/device" + "/" + edgexCommon.Name + "/:" + edgexCommon.Name

	ApiRolePolicyRoute       = edgexCommon.ApiBase + "/rolepolicy"
	ApiAllRolePolicyRoute    = ApiRolePolicyRoute + "/" + edgexCommon.All
	ApiRolePolicyByRoleRoute = ApiRolePolicyRoute + "/" + Role + "/:" + Role

	ApiUserRoute       = edgexCommon.ApiBase + "/user"
	ApiAllUserRoute    = ApiUserRoute + "/" + edgexCommon.All
	ApiUserByNameRoute = ApiUserRoute + "/" + edgexCommon.Name + "/:" + edgexCommon.Name

	ApiLoginRoute        = edgexCommon.ApiBase + "/login"
	ApiLogoutRoute       = edgexCommon.ApiBase + "/logout"
	ApiAuthRoute         = edgexCommon.ApiBase + "/auth"
	ApiAuthRoutesRoute   = edgexCommon.ApiBase + "/auth-routes"
	ApiRefreshTokenRoute = edgexCommon.ApiBase + "/refresh-token"

	ApiKeyRoute                     = edgexCommon.ApiBase + "/key"
	ApiVerificationKeyByIssuerRoute = ApiKeyRoute + "/" + VerificationKeyType + "/" + Issuer + "/:" + Issuer
)

// constants relate to header names
const (
	AuthorizationHeader   = "Authorization"
	ForwardedUriHeader    = "X-Forwarded-Uri"
	ForwardedMethodHeader = "X-Forwarded-Method"
)
