// Copyright (C) 2024-2025 IOTech Ltd

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
	ApiAuthGraphQLRoute  = ApiAuthRoute + "/graphql"
	ApiAuthRoutesRoute   = edgexCommon.ApiBase + "/auth-routes"
	ApiRefreshTokenRoute = edgexCommon.ApiBase + "/refresh-token"

	ApiKeyRoute                     = edgexCommon.ApiBase + "/key"
	ApiVerificationKeyByIssuerRoute = ApiKeyRoute + "/" + VerificationKeyType + "/" + Issuer + "/:" + Issuer

	ApiFilterRoute           = edgexCommon.ApiBase + "/filter"
	ApiFilterIdRoute         = ApiFilterRoute + "/" + edgexCommon.Id + "/:" + edgexCommon.Id
	ApiFilterDeviceNameRoute = ApiFilterRoute + "/" + edgexCommon.DeviceName + "/:" + edgexCommon.DeviceName
	ApiAllFilterRoute        = ApiFilterRoute + "/" + edgexCommon.All

	ApiRetentionPolicyRoute     = edgexCommon.ApiBase + "/retentionpolicy"
	ApiRetentionPolicyByIdRoute = ApiRetentionPolicyRoute + "/" + edgexCommon.Id + "/:" + edgexCommon.Id
	ApiAllRetentionPolicyRoute  = ApiRetentionPolicyRoute + "/" + edgexCommon.All
)

// constants relate to header names
const (
	AuthorizationHeader   = "Authorization"
	ForwardedUriHeader    = "X-Forwarded-Uri"
	ForwardedMethodHeader = "X-Forwarded-Method"
)
