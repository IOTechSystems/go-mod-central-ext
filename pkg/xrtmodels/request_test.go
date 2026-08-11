// Copyright (C) 2021-2026 IOTech Ltd

package xrtmodels

import (
	"encoding/json"
	"reflect"
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

func TestNewDeviceScanRequest_OmitEmptyProfileName(t *testing.T) {
	device := DeviceInfo{}
	request := NewDeviceScanRequest(device, "testClient", nil)

	data, err := json.Marshal(request)
	require.NoError(t, err)

	var raw map[string]any
	err = json.Unmarshal(data, &raw)
	require.NoError(t, err)

	_, hasProfile := raw["profile"]
	assert.False(t, hasProfile, "profile field should be omitted when ProfileName is empty")
}

func TestNewDeviceScanRequest_IncludeProfileName(t *testing.T) {
	device := DeviceInfo{}
	device.ProfileName = "my-profile"
	request := NewDeviceScanRequest(device, "testClient", nil)

	data, err := json.Marshal(request)
	require.NoError(t, err)

	var raw map[string]any
	err = json.Unmarshal(data, &raw)
	require.NoError(t, err)

	profileVal, hasProfile := raw["profile"]
	assert.True(t, hasProfile, "profile field should be present when ProfileName is set")
	assert.Equal(t, "my-profile", profileVal)
}

func TestNewScheduleRequests_BaseFields(t *testing.T) {
	clientName := "testClient"

	readReq := NewScheduleReadRequest("my-schedule", clientName)
	readData, err := json.Marshal(readReq)
	require.NoError(t, err)

	updateReq := NewScheduleUpdateRequest(clientName, Schedule{
		Name:     "my-schedule",
		Device:   "dev1",
		Resource: []string{"temp"},
		Interval: 1000000,
	})
	updateData, err := json.Marshal(updateReq)
	require.NoError(t, err)

	tests := []struct {
		name       string
		data       []byte
		expectedOp string
	}{
		{"NewScheduleReadRequest", readData, ScheduleReadOperation},
		{"NewScheduleUpdateRequest", updateData, ScheduleUpdateOperation},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var base BaseRequest
			err = json.Unmarshal(tc.data, &base)
			require.NoError(t, err)
			assert.Equal(t, clientName, base.Client)
			assert.Equal(t, tc.expectedOp, base.Op)
		})
	}
}

func TestNewScheduleReadRequest_ScheduleField(t *testing.T) {
	req := NewScheduleReadRequest("target-schedule", "client1")
	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"schedule":"target-schedule"`)
}

func TestNewScheduleUpdateRequest_ScheduleField(t *testing.T) {
	s := Schedule{Name: "sched1", Device: "dev1", Resource: []string{"r"}}
	req := NewScheduleUpdateRequest("client1", s)
	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"sched1"`)
}

// The batch selector is one-of and every field is omitempty, because XRT treats a request
// with no selector as "all" — a zero-valued field must not be sent as an empty selector.
func TestNewBatchRequests(t *testing.T) {
	const client = "test-client"

	t.Run("read devices", func(t *testing.T) {
		tests := []struct {
			name    string
			names   []string
			pattern string
			expect  map[string]any
		}{
			{"no argument selects all", nil, "", map[string]any{}},
			{"by names", []string{"a", "b"}, "", map[string]any{"devices": []any{"a", "b"}}},
			{"by pattern", nil, ".*", map[string]any{"pattern": ".*"}},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				request := NewBatchReadDevicesRequest(tc.names, tc.pattern, client)
				assertBatchRequest(t, request, request.BaseRequest, BatchReadDevicesOperation, tc.expect)
			})
		}
	})

	t.Run("delete devices", func(t *testing.T) {
		request := NewBatchDeleteDevicesRequest([]string{"a"}, client)
		assertBatchRequest(t, request, request.BaseRequest, BatchDeleteDevicesOperation,
			map[string]any{"devices": []any{"a"}})
	})

	t.Run("read schedules", func(t *testing.T) {
		request := NewBatchReadSchedulesRequest(nil, "dev-1", "", client)
		assertBatchRequest(t, request, request.BaseRequest, BatchReadSchedulesOperation,
			map[string]any{"device": "dev-1"})
	})

	// Delete offers no pattern and no "all", so the selector cannot express either.
	t.Run("delete schedules by names", func(t *testing.T) {
		request := NewBatchDeleteSchedulesRequest([]string{"s1"}, "", client)
		assertBatchRequest(t, request, request.BaseRequest, BatchDeleteSchedulesOperation,
			map[string]any{"schedules": []any{"s1"}})
	})

	t.Run("add requests carry their payload", func(t *testing.T) {
		devices := NewBatchAddDevicesRequest([]DeviceInfo{{Device: edgexDtos.Device{Name: "a"}}}, client)
		if devices.Op != BatchAddDevicesOperation || len(devices.Devices) != 1 {
			t.Errorf("got %+v, want one device for %s", devices, BatchAddDevicesOperation)
		}
		schedules := NewBatchAddSchedulesRequest([]Schedule{{Name: "s1"}}, client)
		if schedules.Op != BatchAddSchedulesOperation || len(schedules.Schedules) != 1 {
			t.Errorf("got %+v, want one schedule for %s", schedules, BatchAddSchedulesOperation)
		}
	})
}

// assertBatchRequest checks the operation and that exactly the expected selector fields
// are present on the wire. The request is taken as any so it serves both the device and
// schedule request types.
func assertBatchRequest(t *testing.T, request any, base BaseRequest, expectedOp string, expectedFields map[string]any) {
	t.Helper()

	if base.Op != expectedOp {
		t.Errorf("op: got %q, want %q", base.Op, expectedOp)
	}
	if base.RequestId == "" || base.Client == "" {
		t.Errorf("request must carry a client and request id, got %+v", base)
	}

	encoded, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, field := range []string{"devices", "schedules", "device", "pattern"} {
		got, present := decoded[field]
		want, expected := expectedFields[field]
		if expected != present {
			t.Errorf("field %q: present=%v, want present=%v", field, present, expected)
			continue
		}
		if expected && !reflect.DeepEqual(got, want) {
			t.Errorf("field %q: got %v, want %v", field, got, want)
		}
	}
}

// The add_batch item name is derived from the device, so it cannot disagree with the
// device info XRT receives alongside it.
func TestNewBatchAddDevicesRequestDerivesTheName(t *testing.T) {
	devices := []DeviceInfo{
		{Device: edgexDtos.Device{Name: "dev-1", ProfileName: "p"}},
		{Device: edgexDtos.Device{Name: "dev-2", ProfileName: "p"}},
	}

	request := NewBatchAddDevicesRequest(devices, "test-client")

	if len(request.Devices) != len(devices) {
		t.Fatalf("got %d items, want %d", len(request.Devices), len(devices))
	}
	for i, item := range request.Devices {
		if item.DeviceName != devices[i].Name {
			t.Errorf("[%d] name: got %q, want %q", i, item.DeviceName, devices[i].Name)
		}
		if item.DeviceInfo.Name != devices[i].Name {
			t.Errorf("[%d] info name: got %q, want %q", i, item.DeviceInfo.Name, devices[i].Name)
		}
	}
}

// Both add operations must put the same shape on the wire for an empty payload; XRT does
// not define whether null and [] are equivalent, so neither constructor should emit null.
func TestBatchAddRequestsEncodeEmptyPayloadConsistently(t *testing.T) {
	tests := []struct {
		name    string
		request any
		field   string
	}{
		{"devices", NewBatchAddDevicesRequest(nil, "c"), "devices"},
		{"schedules", NewBatchAddSchedulesRequest(nil, "c"), "schedules"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := json.Marshal(tc.request)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			var decoded map[string]json.RawMessage
			if err := json.Unmarshal(encoded, &decoded); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := string(decoded[tc.field]); got != "[]" {
				t.Errorf("%s: got %s, want []", tc.field, got)
			}
		})
	}
}
