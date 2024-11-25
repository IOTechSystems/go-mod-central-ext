// Copyright (C) 2024 IOTech Ltd

package interfaces

import (
	"context"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/requests"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/responses"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

// RolePolicyClient defines the interface for interactions with the role policy API endpoint on the IOTech security-proxy-auth service.
type RolePolicyClient interface {
	// Add adds new a role policy.
	Add(ctx context.Context, reqs requests.AddRolePolicyRequest) (common.BaseWithIdResponse, errors.EdgeX)
	// Update updates a role policy.
	Update(ctx context.Context, reqs requests.AddRolePolicyRequest) (common.BaseResponse, errors.EdgeX)
	// AllRolePolicies returns all role policies.
	// The result can be limited in a certain range by specifying the offset and limit parameters.
	// offset: The number of items to skip before starting to collect the result set. Default is 0.
	// limit: The number of items to return. Specify -1 will return all remaining items after offset. The maximum will be the MaxResultCount as defined in the configuration of service. Default is 20.
	AllRolePolicies(ctx context.Context, offset int, limit int) (responses.MultiRolePolicyResponse, errors.EdgeX)
	// RolePolicyByRole returns a role policy by role.
	RolePolicyByRole(ctx context.Context, name string) (responses.RolePolicyResponse, errors.EdgeX)
	// DeleteRolePolicyByRole deletes a role policy by role.
	DeleteRolePolicyByRole(ctx context.Context, name string) (common.BaseResponse, errors.EdgeX)
}
