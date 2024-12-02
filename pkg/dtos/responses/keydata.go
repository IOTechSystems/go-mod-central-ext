// Copyright (C) 2024 IOTech Ltd

package responses

import (
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
)

// KeyDataResponse defines the Response Content for GET KeyData DTOs.
type KeyDataResponse struct {
	dtoCommon.BaseResponse `json:",inline"`
	KeyData                dtos.KeyData `json:"keyData"`
}

func NewKeyDataResponse(requestId string, message string, statusCode int, keyData dtos.KeyData) KeyDataResponse {
	return KeyDataResponse{
		BaseResponse: dtoCommon.NewBaseResponse(requestId, message, statusCode),
		KeyData:      keyData,
	}
}
