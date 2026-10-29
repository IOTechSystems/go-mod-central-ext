// Copyright (C) 2026 IOTech Ltd

package responses

import (
	"encoding/json"
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"

	"github.com/stretchr/testify/require"
)

func requireAIBase(t *testing.T, requestId, message string, statusCode int) {
	require.Equal(t, common.ExpectedRequestId, requestId)
	require.Equal(t, common.ExpectedMessage, message)
	require.Equal(t, common.ExpectedStatusCode, statusCode)
}

func TestNewAIDashboardResponse(t *testing.T) {
	expected := dtos.AIDashboard{Id: "d", Name: "a"}
	actual := NewAIDashboardResponse(common.ExpectedRequestId, common.ExpectedMessage, common.ExpectedStatusCode, expected)
	requireAIBase(t, actual.RequestId, actual.Message, actual.StatusCode)
	require.Equal(t, expected, actual.Dashboard)
}

func TestNewMultiAIDashboardInfosResponse(t *testing.T) {
	expected := []dtos.AIDashboardInfo{{Id: "d1", Ontology: common.AIOntologyHaystack, PanelCount: 2}, {Id: "d2", Ontology: common.AIOntologyEO66}}
	actual := NewMultiAIDashboardInfosResponse(common.ExpectedRequestId, common.ExpectedMessage, common.ExpectedStatusCode, 2, expected)
	requireAIBase(t, actual.RequestId, actual.Message, actual.StatusCode)
	require.Equal(t, int64(2), actual.TotalCount)
	require.Equal(t, expected, actual.Dashboards)
}

func TestNewAIPanelResponse(t *testing.T) {
	expected := dtos.AIPanel{Id: "p", Title: "t"}
	matched := 0
	actual := NewAIPanelResponse(common.ExpectedRequestId, common.ExpectedMessage, common.ExpectedStatusCode, "d", expected, &matched)
	requireAIBase(t, actual.RequestId, actual.Message, actual.StatusCode)
	require.Equal(t, "d", actual.DashboardId)
	require.Equal(t, expected, actual.Panel)
	require.Equal(t, &matched, actual.MatchedPoints)

	// zero matches is still reported; an unknown count is left out
	b, err := json.Marshal(actual)
	require.NoError(t, err)
	require.Contains(t, string(b), `"matchedPoints":0`)
	b, err = json.Marshal(NewAIPanelResponse(common.ExpectedRequestId, common.ExpectedMessage, common.ExpectedStatusCode, "d", expected, nil))
	require.NoError(t, err)
	require.NotContains(t, string(b), "matchedPoints")
}

func TestNewAIPointsSummaryResponse(t *testing.T) {
	expected := dtos.AIPointsSummary{Ontology: common.AIOntologyHaystack, TotalPoints: 1, TotalDevices: 1, Tags: []dtos.AITagSummary{{Tag: "temp", Points: 1}}}
	actual := NewAIPointsSummaryResponse(common.ExpectedRequestId, common.ExpectedMessage, common.ExpectedStatusCode, expected)
	requireAIBase(t, actual.RequestId, actual.Message, actual.StatusCode)
	require.Equal(t, expected, actual.Summary)
}

func TestNewAIPointsPageResponse(t *testing.T) {
	next := 2
	expected := dtos.AIPointsPage{Ontology: common.AIOntologyEO66, Tag: "temp", Total: 3, NextOffset: &next}
	actual := NewAIPointsPageResponse(common.ExpectedRequestId, common.ExpectedMessage, common.ExpectedStatusCode, expected)
	requireAIBase(t, actual.RequestId, actual.Message, actual.StatusCode)
	require.Equal(t, expected, actual.Page)
}
