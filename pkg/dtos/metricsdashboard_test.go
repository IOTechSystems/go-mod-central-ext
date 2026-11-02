// Copyright (C) 2026 IOTech Ltd

package dtos

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validMetricsDashboardPanel() MetricsDashboardPanel {
	return MetricsDashboardPanel{Title: "Temps", ChartType: "line", Tags: []string{"hs4:zone-air-temp"}}
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
		{"no tag", withPanel(func(p *MetricsDashboardPanel) { p.Tags = nil })},
		{"tag without colon", withPanel(func(p *MetricsDashboardPanel) { p.Tags = []string{"zone-air-temp"} })},
		{"blank key", withPanel(func(p *MetricsDashboardPanel) { p.Tags = []string{" :zone-air-temp"} })},
		{"blank value", withPanel(func(p *MetricsDashboardPanel) { p.Tags = []string{"hs4: "} })},
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
		{"tags of two keys in one panel", func(d *MetricsDashboard) {
			d.Panels[0].Tags = []string{"hs4:zone-air-temp", "eo66PointType:temp"}
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
		{"panels without ids", func(d *MetricsDashboard) {
			d.Panels = []MetricsDashboardPanel{validMetricsDashboardPanel(), validMetricsDashboardPanel()}
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

func TestMetricsDashboardValidate_DedupesTags(t *testing.T) {
	d := validMetricsDashboard()
	d.Panels[0].Tags = []string{"hs4:b", "hs4:a", "hs4:b", " hs4:a "}
	require.NoError(t, d.Validate())
	assert.Equal(t, []string{"hs4:b", "hs4:a", " hs4:a "}, d.Panels[0].Tags)
}

func TestMetricsDashboard_DecodesUnknownField(t *testing.T) {
	doc := `{"name":"temps","range":{"last":"5m"},"future":1,
		"panels":[{"title":"Temps","chartType":"line","tags":["hs4:t"],"rightAxis":true,"layout":{"x":0,"y":0,"w":6,"h":4}}]}`
	var d MetricsDashboard
	require.NoError(t, json.Unmarshal([]byte(doc), &d))
	assert.Equal(t, "temps", d.Name)
	assert.Equal(t, MetricsDashboardLayout{W: 6, H: 4}, d.Panels[0].Layout)
	assert.NoError(t, d.Validate())
}

func TestMetricsDashboardValidate_EmptyPanelsIsArray(t *testing.T) {
	d := MetricsDashboard{Name: "temps", Range: MetricsDashboardRange{Last: "5m"}}
	require.NoError(t, d.Validate())
	b, err := json.Marshal(d)
	require.NoError(t, err)
	assert.Contains(t, string(b), `"panels":[]`)
}

func TestMetricsDashboardPanelValidate_NamesThePanel(t *testing.T) {
	p := MetricsDashboardPanel{Title: "Zone temp", ChartType: "line", Tags: []string{"hs4:t"}, Layout: MetricsDashboardLayout{X: 10, W: 6, H: 4}}
	err := p.Validate()
	require.Error(t, err)
	assert.Equal(t, errors.KindContractInvalid, errors.Kind(err))
	assert.Contains(t, err.Error(), "panel Zone temp")
	assert.Contains(t, err.Error(), "outside the 12-column grid")
}
