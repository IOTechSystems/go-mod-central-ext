// Copyright (C) 2021-2026 IOTech Ltd

package xrtmodels

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
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

// UnmarshalJSON decodes an Int64 or Uint64 value, and each element of their arrays, as an exact int64 or uint64. As
// the float64 a JSON number otherwise becomes, a value beyond 2^53 is rounded, and the conversion back can then
// overflow (MaxInt64 comes back as MinInt64). Every other value decodes as before.
func (r *Reading) UnmarshalJSON(data []byte) error {
	type plain Reading // without this method
	var raw struct {
		plain
		Value json.RawMessage `json:"value"` // shadows plain.Value
	}
	raw.plain = plain(*r)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*r = Reading(raw.plain)
	if len(raw.Value) == 0 {
		return nil
	}
	valueType, _ := edgexCommon.NormalizeValueType(r.Type)
	switch valueType {
	case edgexCommon.ValueTypeInt64, edgexCommon.ValueTypeUint64,
		edgexCommon.ValueTypeInt64Array, edgexCommon.ValueTypeUint64Array:
		decoder := json.NewDecoder(bytes.NewReader(raw.Value))
		decoder.UseNumber()
		if err := decoder.Decode(&r.Value); err != nil {
			return err
		}
		unsigned := valueType == edgexCommon.ValueTypeUint64 || valueType == edgexCommon.ValueTypeUint64Array
		var err error
		r.Value, err = exactIntegers(r.Value, unsigned)
		return err
	}
	return json.Unmarshal(raw.Value, &r.Value)
}

// exactIntegers turns each json.Number in v, a value or a slice of them, into an int64 or uint64. Any other number,
// such as 1.5, becomes the float64 it would have decoded to, or fails as it would have, beyond float64's range.
func exactIntegers(v any, unsigned bool) (any, error) {
	switch value := v.(type) {
	case json.Number:
		return exactInteger(value, unsigned)
	case []any:
		for i, elem := range value {
			var err error
			if value[i], err = exactIntegers(elem, unsigned); err != nil {
				return nil, err
			}
		}
	}
	return v, nil
}

func exactInteger(n json.Number, unsigned bool) (any, error) {
	v, err := parseInteger(n.String(), unsigned)
	if err != nil {
		if digits, ok := integerDigits(n.String()); ok {
			v, err = parseInteger(digits, unsigned)
		}
	}
	if err != nil {
		return n.Float64()
	}
	return v, nil
}

func parseInteger(s string, unsigned bool) (any, error) {
	if unsigned {
		return strconv.ParseUint(s, 10, 64)
	}
	return strconv.ParseInt(s, 10, 64)
}

// integerDigits rewrites a valid JSON number, such as 9007199254740993e0 or 1.5e1, as plain integer digits, which
// ParseInt and ParseUint take but not exponent or fraction forms. ok is false for a non-integer or one over 20 digits,
// beyond uint64; the bound also keeps a huge exponent from building a huge string.
func integerDigits(n string) (digits string, ok bool) {
	sign := ""
	if n[0] == '-' {
		sign, n = "-", n[1:]
	}
	mantissa, exponent := n, "0"
	if i := strings.IndexAny(n, "eE"); i >= 0 {
		mantissa, exponent = n[:i], n[i+1:]
	}
	whole, frac, _ := strings.Cut(mantissa, ".")
	digits = strings.TrimLeft(whole+frac, "0")
	if digits == "" {
		return "0", true // before Atoi, which fails on an exponent beyond int
	}
	exp, err := strconv.Atoi(exponent)
	if err != nil {
		return "", false
	}
	point := len(digits) - len(frac) + exp // digits before the decimal point
	digits = strings.TrimRight(digits, "0")
	if len(digits) > point || point > 20 {
		return "", false
	}
	return sign + digits + strings.Repeat("0", point-len(digits)), true
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

type BatchScheduleResults struct {
	BaseResult `json:",inline"`
	Results    []BatchScheduleResult `json:"schedules"`
}

// BatchScheduleResult names the schedule an outcome belongs to, except on add_batch,
// where XRT omits the name and results map by position.
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
