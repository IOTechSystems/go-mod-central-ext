// Copyright (C) 2026 IOTech Ltd

package dtos

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptr(s string) *string { return &s }

func validMetricsDashboardPanel() MetricsDashboardPanel {
	return MetricsDashboardPanel{Title: "Temps", ChartType: "line", Select: []MetricsDashboardSelector{{Tag: ptr("hs4:zone-air-temp")}}}
}

func selectorsMetricsDashboardPanel() MetricsDashboardPanel {
	return MetricsDashboardPanel{Title: "Temps", ChartType: "line", Select: []MetricsDashboardSelector{
		{Tag: ptr("hs4:co2")},
		{ProfileName: ptr("Johnson-BMS"), ResourceName: ptr("Lobby IAQ Sensor")},
	}}
}

func validMetricsDashboard() MetricsDashboard {
	return MetricsDashboard{Name: "temps", Range: MetricsDashboardRange{Last: "5m"}, Panels: []MetricsDashboardPanel{validMetricsDashboardPanel()}}
}

func TestMetricsDashboardValidate_Refusals(t *testing.T) {
	withPanel := func(edit func(p *MetricsDashboardPanel)) func(d *MetricsDashboard) {
		return func(d *MetricsDashboard) { edit(&d.Panels[0]) }
	}
	tests := []struct {
		name string
		edit func(d *MetricsDashboard)
	}{
		{"empty name", func(d *MetricsDashboard) { d.Name = "" }},
		{"blank name", func(d *MetricsDashboard) { d.Name = "  " }},
		{"neither last nor start/end", func(d *MetricsDashboard) { d.Range = MetricsDashboardRange{} }},
		{"both last and start/end", func(d *MetricsDashboard) { d.Range = MetricsDashboardRange{Last: "5m", Start: 1, End: 2} }},
		{"last outside choices", func(d *MetricsDashboard) { d.Range = MetricsDashboardRange{Last: "7d"} }},
		{"start equals end", func(d *MetricsDashboard) { d.Range = MetricsDashboardRange{Start: 5, End: 5} }},
		{"start after end", func(d *MetricsDashboard) { d.Range = MetricsDashboardRange{Start: 6, End: 5} }},
		{"start without end", func(d *MetricsDashboard) { d.Range = MetricsDashboardRange{Start: 6} }},
		{"negative start without end", func(d *MetricsDashboard) { d.Range = MetricsDashboardRange{Start: -1000} }},
		{"end without start", func(d *MetricsDashboard) { d.Range = MetricsDashboardRange{End: 1000} }},
		{"panel without title", withPanel(func(p *MetricsDashboardPanel) { p.Title = " " })},
		{"blank chartType", withPanel(func(p *MetricsDashboardPanel) { p.ChartType = " " })},
		{"nil select", withPanel(func(p *MetricsDashboardPanel) { p.Select = nil })},
		{"empty select", withPanel(func(p *MetricsDashboardPanel) { p.Select = []MetricsDashboardSelector{} })},
		{"empty selector", withPanel(func(p *MetricsDashboardPanel) { p.Select[0] = MetricsDashboardSelector{} })},
		{"only a deviceName", withPanel(func(p *MetricsDashboardPanel) { p.Select[0] = MetricsDashboardSelector{DeviceName: ptr("J")} })},
		{"only a profileName", withPanel(func(p *MetricsDashboardPanel) { p.Select[0] = MetricsDashboardSelector{ProfileName: ptr("P")} })},
		{"both deviceName and profileName", withPanel(func(p *MetricsDashboardPanel) {
			p.Select[0].DeviceName, p.Select[0].ProfileName = ptr("J"), ptr("P")
		})},
		{"empty tag", withPanel(func(p *MetricsDashboardPanel) { p.Select[0].Tag = ptr("") })},
		{"blank resourceName", withPanel(func(p *MetricsDashboardPanel) { p.Select[0].ResourceName = ptr(" ") })},
		{"empty deviceName", withPanel(func(p *MetricsDashboardPanel) { p.Select[0].DeviceName = ptr("") })},
		{"blank profileName", withPanel(func(p *MetricsDashboardPanel) { p.Select[0].ProfileName = ptr("  ") })},
		{"tag without colon", withPanel(func(p *MetricsDashboardPanel) { p.Select[0].Tag = ptr("zone-air-temp") })},
		{"blank key", withPanel(func(p *MetricsDashboardPanel) { p.Select[0].Tag = ptr(" :zone-air-temp") })},
		{"blank value", withPanel(func(p *MetricsDashboardPanel) { p.Select[0].Tag = ptr("hs4: ") })},
		{"bad tag beside a resourceName", withPanel(func(p *MetricsDashboardPanel) {
			*p = selectorsMetricsDashboardPanel()
			p.Select[1].Tag = ptr("co2")
		})},
		{"duplicate panel id", func(d *MetricsDashboard) {
			p := validMetricsDashboardPanel()
			p.Id = "a"
			d.Panels = []MetricsDashboardPanel{p, p}
		}},
		{"w <= 0", withPanel(func(p *MetricsDashboardPanel) { p.Layout = MetricsDashboardLayout{X: 0, Y: 0, W: 0, H: 4} })},
		{"h <= 0", withPanel(func(p *MetricsDashboardPanel) { p.Layout = MetricsDashboardLayout{X: 0, Y: 0, W: 6, H: 0} })},
		{"x < 0", withPanel(func(p *MetricsDashboardPanel) { p.Layout = MetricsDashboardLayout{X: -1, Y: 0, W: 6, H: 4} })},
		{"y < 0", withPanel(func(p *MetricsDashboardPanel) { p.Layout = MetricsDashboardLayout{X: 0, Y: -1, W: 6, H: 4} })},
		{"x + w > 12", withPanel(func(p *MetricsDashboardPanel) { p.Layout = MetricsDashboardLayout{X: 7, Y: 0, W: 6, H: 4} })},
		{"x + w overflows", withPanel(func(p *MetricsDashboardPanel) { p.Layout = MetricsDashboardLayout{X: math.MaxInt, Y: 0, W: 1, H: 1} })},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := validMetricsDashboard()
			tt.edit(&d)
			err := d.Validate()
			require.Error(t, err)
			assert.Equal(t, errors.KindContractInvalid, errors.Kind(err))
		})
	}
}

func TestMetricsDashboardValidate_Accepts(t *testing.T) {
	tests := []struct {
		name string
		edit func(d *MetricsDashboard)
	}{
		{"fixed range", func(d *MetricsDashboard) { d.Range = MetricsDashboardRange{Start: 1, End: 2} }},
		{"selectors of two keys in one panel", func(d *MetricsDashboard) {
			d.Panels[0].Select = append(d.Panels[0].Select, MetricsDashboardSelector{Tag: ptr("eo66PointType:temp")})
		}},
		{"a tag on a device", func(d *MetricsDashboard) {
			d.Panels[0].Select = []MetricsDashboardSelector{{DeviceName: ptr("AHU-1"), Tag: ptr("hs4:zone-air-temp")}}
		}},
		{"a tag on a profile", func(d *MetricsDashboard) {
			d.Panels[0].Select = []MetricsDashboardSelector{{ProfileName: ptr("AHU"), Tag: ptr("hs4:zone-air-temp")}}
		}},
		{"a resourceName on a device", func(d *MetricsDashboard) {
			d.Panels[0].Select = []MetricsDashboardSelector{{DeviceName: ptr("S"), ResourceName: ptr("Lobby Temperature")}}
		}},
		{"a resourceName on a profile", func(d *MetricsDashboard) {
			d.Panels[0].Select = []MetricsDashboardSelector{{ProfileName: ptr("Johnson-BMS"), ResourceName: ptr("Lobby IAQ Sensor")}}
		}},
		{"a resourceName on every device", func(d *MetricsDashboard) {
			d.Panels[0].Select = []MetricsDashboardSelector{{ResourceName: ptr("Lobby IAQ Sensor")}}
		}},
		{"a resourceName and a tag", func(d *MetricsDashboard) {
			d.Panels[0].Select = []MetricsDashboardSelector{{ResourceName: ptr("Temp"), Tag: ptr("hs4:zone-air-temp")}}
		}},
		{"every field but profileName", func(d *MetricsDashboard) {
			d.Panels[0].Select = []MetricsDashboardSelector{{DeviceName: ptr("S"), ResourceName: ptr("Temp"), Tag: ptr("hs4:zone-air-temp")}}
		}},
		{"a tag of any key", func(d *MetricsDashboard) {
			d.Panels[0].Select = []MetricsDashboardSelector{{Tag: ptr("site:north")}}
		}},
		{"a tag value holding a colon", func(d *MetricsDashboard) {
			d.Panels[0].Select = []MetricsDashboardSelector{{Tag: ptr("site:Zone:1")}}
		}},
		{"any non-blank chartType", func(d *MetricsDashboard) { d.Panels[0].ChartType = "no-such-type" }},
		{"zero layout", func(d *MetricsDashboard) { d.Panels[0].Layout = MetricsDashboardLayout{} }},
		{"layout at the right edge", func(d *MetricsDashboard) { d.Panels[0].Layout = MetricsDashboardLayout{X: 6, Y: 0, W: 6, H: 4} }},
		{"no panels", func(d *MetricsDashboard) { d.Panels = nil }},
		{"100 panels", func(d *MetricsDashboard) {
			d.Panels = make([]MetricsDashboardPanel, 100)
			for i := range d.Panels {
				d.Panels[i] = validMetricsDashboardPanel()
			}
		}},
		{"a rule and a profile's resource", func(d *MetricsDashboard) { d.Panels[0] = selectorsMetricsDashboardPanel() }},
		{"panels without ids", func(d *MetricsDashboard) {
			d.Panels = []MetricsDashboardPanel{validMetricsDashboardPanel(), validMetricsDashboardPanel()}
		}},
		{"valueMapping with blank entries on a gauge", func(d *MetricsDashboard) {
			d.Panels[0].ChartType = "gauge"
			d.Panels[0].ValueMapping = map[string]map[string]string{"Run Status": {"0": "Stopped", "": "Empty", "1": ""}, " ": {"x": "y"}}
		}},
		{"rightAxisUnits with blank and repeated units on a stat", func(d *MetricsDashboard) {
			d.Panels[0].ChartType = "stat"
			d.Panels[0].RightAxisUnits = []string{"%", "", " ", "%"}
		}},
	}
	for _, c := range common.MetricsDashboardLastChoices {
		tests = append(tests, struct {
			name string
			edit func(d *MetricsDashboard)
		}{fmt.Sprintf("last %s", c.Last), func(d *MetricsDashboard) { d.Range = MetricsDashboardRange{Last: c.Last} }})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := validMetricsDashboard()
			tt.edit(&d)
			assert.NoError(t, d.Validate())
		})
	}
}

func TestMetricsDashboardValidate_KeepsSelectorsAsSent(t *testing.T) {
	d := validMetricsDashboard()
	d.Panels[0].Select = []MetricsDashboardSelector{
		{Tag: ptr("hs4:b")}, {DeviceName: ptr(" J "), ResourceName: ptr(" Temp ")}, {Tag: ptr("hs4:b")},
		{DeviceName: ptr(" J "), ResourceName: ptr(" Temp ")},
	}
	want := slices.Clone(d.Panels[0].Select)
	require.NoError(t, d.Validate())
	assert.Equal(t, want, d.Panels[0].Select)
}

func TestMetricsDashboard_DecodesUnknownField(t *testing.T) {
	doc := `{"name":"temps","range":{"last":"5m"},"future":1,
		"panels":[{"title":"Temps","chartType":"line","select":[{"tag":"hs4:t"}],"rightAxis":true,"layout":{"x":0,"y":0,"w":6,"h":4}}]}`
	var d MetricsDashboard
	require.NoError(t, json.Unmarshal([]byte(doc), &d))
	assert.Equal(t, "temps", d.Name)
	assert.Equal(t, MetricsDashboardLayout{W: 6, H: 4}, d.Panels[0].Layout)
	assert.NoError(t, d.Validate())
}

func TestMetricsDashboard_Pinned(t *testing.T) {
	for _, pinned := range []bool{true, false} {
		b, err := json.Marshal(MetricsDashboard{Name: "temps", Pinned: pinned})
		require.NoError(t, err)
		assert.Contains(t, string(b), fmt.Sprintf(`"pinned":%t`, pinned))
		var d MetricsDashboard
		require.NoError(t, json.Unmarshal(b, &d))
		assert.Equal(t, pinned, d.Pinned)

		b, err = json.Marshal(MetricsDashboardInfo{Name: "temps", Pinned: pinned})
		require.NoError(t, err)
		assert.Contains(t, string(b), fmt.Sprintf(`"pinned":%t`, pinned))
		var info MetricsDashboardInfo
		require.NoError(t, json.Unmarshal(b, &info))
		assert.Equal(t, pinned, info.Pinned)
	}

	var d MetricsDashboard
	require.NoError(t, json.Unmarshal([]byte(`{"name":"temps"}`), &d))
	assert.False(t, d.Pinned)
	var info MetricsDashboardInfo
	require.NoError(t, json.Unmarshal([]byte(`{"name":"temps"}`), &info))
	assert.False(t, info.Pinned)
}

func TestMetricsDashboardValidate_EmptyPanelsIsArray(t *testing.T) {
	d := MetricsDashboard{Name: "temps", Range: MetricsDashboardRange{Last: "5m"}}
	require.NoError(t, d.Validate())
	b, err := json.Marshal(d)
	require.NoError(t, err)
	assert.Contains(t, string(b), `"panels":[]`)
}

func TestMetricsDashboardPanelValidate_NamesThePanel(t *testing.T) {
	p := MetricsDashboardPanel{Title: "Zone temp", ChartType: "line", Select: []MetricsDashboardSelector{{Tag: ptr("hs4:t")}}, Layout: MetricsDashboardLayout{X: 10, W: 6, H: 4}}
	err := p.Validate()
	require.Error(t, err)
	assert.Equal(t, errors.KindContractInvalid, errors.Kind(err))
	assert.Contains(t, err.Error(), "panel Zone temp")
	assert.Contains(t, err.Error(), "outside the 12-column grid")
}

func TestMetricsDashboardPanelValidate_Texts(t *testing.T) {
	tests := []struct {
		name string
		edit func(p *MetricsDashboardPanel)
		want string
	}{
		{"no selectors", func(p *MetricsDashboardPanel) { p.Select = nil }, "panel Temps needs at least one selector"},
		{"neither resourceName nor tag", func(p *MetricsDashboardPanel) { p.Select[1].ResourceName = nil },
			"selector 2 of panel Temps needs a resourceName or a tag"},
		{"deviceName and profileName", func(p *MetricsDashboardPanel) { p.Select[1].DeviceName = ptr("J") },
			"selector 2 of panel Temps sets both deviceName and profileName"},
		{"blank deviceName", func(p *MetricsDashboardPanel) { p.Select[0].DeviceName = ptr("") }, "selector 1 of panel Temps has a blank deviceName"},
		{"blank profileName", func(p *MetricsDashboardPanel) { p.Select[1].ProfileName = ptr(" ") }, "selector 2 of panel Temps has a blank profileName"},
		{"blank resourceName", func(p *MetricsDashboardPanel) { p.Select[1].ResourceName = ptr("  ") }, "selector 2 of panel Temps has a blank resourceName"},
		{"blank tag", func(p *MetricsDashboardPanel) { p.Select[0].Tag = ptr("") }, "selector 1 of panel Temps has a blank tag"},
		{"bad tag", func(p *MetricsDashboardPanel) { p.Select[0].Tag = ptr("Lobby IAQ Sensor") },
			`selector 1 of panel Temps: tag "Lobby IAQ Sensor" is not <key>:<value>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := selectorsMetricsDashboardPanel()
			tt.edit(&p)
			err := p.Validate()
			require.Error(t, err)
			assert.Equal(t, errors.KindContractInvalid, errors.Kind(err))
			assert.Equal(t, tt.want, err.Error())
		})
	}
}

func TestValidateMetricsDashboardSelect_Texts(t *testing.T) {
	tests := map[string]struct {
		selectors []MetricsDashboardSelector
		want      string
	}{
		"nil":   {nil, "at least one selector is required"},
		"empty": {[]MetricsDashboardSelector{}, "at least one selector is required"},
		"neither resourceName nor tag": {[]MetricsDashboardSelector{{Tag: ptr("hs4:t")}, {DeviceName: ptr("J")}},
			"selector 2 needs a resourceName or a tag"},
		"deviceName and profileName": {[]MetricsDashboardSelector{{DeviceName: ptr("J"), ProfileName: ptr("P"), Tag: ptr("hs4:t")}},
			"selector 1 sets both deviceName and profileName"},
		"blank deviceName": {[]MetricsDashboardSelector{{Tag: ptr("hs4:t")}, {DeviceName: ptr(" "), Tag: ptr("hs4:t")}},
			"selector 2 has a blank deviceName"},
		"bad tag": {[]MetricsDashboardSelector{{Tag: ptr("hs4")}}, `selector 1: tag "hs4" is not <key>:<value>`},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := ValidateMetricsDashboardSelect(tt.selectors)
			require.Error(t, err)
			assert.Equal(t, errors.KindContractInvalid, errors.Kind(err))
			assert.Equal(t, tt.want, err.Error())
		})
	}
	assert.NoError(t, ValidateMetricsDashboardSelect(selectorsMetricsDashboardPanel().Select))
}

func TestMetricsDashboardPanel_StoredBeforeSelect(t *testing.T) {
	var p MetricsDashboardPanel
	require.NoError(t, json.Unmarshal([]byte(`{"title":"Temps","chartType":"line","tags":["hs4:t"],"devices":["d"]}`), &p))
	assert.Nil(t, p.Select)
	err := p.Validate()
	require.Error(t, err)
	assert.Equal(t, "panel Temps needs at least one selector", err.Error())
}

func TestMetricsDashboardPanel_JSON(t *testing.T) {
	b, err := json.Marshal(selectorsMetricsDashboardPanel())
	require.NoError(t, err)
	assert.JSONEq(t, `{"title":"Temps","chartType":"line","layout":{"x":0,"y":0,"w":0,"h":0},
		"select":[{"tag":"hs4:co2"},{"profileName":"Johnson-BMS","resourceName":"Lobby IAQ Sensor"}]}`, string(b))
	var p MetricsDashboardPanel
	require.NoError(t, json.Unmarshal(b, &p))
	assert.Equal(t, selectorsMetricsDashboardPanel(), p)
	assert.Nil(t, p.Select[0].DeviceName)
	assert.NoError(t, p.Validate())

	s := MetricsDashboardSelector{DeviceName: ptr("S"), ResourceName: ptr("Temp"), Tag: ptr("hs4:t")}
	b, err = json.Marshal(s)
	require.NoError(t, err)
	assert.JSONEq(t, `{"deviceName":"S","resourceName":"Temp","tag":"hs4:t"}`, string(b))
	var got MetricsDashboardSelector
	require.NoError(t, json.Unmarshal(b, &got))
	assert.Equal(t, s, got)

	require.NoError(t, json.Unmarshal([]byte(`{"title":"t","chartType":"line","select":[{"deviceName":"","tag":"hs4:t"}]}`), &p))
	require.NotNil(t, p.Select[0].DeviceName, "an empty deviceName is sent, not absent")
	assert.Error(t, p.Validate())

	withFields := selectorsMetricsDashboardPanel()
	withFields.ValueMapping = map[string]map[string]string{"Run Status": {"0": "Stopped", "1": "Running"}}
	withFields.RightAxisUnits = []string{"%"}
	b, err = json.Marshal(withFields)
	require.NoError(t, err)
	assert.JSONEq(t, `{"title":"Temps","chartType":"line","layout":{"x":0,"y":0,"w":0,"h":0},
		"select":[{"tag":"hs4:co2"},{"profileName":"Johnson-BMS","resourceName":"Lobby IAQ Sensor"}],
		"valueMapping":{"Run Status":{"0":"Stopped","1":"Running"}},"rightAxisUnits":["%"]}`, string(b))
	var gotFields MetricsDashboardPanel
	require.NoError(t, json.Unmarshal(b, &gotFields))
	assert.Equal(t, withFields, gotFields)
}

func TestMetricsDashboardOverview_JSON(t *testing.T) {
	summary := MetricsDashboardOverview{
		PointCount: 180, DeviceCount: 2,
		Tags:     []MetricsDashboardTagSummary{{Tag: "hs4:co2", PointCount: 1, DeviceCount: 1, Units: []string{"ppm"}}},
		Profiles: []MetricsDashboardProfileSummary{{ProfileName: "Schneider-BMS", DeviceCount: 1, ResourceCount: 92}},
	}
	b, err := json.Marshal(summary)
	require.NoError(t, err)
	assert.JSONEq(t, `{"pointCount":180,"deviceCount":2,
		"tags":[{"tag":"hs4:co2","pointCount":1,"deviceCount":1,"units":["ppm"]}],
		"profiles":[{"profileName":"Schneider-BMS","deviceCount":1,"resourceCount":92}]}`, string(b))
	var got MetricsDashboardOverview
	require.NoError(t, json.Unmarshal(b, &got))
	assert.Equal(t, summary, got)
}

func TestMetricsDashboardResource_JSON(t *testing.T) {
	resource := MetricsDashboardResource{ProfileName: "Johnson-BMS", ResourceName: "Lobby IAQ Sensor", DeviceCount: 1, Units: "ppm",
		ValueType: "Float32", Tags: []string{"hs4:co2"}, Description: "Indoor air quality, CO2 ppm"}
	b, err := json.Marshal(resource)
	require.NoError(t, err)
	assert.JSONEq(t, `{"profileName":"Johnson-BMS","resourceName":"Lobby IAQ Sensor","deviceCount":1,"units":"ppm",
		"valueType":"Float32","tags":["hs4:co2"],"description":"Indoor air quality, CO2 ppm"}`, string(b))
	var got MetricsDashboardResource
	require.NoError(t, json.Unmarshal(b, &got))
	assert.Equal(t, resource, got)

	b, err = json.Marshal(MetricsDashboardResource{ProfileName: "p", ResourceName: "r", ValueType: "Int32", Tags: []string{}})
	require.NoError(t, err)
	assert.NotContains(t, string(b), `"units"`)
	assert.NotContains(t, string(b), `"description"`)
}

func TestMetricsDashboardPanel_StoredWithoutResourceLabels(t *testing.T) {
	stored := `{"id":"a1","title":"Temps","chartType":"line","layout":{"x":0,"y":0,"w":6,"h":4},
		"select":[{"deviceName":"Boiler-1","resourceName":"Temperature"}],
		"valueMapping":{"Run Status":{"0":"Stopped"}},"rightAxisUnits":["%"]}`
	var p MetricsDashboardPanel
	require.NoError(t, json.Unmarshal([]byte(stored), &p))
	assert.Nil(t, p.ResourceLabels)
	b, err := json.Marshal(p)
	require.NoError(t, err)
	assert.NotContains(t, string(b), "resourceLabels")
	assert.JSONEq(t, stored, string(b))
}

func TestMetricsDashboardPanel_ResourceLabelsJSON(t *testing.T) {
	p := validMetricsDashboardPanel()
	p.ResourceLabels = map[string]map[string]string{
		"Boiler-1": {"Temperature": "Boiler 1 water"},
		"Boiler-2": {"Temperature": "Boiler 2 water"},
	}
	b, err := json.Marshal(p)
	require.NoError(t, err)
	assert.JSONEq(t, `{"title":"Temps","chartType":"line","layout":{"x":0,"y":0,"w":0,"h":0},
		"select":[{"tag":"hs4:zone-air-temp"}],
		"resourceLabels":{"Boiler-1":{"Temperature":"Boiler 1 water"},"Boiler-2":{"Temperature":"Boiler 2 water"}}}`, string(b))
	var got MetricsDashboardPanel
	require.NoError(t, json.Unmarshal(b, &got))
	assert.Equal(t, p, got)
}

func TestMetricsDashboardPanel_ResourceLabelsNotChecked(t *testing.T) {
	p := validMetricsDashboardPanel()
	p.ResourceLabels = map[string]map[string]string{"Unselected": {"Nothing": ""}, " ": {" ": " "}}
	assert.NoError(t, p.Validate())
}
