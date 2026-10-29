// Copyright (C) 2026 IOTech Ltd

package responses

import (
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
)

// AIDashboardResponse defines the Response Content for GET AIDashboard DTO
type AIDashboardResponse struct {
	dtoCommon.BaseResponse `json:",inline"`
	Dashboard              dtos.AIDashboard `json:"dashboard"`
}

func NewAIDashboardResponse(requestId string, message string, statusCode int, dashboard dtos.AIDashboard) AIDashboardResponse {
	return AIDashboardResponse{
		BaseResponse: dtoCommon.NewBaseResponse(requestId, message, statusCode),
		Dashboard:    dashboard,
	}
}

// MultiAIDashboardInfosResponse defines the Response Content for GET multiple AIDashboardInfo DTOs
type MultiAIDashboardInfosResponse struct {
	dtoCommon.BaseWithTotalCountResponse `json:",inline"`
	Dashboards                           []dtos.AIDashboardInfo `json:"dashboards"`
}

func NewMultiAIDashboardInfosResponse(requestId string, message string, statusCode int, totalCount int64, dashboards []dtos.AIDashboardInfo) MultiAIDashboardInfosResponse {
	return MultiAIDashboardInfosResponse{
		BaseWithTotalCountResponse: dtoCommon.NewBaseWithTotalCountResponse(requestId, message, statusCode, totalCount),
		Dashboards:                 dashboards,
	}
}

// AIPanelResponse defines the Response Content for a written AIPanel DTO.
// MatchedPoints counts the points the panel's tags and devices match; nil when the point catalogue could not be read.
type AIPanelResponse struct {
	dtoCommon.BaseResponse `json:",inline"`
	DashboardId            string       `json:"dashboardId"`
	Panel                  dtos.AIPanel `json:"panel"`
	MatchedPoints          *int         `json:"matchedPoints,omitempty"`
}

func NewAIPanelResponse(requestId string, message string, statusCode int, dashboardId string, panel dtos.AIPanel, matchedPoints *int) AIPanelResponse {
	return AIPanelResponse{
		BaseResponse:  dtoCommon.NewBaseResponse(requestId, message, statusCode),
		DashboardId:   dashboardId,
		Panel:         panel,
		MatchedPoints: matchedPoints,
	}
}

// AIPointsSummaryResponse defines the Response Content for GET AIPointsSummary DTO
type AIPointsSummaryResponse struct {
	dtoCommon.BaseResponse `json:",inline"`
	Summary                dtos.AIPointsSummary `json:"summary"`
}

func NewAIPointsSummaryResponse(requestId string, message string, statusCode int, summary dtos.AIPointsSummary) AIPointsSummaryResponse {
	return AIPointsSummaryResponse{
		BaseResponse: dtoCommon.NewBaseResponse(requestId, message, statusCode),
		Summary:      summary,
	}
}

// AIPointsPageResponse defines the Response Content for GET AIPointsPage DTO
type AIPointsPageResponse struct {
	dtoCommon.BaseResponse `json:",inline"`
	Page                   dtos.AIPointsPage `json:"page"`
}

func NewAIPointsPageResponse(requestId string, message string, statusCode int, page dtos.AIPointsPage) AIPointsPageResponse {
	return AIPointsPageResponse{
		BaseResponse: dtoCommon.NewBaseResponse(requestId, message, statusCode),
		Page:         page,
	}
}
