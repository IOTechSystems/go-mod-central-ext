// Copyright (C) 2025 IOTech Ltd

package responses

import (
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
)

// FilterResponse defines the Response Content for GET Filter DTOs.
type FilterResponse struct {
	dtoCommon.BaseResponse `json:",inline"`
	Filter                 dtos.Filter `json:"filter"`
}

func NewFilterResponse(requestId string, message string, statusCode int, filter dtos.Filter) FilterResponse {
	return FilterResponse{
		BaseResponse: dtoCommon.NewBaseResponse(requestId, message, statusCode),
		Filter:       filter,
	}
}

// MultiFiltersResponse defines the Response Content for GET multiple Filter DTOs.
type MultiFiltersResponse struct {
	dtoCommon.BaseWithTotalCountResponse `json:",inline"`
	Filters                              []dtos.Filter `json:"filters"`
}

func NewMultiFiltersResponse(requestId string, message string, statusCode int, totalCount int64, filters []dtos.Filter) MultiFiltersResponse {
	return MultiFiltersResponse{
		BaseWithTotalCountResponse: dtoCommon.NewBaseWithTotalCountResponse(requestId, message, statusCode, totalCount),
		Filters:                    filters,
	}
}
