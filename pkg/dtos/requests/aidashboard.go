// Copyright (C) 2026 IOTech Ltd

package requests

import (
	"encoding/json"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"
	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

// The AI dashboard requests check only their shape: apiVersion, at least one change, the fields
// that cannot be set together, and a non-blank dashboard name. They do not check title, chartType,
// tags, range or layout.
// Those rules live in dtos.AIDashboard.Validate, which central-ui runs on the whole dashboard
// after it applies the request to the stored one, because a request alone cannot see the
// stored panels, for example whether panel ids are unique or where a new panel is placed.

// AddAIDashboardRequest defines the Request Content for POST AIDashboard DTO.
// Ids and timestamps in the body are dropped, and a missing range becomes last 24h.
type AddAIDashboardRequest struct {
	dtoCommon.BaseRequest `json:",inline"`
	Dashboard             dtos.AIDashboard `json:"dashboard"`
}

// Validate satisfies the Validator interface
func (a *AddAIDashboardRequest) Validate() error {
	return edgexCommon.Validate(a)
}

// UnmarshalJSON implements the Unmarshaler interface for the AddAIDashboardRequest type
func (a *AddAIDashboardRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		dtoCommon.BaseRequest
		Dashboard dtos.AIDashboard
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}

	*a = AddAIDashboardRequest(alias)
	a.Dashboard.Id = ""
	a.Dashboard.DBTimestamp = edgexDtos.DBTimestamp{}
	for i := range a.Dashboard.Panels {
		a.Dashboard.Panels[i].Id = ""
	}
	if a.Dashboard.Range == (dtos.AIRange{}) {
		a.Dashboard.Range = dtos.AIRange{Last: common.DefaultAIRangeLast}
	}

	return a.Validate()
}

// UpdateAIDashboardRequest defines the Request Content for PATCH AIDashboard DTO.
// A nil field is left unchanged.
type UpdateAIDashboardRequest struct {
	dtoCommon.BaseRequest `json:",inline"`
	Name                  *string       `json:"name,omitempty" validate:"omitempty,edgex-dto-none-empty-string"`
	Description           *string       `json:"description,omitempty"`
	Range                 *dtos.AIRange `json:"range,omitempty"`
}

// Validate satisfies the Validator interface
func (u *UpdateAIDashboardRequest) Validate() error {
	if err := edgexCommon.Validate(u); err != nil {
		return err
	}
	if u.Name == nil && u.Description == nil && u.Range == nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "the request needs at least one of name, description or range", nil)
	}
	return nil
}

// UnmarshalJSON implements the Unmarshaler interface for the UpdateAIDashboardRequest type
func (u *UpdateAIDashboardRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		dtoCommon.BaseRequest
		Name        *string
		Description *string
		Range       *dtos.AIRange
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}

	*u = UpdateAIDashboardRequest(alias)
	return u.Validate()
}

// AddAIPanelRequest defines the Request Content for POST AIPanel DTO.
// The panel id in the body is dropped; central-ui places a panel sent with a zero layout.
type AddAIPanelRequest struct {
	dtoCommon.BaseRequest `json:",inline"`
	Panel                 dtos.AIPanel `json:"panel"`
}

// Validate satisfies the Validator interface
func (a *AddAIPanelRequest) Validate() error {
	return edgexCommon.Validate(a)
}

// UnmarshalJSON implements the Unmarshaler interface for the AddAIPanelRequest type
func (a *AddAIPanelRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		dtoCommon.BaseRequest
		Panel dtos.AIPanel
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}

	*a = AddAIPanelRequest(alias)
	a.Panel.Id = ""
	return a.Validate()
}

// UpdateAIPanelRequest defines the Request Content for PATCH AIPanel DTO.
// A nil field is left unchanged. FollowDashboardRange removes the panel's own range; a set Options
// replaces the whole object; a set Layout is used as sent, so a zero one is refused.
type UpdateAIPanelRequest struct {
	dtoCommon.BaseRequest `json:",inline"`
	Title                 *string         `json:"title,omitempty"`
	ChartType             *string         `json:"chartType,omitempty"`
	Tags                  *[]string       `json:"tags,omitempty"`
	Devices               *[]string       `json:"devices,omitempty"`
	Range                 *dtos.AIRange   `json:"range,omitempty"`
	FollowDashboardRange  bool            `json:"followDashboardRange,omitempty"`
	Layout                *dtos.AILayout  `json:"layout,omitempty"`
	Options               *map[string]any `json:"options,omitempty"`
}

// Validate satisfies the Validator interface
func (u *UpdateAIPanelRequest) Validate() error {
	if err := edgexCommon.Validate(u); err != nil {
		return err
	}
	if u.Title == nil && u.ChartType == nil && u.Tags == nil && u.Devices == nil && u.Range == nil &&
		!u.FollowDashboardRange && u.Layout == nil && u.Options == nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "the request needs at least one field to change", nil)
	}
	if u.FollowDashboardRange && u.Range != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "range cannot be set together with followDashboardRange", nil)
	}
	return nil
}

// UnmarshalJSON implements the Unmarshaler interface for the UpdateAIPanelRequest type
func (u *UpdateAIPanelRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		dtoCommon.BaseRequest
		Title                *string
		ChartType            *string
		Tags                 *[]string
		Devices              *[]string
		Range                *dtos.AIRange
		FollowDashboardRange bool
		Layout               *dtos.AILayout
		Options              *map[string]any
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}

	*u = UpdateAIPanelRequest(alias)
	return u.Validate()
}

// UpdateAILayoutsRequest defines the Request Content for PATCH AIDashboard layouts.
// Maps panel id to its new layout.
type UpdateAILayoutsRequest struct {
	dtoCommon.BaseRequest `json:",inline"`
	Layouts               map[string]dtos.AILayout `json:"layouts"`
}

// Validate satisfies the Validator interface
func (u *UpdateAILayoutsRequest) Validate() error {
	if err := edgexCommon.Validate(u); err != nil {
		return err
	}
	if len(u.Layouts) == 0 {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "layouts needs at least one panel", nil)
	}
	return nil
}

// UnmarshalJSON implements the Unmarshaler interface for the UpdateAILayoutsRequest type
func (u *UpdateAILayoutsRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		dtoCommon.BaseRequest
		Layouts map[string]dtos.AILayout
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}

	*u = UpdateAILayoutsRequest(alias)
	return u.Validate()
}
