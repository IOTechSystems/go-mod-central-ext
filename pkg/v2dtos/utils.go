// Copyright (C) 2020-2024 IOTech Ltd

package v2dtos

import (
	"fmt"
	"strings"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

// Constants related to Reading ValueTypes
const (
	ApiVersion     = "v2"
	CommaSeparator = ","

	ReadWrite_RW = "RW"
	ReadWrite_WR = "WR"

	ValueTypeBool         = "Bool"
	ValueTypeString       = "String"
	ValueTypeUint8        = "Uint8"
	ValueTypeUint16       = "Uint16"
	ValueTypeUint32       = "Uint32"
	ValueTypeUint64       = "Uint64"
	ValueTypeInt8         = "Int8"
	ValueTypeInt16        = "Int16"
	ValueTypeInt32        = "Int32"
	ValueTypeInt64        = "Int64"
	ValueTypeFloat32      = "Float32"
	ValueTypeFloat64      = "Float64"
	ValueTypeBinary       = "Binary"
	ValueTypeBoolArray    = "BoolArray"
	ValueTypeStringArray  = "StringArray"
	ValueTypeUint8Array   = "Uint8Array"
	ValueTypeUint16Array  = "Uint16Array"
	ValueTypeUint32Array  = "Uint32Array"
	ValueTypeUint64Array  = "Uint64Array"
	ValueTypeInt8Array    = "Int8Array"
	ValueTypeInt16Array   = "Int16Array"
	ValueTypeInt32Array   = "Int32Array"
	ValueTypeInt64Array   = "Int64Array"
	ValueTypeFloat32Array = "Float32Array"
	ValueTypeFloat64Array = "Float64Array"
	ValueTypeObject       = "Object"
	ValueTypeObjectArray  = "ObjectArray"
)

var valueTypes = []string{
	ValueTypeBool, ValueTypeString,
	ValueTypeUint8, ValueTypeUint16, ValueTypeUint32, ValueTypeUint64,
	ValueTypeInt8, ValueTypeInt16, ValueTypeInt32, ValueTypeInt64,
	ValueTypeFloat32, ValueTypeFloat64,
	ValueTypeBinary,
	ValueTypeBoolArray, ValueTypeStringArray,
	ValueTypeUint8Array, ValueTypeUint16Array, ValueTypeUint32Array, ValueTypeUint64Array,
	ValueTypeInt8Array, ValueTypeInt16Array, ValueTypeInt32Array, ValueTypeInt64Array,
	ValueTypeFloat32Array, ValueTypeFloat64Array,
	ValueTypeObject, ValueTypeObjectArray,
}

// // NormalizeValueType normalizes the valueType to upper camel case
func NormalizeValueType(valueType string) (string, error) {
	for _, v := range valueTypes {
		if strings.EqualFold(valueType, v) {
			return v, nil
		}
	}
	return "", errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("unable to normalize the unknown value type %s", valueType), nil)
}
