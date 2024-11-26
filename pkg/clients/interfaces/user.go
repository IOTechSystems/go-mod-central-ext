// Copyright (C) 2024 IOTech Ltd

package interfaces

import (
	"context"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/requests"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/responses"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

// UserClient defines the interface for interactions with the user API endpoint on the IOTech security-proxy-auth service.
type UserClient interface {
	// Add adds new users.
	Add(ctx context.Context, reqs []requests.AddUserRequest) ([]common.BaseWithIdResponse, errors.EdgeX)
	// Update updates users.
	Update(ctx context.Context, reqs []requests.UpdateUserRequest) ([]common.BaseResponse, errors.EdgeX)
	// AllUsers returns all users.
	// The result can be limited in a certain range by specifying the offset and limit parameters.
	// offset: The number of items to skip before starting to collect the result set. Default is 0.
	// limit: The number of items to return. Specify -1 will return all remaining items after offset. The maximum will be the MaxResultCount as defined in the configuration of service. Default is 20.
	AllUsers(ctx context.Context, offset int, limit int) (responses.MultiUsersResponse, errors.EdgeX)
	// UserByName returns a user by username.
	UserByName(ctx context.Context, name string) (responses.UserResponse, errors.EdgeX)
	// DeleteUserByName deletes a user by username.
	DeleteUserByName(ctx context.Context, name string) (common.BaseResponse, errors.EdgeX)
	// Login logins a user by username and password
	Login(ctx context.Context, req requests.LoginRequest) (responses.TokenResponse, errors.EdgeX)
}
