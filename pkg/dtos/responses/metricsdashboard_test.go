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
		{"overview", NewMetricsDashboardOverviewResponse("", "", http.StatusOK, dtos.MetricsDashboardOverview{Tags: []dtos.MetricsDashboardTagSummary{{Tag: "hs4:t"}}}),
			[]string{`"units":[]`}},
		{"empty overview", NewMetricsDashboardOverviewResponse("", "", http.StatusOK, dtos.MetricsDashboardOverview{}),
			[]string{`"tags":[]`, `"profiles":[]`}},
		{"resources", NewMetricsDashboardResourcesResponse("", "", http.StatusOK, 0, nil), []string{`"totalCount":0`, `"resources":[]`}},
		{"resource without tags", NewMetricsDashboardResourcesResponse("", "", http.StatusOK, 1, []dtos.MetricsDashboardResource{{ProfileName: "p", ResourceName: "r"}}),
			[]string{`"totalCount":1`, `"tags":[]`}},
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

func TestMetricsDashboardResponses_JSON(t *testing.T) {
	summary := NewMetricsDashboardOverviewResponse("", "", http.StatusOK, dtos.MetricsDashboardOverview{
		PointCount: 180, DeviceCount: 2,
		Tags:     []dtos.MetricsDashboardTagSummary{{Tag: "hs4:co2", PointCount: 1, DeviceCount: 1, Units: []string{"ppm"}}},
		Profiles: []dtos.MetricsDashboardProfileSummary{{ProfileName: "Schneider-BMS", DeviceCount: 1, ResourceCount: 92}},
	})
	assert.JSONEq(t, `{"apiVersion":"v3","statusCode":200,"pointCount":180,"deviceCount":2,
		"tags":[{"tag":"hs4:co2","pointCount":1,"deviceCount":1,"units":["ppm"]}],
		"profiles":[{"profileName":"Schneider-BMS","deviceCount":1,"resourceCount":92}]}`, marshalled(t, summary))
	var gotSummary MetricsDashboardOverviewResponse
	require.NoError(t, json.Unmarshal([]byte(marshalled(t, summary)), &gotSummary))
	assert.Equal(t, summary, gotSummary)

	search := NewMetricsDashboardResourcesResponse("", "", http.StatusOK, 2, []dtos.MetricsDashboardResource{{
		ProfileName: "Johnson-BMS", ResourceName: "Lobby IAQ Sensor", DeviceCount: 1, Units: "ppm", ValueType: "Float32",
		Tags: []string{}, Description: "Indoor air quality, CO2 ppm",
	}})
	assert.JSONEq(t, `{"apiVersion":"v3","statusCode":200,"totalCount":2,
		"resources":[{"profileName":"Johnson-BMS","resourceName":"Lobby IAQ Sensor","deviceCount":1,"units":"ppm","valueType":"Float32",
		"tags":[],"description":"Indoor air quality, CO2 ppm"}]}`, marshalled(t, search))
	var gotSearch MetricsDashboardResourcesResponse
	require.NoError(t, json.Unmarshal([]byte(marshalled(t, search)), &gotSearch))
	assert.Equal(t, search, gotSearch)
}
