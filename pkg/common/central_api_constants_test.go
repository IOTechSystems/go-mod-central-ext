// Copyright (C) 2026 IOTech Ltd

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInternalRoute(t *testing.T) {
	tests := []struct {
		route    string
		expected string
	}{
		{ApiMetricsDashboardPath, "/api/v3/internal/metricsdashboard"},
		{ApiMetricsDashboardAllRoute, "/api/v3/internal/metricsdashboard/all"},
		{ApiMetricsDashboardByIdRoute, "/api/v3/internal/metricsdashboard/id/:id"},
		{ApiMetricsDashboardLayoutRoute, "/api/v3/internal/metricsdashboard/id/:id/layout"},
		{ApiMetricsDashboardPanelRoute, "/api/v3/internal/metricsdashboard/id/:id/panel"},
		{ApiMetricsDashboardPanelByIdRoute, "/api/v3/internal/metricsdashboard/id/:id/panel/id/:panelId"},
		{ApiMetricsDashboardOverviewRoute, "/api/v3/internal/metricsdashboard/overview"},
		{ApiMetricsDashboardResourcesRoute, "/api/v3/internal/metricsdashboard/resources"},
		{ApiMetricsDashboardArchiveRoute, "/api/v3/internal/metricsdashboard/archive"},
		{"/api/v3/internal/metricsdashboard/all", "/api/v3/internal/metricsdashboard/all"},
		{"/api/v3/internal", "/api/v3/internal"},
		{"/api/v3/", "/api/v3/"},
		{"/api/v3x/metricsdashboard/all", "/api/v3x/metricsdashboard/all"},
		{"/sseapi/v3/metricsdashboard/id/x", "/sseapi/v3/metricsdashboard/id/x"},
		{"metricsdashboard", "metricsdashboard"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.route, func(t *testing.T) {
			assert.Equal(t, tt.expected, InternalRoute(tt.route))
		})
	}
}
