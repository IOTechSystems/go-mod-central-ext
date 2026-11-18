// Copyright (C) 2022-2025 IOTech Ltd

package xrtmodels

import (
	"encoding/base64"
	"fmt"
	"math"
	"strconv"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/central/dbc"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"

	"github.com/spf13/cast"
)

func toEdgeXProperties(protocol string, protocolProperties map[string]any) map[string]string {
	intProperties, floatProperties, boolProperties := PropertyConversionList(protocol)

	edgexProperties := make(map[string]string)
	for k, v := range protocolProperties {
		edgexProperties[k] = cast.ToString(v)
	}

	for _, p := range intProperties {
		propertyValue, ok := protocolProperties[p]
		if ok {
			// if we use fmt.fmt.Sprintf("%v", propertyValue) to convert the float to string,
			// the 4194148 become 4.194148e+06 and dot(.), plus(+) are invalid for metadata
			// so we can use %.0f to convert the float without the decimal point
			edgexProperties[p] = fmt.Sprintf("%.0f", propertyValue)
		}
	}

	for _, p := range floatProperties {
		propertyValue, ok := protocolProperties[p]
		if ok {
			switch val := propertyValue.(type) {
			case float64:
				// The -1 as the third parameter tells the function to print the fewest digits necessary to accurately represent the float
				// For example:
				//   strconv.FormatFloat(5.2, 'f', -1, 64) -> 5.2
				//   fmt.Sprintf("%f",5.2) -> 5.200000
				edgexProperties[p] = strconv.FormatFloat(val, 'f', -1, 64)
			}
		}
	}

	for _, p := range boolProperties {
		propertyValue, ok := protocolProperties[p]
		if ok {
			edgexProperties[p] = cast.ToString(propertyValue)
		}
	}
	return edgexProperties
}

func PropertyConversionList(protocol string) ([]string, []string, []string) {
	var intProperties []string
	var floatProperties []string
	var boolProperties []string
	switch protocol {
	case common.BacnetIP, common.BacnetMSTP:
		intProperties = []string{common.BacnetDeviceInstance, common.BacnetPort}
	case common.Gps:
		intProperties = []string{common.GpsGpsdPort, common.GpsGpsdRetries, common.GpsGpsdConnTimeout, common.GpsGpsdRequestTimeout}
	case common.ModbusTcp:
		intProperties = []string{common.ModbusUnitID, common.ModbusPort, common.ModbusReadMaxHoldingRegisters,
			common.ModbusReadMaxInputRegisters, common.ModbusReadMaxBitsCoils, common.ModbusReadMaxBitsDiscreteInputs,
			common.ModbusWriteMaxHoldingRegisters, common.ModbusWriteMaxBitsCoils}
	case common.ModbusRtu:
		intProperties = []string{common.ModbusUnitID, common.ModbusBaudRate, common.ModbusDataBits, common.ModbusStopBits,
			common.ModbusReadMaxHoldingRegisters, common.ModbusReadMaxInputRegisters, common.ModbusReadMaxBitsCoils,
			common.ModbusReadMaxBitsDiscreteInputs, common.ModbusWriteMaxHoldingRegisters, common.ModbusWriteMaxBitsCoils}
	case common.Opcua:
		intProperties = []string{common.OpcuaRequestedSessionTimeout, common.OpcuaBrowseDepth, common.OpcuaConnectionReadingPostDelay, common.OpcuaReadBatchSize, common.OpcuaWriteBatchSize, common.OpcuaNodesPerBrowse}
		floatProperties = []string{common.OpcuaBrowsePublishInterval, common.OpcuaSessionKeepAliveInterval}
	case common.S7:
		intProperties = []string{common.S7Rack, common.S7Slot}
	case common.EtherNetIPExplicitConnected:
		intProperties = []string{common.EtherNetIPRPI}
		boolProperties = []string{common.EtherNetIPSaveValue}
	case common.EtherNetIPO2T, common.EtherNetIPT2O:
		intProperties = []string{common.EtherNetIPRPI}
	case common.EtherNetIPKey:
		intProperties = []string{common.EtherNetIPVendorID, common.EtherNetIPDeviceType, common.EtherNetIPProductCode,
			common.EtherNetIPMajorRevision, common.EtherNetIPMinorRevision}
	case dbc.Canbus:
		intProperties = []string{dbc.ID, dbc.DataSize, dbc.Port}
	}
	return intProperties, floatProperties, boolProperties
}

func ToEdgeXV2EventDTO(xrtEvent MultiResourcesResult) (edgexDtos.Event, errors.EdgeX) {
	event := edgexDtos.Event{
		DeviceName:  xrtEvent.Device,
		ProfileName: xrtEvent.Profile,
		SourceName:  xrtEvent.SourceName,
		Tags:        xrtEvent.Tags,
		Readings:    make([]edgexDtos.BaseReading, len(xrtEvent.Readings)),
	}

	index := 0
	for resourceName, reading := range xrtEvent.Readings {
		valueType, err := edgexCommon.NormalizeValueType(reading.Type)
		if err != nil {
			return event, errors.NewCommonEdgeXWrapper(err)
		}
		value, err := ParseXRTReadingValue(valueType, reading.Value)
		if err != nil {
			return event, errors.NewCommonEdgeXWrapper(err)
		}

		switch valueType {
		case edgexCommon.ValueTypeBinary:
			if data, ok := value.([]byte); ok {
				event.Readings[index] = edgexDtos.NewBinaryReading(xrtEvent.Profile, xrtEvent.Device, resourceName, data, "")
			} else {
				return event, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("invalid binary value '%v'", value), nil)
			}
		case edgexCommon.ValueTypeObject:
			event.Readings[index] = edgexDtos.NewObjectReading(xrtEvent.Profile, xrtEvent.Device, resourceName, value)
		case edgexCommon.ValueTypeObjectArray:
			event.Readings[index] = edgexDtos.NewObjectReadingWithArray(xrtEvent.Profile, xrtEvent.Device, resourceName, value)
		default:
			event.Readings[index], err = edgexDtos.NewSimpleReading(xrtEvent.Profile, xrtEvent.Device, resourceName, valueType, value)
			if err != nil {
				return event, errors.NewCommonEdgeXWrapper(err)
			}
		}
		event.Readings[index].Origin = reading.Origin
		event.Readings[index].Tags = reading.Tags
		event.Origin = reading.Origin
		index++
	}

	return event, nil
}

// toInt64 accepts an int64, as Reading.UnmarshalJSON decodes an Int64 value, or a float64, as a JSON number decodes
// into any, within int64's range; the conversion would silently wrap one beyond it.
func toInt64(v any) (int64, bool) {
	switch n := v.(type) {
	case float64:
		return int64(n), n >= -(1<<63) && n < 1<<63
	case int64:
		return n, true
	}
	return 0, false
}

// toUint64 accepts a uint64 or a float64 within uint64's range (see toInt64).
func toUint64(v any) (uint64, bool) {
	switch n := v.(type) {
	case float64:
		return uint64(n), n > -1 && n < 1<<64
	case uint64:
		return n, true
	}
	return 0, false
}

// toInteger accepts a float64 within T's range (see toInt64), checked by converting back, which is exact for T up to
// 32 bits but not beyond: 2^64 could convert to MaxUint64 and back to 2^64.
func toInteger[T int8 | int16 | int32 | uint8 | uint16 | uint32](v any) (T, bool) {
	f, ok := v.(float64)
	if !ok {
		return 0, false
	}
	t := T(f)
	return t, float64(t) == math.Trunc(f)
}

// ParseXRTReadingValue parses the XRT reading value to EdgeX reading value.
func ParseXRTReadingValue(valueType string, reading interface{}) (interface{}, errors.EdgeX) {
	// The reading is as JSON decodes into any: bool, float64, nil, string, []interface{} or map[string]interface{};
	// an Int64 or Uint64 value, or an element of their arrays, may also be int64 or uint64 (see Reading.UnmarshalJSON).
	switch valueType {
	case edgexCommon.ValueTypeString:
		return parseValue(reading, assertType[string], invalidString)
	case edgexCommon.ValueTypeBool:
		return parseValue(reading, assertType[bool], invalidBool)
	case edgexCommon.ValueTypeUint8:
		return parseValue(reading, toInteger[uint8], invalidNumber)
	case edgexCommon.ValueTypeUint16:
		return parseValue(reading, toInteger[uint16], invalidNumber)
	case edgexCommon.ValueTypeUint32:
		return parseValue(reading, toInteger[uint32], invalidNumber)
	case edgexCommon.ValueTypeUint64:
		return parseValue(reading, toUint64, invalidNumber)
	case edgexCommon.ValueTypeInt8:
		return parseValue(reading, toInteger[int8], invalidNumber)
	case edgexCommon.ValueTypeInt16:
		return parseValue(reading, toInteger[int16], invalidNumber)
	case edgexCommon.ValueTypeInt32:
		return parseValue(reading, toInteger[int32], invalidNumber)
	case edgexCommon.ValueTypeInt64:
		return parseValue(reading, toInt64, invalidNumber)
	case edgexCommon.ValueTypeFloat32:
		return parseValue(reading, toFloat32, invalidNumber)
	case edgexCommon.ValueTypeFloat64:
		return parseValue(reading, assertType[float64], invalidNumber)
	case edgexCommon.ValueTypeBinary:
		return parseBinary(reading)
	case edgexCommon.ValueTypeBoolArray:
		return parseArray(reading, assertType[bool], invalidBool)
	case edgexCommon.ValueTypeStringArray:
		return parseArray(reading, assertType[string], invalidString)
	case edgexCommon.ValueTypeUint8Array:
		if _, ok := reading.([]interface{}); !ok {
			return parseBase64Uint8Array(reading)
		}
		return parseArray(reading, toInteger[uint8], invalidNumber)
	case edgexCommon.ValueTypeUint16Array:
		return parseArray(reading, toInteger[uint16], invalidNumber)
	case edgexCommon.ValueTypeUint32Array:
		return parseArray(reading, toInteger[uint32], invalidNumber)
	case edgexCommon.ValueTypeUint64Array:
		return parseArray(reading, toUint64, invalidNumber)
	case edgexCommon.ValueTypeInt8Array:
		return parseArray(reading, toInteger[int8], invalidNumber)
	case edgexCommon.ValueTypeInt16Array:
		return parseArray(reading, toInteger[int16], invalidNumber)
	case edgexCommon.ValueTypeInt32Array:
		return parseArray(reading, toInteger[int32], invalidNumber)
	case edgexCommon.ValueTypeInt64Array:
		return parseArray(reading, toInt64, invalidNumber)
	case edgexCommon.ValueTypeFloat32Array:
		return parseArray(reading, toFloat32, invalidNumber)
	case edgexCommon.ValueTypeFloat64Array:
		return parseArray(reading, assertType[float64], invalidNumber)
	case edgexCommon.ValueTypeObject:
		return reading, nil
	case edgexCommon.ValueTypeObjectArray:
		if _, ok := reading.([]interface{}); !ok {
			return nil, invalidArray(reading)
		}
		return reading, nil
	}
	return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("none supported value type '%s'", valueType), nil)
}

// Formats of the error for a reading value, or an array element, of the wrong type.
const (
	invalidString = "invalid string value '%v'"
	invalidBool   = "invalid bool value '%v'"
	invalidNumber = "invalid number '%v'"
)

func assertType[T any](v any) (T, bool) {
	t, ok := v.(T)
	return t, ok
}

func toFloat32(v any) (float32, bool) {
	f, ok := v.(float64)
	return float32(f), ok
}

func parseValue[T any](reading any, convert func(any) (T, bool), format string) (any, errors.EdgeX) {
	v, ok := convert(reading)
	if !ok {
		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf(format, reading), nil)
	}
	return v, nil
}

func parseArray[T any](reading any, convert func(any) (T, bool), format string) (any, errors.EdgeX) {
	elems, ok := reading.([]interface{})
	if !ok {
		return nil, invalidArray(reading)
	}
	array := make([]T, len(elems))
	for i, elem := range elems {
		if array[i], ok = convert(elem); !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf(format, elem), nil)
		}
	}
	return array, nil
}

func invalidArray(reading any) errors.EdgeX {
	return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("invalid array '%v'", reading), nil)
}

// XRT transfers binary data in a base64 encoded string.
func parseBinary(reading any) (any, errors.EdgeX) {
	strValue, ok := reading.(string)
	if !ok {
		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf(invalidString, reading), nil)
	}
	data, err := base64.StdEncoding.DecodeString(strValue)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("fail to decode the base64 string '%v'", reading), err)
	}
	return data, nil
}

// parseBase64Uint8Array decodes a Uint8Array reading that is not an array, as base64.
func parseBase64Uint8Array(reading any) (any, errors.EdgeX) {
	data, err := base64.StdEncoding.DecodeString(fmt.Sprint(reading))
	if err != nil {
		return nil, invalidArray(reading)
	}
	return data, nil
}

func FromEdgeXToXRTReadings(readings []edgexDtos.BaseReading) (map[string]Reading, errors.EdgeX) {
	xrtReadings := make(map[string]Reading, len(readings))
	for _, r := range readings {
		var val any
		var err errors.EdgeX
		switch r.ValueType {
		case edgexCommon.ValueTypeObject, edgexCommon.ValueTypeObjectArray:
			val = r.ObjectValue
		case edgexCommon.ValueTypeBinary:
			val = r.BinaryValue
		default:
			val, err = common.ParseValueByDeviceResource(r.ValueType, r.Value)
			if err != nil {
				return nil, errors.NewCommonEdgeX(errors.Kind(err), "failed to convert from EdgeX reading to XRT reading", err)
			}
		}
		xrtReadings[r.ResourceName] = Reading{
			Value:  val,
			Type:   r.ValueType,
			Origin: r.Origin,
			Tags:   r.Tags,
		}
	}
	return xrtReadings, nil
}
