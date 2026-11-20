// Copyright (C) 2026 IOTech Ltd

package dtos

import (
	"fmt"
	"strings"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

const MetricsDashboardGridColumns = 12

// MetricsDashboard is a dashboard an agent or a user builds from device resources; created and modified are epoch ms
type MetricsDashboard struct {
	Id          string                  `json:"id,omitempty"`
	Name        string                  `json:"name"`
	Description string                  `json:"description,omitempty"`
	Range       MetricsDashboardRange   `json:"range"`
	Panels      []MetricsDashboardPanel `json:"panels"`
	Pinned      bool                    `json:"pinned"`
	Revision    int64                   `json:"revision,omitempty"`
	Created     int64                   `json:"created,omitempty"`
	Modified    int64                   `json:"modified,omitempty"`
}

// MetricsDashboardRange is either a relative range (Last, one of common.MetricsDashboardLastChoices) or a fixed range of epoch ms
type MetricsDashboardRange struct {
	Last  string `json:"last,omitempty"`
	Start int64  `json:"start,omitempty"`
	End   int64  `json:"end,omitempty"`
}

// MetricsDashboardPanel selects its points one of two ways, never both: by Tags, the points carrying one of them,
// limited to Devices when it is not empty; or by Points, the device resources it lists, in order. Devices only
// narrows Tags.
type MetricsDashboardPanel struct {
	Id        string                       `json:"id,omitempty"`
	Title     string                       `json:"title"`
	ChartType string                       `json:"chartType"`
	Tags      []string                     `json:"tags,omitempty"`
	Devices   []string                     `json:"devices,omitempty"`
	Points    []MetricsDashboardPanelPoint `json:"points,omitempty"`
	Layout    MetricsDashboardLayout       `json:"layout"`
}

// MetricsDashboardPanelPoint is a device resource a panel lists by name
type MetricsDashboardPanelPoint struct {
	DeviceName   string `json:"deviceName"`
	ResourceName string `json:"resourceName"`
}

// MetricsDashboardLayout places a panel on a 12-column grid; the zero layout asks central-ui to place it
type MetricsDashboardLayout struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

// MetricsDashboardInfo is a dashboard in the list of dashboards
type MetricsDashboardInfo struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	PanelCount  int    `json:"panelCount"`
	Pinned      bool   `json:"pinned"`
	Revision    int64  `json:"revision"`
	Modified    int64  `json:"modified"`
}

// MetricsDashboardTagSummary counts the points and devices that carry each tag
type MetricsDashboardTagSummary struct {
	TotalPoints  int                      `json:"totalPoints"`
	TotalDevices int                      `json:"totalDevices"`
	Tags         []MetricsDashboardTagRow `json:"tags"`
}

type MetricsDashboardTagRow struct {
	Tag     string   `json:"tag"`
	Points  int      `json:"points"`
	Devices int      `json:"devices"`
	Units   []string `json:"units"`
}

// MetricsDashboardPoint is a device resource carrying a tag
type MetricsDashboardPoint struct {
	DeviceName   string `json:"deviceName"`
	ResourceName string `json:"resourceName"`
	Tag          string `json:"tag"`
	Units        string `json:"units,omitempty"`
	Description  string `json:"description,omitempty"`
}

// Validate de-duplicates the panels' tags and turns nil panels into an empty list, then checks the shape of the
// dashboard. Tags are kept as sent, spaces included. Which chart types exist and how many panels a dashboard holds
// are central-ui's rules, not checked here.
func (d *MetricsDashboard) Validate() error {
	if d.Panels == nil {
		d.Panels = []MetricsDashboardPanel{}
	}
	if strings.TrimSpace(d.Name) == "" {
		return invalid("dashboard name is required")
	}
	if err := d.Range.Validate(); err != nil {
		return err
	}
	ids := make(map[string]bool, len(d.Panels))
	for i := range d.Panels {
		p := &d.Panels[i]
		if err := p.Validate(); err != nil {
			return err
		}
		if p.Id == "" {
			continue
		}
		if ids[p.Id] {
			return invalid(fmt.Sprintf("panel id %s is used twice", p.Id))
		}
		ids[p.Id] = true
	}
	return nil
}

// Validate checks the range on its own: either a choice of common.MetricsDashboardLastChoices or start before end
func (r MetricsDashboardRange) Validate() error {
	hasLast := r.Last != ""
	hasStart, hasEnd := r.Start != 0, r.End != 0
	if hasStart != hasEnd {
		return invalid("range needs both start and end")
	}
	if hasLast == hasStart {
		return invalid("range needs either last or start and end")
	}
	if hasLast {
		if _, ok := common.MetricsDashboardLastDuration(r.Last); !ok {
			return invalid(fmt.Sprintf("range last %s is not one of %s", r.Last, metricsDashboardLastList()))
		}
		return nil
	}
	if r.Start >= r.End {
		return invalid("range start must be before end")
	}
	return nil
}

// Validate de-duplicates the panel's tags, or drops an empty Devices beside Points, then checks the panel on its own.
// Points are kept as sent, order and repeats included. Panel ids are checked across the dashboard by
// MetricsDashboard.Validate.
func (p *MetricsDashboardPanel) Validate() error {
	if strings.TrimSpace(p.Title) == "" {
		return invalid("panel title is required")
	}
	if strings.TrimSpace(p.ChartType) == "" {
		return invalid(fmt.Sprintf("panel %s needs a chartType", p.Title))
	}
	hasTags, hasPoints := len(p.Tags) > 0, len(p.Points) > 0
	switch {
	case hasTags && hasPoints:
		return invalid(fmt.Sprintf("panel %s has both tags and points", p.Title))
	case hasPoints:
		if err := validateMetricsDashboardPoints(p.Points, fmt.Sprintf(" of panel %s", p.Title)); err != nil {
			return err
		}
		if len(p.Devices) > 0 {
			return invalid(fmt.Sprintf("panel %s has devices beside points", p.Title))
		}
		p.Tags, p.Devices = nil, nil
	case hasTags:
		tags, err := ValidateMetricsDashboardTags(p.Tags)
		if err != nil {
			return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("panel %s", p.Title), err)
		}
		p.Tags = tags
	default:
		return invalid(fmt.Sprintf("panel %s needs tags or points", p.Title))
	}
	if err := p.Layout.Validate(); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("panel %s", p.Title), err)
	}
	return nil
}

// ValidateMetricsDashboardTags returns the tags de-duplicated in order, or an error when none is left or one is not
// <key>:<value> with a non-blank key and value. Each tag is kept as sent, spaces included.
func ValidateMetricsDashboardTags(tags []string) ([]string, error) {
	tags = dedupeTags(tags)
	if len(tags) == 0 {
		return nil, invalid("at least one tag is required")
	}
	for _, tag := range tags {
		key, value, ok := strings.Cut(tag, ":")
		if !ok || strings.TrimSpace(key) == "" || strings.TrimSpace(value) == "" {
			return nil, invalid(fmt.Sprintf("tag %q is not <key>:<value>", tag))
		}
	}
	return tags, nil
}

// ValidateMetricsDashboardPoints returns an error when the list is empty or a point has a blank deviceName or
// resourceName. Points are counted from 1 in the error.
func ValidateMetricsDashboardPoints(points []MetricsDashboardPanelPoint) error {
	if len(points) == 0 {
		return invalid("at least one point is required")
	}
	return validateMetricsDashboardPoints(points, "")
}

// validateMetricsDashboardPoints checks each point's names; where follows "point <i>" in the error
func validateMetricsDashboardPoints(points []MetricsDashboardPanelPoint, where string) error {
	for i, pt := range points {
		if strings.TrimSpace(pt.DeviceName) == "" || strings.TrimSpace(pt.ResourceName) == "" {
			return invalid(fmt.Sprintf("point %d%s needs a deviceName and a resourceName", i+1, where))
		}
	}
	return nil
}

// Validate checks the layout on its own; a zero layout passes and means "place it for me"
func (l MetricsDashboardLayout) Validate() error {
	if l == (MetricsDashboardLayout{}) {
		return nil
	}
	if l.W <= 0 || l.H <= 0 || l.X < 0 || l.Y < 0 || l.X > MetricsDashboardGridColumns-l.W {
		return invalid(fmt.Sprintf("layout %+v is outside the %d-column grid", l, MetricsDashboardGridColumns))
	}
	return nil
}

func dedupeTags(tags []string) []string {
	seen := make(map[string]bool, len(tags))
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		if seen[t] {
			continue
		}
		seen[t] = true
		out = append(out, t)
	}
	return out
}

func metricsDashboardLastList() string {
	lasts := make([]string, len(common.MetricsDashboardLastChoices))
	for i, c := range common.MetricsDashboardLastChoices {
		lasts[i] = c.Last
	}
	return strings.Join(lasts, ", ")
}

func invalid(msg string) error {
	return errors.NewCommonEdgeX(errors.KindContractInvalid, msg, nil)
}
