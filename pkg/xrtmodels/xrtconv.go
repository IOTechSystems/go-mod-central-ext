// Copyright (C) 2022-2025 IOTech Ltd

package xrtmodels

import (
	"encoding/base64"
	"fmt"
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

// ParseXRTReadingValue parses the XRT reading value to EdgeX reading value
func ParseXRTReadingValue(valueType string, reading interface{}) (interface{}, errors.EdgeX) {
	// Since we receive the reading in JSON format, the JSON lib will unmarshal the reading to specified data type:
	// bool for JSON booleans,  float64 for JSON numbers, nil for JSON null
	// string for JSON strings, []interface{} for JSON arrays, map[string]interface{} for JSON objects
	var err error
	var val interface{}
	numberErrMsg := fmt.Sprintf("invalid number '%v'", reading)
	arraryErrMsg := fmt.Sprintf("invalid array '%v'", reading)
	switch valueType {
	case edgexCommon.ValueTypeString:
		strValue, ok := reading.(string)
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("invalid string value '%v'", reading), nil)
		}
		val = strValue
	case edgexCommon.ValueTypeBool:
		boolValue, ok := reading.(bool)
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("invalid bool value '%v'", reading), nil)
		}
		val = boolValue
	case edgexCommon.ValueTypeUint8:
		float64Value, ok := reading.(float64)
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("invalid numbers '%v'", reading), nil)
		}
		val = uint8(float64Value)
	case edgexCommon.ValueTypeUint16:
		float64Value, ok := reading.(float64)
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("invalid numbers '%v'", reading), nil)
		}
		val = uint16(float64Value)
	case edgexCommon.ValueTypeUint32:
		float64Value, ok := reading.(float64)
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, numberErrMsg, nil)
		}
		val = uint32(float64Value)
	case edgexCommon.ValueTypeUint64:
		float64Value, ok := reading.(float64)
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, numberErrMsg, nil)
		}
		val = uint64(float64Value)
	case edgexCommon.ValueTypeInt8:
		float64Value, ok := reading.(float64)
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, numberErrMsg, nil)
		}
		val = int8(float64Value)
	case edgexCommon.ValueTypeInt16:
		float64Value, ok := reading.(float64)
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, numberErrMsg, nil)
		}
		val = int16(float64Value)
	case edgexCommon.ValueTypeInt32:
		float64Value, ok := reading.(float64)
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, numberErrMsg, nil)
		}
		val = int32(float64Value)
	case edgexCommon.ValueTypeInt64:
		float64Value, ok := reading.(float64)
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, numberErrMsg, nil)
		}
		val = int64(float64Value)
	case edgexCommon.ValueTypeFloat32:
		float64Value, ok := reading.(float64)
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, numberErrMsg, nil)
		}
		val = float32(float64Value)
	case edgexCommon.ValueTypeFloat64:
		float64Value, ok := reading.(float64)
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, numberErrMsg, nil)
		}
		val = float64Value
	case edgexCommon.ValueTypeBinary:
		// XRT transfer binary data in base64 encoded string
		strValue, ok := reading.(string)
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("invalid string value '%v'", reading), nil)
		}
		val, err = base64.StdEncoding.DecodeString(strValue)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("fail to decode the base64 string '%v'", reading), err)
		}
	case edgexCommon.ValueTypeBoolArray:
		interfaceArray, ok := reading.([]interface{})
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, arraryErrMsg, nil)
		}
		boolArray := make([]bool, len(interfaceArray))
		for i, v := range interfaceArray {
			boolValue, ok := v.(bool)
			if !ok {
				return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("invalid bool value '%v'", v), nil)
			}
			boolArray[i] = boolValue
		}
		val = boolArray
	case edgexCommon.ValueTypeStringArray:
		interfaceArray, ok := reading.([]interface{})
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, arraryErrMsg, nil)
		}
		stringArray := make([]string, len(interfaceArray))
		for i, v := range interfaceArray {
			strValue, ok := v.(string)
			if !ok {
				return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("invalid string value '%v'", v), nil)
			}
			stringArray[i] = strValue
		}
		val = stringArray
	case edgexCommon.ValueTypeUint8Array:
		interfaceArray, ok := reading.([]interface{})
		if !ok {
			// try to decode from the base64 string
			data, decodeErr := base64.StdEncoding.DecodeString(fmt.Sprint(reading))
			if decodeErr != nil {
				return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, arraryErrMsg, nil)
			}
			val = data
			break
		}
		uint8Array := make([]uint8, len(interfaceArray))
		for i, v := range interfaceArray {
			float64Value, ok := v.(float64)
			if !ok {
				return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("invalid nunmber '%v'", v), nil)
			}
			uint8Array[i] = uint8(float64Value)
		}
		val = uint8Array
	case edgexCommon.ValueTypeUint16Array:
		interfaceArray, ok := reading.([]interface{})
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, arraryErrMsg, nil)
		}
		uint16Array := make([]uint16, len(interfaceArray))
		for i, v := range interfaceArray {
			float64Value, ok := v.(float64)
			if !ok {
				return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("invalid nunmber '%v'", v), nil)
			}
			uint16Array[i] = uint16(float64Value)
		}
		val = uint16Array
	case edgexCommon.ValueTypeUint32Array:
		interfaceArray, ok := reading.([]interface{})
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, arraryErrMsg, nil)
		}
		uint32Array := make([]uint32, len(interfaceArray))
		for i, v := range interfaceArray {
			float64Value, ok := v.(float64)
			if !ok {
				return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("invalid nunmber '%v'", v), nil)
			}
			uint32Array[i] = uint32(float64Value)
		}
		val = uint32Array
	case edgexCommon.ValueTypeUint64Array:
		interfaceArray, ok := reading.([]interface{})
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, arraryErrMsg, nil)
		}
		uint64Array := make([]uint64, len(interfaceArray))
		for i, v := range interfaceArray {
			float64Value, ok := v.(float64)
			if !ok {
				return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("invalid nunmber '%v'", v), nil)
			}
			uint64Array[i] = uint64(float64Value)
		}
		val = uint64Array
	case edgexCommon.ValueTypeInt8Array:
		interfaceArray, ok := reading.([]interface{})
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, arraryErrMsg, nil)
		}
		int8Array := make([]int8, len(interfaceArray))
		for i, v := range interfaceArray {
			float64Value, ok := v.(float64)
			if !ok {
				return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("invalid nunmber '%v'", v), nil)
			}
			int8Array[i] = int8(float64Value)
		}
		val = int8Array
	case edgexCommon.ValueTypeInt16Array:
		interfaceArray, ok := reading.([]interface{})
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, arraryErrMsg, nil)
		}
		int16Array := make([]int16, len(interfaceArray))
		for i, v := range interfaceArray {
			float64Value, ok := v.(float64)
			if !ok {
				return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("invalid number '%v'", v), nil)
			}
			int16Array[i] = int16(float64Value)
		}
		val = int16Array
	case edgexCommon.ValueTypeInt32Array:
		interfaceArray, ok := reading.([]interface{})
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, arraryErrMsg, nil)
		}
		int32Array := make([]int32, len(interfaceArray))
		for i, v := range interfaceArray {
			float64Value, ok := v.(float64)
			if !ok {
				return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("invalid number '%v'", v), nil)
			}
			int32Array[i] = int32(float64Value)
		}
		val = int32Array
	case edgexCommon.ValueTypeInt64Array:
		interfaceArray, ok := reading.([]interface{})
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, arraryErrMsg, nil)
		}
		int64Array := make([]int64, len(interfaceArray))
		for i, v := range interfaceArray {
			float64Value, ok := v.(float64)
			if !ok {
				return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("invalid number '%v'", v), nil)
			}
			int64Array[i] = int64(float64Value)
		}
		val = int64Array
	case edgexCommon.ValueTypeFloat32Array:
		interfaceArray, ok := reading.([]interface{})
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, arraryErrMsg, nil)
		}
		float32Array := make([]float32, len(interfaceArray))
		for i, v := range interfaceArray {
			float64Value, ok := v.(float64)
			if !ok {
				return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("invalid number '%v'", v), nil)
			}
			float32Array[i] = float32(float64Value)
		}
		val = float32Array
	case edgexCommon.ValueTypeFloat64Array:
		interfaceArray, ok := reading.([]interface{})
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, arraryErrMsg, nil)
		}
		float64Array := make([]float64, len(interfaceArray))
		for i, v := range interfaceArray {
			float64Value, ok := v.(float64)
			if !ok {
				return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("invalid number '%v'", v), nil)
			}
			float64Array[i] = float64Value
		}
		val = float64Array
	case edgexCommon.ValueTypeObject:
		val = reading
	case edgexCommon.ValueTypeObjectArray:
		_, ok := reading.([]interface{})
		if !ok {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, arraryErrMsg, nil)
		}
		val = reading
	default:
		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("none supported value type '%s'", valueType), nil)
	}

	return val, nil
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
