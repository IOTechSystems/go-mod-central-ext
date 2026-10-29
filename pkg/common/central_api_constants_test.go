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
		{ApiAIDashboardRoute, "/api/v3/internal/ai-dashboard"},
		{ApiAIDashboardAllRoute, "/api/v3/internal/ai-dashboard/all"},
		{ApiAIDashboardByIdPath, "/api/v3/internal/ai-dashboard/id"},
		{ApiAIDashboardByIdRoute, "/api/v3/internal/ai-dashboard/id/:id"},
		{ApiAIDashboardPanelRoute, "/api/v3/internal/ai-dashboard/id/:id/panel"},
		{ApiAIDashboardPanelByIdRoute, "/api/v3/internal/ai-dashboard/id/:id/panel/id/:panelId"},
		{ApiAIDashboardPointsRoute, "/api/v3/internal/ai-dashboard/points"},
		{ApiInternalResourceTagJobRoute, ApiInternalResourceTagJobRoute},
		{"/api/v3/internal", "/api/v3/internal"},
		{"/api/v3x/ai-dashboard", "/api/v3x/ai-dashboard"},
		{"/api/v3", "/api/v3"},
		{"/ai-dashboard", "/ai-dashboard"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.route, func(t *testing.T) {
			actual := InternalRoute(tt.route)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, actual, InternalRoute(actual), "idempotent")
		})
	}
}
