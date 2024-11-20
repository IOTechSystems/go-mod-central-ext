// Copyright (C) 2021-2024 IOTech Ltd

package xrtmodels

import (
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"

	"github.com/google/uuid"
)

const (
	ProfileAddOperation    = "profile:add"
	ProfileUpdateOperation = "profile:update"
	ProfileListOperation   = "profile:list"
	ProfileGetOperation    = "profile:read"
	ProfileDeleteOperation = "profile:delete"

	DeviceAddOperation         = "device:add"
	DeviceUpdateOperation      = "device:update"
	DeviceResourceGetOperation = "device:get"
	DeviceGetOperation         = "device:read"
	DeviceResourceSetOperation = "device:put"
	DeviceDeleteOperation      = "device:delete"
	DeviceListOperation        = "device:list"
	DeviceScanOperation        = "device:scan"

	ScheduleAddOperation    = "schedule:add"
	ScheduleListOperation   = "schedule:list"
	ScheduleDeleteOperation = "schedule:delete"

	DiscoveryTriggerOperation = "discovery:trigger"

	ComponentUpdateOperation    = "component:update"
	ComponentDiscoverOperation  = "component:discover" // Discover state and configuration of all components in an instance
	ComponentDiscoveryOperation = "discovery:discover" // Discover state and configuration of all components from all Xrt instances
)

type BaseRequest struct {
	Client    string `json:"client"`
	RequestId string `json:"request_id"`
	Op        string `json:"op"`
}

type AddProfileRequest struct {
	BaseRequest `json:",inline"`
	Profile     edgexDtos.DeviceProfile `json:"profile"`
}

type UpdateProfileRequest struct {
	BaseRequest `json:",inline"`
	Profile     edgexDtos.DeviceProfile `json:"profile"`
}

type ProfileRequest struct {
	BaseRequest `json:",inline"`
	Profile     string `json:"profile"`
}

type AddDeviceRequest struct {
	BaseRequest `json:",inline"`
	DeviceName  string     `json:"device"`
	DeviceInfo  DeviceInfo `json:"device_info"`
}

type AddDiscoveredDeviceRequest struct {
	BaseRequest `json:",inline"`
	DeviceName  string               `json:"device"`
	DeviceInfo  DiscoveredDeviceInfo `json:"device_info"`
}

type DiscoveredDeviceInfo struct {
	Protocols map[string]edgexDtos.ProtocolProperties `json:"protocols"`
}

type ScanDeviceRequest struct {
	BaseRequest `json:",inline"`
	DeviceName  string         `json:"device"`
	ProfileName string         `json:"profile"`
	Options     map[string]any `json:"options"`
}

type UpdateDeviceRequest struct {
	BaseRequest `json:",inline"`
	DeviceName  string     `json:"device"`
	DeviceInfo  DeviceInfo `json:"device_info"`
}

type DeviceRequest struct {
	BaseRequest `json:",inline"`
	Device      string `json:"device"`
}

type GetResourcesRequest struct {
	BaseRequest `json:",inline"`
	DeviceName  string   `json:"device"`
	Resource    []string `json:"resource"`
}

type GetResourceRequest struct {
	BaseRequest `json:",inline"`
	DeviceName  string `json:"device"`
	Resource    string `json:"resource"`
}

type PutResourceRequest struct {
	BaseRequest `json:",inline"`
	DeviceName  string                 `json:"device"`
	Values      map[string]interface{} `json:"values"`
	Options     map[string]interface{} `json:"options"`
}

type ScheduleRequest struct {
	BaseRequest `json:",inline"`
	Schedule    string `json:"schedule"`
}

type AddScheduleRequest struct {
	BaseRequest `json:",inline"`
	Schedule    Schedule `json:"schedule"`
}

type UpdateComponentRequest struct {
	BaseRequest `json:",inline"`
	Component   string                 `json:"component"`
	Config      map[string]interface{} `json:"config"`
}

type DiscoverComponentRequest struct {
	BaseRequest `json:",inline"`
	Category    string `json:"category,omitempty"`
}

type DiscoveryRequest struct {
	BaseRequest `json:",inline"`
	Options     map[string]interface{} `json:"options"`
}

func NewBaseRequest(op string, clientName string) BaseRequest {
	return BaseRequest{
		Client:    clientName,
		RequestId: uuid.New().String(),
		Op:        op,
	}
}

func NewAllProfilesRequest(clientName string) BaseRequest {
	req := BaseRequest{
		Client:    clientName,
		RequestId: uuid.New().String(),
		Op:        ProfileListOperation,
	}
	return req
}

// NewProfileAddRequest creates request with device profile
func NewProfileAddRequest(profile edgexDtos.DeviceProfile, clientName string) AddProfileRequest {
	req := AddProfileRequest{
		BaseRequest: BaseRequest{
			Client:    clientName,
			RequestId: uuid.New().String(),
			Op:        ProfileAddOperation,
		},
		Profile: profile,
	}

	return req
}

// NewProfileUpdateRequest creates request with v1 device profile
func NewProfileUpdateRequest(profile edgexDtos.DeviceProfile, clientName string) UpdateProfileRequest {
	req := UpdateProfileRequest{
		BaseRequest: BaseRequest{
			Client:    clientName,
			RequestId: uuid.New().String(),
			Op:        ProfileUpdateOperation,
		},
		Profile: profile,
	}

	return req
}

func NewProfileGetRequest(profileName string, clientName string) ProfileRequest {
	return ProfileRequest{
		BaseRequest: BaseRequest{
			Client:    clientName,
			RequestId: uuid.New().String(),
			Op:        ProfileGetOperation,
		},
		Profile: profileName,
	}
}

func NewProfileDeleteRequest(profileName string, clientName string) ProfileRequest {
	req := ProfileRequest{
		BaseRequest: BaseRequest{
			Client:    clientName,
			RequestId: uuid.New().String(),
			Op:        ProfileDeleteOperation,
		},
		Profile: profileName,
	}
	return req
}

func NewDeviceAddRequest(device DeviceInfo, clientName string) AddDeviceRequest {
	deviceRequest := AddDeviceRequest{
		BaseRequest: BaseRequest{
			Client:    clientName,
			RequestId: uuid.New().String(),
			Op:        DeviceAddOperation,
		},
		DeviceName: device.Name,
		DeviceInfo: device,
	}
	return deviceRequest
}

// NewDiscoveredDeviceAddRequest creates a request to add the discovered device without profile before sending the device:scan to generate the profile
func NewDiscoveredDeviceAddRequest(device DeviceInfo, clientName string) AddDiscoveredDeviceRequest {
	deviceRequest := AddDiscoveredDeviceRequest{
		BaseRequest: BaseRequest{
			Client:    clientName,
			RequestId: uuid.New().String(),
			Op:        DeviceAddOperation,
		},
		DeviceName: device.Name,
		DeviceInfo: DiscoveredDeviceInfo{Protocols: device.Protocols},
	}
	return deviceRequest
}

// NewDeviceScanRequest creates a request to scan the device and generate the profile
func NewDeviceScanRequest(device DeviceInfo, clientName string, options map[string]any) ScanDeviceRequest {
	deviceRequest := ScanDeviceRequest{
		BaseRequest: BaseRequest{
			Client:    clientName,
			RequestId: uuid.New().String(),
			Op:        DeviceScanOperation,
		},
		DeviceName:  device.Name,
		ProfileName: device.ProfileName,
		Options:     options,
	}
	return deviceRequest
}

func NewDeviceUpdateRequest(device DeviceInfo, clientName string) UpdateDeviceRequest {
	deviceRequest := UpdateDeviceRequest{
		BaseRequest: BaseRequest{
			Client:    clientName,
			RequestId: uuid.New().String(),
			Op:        DeviceUpdateOperation,
		},
		DeviceName: device.Name,
		DeviceInfo: device,
	}
	return deviceRequest
}

func NewAllDevicesRequest(clientName string) BaseRequest {
	req := BaseRequest{
		Client:    clientName,
		RequestId: uuid.New().String(),
		Op:        DeviceListOperation,
	}
	return req
}

func NewDeviceGetRequest(deviceName string, clientName string) DeviceRequest {
	req := DeviceRequest{
		BaseRequest: BaseRequest{
			Client:    clientName,
			RequestId: uuid.New().String(),
			Op:        DeviceGetOperation,
		},
		Device: deviceName,
	}
	return req
}

func NewDeviceDeleteRequest(deviceName string, clientName string) DeviceRequest {
	req := DeviceRequest{
		BaseRequest: BaseRequest{
			Client:    clientName,
			RequestId: uuid.New().String(),
			Op:        DeviceDeleteOperation,
		},
		Device: deviceName,
	}
	return req
}

func NewDeviceResourceGetRequest(deviceName string, clientName string, resources []string) GetResourcesRequest {
	req := GetResourcesRequest{
		BaseRequest: BaseRequest{
			Client:    clientName,
			RequestId: uuid.New().String(),
			Op:        DeviceResourceGetOperation,
		},
		DeviceName: deviceName,
		Resource:   nil,
	}
	req.Resource = resources
	return req
}

func NewDeviceResourceSetRequest(deviceName string, clientName string, values map[string]interface{}, options map[string]interface{}) PutResourceRequest {
	req := PutResourceRequest{
		BaseRequest: BaseRequest{
			Client:    clientName,
			RequestId: uuid.New().String(),
			Op:        DeviceResourceSetOperation,
		},
		DeviceName: deviceName,
		Values:     values,
		Options:    options,
	}
	return req
}

func NewAllSchedulesRequest(clientName string) BaseRequest {
	req := BaseRequest{
		Client:    clientName,
		RequestId: uuid.New().String(),
		Op:        ScheduleListOperation,
	}
	return req
}

func NewScheduleAddRequest(clientName string, schedule Schedule) AddScheduleRequest {
	req := AddScheduleRequest{
		BaseRequest: BaseRequest{
			Client:    clientName,
			RequestId: uuid.New().String(),
			Op:        ScheduleAddOperation,
		},
		Schedule: schedule,
	}
	return req
}

func NewScheduleDeleteRequest(scheduleName string, clientName string) ScheduleRequest {
	req := ScheduleRequest{
		BaseRequest: BaseRequest{
			Client:    clientName,
			RequestId: uuid.New().String(),
			Op:        ScheduleDeleteOperation,
		},
		Schedule: scheduleName,
	}
	return req
}

func NewComponentUpdateRequest(Component string, clientName string, config map[string]interface{}) UpdateComponentRequest {
	req := UpdateComponentRequest{
		BaseRequest: BaseRequest{
			Client:    clientName,
			RequestId: uuid.New().String(),
			Op:        ComponentUpdateOperation,
		},
		Component: Component,
		Config:    config,
	}
	return req
}

func NewComponentDiscoverRequest(clientName string, category string) DiscoverComponentRequest {
	req := DiscoverComponentRequest{
		BaseRequest: BaseRequest{
			Client:    clientName,
			RequestId: uuid.New().String(),
			Op:        ComponentDiscoverOperation,
		},
		Category: category,
	}
	return req
}

func NewDiscoveryRequest(clientName string, options map[string]interface{}) DiscoveryRequest {
	req := DiscoveryRequest{
		BaseRequest: BaseRequest{
			Client:    clientName,
			RequestId: uuid.New().String(),
			Op:        DiscoveryTriggerOperation,
		},
		Options: options,
	}
	return req
}
