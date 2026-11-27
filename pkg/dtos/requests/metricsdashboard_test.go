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
		"add panel":        `{"panel":{"title":"t","chartType":"line","select":[{"tag":"hs4:t"}]}}`,
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
	assert.Equal(t, "the update sets none of title, chartType, select, layout, valueMapping, rightAxisUnits and resourceLabels", err.Error())

	for _, body := range []string{
		`{"apiVersion":"v3","title":"t"}`,
		`{"apiVersion":"v3","chartType":"bar"}`,
		`{"apiVersion":"v3","select":[{"tag":"hs4:t"}]}`,
		`{"apiVersion":"v3","layout":{"x":0,"y":0,"w":6,"h":4}}`,
		`{"apiVersion":"v3","valueMapping":{"Run Status":{"0":"Stopped","":"Empty","1":""}," ":{}}}`,
		`{"apiVersion":"v3","rightAxisUnits":["%",""," ","%"]}`,
		`{"apiVersion":"v3","valueMapping":{}}`,
		`{"apiVersion":"v3","rightAxisUnits":[]}`,
		`{"apiVersion":"v3","resourceLabels":{}}`,
	} {
		assert.NoError(t, json.Unmarshal([]byte(body), &UpdateMetricsDashboardPanelRequest{}), body)
	}
}

func TestUpdateMetricsDashboardPanelRequest_SetsAndClears(t *testing.T) {
	var set UpdateMetricsDashboardPanelRequest
	require.NoError(t, json.Unmarshal([]byte(`{"apiVersion":"v3","valueMapping":{"Run Status":{"0":"Stopped"}},"rightAxisUnits":["%"]}`), &set))
	require.NotNil(t, set.ValueMapping)
	require.NotNil(t, set.RightAxisUnits)
	assert.Equal(t, map[string]map[string]string{"Run Status": {"0": "Stopped"}}, *set.ValueMapping)
	assert.Equal(t, []string{"%"}, *set.RightAxisUnits)

	var cleared UpdateMetricsDashboardPanelRequest
	require.NoError(t, json.Unmarshal([]byte(`{"apiVersion":"v3","valueMapping":{},"rightAxisUnits":[]}`), &cleared))
	require.NotNil(t, cleared.ValueMapping, "{} clears, it is not absent")
	require.NotNil(t, cleared.RightAxisUnits, "[] clears, it is not absent")
	assert.Empty(t, *cleared.ValueMapping)
	assert.Empty(t, *cleared.RightAxisUnits)
	out, err := json.Marshal(cleared)
	require.NoError(t, err)
	assert.Contains(t, string(out), `"valueMapping":{}`)
	assert.Contains(t, string(out), `"rightAxisUnits":[]`)

	var untouched UpdateMetricsDashboardPanelRequest
	require.NoError(t, json.Unmarshal([]byte(`{"apiVersion":"v3","title":"t"}`), &untouched))
	assert.Nil(t, untouched.ValueMapping)
	assert.Nil(t, untouched.RightAxisUnits)
	out, err = json.Marshal(untouched)
	require.NoError(t, err)
	assert.NotContains(t, string(out), `"valueMapping"`)
	assert.NotContains(t, string(out), `"rightAxisUnits"`)
}

func TestUpdateMetricsDashboardPanelRequest_OldFieldsAreNoChange(t *testing.T) {
	for _, body := range []string{
		`{"apiVersion":"v3","tags":["hs4:t"]}`,
		`{"apiVersion":"v3","devices":["d"]}`,
		`{"apiVersion":"v3","points":[{"deviceName":"d","resourceName":"r"}]}`,
	} {
		err := json.Unmarshal([]byte(body), &UpdateMetricsDashboardPanelRequest{})
		require.Error(t, err, body)
		assert.Equal(t, "the update sets none of title, chartType, select, layout, valueMapping, rightAxisUnits and resourceLabels", err.Error(), body)
	}
}

func TestAddMetricsDashboardPanelRequest_DropsId(t *testing.T) {
	var req AddMetricsDashboardPanelRequest
	body := `{"apiVersion":"v3","panel":{"id":"mine","title":"t","chartType":"line","select":[{"tag":"hs4:t"}]}}`
	require.NoError(t, json.Unmarshal([]byte(body), &req))
	assert.Empty(t, req.Panel.Id)
	tag := "hs4:t"
	assert.Equal(t, dtos.MetricsDashboardPanel{Title: "t", ChartType: "line", Select: []dtos.MetricsDashboardSelector{{Tag: &tag}}}, req.Panel)
}

func TestMetricsDashboardRequests_CheckTheFieldsSent(t *testing.T) {
	refused := map[string][]string{
		"add dashboard": {
			`{"apiVersion":"v3","dashboard":{"name":" "}}`,
			`{"apiVersion":"v3","dashboard":{"name":"d","range":{"last":"7d"}}}`,
			`{"apiVersion":"v3","dashboard":{"name":"d","range":{"end":1000}}}`,
			`{"apiVersion":"v3","dashboard":{"name":"d","panels":[{"title":"t","chartType":"line","select":[]}]}}`,
			`{"apiVersion":"v3","dashboard":{"name":"d","panels":[{"title":"t","chartType":"line","tags":["hs4:t"]}]}}`,
			`{"apiVersion":"v3","dashboard":{"name":"d","panels":[{"id":"p","title":"a","chartType":"line","select":[{"tag":"hs4:t"}]},{"id":"p","title":"b","chartType":"line","select":[{"tag":"hs4:t"}]}]}}`,
		},
		"update dashboard": {
			`{"apiVersion":"v3","name":""}`,
			`{"apiVersion":"v3","range":{}}`,
			`{"apiVersion":"v3","range":{"last":"7d"}}`,
			`{"apiVersion":"v3","range":{"start":1000}}`,
			`{"apiVersion":"v3","range":{"start":2000,"end":1000}}`,
		},
		"add panel": {
			`{"apiVersion":"v3","panel":{"chartType":"line","select":[{"tag":"hs4:t"}]}}`,
			`{"apiVersion":"v3","panel":{"title":"t","select":[{"tag":"hs4:t"}]}}`,
			`{"apiVersion":"v3","panel":{"title":"t","chartType":"line"}}`,
			`{"apiVersion":"v3","panel":{"title":"t","chartType":"line","select":[{"tag":"no-colon"}]}}`,
			`{"apiVersion":"v3","panel":{"title":"t","chartType":"line","select":[{"deviceName":"d"}]}}`,
			`{"apiVersion":"v3","panel":{"title":"t","chartType":"line","select":[{"deviceName":"d","profileName":"p","tag":"hs4:t"}]}}`,
			`{"apiVersion":"v3","panel":{"title":"t","chartType":"line","select":[{"deviceName":"","tag":"hs4:t"}]}}`,
			`{"apiVersion":"v3","panel":{"title":"t","chartType":"line","select":[{"profileName":" ","resourceName":"r"}]}}`,
			`{"apiVersion":"v3","panel":{"title":"t","chartType":"line","select":[{"resourceName":""}]}}`,
			`{"apiVersion":"v3","panel":{"title":"t","chartType":"line","select":[{"tag":"hs4:t"}],"layout":{"x":10,"y":0,"w":6,"h":4}}}`,
		},
		"update panel": {
			`{"apiVersion":"v3","title":" "}`,
			`{"apiVersion":"v3","chartType":""}`,
			`{"apiVersion":"v3","select":[]}`,
			`{"apiVersion":"v3","select":[{"tag":"hs4:"}]}`,
			`{"apiVersion":"v3","select":[{"deviceName":"","tag":"hs4:t"}]}`,
			`{"apiVersion":"v3","select":[{"profileName":"p"}]}`,
			`{"apiVersion":"v3","select":[{"deviceName":"d","profileName":"p","resourceName":"r"}]}`,
			`{"apiVersion":"v3","layout":{"x":-1,"y":0,"w":6,"h":4}}`,
		},
	}
	accepted := map[string][]string{
		"add dashboard": {
			`{"apiVersion":"v3","dashboard":{"name":"d"}}`,
			`{"apiVersion":"v3","dashboard":{"name":"d","range":{"start":1000,"end":2000},"panels":[{"title":"t","chartType":"line","select":[{"tag":"hs4:t"},{"profileName":"Johnson-BMS","resourceName":"Lobby IAQ Sensor"}]}]}}`,
		},
		"add panel": {
			`{"apiVersion":"v3","panel":{"title":"t","chartType":"line","select":[{"deviceName":"J","resourceName":"Lobby IAQ Sensor"}]}}`,
			`{"apiVersion":"v3","panel":{"title":"t","chartType":"line","select":[{"resourceName":"Lobby IAQ Sensor","tag":"hs4:co2"}]}}`,
			`{"apiVersion":"v3","panel":{"title":"t","chartType":"line","select":[{"tag":"hs4:co2"},{"tag":"hs4:co2"}]}}`,
		},
		"update panel": {
			`{"apiVersion":"v3","layout":{"x":0,"y":0,"w":0,"h":0}}`,
			`{"apiVersion":"v3","select":[{"tag":"hs4:t"},{"tag":"hs4:t"}]}`,
			`{"apiVersion":"v3","select":[{"deviceName":"d","tag":"hs4:t"}]}`,
			`{"apiVersion":"v3","select":[{"profileName":"p","resourceName":"r"}]}`,
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

func TestUpdateMetricsDashboardPanelRequest_SelectTexts(t *testing.T) {
	tests := map[string]string{
		`{"apiVersion":"v3","select":[]}`:                                                        "at least one selector is required",
		`{"apiVersion":"v3","select":[{"tag":"hs4:t"},{"deviceName":"","tag":"hs4:t"}]}`:         "selector 2 has a blank deviceName",
		`{"apiVersion":"v3","select":[{"profileName":"p"}]}`:                                     "selector 1 needs a resourceName or a tag",
		`{"apiVersion":"v3","select":[{"deviceName":"d","profileName":"p","resourceName":"r"}]}`: "selector 1 sets both deviceName and profileName",
		`{"apiVersion":"v3","select":[{"tag":"x"}]}`:                                             `selector 1: tag "x" is not <key>:<value>`,
	}
	for body, want := range tests {
		err := json.Unmarshal([]byte(body), &UpdateMetricsDashboardPanelRequest{})
		require.Error(t, err, body)
		assert.Equal(t, errors.KindContractInvalid, errors.Kind(err), body)
		assert.Equal(t, want, err.Error(), body)
	}
}

func TestUpdateMetricsDashboardPanelRequest_KeepsSelectAsSent(t *testing.T) {
	var req UpdateMetricsDashboardPanelRequest
	body := `{"apiVersion":"v3","select":[{"deviceName":"b","tag":"hs4:t"},{"resourceName":"r"},{"deviceName":"b","tag":"hs4:t"}]}`
	require.NoError(t, json.Unmarshal([]byte(body), &req))
	require.NotNil(t, req.Select)
	b, tag, r := "b", "hs4:t", "r"
	assert.Equal(t, []dtos.MetricsDashboardSelector{
		{DeviceName: &b, Tag: &tag}, {ResourceName: &r}, {DeviceName: &b, Tag: &tag},
	}, *req.Select)
	assert.Nil(t, req.Title)

	out, err := json.Marshal(UpdateMetricsDashboardPanelRequest{Select: req.Select})
	require.NoError(t, err)
	assert.Contains(t, string(out), `"select":[{"deviceName":"b","tag":"hs4:t"},{"resourceName":"r"},{"deviceName":"b","tag":"hs4:t"}]`)
	assert.NotContains(t, string(out), `"title"`)
}
