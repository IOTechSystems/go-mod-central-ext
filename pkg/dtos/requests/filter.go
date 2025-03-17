// Copyright (C) 2025 IOTech Ltd

package requests

import (
	"encoding/json"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

// FilterRequest defines the Request Content for POST Filter DTO.
type FilterRequest struct {
	dtoCommon.BaseRequest `json:",inline"`
	Filter                dtos.Filter `json:"filter"`
}

// Validate satisfies the Validator interface
func (d *FilterRequest) Validate() error {
	err := common.Validate(d)
	return err
}

// UnmarshalJSON implements the Unmarshaler interface for the FilterRequest type
func (d *FilterRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		dtoCommon.BaseRequest
		Filter dtos.Filter
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}

	*d = FilterRequest(alias)

	// validate FilterRequest DTO
	if err := d.Validate(); err != nil {
		return err
	}
	return nil
}

// AddFilterReqToFilterModels transforms the FilterRequest DTO array to the Filter model array
func AddFilterReqToFilterModels(addRequests []FilterRequest) (Filters []models.Filter) {
	for _, req := range addRequests {
		d := dtos.ToFilterModel(req.Filter)
		Filters = append(Filters, d)
	}
	return Filters
}
