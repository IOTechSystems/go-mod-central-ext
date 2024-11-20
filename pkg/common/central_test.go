// Copyright (C) 2023 IOTech Ltd

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
)

func TestParseValueByDeviceResource(t *testing.T) {
	tests := []struct {
		name      string
		valueType string
		value     any
		expected  any
	}{
		{"string", edgexCommon.ValueTypeString, "foo and bar", "foo and bar"},
		{"uint8", edgexCommon.ValueTypeUint8, "127", uint8(127)},
		{"uint16", edgexCommon.ValueTypeUint16, "127", uint16(127)},
		{"uint32", edgexCommon.ValueTypeUint32, "127", uint32(127)},
		{"uint64", edgexCommon.ValueTypeUint64, "127", uint64(127)},
		{"int8", edgexCommon.ValueTypeInt8, "-127", int8(-127)},
		{"int16", edgexCommon.ValueTypeInt16, "-127", int16(-127)},
		{"int32", edgexCommon.ValueTypeInt32, "-127", int32(-127)},
		{"int64", edgexCommon.ValueTypeInt64, "-127", int64(-127)},
		{"float32", edgexCommon.ValueTypeFloat32, "0.123", float32(0.123)},
		{"float64", edgexCommon.ValueTypeFloat64, "0.123", 0.123},
		{"string array - EdgeX readings", edgexCommon.ValueTypeStringArray, "[foo, bar]", []string{"foo", "bar"}},
		{"string array - Set command payload ", edgexCommon.ValueTypeStringArray, "[\"foo\",\"bar\"]", []string{"foo", "bar"}},
		{"bool array", edgexCommon.ValueTypeBoolArray, "[true, false]", []bool{true, false}},
		{"uint8 array", edgexCommon.ValueTypeUint8Array, "[100, 127]", []uint8{100, 127}},
		{"uint16 array", edgexCommon.ValueTypeUint16Array, "[100, 127]", []uint16{100, 127}},
		{"uint32 array", edgexCommon.ValueTypeUint32Array, "[100, 127]", []uint32{100, 127}},
		{"uint64 array", edgexCommon.ValueTypeUint64Array, "[100, 127]", []uint64{100, 127}},
		{"int8 array", edgexCommon.ValueTypeInt8Array, "[-127, 127]", []int8{-127, 127}},
		{"int16 array", edgexCommon.ValueTypeInt16Array, "[-127, 127]", []int16{-127, 127}},
		{"int32 array", edgexCommon.ValueTypeInt32Array, "[-127, 127]", []int32{-127, 127}},
		{"int64 array", edgexCommon.ValueTypeInt64Array, "[-127, 127]", []int64{-127, 127}},
		{"float32 array", edgexCommon.ValueTypeFloat32Array, "[-0.123, 0.123]", []float32{-0.123, 0.123}},
		{"float64 array", edgexCommon.ValueTypeFloat64Array, "[-0.123, 0.123]", []float64{-0.123, 0.123}},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := ParseValueByDeviceResource(testCase.valueType, testCase.value)
			require.NoError(t, err)
			assert.Equal(t, testCase.expected, res)
		})
	}
}
