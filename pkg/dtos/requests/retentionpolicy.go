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

// RetentionPolicyRequest defines the Request Content for POST/PUT RetentionPolicy DTO
type RetentionPolicyRequest struct {
	dtoCommon.BaseRequest `json:",inline"`
	RetentionPolicy       dtos.RetentionPolicy `json:"retentionPolicy"`
}

// Validate satisfies the Validator interface
func (r *RetentionPolicyRequest) Validate() error {
	err := common.Validate(r)
	return err
}

// UnmarshalJSON implements the Unmarshaler interface for the FilterRequest type
func (r *RetentionPolicyRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		dtoCommon.BaseRequest
		RetentionPolicy dtos.RetentionPolicy
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}

	*r = RetentionPolicyRequest(alias)

	// Validate RetentionPolicyRequest DTO
	if err := r.Validate(); err != nil {
		return err
	}
	return nil
}

// RtPolicyReqToRtPolicyModels transforms the RetentionPolicyRequest DTO array to the RetentionPolicy model array
func RtPolicyReqToRtPolicyModels(addRequests []RetentionPolicyRequest) (policies []models.RetentionPolicy) {
	for _, req := range addRequests {
		rtPolicy := dtos.ToRetentionPolicyModel(req.RetentionPolicy)
		policies = append(policies, rtPolicy)
	}
	return policies
}
