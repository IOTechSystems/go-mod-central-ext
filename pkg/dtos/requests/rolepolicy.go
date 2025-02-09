// Copyright (C) 2024-2025 IOTech Ltd

package requests

import (
	"encoding/json"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/models"
	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

// AddRolePolicyRequest defines the Request Content for POST RolePolicy DTO
type AddRolePolicyRequest struct {
	dtoCommon.BaseRequest `json:",inline"`
	RolePolicy            dtos.RolePolicy `json:"rolePolicy"`
}

// Validate satisfies the Validator interface
func (a *AddRolePolicyRequest) Validate() error {
	err := edgexCommon.Validate(a)
	return err
}

// UnmarshalJSON implements the Unmarshaler interface for the AddRolePolicyRequest type
func (a *AddRolePolicyRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		dtoCommon.BaseRequest
		RolePolicy dtos.RolePolicy
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}

	*a = AddRolePolicyRequest(alias)

	// Validate AddRolePolicyRequest DTO
	if err := a.Validate(); err != nil {
		return err
	}
	return nil
}

// AddRolePolicyReqToRolePolicyModels transforms the AddRolePolicyRequest DTO array to the RolePolicy model array
func AddRolePolicyReqToRolePolicyModels(addRequests []AddRolePolicyRequest) (rps []models.RolePolicy) {
	for _, req := range addRequests {
		d := dtos.ToRolePolicyModel(req.RolePolicy)
		rps = append(rps, d)
	}
	return rps
}

// AuthRouteRequest defines the Request Content for POST AuthRoute DTO
type AuthRouteRequest struct {
	dtoCommon.BaseRequest `json:",inline"`
	AuthRoute             dtos.AuthRoute `json:"authRoute"`
}

// Validate satisfies the Validator interface
func (a *AuthRouteRequest) Validate() error {
	err := edgexCommon.Validate(a)
	return err
}

// UnmarshalJSON implements the Unmarshaler interface for the AuthRouteRequest type
func (a *AuthRouteRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		dtoCommon.BaseRequest
		AuthRoute dtos.AuthRoute
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}

	*a = AuthRouteRequest(alias)
	// Validate AuthRouteRequest DTO
	if err := a.Validate(); err != nil {
		return err
	}
	return nil
}

// AuthGraphQLRequest defines the Request Content for POST AuthGraphQL DTO
type AuthGraphQLRequest struct {
	dtoCommon.BaseRequest `json:",inline"`
	dtos.AuthGraphQL      `json:",inline"`
}

// Validate satisfies the Validator interface
func (a *AuthGraphQLRequest) Validate() error {
	err := edgexCommon.Validate(a)
	return err
}

// UnmarshalJSON implements the Unmarshaler interface for the AuthRouteRequest type
func (a *AuthGraphQLRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		dtoCommon.BaseRequest
		dtos.AuthGraphQL
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}

	*a = AuthGraphQLRequest(alias)
	// Validate AuthGraphQLRequest DTO
	if err := a.Validate(); err != nil {
		return err
	}
	return nil
}
