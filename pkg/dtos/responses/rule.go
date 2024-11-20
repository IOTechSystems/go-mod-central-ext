// Copyright (C) 2023-2024 IOTech Ltd

package responses

import (
	"github.com/edgexfoundry/go-mod-central-ext/v4/pkg/dtos"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
)

// RuleResponse defines the Response Content for GET rule DTO.
type RuleResponse struct {
	dtoCommon.BaseResponse `json:",inline"`
	Rule                   dtos.Rule `json:"rule"`
}

// MultiRulesResponse defines the Response Content for GET multiple rule DTO.
type MultiRulesResponse struct {
	dtoCommon.BaseWithTotalCountResponse `json:",inline"`
	Rules                                []dtos.Rule `json:"rules"`
}

func NewMultiRulesResponse(requestId string, message string, statusCode int, totalCount uint32, rules []dtos.Rule) MultiRulesResponse {
	return MultiRulesResponse{
		BaseWithTotalCountResponse: dtoCommon.NewBaseWithTotalCountResponse(requestId, message, statusCode, totalCount),
		Rules:                      rules,
	}
}
