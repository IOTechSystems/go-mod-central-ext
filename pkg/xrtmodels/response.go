// Copyright (C) 2021-2024 IOTech Ltd

package xrtmodels

import (
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

const (
	// https://docs.iotechsys.com/edge-xrt21/mqtt-management/mqtt-management.html#general-result-format
	XrtSdkStatusOk               = 0
	XrtSdkStatusNotFound         = 1
	XrtSdkStatusNotSupported     = 2
	XrtSdkStatusInvalidOperation = 3
	XrtSdkStatusAlreadyExists    = 7
	XrtSdkStatusServerError      = 500 // server error code for uncovered XRT error
)

type BaseResponse struct {
	Client    string `json:"client"`
	RequestId string `json:"request_id"`
	Type      string `json:"type"`
}

type CommonResponse struct {
	BaseResponse `json:",inline"`
	Result       BaseResult `json:"result"`
}

type BaseResult struct {
	Status       int    `json:"status"`
	ErrorMessage string `json:"error,omitempty"`
}

func (result BaseResult) Error() errors.EdgeX {
	switch result.Status {
	case XrtSdkStatusOk:
		return nil
	case XrtSdkStatusNotFound:
		return errors.NewCommonEdgeX(errors.KindEntityDoesNotExist, result.ErrorMessage, nil)
	case XrtSdkStatusNotSupported:
		return errors.NewCommonEdgeX(errors.KindNotImplemented, result.ErrorMessage, nil)
	case XrtSdkStatusInvalidOperation:
		return errors.NewCommonEdgeX(errors.KindInvalidId, result.ErrorMessage, nil)
	case XrtSdkStatusAlreadyExists:
		return errors.NewCommonEdgeX(errors.KindDuplicateName, result.ErrorMessage, nil)
	default:
		return errors.NewCommonEdgeX(errors.KindServerError, result.ErrorMessage, nil)
	}
}

// XrtErrorCode returns the XRT error code from EdgeX error
func XrtErrorCode(err errors.EdgeX) int {
	switch errors.Kind(err) {
	case errors.KindEntityDoesNotExist:
		return XrtSdkStatusNotFound
	case errors.KindNotImplemented:
		return XrtSdkStatusNotSupported
	case errors.KindInvalidId:
		return XrtSdkStatusInvalidOperation
	case errors.KindDuplicateName:
		return XrtSdkStatusAlreadyExists
	default:
		return XrtSdkStatusServerError
	}
}

type MultiResourcesResponse struct {
	BaseResponse `json:",inline"`
	Result       MultiResourcesResult `json:"result"`
}

type MultiResourcesResult struct {
	BaseResult `json:",inline"`
	Device     string                 `json:"device"`
	Profile    string                 `json:"profile"`
	SourceName string                 `json:"sourceName"`
	Readings   map[string]Reading     `json:"readings"`
	Tags       map[string]interface{} `json:"tags"`
	Type       string                 `json:"type"`
}

type Reading struct {
	Value  interface{}            `json:"value"`
	Type   string                 `json:"type"`
	Origin int64                  `json:"origin"`
	Tags   map[string]interface{} `json:"tags"`
}

type MultiDevicesResponse struct {
	BaseResponse `json:",inline"`
	Result       MultiDevicesResult `json:"result"`
}

type MultiDevicesResult struct {
	BaseResult `json:",inline"`
	Devices    []string `json:"devices"`
}

type DeviceResponse struct {
	BaseResponse `json:",inline"`
	Result       DeviceResult `json:"result"`
}

type DeviceResult struct {
	BaseResult `json:",inline"`
	Device     DeviceInfo `json:"device"`
}

type DiscoveredDevicesResult struct {
	BaseResult `json:",inline"`
	Devices    map[string]DeviceInfo `json:"devices"`
}

type MultiProfilesResponse struct {
	BaseResponse `json:",inline"`
	Result       MultiProfilesResult `json:"result"`
}

type MultiProfilesResult struct {
	BaseResult `json:",inline"`
	Profiles   []string `json:"profiles"`
}

type ProfileResponse struct {
	BaseResponse `json:",inline"`
	Result       ProfileResult `json:"result"`
}

type ProfileResult struct {
	BaseResult `json:",inline"`
	Profile    edgexDtos.DeviceProfile `json:"profile"`
}

type MultiSchedulesResponse struct {
	BaseResponse `json:",inline"`
	Result       MultiSchedulesResult `json:"result"`
}

type MultiSchedulesResult struct {
	BaseResult `json:",inline"`
	Schedules  []string `json:"schedules"`
}

type MultiComponentsResponse struct {
	BaseResponse `json:",inline"`
	Result       ComponentsDiscoveryResponse `json:"result"`
}

// ComponentsDiscoveryResponse is used to reply the discovery:discover operation
type ComponentsDiscoveryResponse struct {
	Components []Component `json:"components"`
	NodeID     string      `json:"node_id"`
	ServerID   string      `json:"server_id"`
	Type       string      `json:"type"`
}
