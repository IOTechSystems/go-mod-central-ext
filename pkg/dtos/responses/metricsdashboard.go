// Copyright (C) 2026 IOTech Ltd

package responses

import (
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"

	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
)

// MetricsDashboardResponse defines the Response Content for a single MetricsDashboard
type MetricsDashboardResponse struct {
	dtoCommon.BaseResponse `json:",inline"`
	Dashboard              dtos.MetricsDashboard `json:"dashboard"`
}

func NewMetricsDashboardResponse(requestId string, message string, statusCode int, dashboard dtos.MetricsDashboard) MetricsDashboardResponse {
	dashboard.Panels = orEmpty(dashboard.Panels)
	return MetricsDashboardResponse{
		BaseResponse: dtoCommon.NewBaseResponse(requestId, message, statusCode),
		Dashboard:    dashboard,
	}
}

// MultiMetricsDashboardInfosResponse defines the Response Content for the list of MetricsDashboards
type MultiMetricsDashboardInfosResponse struct {
	dtoCommon.BaseWithTotalCountResponse `json:",inline"`
	Dashboards                           []dtos.MetricsDashboardInfo `json:"dashboards"`
}

func NewMultiMetricsDashboardInfosResponse(requestId string, message string, statusCode int, totalCount int64, dashboards []dtos.MetricsDashboardInfo) MultiMetricsDashboardInfosResponse {
	return MultiMetricsDashboardInfosResponse{
		BaseWithTotalCountResponse: dtoCommon.NewBaseWithTotalCountResponse(requestId, message, statusCode, totalCount),
		Dashboards:                 orEmpty(dashboards),
	}
}

// MetricsDashboardPanelResponse defines the Response Content for a panel write; Panel is nil after a delete
type MetricsDashboardPanelResponse struct {
	dtoCommon.BaseResponse `json:",inline"`
	DashboardId            string                      `json:"dashboardId"`
	Revision               int64                       `json:"revision"`
	Panel                  *dtos.MetricsDashboardPanel `json:"panel,omitempty"`
}

func NewMetricsDashboardPanelResponse(requestId string, message string, statusCode int, dashboardId string, revision int64, panel *dtos.MetricsDashboardPanel) MetricsDashboardPanelResponse {
	return MetricsDashboardPanelResponse{
		BaseResponse: dtoCommon.NewBaseResponse(requestId, message, statusCode),
		DashboardId:  dashboardId,
		Revision:     revision,
		Panel:        panel,
	}
}

// MetricsDashboardTagSummaryResponse defines the Response Content for GET metrics dashboard tags
type MetricsDashboardTagSummaryResponse struct {
	dtoCommon.BaseResponse          `json:",inline"`
	dtos.MetricsDashboardTagSummary `json:",inline"`
}

func NewMetricsDashboardTagSummaryResponse(requestId string, message string, statusCode int, summary dtos.MetricsDashboardTagSummary) MetricsDashboardTagSummaryResponse {
	rows := make([]dtos.MetricsDashboardTagRow, len(summary.Tags))
	for i, row := range summary.Tags {
		row.Units = orEmpty(row.Units)
		rows[i] = row
	}
	summary.Tags = rows
	return MetricsDashboardTagSummaryResponse{
		BaseResponse:               dtoCommon.NewBaseResponse(requestId, message, statusCode),
		MetricsDashboardTagSummary: summary,
	}
}

// MetricsDashboardPointsResponse defines the Response Content for GET metrics dashboard points
type MetricsDashboardPointsResponse struct {
	dtoCommon.BaseWithTotalCountResponse `json:",inline"`
	Points                               []dtos.MetricsDashboardPoint `json:"points"`
}

func NewMetricsDashboardPointsResponse(requestId string, message string, statusCode int, totalCount int64, points []dtos.MetricsDashboardPoint) MetricsDashboardPointsResponse {
	return MetricsDashboardPointsResponse{
		BaseWithTotalCountResponse: dtoCommon.NewBaseWithTotalCountResponse(requestId, message, statusCode, totalCount),
		Points:                     orEmpty(points),
	}
}

// orEmpty turns a nil list into an empty one, so it is sent as [] and not null
func orEmpty[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}
