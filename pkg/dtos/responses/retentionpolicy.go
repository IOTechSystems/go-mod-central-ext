// Copyright (C) 2025 IOTech Ltd

package responses

import (
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"

	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
)

// RetentionPolicyResponse defines the Response Content for GET RetentionPolicy DTO
type RetentionPolicyResponse struct {
	dtoCommon.BaseResponse `json:",inline"`
	RetentionPolicy        dtos.RetentionPolicy `json:"retentionPolicy"`
}

func NewRetentionPolicyResponse(requestId string, message string, statusCode int, rp dtos.RetentionPolicy) RetentionPolicyResponse {
	return RetentionPolicyResponse{
		BaseResponse:    dtoCommon.NewBaseResponse(requestId, message, statusCode),
		RetentionPolicy: rp,
	}
}

// MultiRetentionPolicyResponse defines the Response Content for GET multiple RetentionPolicy DTOs.
type MultiRetentionPolicyResponse struct {
	dtoCommon.BaseWithTotalCountResponse `json:",inline"`
	RetentionPolicies                    []dtos.RetentionPolicy `json:"retentionPolicies"`
}

func NewMultiRetentionPolicyResponse(requestId string, message string, statusCode int, totalCount int64, rps []dtos.RetentionPolicy) MultiRetentionPolicyResponse {
	return MultiRetentionPolicyResponse{
		BaseWithTotalCountResponse: dtoCommon.NewBaseWithTotalCountResponse(requestId, message, statusCode, totalCount),
		RetentionPolicies:          rps,
	}
}
