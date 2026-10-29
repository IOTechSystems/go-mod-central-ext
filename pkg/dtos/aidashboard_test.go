// Copyright (C) 2026 IOTech Ltd

package dtos

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validAIPanel() AIPanel {
	return AIPanel{
		Id:        "p",
		Title:     "Supply air",
		ChartType: common.AIChartLine,
		Tags:      []string{"t"},
		Layout:    AILayout{X: 0, Y: 0, W: 6, H: 4},
	}
}

func TestAIDashboardValidate(t *testing.T) {
	valid := func() AIDashboard {
		return AIDashboard{Name: "a", Ontology: common.AIOntologyHaystack, Range: AIRange{Last: "24h"}, Panels: []AIPanel{validAIPanel()}}
	}
	panel := func(edit func(p *AIPanel)) func(d *AIDashboard) {
		return func(d *AIDashboard) { edit(&d.Panels[0]) }
	}
	second := func(edit func(p *AIPanel)) func(d *AIDashboard) {
		return func(d *AIDashboard) {
			p := validAIPanel()
			p.Id = "q"
			edit(&p)
			d.Panels = append(d.Panels, p)
		}
	}
	tests := []struct {
		name    string
		edit    func(d *AIDashboard)
		wantErr string
	}{
		{"valid", func(d *AIDashboard) {}, ""},
		{"no panels", func(d *AIDashboard) { d.Panels = nil }, ""},
		{"no name", func(d *AIDashboard) { d.Name = "" }, "Name"},
		{"blank name", func(d *AIDashboard) { d.Name = " " }, "Name"},
		{"long name", func(d *AIDashboard) { d.Name = strings.Repeat("é", 1000) }, ""},
		{"ontology eo66", func(d *AIDashboard) { d.Ontology = common.AIOntologyEO66 }, ""},
		{"no ontology", func(d *AIDashboard) { d.Ontology = "" }, "ontology"},
		{"unknown ontology", func(d *AIDashboard) { d.Ontology = "brick" }, "ontology"},
		{"range last 2x", func(d *AIDashboard) { d.Range = AIRange{Last: "2x"} }, "range.last"},
		{"range last 7d", func(d *AIDashboard) { d.Range = AIRange{Last: "7d"} }, ""},
		{"range start after end", func(d *AIDashboard) { d.Range = AIRange{Start: 2, End: 1} }, "range.start"},
		{"range empty", func(d *AIDashboard) { d.Range = AIRange{} }, "range"},
		{"blank panel id", panel(func(p *AIPanel) { p.Id = " " }), "panels[0]: id"},
		{"two panels with id p", second(func(p *AIPanel) { p.Id = "p" }), "panels[1]: id"},
		{"blank title", panel(func(p *AIPanel) { p.Title = "  " }), "panels[0]: title"},
		{"kind pie", panel(func(p *AIPanel) { p.ChartType = "pie" }), "panels[0]: chartType"},
		{"no tags", panel(func(p *AIPanel) { p.Tags = nil }), "panels[0]: tags"},
		{"blank tag", second(func(p *AIPanel) { p.Tags = []string{"t", " "} }), "panels[1]: tags"},
		{"panel range start and end", panel(func(p *AIPanel) { p.Range = &AIRange{Start: 1, End: 2} }), "panels[0]: range"},
		{"panel range last 0d", panel(func(p *AIPanel) { p.Range = &AIRange{Last: "0d"} }), "panels[0]: range.last"},
		{"panel range last 1h", panel(func(p *AIPanel) { p.Range = &AIRange{Last: "1h"} }), ""},
		{"zero layout", panel(func(p *AIPanel) { p.Layout = AILayout{} }), "panels[0]: layout"},
		{"layout off the grid", second(func(p *AIPanel) { p.Layout = AILayout{X: 8, Y: 0, W: 6, H: 4} }), "panels[1]: layout"},
		{"options of line on line", panel(func(p *AIPanel) { p.Options = map[string]any{"rightAxisTags": []any{"t"}} }), ""},
		{"options with an unknown key", panel(func(p *AIPanel) { p.Options = map[string]any{"colour": "red"} }), ""},
		{"options with a wrong type", panel(func(p *AIPanel) {
			p.ChartType = common.AIChartState
			p.Options = map[string]any{"stateLabels": 3, "rightAxisTags": "x"}
		}), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := valid()
			tt.edit(&d)
			err := d.Validate()
			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestParseAIDuration(t *testing.T) {
	tests := []struct {
		in      string
		want    time.Duration
		wantErr bool
	}{
		{"15m", 15 * time.Minute, false},
		{"24h", 24 * time.Hour, false},
		{"7d", 168 * time.Hour, false},
		{"0d", 0, true},
		{"-1h", 0, true},
		{"x", 0, true},
		{"NaNd", 0, true},
		{"Infd", 0, true},
		{"1e300d", 0, true},
		{"1e-20d", 0, true},
		{"1.5d", 36 * time.Hour, false},
		{"1000000000d", 0, true},
		{"0.00000000000000000001d", 0, true},
		{"0x1p3d", 0, true},
		{"1_0d", 0, true},
		{"+7d", 0, true},
		{"1e1d", 0, true},
		{".d", 0, true},
		{"d", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := ParseAIDuration(tt.in)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAILayoutValidate(t *testing.T) {
	tests := []struct {
		name    string
		layout  AILayout
		wantErr bool
	}{
		{"full width", AILayout{X: 0, Y: 0, W: 12, H: 4}, false},
		{"right edge", AILayout{X: 6, Y: 9, W: 6, H: 4}, false},
		{"past the right edge", AILayout{X: 7, Y: 0, W: 6, H: 4}, true},
		{"huge x", AILayout{X: math.MaxInt, Y: 0, W: 1, H: 1}, true},
		{"huge w", AILayout{X: 1, Y: 0, W: math.MaxInt, H: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.layout.validate()
			if !tt.wantErr {
				assert.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), "layout {x:")
		})
	}
}

func TestAIRangeValidate(t *testing.T) {
	tests := []struct {
		name      string
		r         AIRange
		wantField string
	}{
		{"last", AIRange{Last: "1h"}, ""},
		{"start and end", AIRange{Start: 1, End: 2}, ""},
		{"empty", AIRange{}, "range"},
		{"last with start and end", AIRange{Last: "1h", Start: 1, End: 2}, "range"},
		{"end without start", AIRange{End: 1000}, "range.start"},
		{"start without end", AIRange{Start: 1000}, "range.start"},
		{"negative epochs", AIRange{Start: -100, End: -1}, "range.start"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.r.validate("range")
			if tt.wantField == "" {
				assert.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantField)
		})
	}
}
