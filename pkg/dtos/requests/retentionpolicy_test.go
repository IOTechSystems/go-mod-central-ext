// Copyright (C) 2025 IOTech Ltd

package requests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/models"

	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"

	"github.com/stretchr/testify/require"
)

var (
	testDeviceName      = "testDevice"
	testSourceName      = "testSource"
	testDuration        = "10h"
	testRetentionPolicy = RetentionPolicyRequest{
		BaseRequest: dtoCommon.BaseRequest{
			RequestId:   common.ExampleUUID,
			Versionable: dtoCommon.NewVersionable(),
		},
		RetentionPolicy: dtos.RetentionPolicy{
			DeviceName: testDeviceName,
			SourceName: testSourceName,
			Duration:   testDuration,
		},
	}
)

func TestRetentionPolicyRequest_Validate(t *testing.T) {
	valid := testRetentionPolicy

	noSourceName := valid
	noSourceName.RetentionPolicy.SourceName = ""

	noDeviceName := valid
	noDeviceName.RetentionPolicy.DeviceName = ""

	withDayDuration := valid
	withDayDuration.RetentionPolicy.Duration = "1d20h"

	noDeviceSource := noSourceName
	noDeviceSource.RetentionPolicy.DeviceName = ""

	noDuration := valid
	noDuration.RetentionPolicy.Duration = ""

	invalidDuration := noDuration
	invalidDuration.RetentionPolicy.Duration = "xxx"

	tests := []struct {
		name            string
		RetentionPolicy RetentionPolicyRequest
		expectError     bool
	}{
		{"valid - RetentionPolicyRequest", valid, false},
		{"valid - RetentionPolicyRequest with no source name", noSourceName, false},
		{"valid - RetentionPolicyRequest with no device name", noDeviceName, false},
		{"valid - RetentionPolicyRequest with day duration", withDayDuration, false},
		{"invalid - RetentionPolicyRequest with neither device name nor source name", noDeviceSource, true},
		{"invalid - RetentionPolicyRequest with no duration", noDuration, true},
		{"invalid - RetentionPolicyRequest with invalid duration string", invalidDuration, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.RetentionPolicy.Validate()
			if tt.expectError {
				require.Error(t, err, fmt.Sprintf("expect error but not : %s", tt.name))
			} else {
				require.NoError(t, err, fmt.Sprintf("unexpected error occurs : %s", tt.name))
			}
		})
	}
}

func TestRetentionPolicyRequest_UnmarshalJSON(t *testing.T) {
	valid := testRetentionPolicy
	resultTestBytes, err := json.Marshal(testRetentionPolicy)
	require.NoError(t, err)

	type args struct {
		data []byte
	}

	result := RetentionPolicyRequest{}

	tests := []struct {
		name    string
		rp      *RetentionPolicyRequest
		args    args
		wantErr bool
	}{
		{"unmarshal RetentionPolicyRequest with success", &valid, args{resultTestBytes}, false},
		{"unmarshal invalid RetentionPolicyRequest, empty data", &result, args{[]byte{}}, true},
		{"unmarshal invalid RetentionPolicyRequest, string data", &result, args{[]byte("Invalid RetentionPolicyRequest")}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expected = tt.rp
			err := tt.rp.UnmarshalJSON(tt.args.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, expected, tt.rp, "Unmarshal did not result in expected RetentionPolicyRequest")
			}
		})
	}
}

func TestRtPolicyReqToRtPolicyModels(t *testing.T) {
	requests := []RetentionPolicyRequest{testRetentionPolicy}
	expectedRtPolicyModels := []models.RetentionPolicy{
		{
			DeviceName: testDeviceName,
			SourceName: testSourceName,
			Duration:   testDuration,
		},
	}
	resultModels := RtPolicyReqToRtPolicyModels(requests)
	require.Equal(t, expectedRtPolicyModels, resultModels, "RetentionPolicyRequest did not result in expected RetentionPolicy model")
}
