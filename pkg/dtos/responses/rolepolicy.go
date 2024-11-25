// Copyright (C) 2024 IOTech Ltd

package responses

import (
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
)

// RolePolicyResponse defines the Response Content for GET RolePolicy DTOs
type RolePolicyResponse struct {
	dtoCommon.BaseResponse `json:",inline"`
	RolePolicy             dtos.RolePolicy `json:"rolePolicy"`
}

func NewRolePolicyResponse(requestId string, message string, statusCode int, rolePolicy dtos.RolePolicy) RolePolicyResponse {
	return RolePolicyResponse{
		BaseResponse: dtoCommon.NewBaseResponse(requestId, message, statusCode),
		RolePolicy:   rolePolicy,
	}
}

// MultiRolePolicyResponse defines the Response Content for GET multiple RolePolicy DTOs
type MultiRolePolicyResponse struct {
	dtoCommon.BaseWithTotalCountResponse `json:",inline"`
	RolePolicies                         []dtos.RolePolicy `json:"rolePolicies"`
}

func NewMultiRolePolicyResponse(requestId string, message string, statusCode int, totalCount uint32, rolePolicies []dtos.RolePolicy) MultiRolePolicyResponse {
	return MultiRolePolicyResponse{
		BaseWithTotalCountResponse: dtoCommon.NewBaseWithTotalCountResponse(requestId, message, statusCode, totalCount),
		RolePolicies:               rolePolicies,
	}
}

// AuthRouteResponse defines the Response Content for POST AuthRoute DTO
type AuthRouteResponse struct {
	dtoCommon.BaseResponse `json:",inline"`
	AuthResponses          []dtos.AuthRouteResult `json:"authResponses"`
}

func NewAuthRouteResponse(requestId string, message string, statusCode int, authRouteResults []dtos.AuthRouteResult) AuthRouteResponse {
	return AuthRouteResponse{
		BaseResponse:  dtoCommon.NewBaseResponse(requestId, message, statusCode),
		AuthResponses: authRouteResults,
	}
}
