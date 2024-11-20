// Copyright (C) 2024 IOTech Ltd

package responses

import (
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
)

type TokenResponse struct {
	dtoCommon.BaseResponse `json:",inline"`
	JWT                    string `json:"jwt"`
}

// NewTokenResponse returns the JWT
func NewTokenResponse(requestId string, message string, statusCode int, token string) TokenResponse {
	return TokenResponse{
		BaseResponse: dtoCommon.NewBaseResponse(requestId, message, statusCode),
		JWT:          token,
	}
}
