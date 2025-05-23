// Copyright (C) 2021 IOTech Ltd

package requests

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
)

var testOperationRequest = OperationRequest{
	BaseRequest: dtoCommon.BaseRequest{
		RequestId:   common.ExampleUUID,
		Versionable: dtoCommon.NewVersionable(),
	},
	ServiceName: common.TestServiceName,
	Action:      common.TestActionName,
}

func TestOperationRequest_Validate(t *testing.T) {
	valid := testOperationRequest
	noReqId := testOperationRequest
	noReqId.RequestId = ""
	invalidReqId := testOperationRequest
	invalidReqId.RequestId = "abc"
	noServiceName := testOperationRequest
	noServiceName.ServiceName = ""
	noAction := testOperationRequest
	noAction.Action = ""

	tests := []struct {
		name        string
		request     OperationRequest
		expectedErr bool
	}{
		{"valid", valid, false},
		{"valid - no Request Id", noReqId, false},
		{"invalid - RequestId is not an uuid", invalidReqId, true},
		{"invalid - no ServiceName", noServiceName, true},
		{"invalid - no Action", noAction, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOperationRequest_UnmarshalJSON(t *testing.T) {
	valid := testOperationRequest
	resultTestBytes, _ := json.Marshal(testOperationRequest)
	type args struct {
		data []byte
	}
	tests := []struct {
		name        string
		request     OperationRequest
		args        args
		expectedErr bool
	}{
		{"valid", valid, args{resultTestBytes}, false},
		{"invalid - empty data", OperationRequest{}, args{[]byte{}}, true},
		{"invalid - string data", OperationRequest{}, args{[]byte("Invalid OperationRequest")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expected = tt.request
			err := tt.request.UnmarshalJSON(tt.args.data)
			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, expected, tt.request, "Unmarshal did not result in expected AddProvisionWatcherRequest.")
			}
		})
	}
}
