// Copyright (C) 2021-2024 IOTech Ltd

package xrtmodels

import (
	"encoding/json"
	"testing"

	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRequest(t *testing.T) {
	profile := edgexDtos.DeviceProfile{}
	device := DeviceInfo{}
	clientName := "testClient"
	addProfileRequest, err := json.Marshal(NewProfileAddRequest(profile, clientName))
	require.NoError(t, err)
	updateProfileRequest, err := json.Marshal(NewProfileUpdateRequest(profile, clientName))
	require.NoError(t, err)
	getProfileRequest, err := json.Marshal(NewProfileGetRequest(profile.Name, clientName))
	require.NoError(t, err)
	deleteProfileRequest, err := json.Marshal(NewProfileDeleteRequest(profile.Name, clientName))
	require.NoError(t, err)
	addDeviceRequest, err := json.Marshal(NewDeviceAddRequest(device, clientName))
	require.NoError(t, err)
	updateDeviceRequest, err := json.Marshal(NewDeviceUpdateRequest(device, clientName))
	require.NoError(t, err)
	getDeviceRequest, err := json.Marshal(NewDeviceGetRequest(device.Name, clientName))
	require.NoError(t, err)
	deleteDeviceRequest, err := json.Marshal(NewDeviceDeleteRequest(device.Name, clientName))
	require.NoError(t, err)
	getDeviceResourceRequest, err := json.Marshal(NewDeviceResourceGetRequest(device.Name, clientName, []string{}))
	require.NoError(t, err)
	setDeviceResourceRequest, err := json.Marshal(NewDeviceResourceSetRequest(device.Name, clientName, map[string]any{}, map[string]any{}))
	require.NoError(t, err)
	categoryName := "IOT::Core"
	componentDiscoverReq, err := json.Marshal(NewComponentDiscoverRequest(clientName, categoryName))
	require.NoError(t, err)
	scanDeviceReq, err := json.Marshal(NewDeviceScanRequest(device, clientName, map[string]any{}))
	require.NoError(t, err)

	var tests = []struct {
		name       string
		data       []byte
		expectedOp string
	}{
		{"new AddProfileRequest", addProfileRequest, ProfileAddOperation},
		{"new UpdateProfileRequest", updateProfileRequest, ProfileUpdateOperation},
		{"new GetProfileRequest", getProfileRequest, ProfileGetOperation},
		{"new DeleteProfileRequest", deleteProfileRequest, ProfileDeleteOperation},
		{"new AddDeviceRequest", addDeviceRequest, DeviceAddOperation},
		{"new UpdateDeviceRequest", updateDeviceRequest, DeviceUpdateOperation},
		{"new GetDeviceRequest", getDeviceRequest, DeviceGetOperation},
		{"new DeleteDeviceRequest", deleteDeviceRequest, DeviceDeleteOperation},
		{"new GetDeviceResourceRequest", getDeviceResourceRequest, DeviceResourceGetOperation},
		{"new SetDeviceResourceRequest", setDeviceResourceRequest, DeviceResourceSetOperation},
		{"new ComponentDiscoverRequest", componentDiscoverReq, ComponentDiscoverOperation},
		{"new DeviceScanRequest", scanDeviceReq, DeviceScanOperation},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var request BaseRequest
			err = json.Unmarshal(testCase.data, &request)
			require.NoError(t, err)
			assert.Equal(t, clientName, request.Client)
			assert.Equal(t, testCase.expectedOp, request.Op)
		})
	}
}
