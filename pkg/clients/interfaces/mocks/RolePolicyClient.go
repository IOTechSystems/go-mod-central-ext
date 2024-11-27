// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	common "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"

	errors "github.com/edgexfoundry/go-mod-core-contracts/v4/errors"

	mock "github.com/stretchr/testify/mock"

	requests "github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/requests"

	responses "github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/responses"
)

// RolePolicyClient is an autogenerated mock type for the RolePolicyClient type
type RolePolicyClient struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, reqs
func (_m *RolePolicyClient) Add(ctx context.Context, reqs requests.AddRolePolicyRequest) (common.BaseWithIdResponse, errors.EdgeX) {
	ret := _m.Called(ctx, reqs)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 common.BaseWithIdResponse
	var r1 errors.EdgeX
	if rf, ok := ret.Get(0).(func(context.Context, requests.AddRolePolicyRequest) (common.BaseWithIdResponse, errors.EdgeX)); ok {
		return rf(ctx, reqs)
	}
	if rf, ok := ret.Get(0).(func(context.Context, requests.AddRolePolicyRequest) common.BaseWithIdResponse); ok {
		r0 = rf(ctx, reqs)
	} else {
		r0 = ret.Get(0).(common.BaseWithIdResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, requests.AddRolePolicyRequest) errors.EdgeX); ok {
		r1 = rf(ctx, reqs)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(errors.EdgeX)
		}
	}

	return r0, r1
}

// AllRolePolicies provides a mock function with given fields: ctx, offset, limit
func (_m *RolePolicyClient) AllRolePolicies(ctx context.Context, offset int, limit int) (responses.MultiRolePolicyResponse, errors.EdgeX) {
	ret := _m.Called(ctx, offset, limit)

	if len(ret) == 0 {
		panic("no return value specified for AllRolePolicies")
	}

	var r0 responses.MultiRolePolicyResponse
	var r1 errors.EdgeX
	if rf, ok := ret.Get(0).(func(context.Context, int, int) (responses.MultiRolePolicyResponse, errors.EdgeX)); ok {
		return rf(ctx, offset, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int) responses.MultiRolePolicyResponse); ok {
		r0 = rf(ctx, offset, limit)
	} else {
		r0 = ret.Get(0).(responses.MultiRolePolicyResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int) errors.EdgeX); ok {
		r1 = rf(ctx, offset, limit)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(errors.EdgeX)
		}
	}

	return r0, r1
}

// DeleteRolePolicyByRole provides a mock function with given fields: ctx, name
func (_m *RolePolicyClient) DeleteRolePolicyByRole(ctx context.Context, name string) (common.BaseResponse, errors.EdgeX) {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for DeleteRolePolicyByRole")
	}

	var r0 common.BaseResponse
	var r1 errors.EdgeX
	if rf, ok := ret.Get(0).(func(context.Context, string) (common.BaseResponse, errors.EdgeX)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) common.BaseResponse); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Get(0).(common.BaseResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) errors.EdgeX); ok {
		r1 = rf(ctx, name)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(errors.EdgeX)
		}
	}

	return r0, r1
}

// RolePolicyByRole provides a mock function with given fields: ctx, name
func (_m *RolePolicyClient) RolePolicyByRole(ctx context.Context, name string) (responses.RolePolicyResponse, errors.EdgeX) {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for RolePolicyByRole")
	}

	var r0 responses.RolePolicyResponse
	var r1 errors.EdgeX
	if rf, ok := ret.Get(0).(func(context.Context, string) (responses.RolePolicyResponse, errors.EdgeX)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) responses.RolePolicyResponse); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Get(0).(responses.RolePolicyResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) errors.EdgeX); ok {
		r1 = rf(ctx, name)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(errors.EdgeX)
		}
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, reqs
func (_m *RolePolicyClient) Update(ctx context.Context, reqs requests.AddRolePolicyRequest) (common.BaseResponse, errors.EdgeX) {
	ret := _m.Called(ctx, reqs)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 common.BaseResponse
	var r1 errors.EdgeX
	if rf, ok := ret.Get(0).(func(context.Context, requests.AddRolePolicyRequest) (common.BaseResponse, errors.EdgeX)); ok {
		return rf(ctx, reqs)
	}
	if rf, ok := ret.Get(0).(func(context.Context, requests.AddRolePolicyRequest) common.BaseResponse); ok {
		r0 = rf(ctx, reqs)
	} else {
		r0 = ret.Get(0).(common.BaseResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, requests.AddRolePolicyRequest) errors.EdgeX); ok {
		r1 = rf(ctx, reqs)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(errors.EdgeX)
		}
	}

	return r0, r1
}

// NewRolePolicyClient creates a new instance of RolePolicyClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRolePolicyClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *RolePolicyClient {
	mock := &RolePolicyClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
