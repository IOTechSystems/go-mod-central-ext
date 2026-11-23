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

// MetricsDashboardPanel charts the union of what its selectors select
type MetricsDashboardPanel struct {
	Id        string                     `json:"id,omitempty"`
	Title     string                     `json:"title"`
	ChartType string                     `json:"chartType"`
	Select    []MetricsDashboardSelector `json:"select"`
	Layout    MetricsDashboardLayout     `json:"layout"`
}

// MetricsDashboardSelector is one entry of a panel's select; a panel charts every point its selectors pick.
// The set fields must all match: set ResourceName or Tag, and at most one of DeviceName and ProfileName.
// Without DeviceName it is a rule that also picks devices added later.
// Selectors may pick the same point, so whoever resolves them must remove the duplicates.
// For example:
//
//	{"deviceName": "AHU-1", "resourceName": "Supply Air Temp"}   AHU-1's Supply Air Temp
//	{"deviceName": "AHU-1", "tag": "hs4:zone-air-temp"}         AHU-1's resources tagged hs4:zone-air-temp
//	{"profileName": "AHU", "tag": "hs4:zone-air-temp"}          the same on every device of profile AHU
//	{"resourceName": "Lobby IAQ Sensor"}                        that resource on every device that has it
//	{"tag": "hs4:co2"}                                          every resource tagged hs4:co2, on any device
//
// A nil field is absent; a field sent as "" is refused as blank.
type MetricsDashboardSelector struct {
	DeviceName   *string `json:"deviceName,omitempty"`
	ProfileName  *string `json:"profileName,omitempty"`
	ResourceName *string `json:"resourceName,omitempty"`
	Tag          *string `json:"tag,omitempty"`
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

// MetricsDashboardOverview is what an agent reads first of what a panel can select: the tags stored on readable
// resources and the profiles, with PointCount and DeviceCount over every readable resource on every device. Its
// length does not grow with the number of resources; the search finds resources by keyword.
type MetricsDashboardOverview struct {
	PointCount  int                              `json:"pointCount"`
	DeviceCount int                              `json:"deviceCount"`
	Tags        []MetricsDashboardTagSummary     `json:"tags"`
	Profiles    []MetricsDashboardProfileSummary `json:"profiles"`
}

// MetricsDashboardTagSummary is one tag in the overview: how many points and devices carry it, and their units
type MetricsDashboardTagSummary struct {
	Tag         string   `json:"tag"`
	PointCount  int      `json:"pointCount"`
	DeviceCount int      `json:"deviceCount"`
	Units       []string `json:"units"`
}

// MetricsDashboardProfileSummary is one profile in the overview: how many devices it has and how many readable
// resources, so an agent sees where untagged resources live and which profile to search
type MetricsDashboardProfileSummary struct {
	ProfileName   string `json:"profileName"`
	DeviceCount   int    `json:"deviceCount"`
	ResourceCount int    `json:"resourceCount"`
}

// MetricsDashboardResource is one search result: a profile's readable resource, which every device of that profile
// has, DeviceCount counting them; its ProfileName and ResourceName make a selector as they stand. Tags holds the tags
// stored on it. It is not a point, which is one device's resource.
type MetricsDashboardResource struct {
	ProfileName  string   `json:"profileName"`
	ResourceName string   `json:"resourceName"`
	DeviceCount  int      `json:"deviceCount"`
	Units        string   `json:"units,omitempty"`
	ValueType    string   `json:"valueType"`
	Tags         []string `json:"tags"`
	Description  string   `json:"description,omitempty"`
}

// Validate turns nil panels into an empty list, then checks the shape of the dashboard. Which chart types exist and
// how many panels a dashboard holds are central-ui's rules, not checked here.
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

// Validate checks the panel on its own. Selectors are kept as sent, order, repeats and spaces included. Panel ids
// are checked across the dashboard by MetricsDashboard.Validate.
func (p *MetricsDashboardPanel) Validate() error {
	if strings.TrimSpace(p.Title) == "" {
		return invalid("panel title is required")
	}
	if strings.TrimSpace(p.ChartType) == "" {
		return invalid(fmt.Sprintf("panel %s needs a chartType", p.Title))
	}
	if len(p.Select) == 0 {
		return invalid(fmt.Sprintf("panel %s needs at least one selector", p.Title))
	}
	if err := validateMetricsDashboardSelectors(p.Select, fmt.Sprintf(" of panel %s", p.Title)); err != nil {
		return err
	}
	if err := p.Layout.Validate(); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("panel %s", p.Title), err)
	}
	return nil
}

// ValidateMetricsDashboardSelect returns an error when the list is empty or a selector is not valid: a field set but
// blank, neither resourceName nor tag, both deviceName and profileName, or a tag that is not <key>:<value> with a
// non-blank key and value. Selectors are counted from 1 in the error.
func ValidateMetricsDashboardSelect(selectors []MetricsDashboardSelector) error {
	if len(selectors) == 0 {
		return invalid("at least one selector is required")
	}
	return validateMetricsDashboardSelectors(selectors, "")
}

// validateMetricsDashboardSelectors checks each selector; where follows "selector <i>" in the error
func validateMetricsDashboardSelectors(selectors []MetricsDashboardSelector, where string) error {
	for i, s := range selectors {
		if err := s.validate(fmt.Sprintf("selector %d%s", i+1, where)); err != nil {
			return err
		}
	}
	return nil
}

// validate checks the selector on its own; name starts the error
func (s MetricsDashboardSelector) validate(name string) error {
	for _, f := range []struct {
		field string
		value *string
	}{{"deviceName", s.DeviceName}, {"profileName", s.ProfileName}, {"resourceName", s.ResourceName}, {"tag", s.Tag}} {
		if f.value != nil && strings.TrimSpace(*f.value) == "" {
			return invalid(fmt.Sprintf("%s has a blank %s", name, f.field))
		}
	}
	if s.ResourceName == nil && s.Tag == nil {
		return invalid(fmt.Sprintf("%s needs a resourceName or a tag", name))
	}
	if s.DeviceName != nil && s.ProfileName != nil {
		return invalid(fmt.Sprintf("%s sets both deviceName and profileName", name))
	}
	if s.Tag != nil {
		key, value, ok := strings.Cut(*s.Tag, ":")
		if !ok || strings.TrimSpace(key) == "" || strings.TrimSpace(value) == "" {
			return invalid(fmt.Sprintf("%s: tag %q is not <key>:<value>", name, *s.Tag))
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
