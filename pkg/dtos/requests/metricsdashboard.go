// Copyright (C) 2026 IOTech Ltd

package requests

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

// AddMetricsDashboardRequest defines the Request Content for POST MetricsDashboard
type AddMetricsDashboardRequest struct {
	dtoCommon.BaseRequest `json:",inline"`
	Dashboard             dtos.MetricsDashboard `json:"dashboard"`
}

// Validate satisfies the Validator interface. It checks the dashboard as sent, a missing range read as the default;
// central-ui checks it again once it fills ids, the range and the layouts.
func (a *AddMetricsDashboardRequest) Validate() error {
	if err := edgexCommon.Validate(a); err != nil {
		return err
	}
	d := a.Dashboard
	if d.Range == (dtos.MetricsDashboardRange{}) {
		d.Range = dtos.MetricsDashboardRange{Last: common.DefaultMetricsDashboardRangeLast}
	}
	d.Panels = slices.Clone(d.Panels)
	return d.Validate()
}

// UnmarshalJSON implements the Unmarshaler interface for the AddMetricsDashboardRequest type
func (a *AddMetricsDashboardRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		dtoCommon.BaseRequest
		Dashboard dtos.MetricsDashboard
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}
	*a = AddMetricsDashboardRequest(alias)
	return a.Validate()
}

// UpdateMetricsDashboardRequest defines the Request Content for PATCH MetricsDashboard; only the fields set are changed.
// Validate checks each field set on its own; rules across the dashboard are central-ui's, on the dashboard it stores.
type UpdateMetricsDashboardRequest struct {
	dtoCommon.BaseRequest `json:",inline"`
	Name                  *string                     `json:"name,omitempty"`
	Description           *string                     `json:"description,omitempty"`
	Range                 *dtos.MetricsDashboardRange `json:"range,omitempty"`
}

// Validate satisfies the Validator interface
func (u *UpdateMetricsDashboardRequest) Validate() error {
	if err := edgexCommon.Validate(u); err != nil {
		return err
	}
	if u.Name == nil && u.Description == nil && u.Range == nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "the update sets none of name, description and range", nil)
	}
	if u.Name != nil && strings.TrimSpace(*u.Name) == "" {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "dashboard name is required", nil)
	}
	if u.Range != nil {
		return u.Range.Validate()
	}
	return nil
}

// UnmarshalJSON implements the Unmarshaler interface for the UpdateMetricsDashboardRequest type
func (u *UpdateMetricsDashboardRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		dtoCommon.BaseRequest
		Name        *string
		Description *string
		Range       *dtos.MetricsDashboardRange
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}
	*u = UpdateMetricsDashboardRequest(alias)
	return u.Validate()
}

// AddMetricsDashboardPanelRequest defines the Request Content for POST MetricsDashboardPanel; central-ui assigns the panel id.
// Validate checks the panel on its own; rules across the dashboard are central-ui's, on the dashboard it stores.
type AddMetricsDashboardPanelRequest struct {
	dtoCommon.BaseRequest `json:",inline"`
	Panel                 dtos.MetricsDashboardPanel `json:"panel"`
}

// Validate satisfies the Validator interface
func (a *AddMetricsDashboardPanelRequest) Validate() error {
	if err := edgexCommon.Validate(a); err != nil {
		return err
	}
	return a.Panel.Validate()
}

// UnmarshalJSON implements the Unmarshaler interface for the AddMetricsDashboardPanelRequest type and drops a supplied panel id
func (a *AddMetricsDashboardPanelRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		dtoCommon.BaseRequest
		Panel dtos.MetricsDashboardPanel
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}
	*a = AddMetricsDashboardPanelRequest(alias)
	a.Panel.Id = ""
	return a.Validate()
}

// UpdateMetricsDashboardPanelRequest defines the Request Content for PATCH MetricsDashboardPanel; only the fields set are
// changed, and a set empty Devices means every device. Validate checks each field set on its own; rules across the
// dashboard are central-ui's, on the dashboard it stores.
type UpdateMetricsDashboardPanelRequest struct {
	dtoCommon.BaseRequest `json:",inline"`
	Title                 *string                      `json:"title,omitempty"`
	ChartType             *string                      `json:"chartType,omitempty"`
	Tags                  *[]string                    `json:"tags,omitempty"`
	Devices               *[]string                    `json:"devices,omitempty"`
	Layout                *dtos.MetricsDashboardLayout `json:"layout,omitempty"`
}

// Validate satisfies the Validator interface
func (u *UpdateMetricsDashboardPanelRequest) Validate() error {
	if err := edgexCommon.Validate(u); err != nil {
		return err
	}
	if u.Title == nil && u.ChartType == nil && u.Tags == nil && u.Devices == nil && u.Layout == nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "the update sets none of title, chartType, tags, devices and layout", nil)
	}
	if u.Title != nil && strings.TrimSpace(*u.Title) == "" {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "panel title is required", nil)
	}
	if u.ChartType != nil && strings.TrimSpace(*u.ChartType) == "" {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "panel chartType is required", nil)
	}
	if u.Tags != nil {
		tags, err := dtos.ValidateMetricsDashboardTags(*u.Tags)
		if err != nil {
			return err
		}
		*u.Tags = tags
	}
	if u.Layout != nil {
		return u.Layout.Validate()
	}
	return nil
}

// UnmarshalJSON implements the Unmarshaler interface for the UpdateMetricsDashboardPanelRequest type
func (u *UpdateMetricsDashboardPanelRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		dtoCommon.BaseRequest
		Title     *string
		ChartType *string
		Tags      *[]string
		Devices   *[]string
		Layout    *dtos.MetricsDashboardLayout
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}
	*u = UpdateMetricsDashboardPanelRequest(alias)
	return u.Validate()
}
