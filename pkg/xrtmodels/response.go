// Copyright (C) 2021-2026 IOTech Ltd

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

type ScheduleReadResponse struct {
	BaseResponse `json:",inline"`
	Result       ScheduleReadResult `json:"result"`
}

type ScheduleReadResult struct {
	BaseResult `json:",inline"`
	Schedule   Schedule `json:"schedule"`
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

// BatchDevicesResponse is the reply to device:read_batch. Entries are null for devices
// that do not exist, so the element type must be a pointer to keep the result aligned
// with the requested names.
type BatchDevicesResponse struct {
	BaseResponse `json:",inline"`
	Result       BatchDevicesResult `json:"result"`
}

type BatchDevicesResult struct {
	BaseResult `json:",inline"`
	Devices    []*DeviceInfo `json:"devices"`
}

// BatchDeviceResultsResponse is the reply to device:add_batch and device:delete_batch.
type BatchDeviceResultsResponse struct {
	BaseResponse `json:",inline"`
	Result       BatchDeviceResults `json:"result"`
}

// BatchDeviceResults holds the per-item outcomes. The embedded status reports whether the
// request was accepted and stays 0 when individual items fail.
type BatchDeviceResults struct {
	BaseResult `json:",inline"`
	Results    []BatchDeviceResult `json:"device_results"`
}

type BatchDeviceResult struct {
	BaseResult `json:",inline"`
	Device     string `json:"device"`
}

// BatchSchedulesResponse is the reply to schedule:read_batch, shaped as
// BatchDevicesResponse is.
type BatchSchedulesResponse struct {
	BaseResponse `json:",inline"`
	Result       BatchSchedulesResult `json:"result"`
}

type BatchSchedulesResult struct {
	BaseResult `json:",inline"`
	Schedules  []*Schedule `json:"schedules"`
}

// BatchScheduleResultsResponse is the reply to schedule:add_batch and
// schedule:delete_batch, with the same per-item reporting as the device operations.
type BatchScheduleResultsResponse struct {
	BaseResponse `json:",inline"`
	Result       BatchScheduleResults `json:"result"`
}

// BatchScheduleResults holds the per-item outcomes, as BatchDeviceResults does.
type BatchScheduleResults struct {
	BaseResult `json:",inline"`
	Results    []BatchScheduleResult `json:"schedule_results"`
}

type BatchScheduleResult struct {
	BaseResult `json:",inline"`
	Schedule   string `json:"schedule"`
}

// BatchItemResult is the outcome of one item in a batch operation. A nil error from the
// call means the request was processed, not that every item succeeded.
type BatchItemResult struct {
	Name   string
	Status int
	Err    errors.EdgeX // set when Status is non-zero
}

// Failed reports whether this item's operation failed.
func (result BatchItemResult) Failed() bool {
	return result.Status != 0
}

// NewBatchItemResult converts a per-item result, reusing BaseResult.Error so a failure
// carries the same error kind as the equivalent single-item operation.
func NewBatchItemResult(name string, result BaseResult) BatchItemResult {
	return BatchItemResult{
		Name:   name,
		Status: result.Status,
		Err:    result.Error(),
	}
}
