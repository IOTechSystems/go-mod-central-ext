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
)

const (
	ApiRuleRoute       = edgexCommon.ApiBase + "/rule"
	ApiAllRulesRoute   = ApiRuleRoute + "/" + edgexCommon.All
	ApiRuleByNameRoute = ApiRuleRoute + "/" + edgexCommon.Name + "/:" + edgexCommon.Name

	ApiCoreCommandsByDeviceNameRoute = edgexCommon.ApiBase + "/command/device" + "/" + edgexCommon.Name + "/:" + edgexCommon.Name
)
