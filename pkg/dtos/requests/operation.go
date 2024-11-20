//
// Copyright (C) 2021-2023 IOTech Ltd
//

package requests

import (
	"encoding/json"

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

// OperationRequest defines the Request Content for SMA POST Operation.
type OperationRequest struct {
	dtoCommon.BaseRequest `json:",inline"`
	ServiceName           string   `json:"serviceName" validate:"required"`
	Action                string   `json:"action" validate:"required"`
	CmdFlags              []string `json:"cmdFlags"` // to identify the additional CMD flags or argument of a CRI operation
}

// Validate satisfies the Validator interface
func (o *OperationRequest) Validate() error {
	err := edgexCommon.Validate(o)
	return err
}

// UnmarshalJSON implements the Unmarshaler interface for the OperationRequest type
func (o *OperationRequest) UnmarshalJSON(b []byte) error {
	alias := struct {
		dtoCommon.BaseRequest
		ServiceName string
		Action      string
		CmdFlags    []string
	}{}

	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}
	*o = OperationRequest(alias)

	if err := o.Validate(); err != nil {
		return err
	}

	return nil
}
