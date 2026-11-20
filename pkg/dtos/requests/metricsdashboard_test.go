// Copyright (C) 2026 IOTech Ltd

package requests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricsDashboardRequests_NeedApiVersion(t *testing.T) {
	bodies := map[string]string{
		"add dashboard":    `{"dashboard":{"name":"d"}}`,
		"update dashboard": `{"name":"d"}`,
		"add panel":        `{"panel":{"title":"t","chartType":"line","tags":["hs4:t"]}}`,
		"update panel":     `{"title":"t"}`,
	}
	targets := map[string]func() json.Unmarshaler{
		"add dashboard":    func() json.Unmarshaler { return &AddMetricsDashboardRequest{} },
		"update dashboard": func() json.Unmarshaler { return &UpdateMetricsDashboardRequest{} },
		"add panel":        func() json.Unmarshaler { return &AddMetricsDashboardPanelRequest{} },
		"update panel":     func() json.Unmarshaler { return &UpdateMetricsDashboardPanelRequest{} },
	}
	for name, body := range bodies {
		t.Run(name, func(t *testing.T) {
			err := json.Unmarshal([]byte(body), targets[name]())
			require.Error(t, err)
			assert.Equal(t, errors.KindContractInvalid, errors.Kind(err))

			withVersion := `{"apiVersion":"v3",` + body[1:]
			assert.NoError(t, json.Unmarshal([]byte(withVersion), targets[name]()))
		})
	}
}

func TestUpdateMetricsDashboardRequest_NeedsAChange(t *testing.T) {
	var req UpdateMetricsDashboardRequest
	err := json.Unmarshal([]byte(`{"apiVersion":"v3"}`), &req)
	require.Error(t, err)
	assert.Equal(t, errors.KindContractInvalid, errors.Kind(err))
	assert.Contains(t, err.Error(), "pinned")

	for _, body := range []string{
		`{"apiVersion":"v3","name":"n"}`,
		`{"apiVersion":"v3","description":""}`,
		`{"apiVersion":"v3","range":{"last":"1h"}}`,
		`{"apiVersion":"v3","pinned":true}`,
		`{"apiVersion":"v3","pinned":false}`,
	} {
		assert.NoError(t, json.Unmarshal([]byte(body), &UpdateMetricsDashboardRequest{}), body)
	}
}

func TestUpdateMetricsDashboardRequest_PinnedAlone(t *testing.T) {
	for _, pinned := range []bool{true, false} {
		var req UpdateMetricsDashboardRequest
		require.NoError(t, json.Unmarshal([]byte(fmt.Sprintf(`{"apiVersion":"v3","pinned":%t}`, pinned)), &req))
		require.NotNil(t, req.Pinned)
		assert.Equal(t, pinned, *req.Pinned)
		assert.Nil(t, req.Name)
	}
}

func TestUpdateMetricsDashboardPanelRequest_NeedsAChange(t *testing.T) {
	var req UpdateMetricsDashboardPanelRequest
	err := json.Unmarshal([]byte(`{"apiVersion":"v3"}`), &req)
	require.Error(t, err)
	assert.Equal(t, errors.KindContractInvalid, errors.Kind(err))
	assert.Equal(t, "the update sets none of title, chartType, tags, points, devices and layout", err.Error())

	for _, body := range []string{
		`{"apiVersion":"v3","title":"t"}`,
		`{"apiVersion":"v3","chartType":"bar"}`,
		`{"apiVersion":"v3","tags":["hs4:t"]}`,
		`{"apiVersion":"v3","devices":[]}`,
		`{"apiVersion":"v3","points":[{"deviceName":"d","resourceName":"r"}]}`,
		`{"apiVersion":"v3","layout":{"x":0,"y":0,"w":6,"h":4}}`,
	} {
		assert.NoError(t, json.Unmarshal([]byte(body), &UpdateMetricsDashboardPanelRequest{}), body)
	}
}

func TestUpdateMetricsDashboardPanelRequest_EmptyDevicesIsSet(t *testing.T) {
	var req UpdateMetricsDashboardPanelRequest
	require.NoError(t, json.Unmarshal([]byte(`{"apiVersion":"v3","devices":[]}`), &req))
	require.NotNil(t, req.Devices)
	assert.Empty(t, *req.Devices)

	devices := []string{}
	b, err := json.Marshal(UpdateMetricsDashboardPanelRequest{Devices: &devices})
	require.NoError(t, err)
	assert.Contains(t, string(b), `"devices":[]`)
	assert.NotContains(t, string(b), `"title"`)
}

func TestAddMetricsDashboardPanelRequest_DropsId(t *testing.T) {
	var req AddMetricsDashboardPanelRequest
	body := `{"apiVersion":"v3","panel":{"id":"mine","title":"t","chartType":"line","tags":["hs4:t"]}}`
	require.NoError(t, json.Unmarshal([]byte(body), &req))
	assert.Empty(t, req.Panel.Id)
	assert.Equal(t, dtos.MetricsDashboardPanel{Title: "t", ChartType: "line", Tags: []string{"hs4:t"}}, req.Panel)
}

func TestMetricsDashboardRequests_CheckTheFieldsSent(t *testing.T) {
	refused := map[string][]string{
		"add dashboard": {
			`{"apiVersion":"v3","dashboard":{"name":" "}}`,
			`{"apiVersion":"v3","dashboard":{"name":"d","range":{"last":"7d"}}}`,
			`{"apiVersion":"v3","dashboard":{"name":"d","range":{"end":1000}}}`,
			`{"apiVersion":"v3","dashboard":{"name":"d","panels":[{"title":"t","chartType":"line","tags":[]}]}}`,
			`{"apiVersion":"v3","dashboard":{"name":"d","panels":[{"id":"p","title":"a","chartType":"line","tags":["hs4:t"]},{"id":"p","title":"b","chartType":"line","tags":["hs4:t"]}]}}`,
			`{"apiVersion":"v3","dashboard":{"name":"d","panels":[{"title":"t","chartType":"line","tags":["hs4:t"],"points":[{"deviceName":"d","resourceName":"r"}]}]}}`,
		},
		"update dashboard": {
			`{"apiVersion":"v3","name":""}`,
			`{"apiVersion":"v3","range":{}}`,
			`{"apiVersion":"v3","range":{"last":"7d"}}`,
			`{"apiVersion":"v3","range":{"start":1000}}`,
			`{"apiVersion":"v3","range":{"start":2000,"end":1000}}`,
		},
		"add panel": {
			`{"apiVersion":"v3","panel":{"chartType":"line","tags":["hs4:t"]}}`,
			`{"apiVersion":"v3","panel":{"title":"t","tags":["hs4:t"]}}`,
			`{"apiVersion":"v3","panel":{"title":"t","chartType":"line"}}`,
			`{"apiVersion":"v3","panel":{"title":"t","chartType":"line","tags":["no-colon"]}}`,
			`{"apiVersion":"v3","panel":{"title":"t","chartType":"line","tags":["hs4:t"],"layout":{"x":10,"y":0,"w":6,"h":4}}}`,
			`{"apiVersion":"v3","panel":{"title":"t","chartType":"line","points":[{"deviceName":"d"}]}}`,
			`{"apiVersion":"v3","panel":{"title":"t","chartType":"line","points":[{"deviceName":"d","resourceName":"r"}],"devices":["d"]}}`,
		},
		"update panel": {
			`{"apiVersion":"v3","title":" "}`,
			`{"apiVersion":"v3","chartType":""}`,
			`{"apiVersion":"v3","tags":[]}`,
			`{"apiVersion":"v3","tags":["hs4:"]}`,
			`{"apiVersion":"v3","layout":{"x":-1,"y":0,"w":6,"h":4}}`,
			`{"apiVersion":"v3","points":[]}`,
			`{"apiVersion":"v3","tags":["hs4:t"],"points":[{"deviceName":"d","resourceName":"r"}]}`,
			`{"apiVersion":"v3","points":[{"deviceName":" ","resourceName":"r"}]}`,
			`{"apiVersion":"v3","points":[{"deviceName":"d","resourceName":""}]}`,
			`{"apiVersion":"v3","points":[{"deviceName":"d","resourceName":"r"}],"devices":["d"]}`,
		},
	}
	accepted := map[string][]string{
		"add dashboard": {
			`{"apiVersion":"v3","dashboard":{"name":"d"}}`,
			`{"apiVersion":"v3","dashboard":{"name":"d","range":{"start":1000,"end":2000},"panels":[{"title":"t","chartType":"line","tags":["hs4:t"]}]}}`,
			`{"apiVersion":"v3","dashboard":{"name":"d","panels":[{"title":"t","chartType":"line","points":[{"deviceName":"d","resourceName":"r"}]}]}}`,
		},
		"add panel": {
			`{"apiVersion":"v3","panel":{"title":"t","chartType":"line","points":[{"deviceName":"d","resourceName":"r"}],"devices":[]}}`,
		},
		"update panel": {
			`{"apiVersion":"v3","layout":{"x":0,"y":0,"w":0,"h":0}}`,
			`{"apiVersion":"v3","tags":["hs4:t","hs4:t"]}`,
			`{"apiVersion":"v3","points":[{"deviceName":"d","resourceName":"r"}],"devices":[]}`,
			`{"apiVersion":"v3","tags":["hs4:t"],"devices":["d"]}`,
		},
	}
	targets := map[string]func() json.Unmarshaler{
		"add dashboard":    func() json.Unmarshaler { return &AddMetricsDashboardRequest{} },
		"update dashboard": func() json.Unmarshaler { return &UpdateMetricsDashboardRequest{} },
		"add panel":        func() json.Unmarshaler { return &AddMetricsDashboardPanelRequest{} },
		"update panel":     func() json.Unmarshaler { return &UpdateMetricsDashboardPanelRequest{} },
	}
	for name, bodies := range refused {
		for _, body := range bodies {
			err := json.Unmarshal([]byte(body), targets[name]())
			require.Error(t, err, body)
			assert.Equal(t, errors.KindContractInvalid, errors.Kind(err), body)
		}
	}
	for name, bodies := range accepted {
		for _, body := range bodies {
			assert.NoError(t, json.Unmarshal([]byte(body), targets[name]()), body)
		}
	}
}

func TestAddMetricsDashboardRequest_LeavesTheDashboardAsSent(t *testing.T) {
	var req AddMetricsDashboardRequest
	require.NoError(t, json.Unmarshal([]byte(`{"apiVersion":"v3","dashboard":{"name":"d"}}`), &req))
	assert.Equal(t, dtos.MetricsDashboardRange{}, req.Dashboard.Range)
	assert.Nil(t, req.Dashboard.Panels)
}

func TestUpdateMetricsDashboardPanelRequest_DedupesTags(t *testing.T) {
	var req UpdateMetricsDashboardPanelRequest
	require.NoError(t, json.Unmarshal([]byte(`{"apiVersion":"v3","tags":["hs4:t","hs4:u","hs4:t"]}`), &req))
	assert.Equal(t, []string{"hs4:t", "hs4:u"}, *req.Tags)
}

func TestUpdateMetricsDashboardPanelRequest_PointsTexts(t *testing.T) {
	tests := map[string]string{
		`{"apiVersion":"v3","points":[]}`: "at least one point is required",
		`{"apiVersion":"v3","tags":["hs4:t"],"points":[{"deviceName":"d","resourceName":"r"}]}`:          "the update sets both tags and points",
		`{"apiVersion":"v3","points":[{"deviceName":"d","resourceName":"r"},{"deviceName":"d"}]}`:        "point 2 needs a deviceName and a resourceName",
		`{"apiVersion":"v3","points":[{"deviceName":"d","resourceName":"r"}],"devices":["d"]}`:           "the update sets devices beside points",
		`{"apiVersion":"v3","points":[{"deviceName":"d","resourceName":"r"}],"tags":[],"devices":["d"]}`: "the update sets both tags and points",
	}
	for body, want := range tests {
		err := json.Unmarshal([]byte(body), &UpdateMetricsDashboardPanelRequest{})
		require.Error(t, err, body)
		assert.Equal(t, errors.KindContractInvalid, errors.Kind(err), body)
		assert.Equal(t, want, err.Error(), body)
	}
}

func TestUpdateMetricsDashboardPanelRequest_KeepsPointsAsSent(t *testing.T) {
	var req UpdateMetricsDashboardPanelRequest
	body := `{"apiVersion":"v3","points":[{"deviceName":"b","resourceName":"r"},{"deviceName":"a","resourceName":"r"},{"deviceName":"b","resourceName":"r"}]}`
	require.NoError(t, json.Unmarshal([]byte(body), &req))
	require.NotNil(t, req.Points)
	assert.Equal(t, []dtos.MetricsDashboardPanelPoint{
		{DeviceName: "b", ResourceName: "r"}, {DeviceName: "a", ResourceName: "r"}, {DeviceName: "b", ResourceName: "r"},
	}, *req.Points)
	assert.Nil(t, req.Tags)
}
