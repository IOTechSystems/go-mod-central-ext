// Copyright (C) 2023-2024 IOTech Ltd

package common

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

// ParseValueByDeviceResource parses value for the specified value type
func ParseValueByDeviceResource(valueType string, value any) (any, errors.EdgeX) {
	var err error

	// Support writing the null value for specific protocol like BACnet.
	// For example, the user send a put request with JSON body {"test-resource": null}, then the device service will receive the nil value and marshal to null before sending to the XRT.
	if value == nil {
		return nil, nil
	}

	v := fmt.Sprint(value)

	if valueType != edgexCommon.ValueTypeString && strings.TrimSpace(v) == "" {
		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("empty string is invalid for %v value type", valueType), nil)
	}

	errMsg := fmt.Sprintf("failed to convert set parameter %s to ValueType %s", v, valueType)
	switch valueType {
	case edgexCommon.ValueTypeString:
		return value, nil
	case edgexCommon.ValueTypeStringArray:
		var arr []string
		err = json.Unmarshal([]byte(v), &arr)
		if err != nil {
			// try to handle the nonstandard json format, for example, [foo, bar]
			strArr := strings.Split(strings.Trim(v, "[]"), ",")
			for _, u := range strArr {
				arr = append(arr, strings.TrimSpace(u))
			}
			return arr, nil
		}
		return arr, nil
	case edgexCommon.ValueTypeBool:
		boolVal, err := strconv.ParseBool(v)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
		}
		return boolVal, nil
	case edgexCommon.ValueTypeBoolArray:
		var arr []bool
		err = json.Unmarshal([]byte(v), &arr)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
		}
		return arr, nil
	case edgexCommon.ValueTypeUint8:
		var n uint64
		n, err = strconv.ParseUint(v, 10, 8)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
		}
		return uint8(n), nil
	case edgexCommon.ValueTypeUint8Array:
		var arr []uint8
		strArr := strings.Split(strings.Trim(v, "[]"), ",")
		for _, u := range strArr {
			n, err := strconv.ParseUint(strings.TrimSpace(u), 10, 8)
			if err != nil {
				return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
			}
			arr = append(arr, uint8(n))
		}
		return arr, nil
	case edgexCommon.ValueTypeUint16:
		var n uint64
		n, err = strconv.ParseUint(v, 10, 16)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
		}
		return uint16(n), nil
	case edgexCommon.ValueTypeUint16Array:
		var arr []uint16
		strArr := strings.Split(strings.Trim(v, "[]"), ",")
		for _, u := range strArr {
			n, err := strconv.ParseUint(strings.TrimSpace(u), 10, 16)
			if err != nil {
				return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
			}
			arr = append(arr, uint16(n))
		}
		return arr, nil
	case edgexCommon.ValueTypeUint32:
		var n uint64
		n, err = strconv.ParseUint(v, 10, 32)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
		}
		return uint32(n), nil
	case edgexCommon.ValueTypeUint32Array:
		var arr []uint32
		strArr := strings.Split(strings.Trim(v, "[]"), ",")
		for _, u := range strArr {
			n, err := strconv.ParseUint(strings.TrimSpace(u), 10, 32)
			if err != nil {
				return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
			}
			arr = append(arr, uint32(n))
		}
		return arr, nil
	case edgexCommon.ValueTypeUint64:
		var n uint64
		n, err = strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
		}
		return n, nil
	case edgexCommon.ValueTypeUint64Array:
		var arr []uint64
		strArr := strings.Split(strings.Trim(v, "[]"), ",")
		for _, u := range strArr {
			n, err := strconv.ParseUint(strings.TrimSpace(u), 10, 64)
			if err != nil {
				return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
			}
			arr = append(arr, n)
		}
		return arr, nil
	case edgexCommon.ValueTypeInt8:
		var n int64
		n, err = strconv.ParseInt(v, 10, 8)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
		}
		return int8(n), nil
	case edgexCommon.ValueTypeInt8Array:
		var arr []int8
		err = json.Unmarshal([]byte(v), &arr)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
		}
		return arr, nil
	case edgexCommon.ValueTypeInt16:
		var n int64
		n, err = strconv.ParseInt(v, 10, 16)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
		}
		return int16(n), nil
	case edgexCommon.ValueTypeInt16Array:
		var arr []int16
		err = json.Unmarshal([]byte(v), &arr)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
		}
		return arr, nil
	case edgexCommon.ValueTypeInt32:
		var n int64
		n, err = strconv.ParseInt(v, 10, 32)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
		}
		return int32(n), nil
	case edgexCommon.ValueTypeInt32Array:
		var arr []int32
		err = json.Unmarshal([]byte(v), &arr)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
		}
		return arr, nil
	case edgexCommon.ValueTypeInt64:
		var n int64
		n, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
		}
		return n, nil
	case edgexCommon.ValueTypeInt64Array:
		var arr []int64
		err = json.Unmarshal([]byte(v), &arr)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
		}
		return arr, nil
	case edgexCommon.ValueTypeFloat32:
		var val float64
		val, err = strconv.ParseFloat(v, 32)
		if err == nil {
			return float32(val), nil
		}
		if numError, ok := err.(*strconv.NumError); ok {
			if numError.Err == strconv.ErrRange {
				return nil, errors.NewCommonEdgeX(errors.KindServerError, "NumError", err)
			}
		}
		var decodedToBytes []byte
		decodedToBytes, err = base64.StdEncoding.DecodeString(v)
		if err == nil {
			var val float32
			val, err = float32FromBytes(decodedToBytes)
			if err != nil {
				return nil, errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("fail to parse %v to float32", v), err)
			} else if math.IsNaN(float64(val)) {
				return nil, errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("fail to parse %v to float32, unexpected result %v", v, val), nil)
			} else {
				return val, nil
			}
		}
	case edgexCommon.ValueTypeFloat32Array:
		var arr []float32
		err = json.Unmarshal([]byte(v), &arr)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
		}
		return arr, nil
	case edgexCommon.ValueTypeFloat64:
		var val float64
		val, err = strconv.ParseFloat(v, 64)
		if err == nil {
			return val, nil
		}
		if numError, ok := err.(*strconv.NumError); ok {
			if numError.Err == strconv.ErrRange {
				return nil, errors.NewCommonEdgeX(errors.KindServerError, "NumError", err)
			}
		}
		var decodedToBytes []byte
		decodedToBytes, err = base64.StdEncoding.DecodeString(v)
		if err == nil {
			val, err = float64FromBytes(decodedToBytes)
			if err != nil {
				return nil, errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("fail to parse %v to float64", v), err)
			} else if math.IsNaN(val) {
				return nil, errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("fail to parse %v to float64, unexpected result %v", v, val), nil)
			} else {
				return val, nil
			}
		}
	case edgexCommon.ValueTypeFloat64Array:
		var arr []float64
		err = json.Unmarshal([]byte(v), &arr)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
		}
		return arr, nil
	case edgexCommon.ValueTypeObject:
		return value, nil
	case edgexCommon.ValueTypeObjectArray:
		return value, nil
	default:
		return nil, errors.NewCommonEdgeX(errors.KindServerError, "unrecognized value type", nil)
	}
	return value, nil
}

func float32FromBytes(numericValue []byte) (float32, error) {
	var res float32
	reader := bytes.NewReader(numericValue)
	err := binary.Read(reader, binary.BigEndian, &res)
	return res, err
}

func float64FromBytes(numericValue []byte) (float64, error) {
	var res float64
	reader := bytes.NewReader(numericValue)
	err := binary.Read(reader, binary.BigEndian, &res)
	return res, err
}
