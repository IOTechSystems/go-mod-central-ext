// Copyright (C) 2026 IOTech Ltd

package responses

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func marshalled(t *testing.T, v any) string {
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return string(b)
}

func TestMetricsDashboardResponses_EmptyListsAreArrays(t *testing.T) {
	tests := []struct {
		name     string
		response any
		expected []string
	}{
		{"dashboard", NewMetricsDashboardResponse("", "", http.StatusOK, dtos.MetricsDashboard{Name: "d"}), []string{`"panels":[]`}},
		{"dashboard infos", NewMultiMetricsDashboardInfosResponse("", "", http.StatusOK, 0, nil), []string{`"dashboards":[]`}},
		{"tag summary", NewMetricsDashboardTagSummaryResponse("", "", http.StatusOK, dtos.MetricsDashboardTagSummary{Tags: []dtos.MetricsDashboardTagRow{{Tag: "hs4:t"}}}),
			[]string{`"units":[]`}},
		{"empty tag summary", NewMetricsDashboardTagSummaryResponse("", "", http.StatusOK, dtos.MetricsDashboardTagSummary{}), []string{`"tags":[]`}},
		{"points", NewMetricsDashboardPointsResponse("", "", http.StatusOK, 0, nil), []string{`"totalCount":0`, `"points":[]`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := marshalled(t, tt.response)
			for _, e := range tt.expected {
				assert.Contains(t, out, e)
			}
			assert.NotContains(t, out, "null")
		})
	}
}
