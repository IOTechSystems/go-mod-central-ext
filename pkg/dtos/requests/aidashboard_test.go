// Copyright (C) 2026 IOTech Ltd

package requests

import (
	"encoding/json"
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const aiPanelJSON = `{"title":"t","chartType":"line","tags":["temp"]}`

func TestAddAIDashboardRequest_Unmarshal(t *testing.T) {
	t.Run("ok, range defaults to last 24h", func(t *testing.T) {
		var req AddAIDashboardRequest
		err := json.Unmarshal([]byte(`{"apiVersion":"v3","dashboard":{"name":"a","panels":[`+aiPanelJSON+`,`+aiPanelJSON+`]}}`), &req)
		require.NoError(t, err)
		assert.Equal(t, dtos.AIRange{Last: "24h"}, req.Dashboard.Range)
		assert.Len(t, req.Dashboard.Panels, 2)
	})
	t.Run("ids and timestamps dropped, rules left to the server", func(t *testing.T) {
		var req AddAIDashboardRequest
		body := `{"apiVersion":"v3","dashboard":{"id":"d","name":"a","created":5,"modified":6,"range":{"last":"1h"},` +
			`"panels":[{"id":"p","title":"","chartType":"pie","layout":{"x":0,"y":0,"w":0,"h":0}}]}}`
		require.NoError(t, json.Unmarshal([]byte(body), &req))
		assert.Empty(t, req.Dashboard.Id)
		assert.Zero(t, req.Dashboard.Created)
		assert.Zero(t, req.Dashboard.Modified)
		assert.Empty(t, req.Dashboard.Panels[0].Id)
		assert.Equal(t, dtos.AIRange{Last: "1h"}, req.Dashboard.Range)
		// ontology is an A1 rule, so a create without it still decodes
		assert.Empty(t, req.Dashboard.Ontology)
	})
	t.Run("ontology kept", func(t *testing.T) {
		var req AddAIDashboardRequest
		require.NoError(t, json.Unmarshal([]byte(`{"apiVersion":"v3","dashboard":{"name":"a","ontology":"eo66"}}`), &req))
		assert.Equal(t, common.AIOntologyEO66, req.Dashboard.Ontology)
	})
	t.Run("blank name", func(t *testing.T) {
		var req AddAIDashboardRequest
		err := json.Unmarshal([]byte(`{"apiVersion":"v3","dashboard":{"name":" "}}`), &req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Name")
	})
	t.Run("no apiVersion", func(t *testing.T) {
		var req AddAIDashboardRequest
		err := json.Unmarshal([]byte(`{"dashboard":{"name":"a","panels":[`+aiPanelJSON+`,`+aiPanelJSON+`]}}`), &req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ApiVersion")
	})
}

func TestUpdateAIPanelRequest_Unmarshal(t *testing.T) {
	t.Run("only devices set", func(t *testing.T) {
		var req UpdateAIPanelRequest
		require.NoError(t, json.Unmarshal([]byte(`{"apiVersion":"v3","devices":["d1"]}`), &req))
		want := UpdateAIPanelRequest{BaseRequest: req.BaseRequest, Devices: &[]string{"d1"}}
		assert.Equal(t, want, req)
	})
	t.Run("only options set, unknown key and wrong type kept", func(t *testing.T) {
		var req UpdateAIPanelRequest
		require.NoError(t, json.Unmarshal([]byte(`{"apiVersion":"v3","options":{"rightAxisTags":"x","colour":1}}`), &req))
		require.NotNil(t, req.Options)
		assert.Equal(t, map[string]any{"rightAxisTags": "x", "colour": float64(1)}, *req.Options)
	})
	t.Run("empty options is a change", func(t *testing.T) {
		var req UpdateAIPanelRequest
		require.NoError(t, json.Unmarshal([]byte(`{"apiVersion":"v3","options":{}}`), &req))
		require.NotNil(t, req.Options)
		assert.Empty(t, *req.Options)
	})
	tests := []struct {
		name    string
		body    string
		wantErr string
	}{
		{"no change", `{"apiVersion":"v3"}`, "at least one"},
		{"range with followDashboardRange", `{"apiVersion":"v3","followDashboardRange":true,"range":{"last":"1h"}}`, "range"},
		{"no apiVersion", `{"devices":["d1"]}`, "ApiVersion"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req UpdateAIPanelRequest
			err := json.Unmarshal([]byte(tt.body), &req)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
	// the dashboard rules are the server's (AIDashboard.Validate), so a well-formed body decodes
	for name, body := range map[string]string{
		"blank title":         `{"apiVersion":"v3","title":""}`,
		"unknown chartType":   `{"apiVersion":"v3","chartType":"pie"}`,
		"zero layout":         `{"apiVersion":"v3","layout":{"x":0,"y":0,"w":0,"h":0}}`,
		"layout off the grid": `{"apiVersion":"v3","layout":{"x":8,"y":0,"w":6,"h":4}}`,
		"fixed range":         `{"apiVersion":"v3","range":{"start":1,"end":2}}`,
	} {
		t.Run(name+" decodes", func(t *testing.T) {
			var req UpdateAIPanelRequest
			assert.NoError(t, json.Unmarshal([]byte(body), &req))
		})
	}
}

func TestUpdateAIDashboardRequest_Unmarshal(t *testing.T) {
	var req UpdateAIDashboardRequest
	err := json.Unmarshal([]byte(`{"apiVersion":"v3"}`), &req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "at least one")

	req = UpdateAIDashboardRequest{}
	require.Error(t, json.Unmarshal([]byte(`{"name":"a"}`), &req))

	req = UpdateAIDashboardRequest{}
	require.NoError(t, json.Unmarshal([]byte(`{"apiVersion":"v3","name":"a"}`), &req))
	require.NotNil(t, req.Name)
	assert.Equal(t, "a", *req.Name)

	req = UpdateAIDashboardRequest{}
	err = json.Unmarshal([]byte(`{"apiVersion":"v3","name":" "}`), &req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Name")

	// range rules are the server's, so a bad range decodes
	req = UpdateAIDashboardRequest{}
	require.NoError(t, json.Unmarshal([]byte(`{"apiVersion":"v3","range":{"last":"2x"}}`), &req))
}

func TestUpdateAILayoutsRequest_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr string
	}{
		{"empty layouts", `{"apiVersion":"v3","layouts":{}}`, "layouts"},
		{"no layouts", `{"apiVersion":"v3"}`, "layouts"},
		{"no apiVersion", `{"layouts":{"p":{"x":0,"y":0,"w":6,"h":4}}}`, "ApiVersion"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req UpdateAILayoutsRequest
			err := json.Unmarshal([]byte(tt.body), &req)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}

	// a zero or off-grid layout decodes: the server refuses it
	var req UpdateAILayoutsRequest
	require.NoError(t, json.Unmarshal([]byte(`{"apiVersion":"v3","layouts":{"p":{"x":0,"y":0,"w":0,"h":0},"q":{"x":8,"y":0,"w":6,"h":4}}}`), &req))
	assert.Len(t, req.Layouts, 2)
}

func TestAddAIPanelRequest_Unmarshal(t *testing.T) {
	var req AddAIPanelRequest
	require.NoError(t, json.Unmarshal([]byte(`{"apiVersion":"v3","panel":{"id":"x","title":"t","chartType":"line","tags":["temp"]}}`), &req))
	assert.Empty(t, req.Panel.Id)
	assert.Equal(t, []string{"temp"}, req.Panel.Tags)

	req = AddAIPanelRequest{}
	require.NoError(t, json.Unmarshal([]byte(`{"apiVersion":"v3","panel":{"title":"t","chartType":"bar","tags":["temp"],"options":{"stateLabels":[1]}}}`), &req))
	assert.Equal(t, map[string]any{"stateLabels": []any{float64(1)}}, req.Panel.Options)

	// panel rules are the server's, so an empty panel decodes
	req = AddAIPanelRequest{}
	require.NoError(t, json.Unmarshal([]byte(`{"apiVersion":"v3","panel":{"title":"","chartType":"pie"}}`), &req))

	req = AddAIPanelRequest{}
	err := json.Unmarshal([]byte(`{"panel":{"title":"t","chartType":"line","tags":["temp"]}}`), &req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ApiVersion")
}
