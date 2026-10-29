// Copyright (C) 2026 IOTech Ltd

package dtos

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
)

// AIDashboard is a dashboard drawn by an agent through central-mcp and arranged by people in Central UI.
// Ontology names the one tag ontology of every panel's tags; it is set at create and never changes.
type AIDashboard struct {
	Id          string    `json:"id,omitempty"`
	Name        string    `json:"name" validate:"required,edgex-dto-none-empty-string"`
	Description string    `json:"description,omitempty"`
	Ontology    string    `json:"ontology"`
	Range       AIRange   `json:"range"`
	Panels      []AIPanel `json:"panels"`
	edgexDtos.DBTimestamp
}

// AIDashboardInfo is an AIDashboard without its panels
type AIDashboardInfo struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Ontology    string `json:"ontology"`
	PanelCount  int    `json:"panelCount"`
	Created     int64  `json:"created,omitempty"`
	Modified    int64  `json:"modified,omitempty"`
}

// AIPanel reads every point whose tag is in Tags, on the devices in Devices (empty = all devices)
type AIPanel struct {
	Id        string   `json:"id,omitempty"`
	Title     string   `json:"title"`
	ChartType string   `json:"chartType"`
	Tags      []string `json:"tags"`
	Devices   []string `json:"devices,omitempty"`
	Range     *AIRange `json:"range,omitempty"`
	Layout    AILayout `json:"layout"`
	// Options holds the settings only one chart kind reads, such as rightAxisTags on line; it is not validated here
	Options map[string]any `json:"options,omitempty"`
}

// AIRange is either Last (a Go duration or "<n>d") or a fixed Start/End in epoch ms
type AIRange struct {
	Last  string `json:"last,omitempty"`
	Start int64  `json:"start,omitempty"`
	End   int64  `json:"end,omitempty"`
}

// AILayout is a panel's place on the grid of common.AIGridColumns columns
type AILayout struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

// ParseAIDuration accepts a Go duration or a number of days such as "7d"; the result must be positive
func ParseAIDuration(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if d, ok := parseAIDuration(s); ok && d > 0 {
		return d, nil
	}
	return 0, fmt.Errorf("%q is not a positive duration such as 15m, 24h or 7d", s)
}

func parseAIDuration(s string) (time.Duration, bool) {
	n, ok := strings.CutSuffix(s, "d")
	if !ok {
		d, err := time.ParseDuration(s)
		return d, err == nil
	}
	// decimal days only: ParseFloat alone would also take signs, exponents, hex, underscores, NaN and Inf
	if n == "" || strings.Trim(n, "0123456789.") != "" {
		return 0, false
	}
	days, err := strconv.ParseFloat(n, 64)
	// day counts past the int64 range would convert to an undefined duration
	if err != nil || days*float64(24*time.Hour) >= math.MaxInt64 {
		return 0, false
	}
	return time.Duration(days * float64(24*time.Hour)), true
}

// validate checks that a range has either last or a start before its end, not both; field prefixes the JSON path in the error text
func (r AIRange) validate(field string) error {
	hasWindow := r.Start != 0 || r.End != 0
	switch {
	case r.Last != "" && hasWindow:
		return fmt.Errorf("%s needs either last or start and end, not both", field)
	case r.Last != "":
		if _, err := ParseAIDuration(r.Last); err != nil {
			return fmt.Errorf("%s.last: %w", field, err)
		}
	case !hasWindow:
		return fmt.Errorf("%s needs either last or start and end", field)
	// a zero start reads as absent, so a fixed window needs both ends as positive epochs
	case r.Start <= 0:
		return fmt.Errorf("%s.start %d must be a positive epoch in ms, given with %s.end", field, r.Start, field)
	case r.Start >= r.End:
		return fmt.Errorf("%s.start %d must be before %s.end %d", field, r.Start, field, r.End)
	}
	return nil
}

// IsZero reports whether the layout is unset, in which case the server places the panel
func (l AILayout) IsZero() bool {
	return l == AILayout{}
}

func (l AILayout) validate() error {
	// x > columns-w rather than x+w > columns, which a huge x or w would overflow past
	if l.W < 1 || l.H < 1 || l.X < 0 || l.Y < 0 || l.W > common.AIGridColumns || l.X > common.AIGridColumns-l.W {
		return fmt.Errorf("layout {x:%d, y:%d, w:%d, h:%d} is not a place on the %d-column grid: w and h at least 1, x and y not negative, x+w at most %d",
			l.X, l.Y, l.W, l.H, common.AIGridColumns, common.AIGridColumns)
	}
	return nil
}

// Validate checks a dashboard as it will be stored, with every panel rule. It is the only place
// these rules live: central-ui calls it after each edit, and the requests do not repeat them.
func (d AIDashboard) Validate() error {
	if err := edgexCommon.Validate(d); err != nil {
		return err
	}
	if !slices.Contains(common.AIOntologies, d.Ontology) {
		return fmt.Errorf("ontology %q is not supported; use one of %s", d.Ontology, strings.Join(common.AIOntologies, ", "))
	}
	if err := d.Range.validate("range"); err != nil {
		return err
	}
	ids := make(map[string]bool, len(d.Panels))
	for i, p := range d.Panels {
		if err := p.validate(ids); err != nil {
			return fmt.Errorf("panels[%d]: %w", i, err)
		}
	}
	return nil
}

// validate checks a stored panel and records its id in ids, which holds the ids of the panels before it
func (p AIPanel) validate(ids map[string]bool) error {
	switch {
	case strings.TrimSpace(p.Id) == "":
		return errors.New("id is required")
	case ids[p.Id]:
		return fmt.Errorf("id %q is used by another panel", p.Id)
	case strings.TrimSpace(p.Title) == "":
		return errors.New("title is required")
	case !slices.Contains(common.AIChartKinds, p.ChartType):
		return fmt.Errorf("chartType %q is not supported; use one of %s", p.ChartType, strings.Join(common.AIChartKinds, ", "))
	case len(p.Tags) == 0:
		return errors.New("tags needs at least one tag")
	case slices.ContainsFunc(p.Tags, func(t string) bool { return strings.TrimSpace(t) == "" }):
		return errors.New("tags must not contain a blank tag")
	}
	ids[p.Id] = true
	if p.Range != nil {
		// a panel never has a fixed window
		if p.Range.Last == "" || p.Range.Start != 0 || p.Range.End != 0 {
			return errors.New("range of a panel may only be last N, such as 1h")
		}
		if err := p.Range.validate("range"); err != nil {
			return err
		}
	}
	// the server places a new panel before it validates, so a stored panel always has a real place
	return p.Layout.validate()
}
