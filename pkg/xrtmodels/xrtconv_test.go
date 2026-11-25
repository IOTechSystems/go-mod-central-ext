// Copyright (C) 2022-2026 IOTech Ltd

package xrtmodels

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
	edgexModels "github.com/edgexfoundry/go-mod-core-contracts/v4/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToXrtDevice(t *testing.T) {
	deviceName := "test-device"
	profileName := "test-profile"
	serviceName := "device-bacnet-ip"
	device := edgexDtos.Device{
		Name: deviceName,
		Protocols: map[string]edgexDtos.ProtocolProperties{
			common.BacnetIP: {
				common.BacnetDeviceInstance: 1234,
			},
		},
		ProfileName:    profileName,
		ServiceName:    serviceName,
		AdminState:     edgexModels.Unlocked,
		OperatingState: edgexModels.Up,
		Properties: map[string]any{
			common.ProtocolName: common.BacnetIP,
		},
	}
	xrtDevice, err := ToXrtDevice(device)
	require.NoError(t, err)

	assert.Equal(t, deviceName, xrtDevice.Name)
	assert.Equal(t, common.BacnetIP, xrtDevice.Properties[common.ProtocolName])
	assert.Equal(t, float64(1234), xrtDevice.Protocols[common.BacnetIP][common.BacnetDeviceInstance])
	assert.Equal(t, profileName, xrtDevice.ProfileName)
	assert.Equal(t, serviceName, xrtDevice.ServiceName)
	assert.Equal(t, edgexModels.Unlocked, xrtDevice.AdminState)
	assert.Equal(t, edgexModels.Up, xrtDevice.OperatingState)
}

func TestToEdgeXProperties(t *testing.T) {
	tests := []struct {
		protocol   string
		properties map[string]interface{}
		expected   map[string]string
	}{
		{
			protocol: common.BacnetIP,
			properties: map[string]interface{}{
				common.BacnetDeviceInstance: float64(1234),
			},
			expected: map[string]string{
				common.BacnetDeviceInstance: "1234",
			},
		},
		{
			protocol: common.BacnetIP,
			properties: map[string]interface{}{
				common.BacnetDeviceInstance: float64(4194302),
			},
			expected: map[string]string{
				common.BacnetDeviceInstance: "4194302",
			},
		},
		{
			protocol: common.Opcua,
			properties: map[string]interface{}{
				common.OpcuaRequestedSessionTimeout:    float64(1200000),
				common.OpcuaBrowseDepth:                float64(0),
				common.OpcuaBrowsePublishInterval:      5.2,
				common.OpcuaConnectionReadingPostDelay: float64(0),
				common.OpcuaIDType:                     float64(1),
			},
			expected: map[string]string{
				common.OpcuaRequestedSessionTimeout:    "1200000",
				common.OpcuaBrowseDepth:                "0",
				common.OpcuaBrowsePublishInterval:      "5.2",
				common.OpcuaConnectionReadingPostDelay: "0",
				common.OpcuaIDType:                     "1",
			},
		},
		{
			protocol: common.ModbusTcp,
			properties: map[string]interface{}{
				common.ModbusUnitID: float64(1),
				common.ModbusPort:   float64(1234),
			},
			expected: map[string]string{
				common.ModbusUnitID: "1",
				common.ModbusPort:   "1234",
			},
		},
		{
			protocol: common.ModbusRtu,
			properties: map[string]interface{}{
				common.ModbusUnitID:   float64(1),
				common.ModbusBaudRate: float64(0),
				common.ModbusDataBits: float64(5),
				common.ModbusStopBits: float64(1),
			},
			expected: map[string]string{
				common.ModbusUnitID:   "1",
				common.ModbusBaudRate: "0",
				common.ModbusDataBits: "5",
				common.ModbusStopBits: "1",
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.protocol, func(t *testing.T) {
			result := toEdgeXProperties(testCase.protocol, testCase.properties)
			assert.EqualValues(t, testCase.expected, result)
		})
	}
}

// A JSON reader gives float64; the Sparkplug client gives int64 and uint64, exact beyond 2^53.
func TestParseXRTReadingValue_64BitIntegers(t *testing.T) {
	for _, tc := range []struct {
		name      string
		valueType string
		reading   any
		want      any
	}{
		{"Int64 from JSON", edgexCommon.ValueTypeInt64, float64(-42), int64(-42)},
		{"Int64 from JSON at its minimum", edgexCommon.ValueTypeInt64, float64(-(1 << 63)), int64(math.MinInt64)},
		{"Uint64 from JSON truncated to zero", edgexCommon.ValueTypeUint64, -0.5, uint64(0)},
		{"Int64 exact", edgexCommon.ValueTypeInt64, int64(math.MinInt64), int64(math.MinInt64)},
		{"Int64 beyond 2^53", edgexCommon.ValueTypeInt64, int64(1<<53 + 1), int64(1<<53 + 1)},
		{"Uint64 from JSON", edgexCommon.ValueTypeUint64, float64(42), uint64(42)},
		{"Uint64 exact", edgexCommon.ValueTypeUint64, uint64(math.MaxUint64), uint64(math.MaxUint64)},
		{"Int64Array mixed", edgexCommon.ValueTypeInt64Array, []any{float64(1), int64(math.MaxInt64)},
			[]int64{1, math.MaxInt64}},
		{"Uint64Array mixed", edgexCommon.ValueTypeUint64Array, []any{float64(1), uint64(18244850924796673191)},
			[]uint64{1, 18244850924796673191}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseXRTReadingValue(tc.valueType, tc.reading)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestParseXRTReadingValue_Rejected(t *testing.T) {
	for _, tc := range []struct {
		valueType string
		reading   any
	}{
		// Each 64-bit type takes only its own Go integer type, besides float64.
		{edgexCommon.ValueTypeInt64, uint64(1)},
		{edgexCommon.ValueTypeUint64, int64(1)},
		{edgexCommon.ValueTypeInt64, "1"},
		{edgexCommon.ValueTypeInt64Array, []any{uint64(1)}},
		{edgexCommon.ValueTypeUint64Array, []any{int64(1)}},
		{edgexCommon.ValueTypeInt32, int64(1)},
		// Beyond the range, as a JSON number Reading.UnmarshalJSON cannot decode exactly.
		{edgexCommon.ValueTypeInt64, float64(1 << 63)},
		{edgexCommon.ValueTypeInt64Array, []any{float64(1 << 63)}},
		{edgexCommon.ValueTypeUint64, float64(-1)},
		{edgexCommon.ValueTypeUint64, float64(1 << 64)},
		{edgexCommon.ValueTypeUint64Array, []any{float64(-1)}},
		// A json.Number beyond its type's range, or for a type that is no number.
		{edgexCommon.ValueTypeUint8, json.Number("256")},
		{edgexCommon.ValueTypeUint8Array, []any{json.Number("256")}},
		{edgexCommon.ValueTypeInt8, json.Number("-129")},
		{edgexCommon.ValueTypeInt32, json.Number("2147483648")},
		{edgexCommon.ValueTypeInt64, json.Number("9223372036854775808")},
		{edgexCommon.ValueTypeInt64Array, []any{json.Number("9223372036854775808")}},
		{edgexCommon.ValueTypeInt64, json.Number("-9223372036854775809")}, // rounds to MinInt64 as a float64
		{edgexCommon.ValueTypeInt64Array, []any{json.Number("-9223372036854775809")}},
		{edgexCommon.ValueTypeUint64, json.Number("-1")},
		{edgexCommon.ValueTypeUint64Array, []any{json.Number("-1")}},
		{edgexCommon.ValueTypeUint64, json.Number("18446744073709551616")},
		{edgexCommon.ValueTypeInt64, json.Number("1e400")},
		{edgexCommon.ValueTypeFloat64, json.Number("1e400")},
		{edgexCommon.ValueTypeFloat64Array, []any{json.Number("1e400")}},
		{edgexCommon.ValueTypeFloat32, float64(1e39)},
		{edgexCommon.ValueTypeFloat32Array, []any{json.Number("-1e39")}},
		// A number for a Uint8Array, whose digits would decode as base64.
		{edgexCommon.ValueTypeUint8Array, json.Number("12345678")},
		{edgexCommon.ValueTypeUint8Array, float64(1234)},
		{edgexCommon.ValueTypeInt64, json.Number("")}, // not from a decoder, but must not panic
		{edgexCommon.ValueTypeString, json.Number("1")},
		{edgexCommon.ValueTypeBool, json.Number("1")},
	} {
		_, err := ParseXRTReadingValue(tc.valueType, tc.reading)
		assert.Error(t, err, "%s %#v", tc.valueType, tc.reading)
	}
}

func TestToEdgeXV2EventDTO_ExactUint64(t *testing.T) {
	event, err := ToEdgeXV2EventDTO(MultiResourcesResult{Device: "d", Profile: "p", Readings: map[string]Reading{
		"count": {Value: uint64(math.MaxUint64), Type: edgexCommon.ValueTypeUint64},
	}})

	require.NoError(t, err)
	require.Len(t, event.Readings, 1)
	assert.Equal(t, "18446744073709551615", event.Readings[0].Value)
}

func TestParseXRTReadingValue_SmallIntegers(t *testing.T) {
	for _, tc := range []struct {
		valueType string
		min, max  float64
		want      [2]any
	}{
		{edgexCommon.ValueTypeUint8, 0, math.MaxUint8, [2]any{uint8(0), uint8(math.MaxUint8)}},
		{edgexCommon.ValueTypeUint16, 0, math.MaxUint16, [2]any{uint16(0), uint16(math.MaxUint16)}},
		{edgexCommon.ValueTypeUint32, 0, math.MaxUint32, [2]any{uint32(0), uint32(math.MaxUint32)}},
		{edgexCommon.ValueTypeInt8, math.MinInt8, math.MaxInt8, [2]any{int8(math.MinInt8), int8(math.MaxInt8)}},
		{edgexCommon.ValueTypeInt16, math.MinInt16, math.MaxInt16, [2]any{int16(math.MinInt16), int16(math.MaxInt16)}},
		{edgexCommon.ValueTypeInt32, math.MinInt32, math.MaxInt32, [2]any{int32(math.MinInt32), int32(math.MaxInt32)}},
	} {
		t.Run(tc.valueType, func(t *testing.T) {
			assertSmallInteger(t, tc.valueType, tc.min, tc.max, tc.want)
		})
	}
}

// assertSmallInteger checks valueType, and its array type, at both ends of its range and just beyond them.
func assertSmallInteger(t *testing.T, valueType string, min, max float64, want [2]any) {
	arrayType := valueType + "Array"
	for i, f := range []float64{min, max} {
		got, err := ParseXRTReadingValue(valueType, f)
		require.NoError(t, err, f)
		assert.Equal(t, want[i], got, f)

		got, err = ParseXRTReadingValue(arrayType, []any{f})
		require.NoError(t, err, f)
		assert.Equal(t, want[i], reflect.ValueOf(got).Index(0).Interface(), f)
	}
	for _, f := range []float64{min - 1, max + 1} {
		_, err := ParseXRTReadingValue(valueType, f)
		assert.Error(t, err, f)
		_, err = ParseXRTReadingValue(arrayType, []any{f})
		assert.Error(t, err, f)
	}
}

// A fraction still truncates toward zero, as before the range check.
func TestParseXRTReadingValue_SmallIntegerFraction(t *testing.T) {
	got, err := ParseXRTReadingValue(edgexCommon.ValueTypeUint8, 255.9)
	require.NoError(t, err)
	assert.Equal(t, uint8(255), got)

	got, err = ParseXRTReadingValue(edgexCommon.ValueTypeInt8, -128.9)
	require.NoError(t, err)
	assert.Equal(t, int8(-128), got)

	got, err = ParseXRTReadingValue(edgexCommon.ValueTypeUint8, -0.5)
	require.NoError(t, err)
	assert.Equal(t, uint8(0), got)
}

// Every numeric type takes a json.Number, as a device:put request decodes one (see PutResourceRequest.UnmarshalJSON).
func TestParseXRTReadingValue_JSONNumber(t *testing.T) {
	for _, tc := range []struct {
		valueType string
		number    json.Number
		want      any
	}{
		{edgexCommon.ValueTypeUint8, "255", uint8(255)},
		{edgexCommon.ValueTypeUint16, "65535", uint16(65535)},
		{edgexCommon.ValueTypeUint32, "4294967295", uint32(4294967295)},
		{edgexCommon.ValueTypeUint64, "18446744073709551615", uint64(math.MaxUint64)},
		{edgexCommon.ValueTypeInt8, "-128", int8(-128)},
		{edgexCommon.ValueTypeInt16, "-32768", int16(-32768)},
		{edgexCommon.ValueTypeInt32, "-2147483648", int32(-2147483648)},
		{edgexCommon.ValueTypeInt64, "-9223372036854775808", int64(math.MinInt64)},
		{edgexCommon.ValueTypeInt64, "9223372036854775807", int64(math.MaxInt64)},
		{edgexCommon.ValueTypeFloat32, "1.5", float32(1.5)},
		{edgexCommon.ValueTypeFloat32, "3.4028234663852886e38", float32(math.MaxFloat32)},
		{edgexCommon.ValueTypeFloat64, "1.5", 1.5},
		// Exponent and fraction forms, as before as float64.
		{edgexCommon.ValueTypeInt8, "1e2", int8(100)},
		{edgexCommon.ValueTypeInt64, "9007199254740993e0", int64(1<<53 + 1)},
		{edgexCommon.ValueTypeUint64, "1.8446744073709551615e19", uint64(math.MaxUint64)},
		{edgexCommon.ValueTypeUint8, "255.9", uint8(255)},
		{edgexCommon.ValueTypeInt64, "-1.5", int64(-1)},
		{edgexCommon.ValueTypeUint64, "-0.5", uint64(0)},
	} {
		t.Run(tc.valueType+" "+tc.number.String(), func(t *testing.T) {
			got, err := ParseXRTReadingValue(tc.valueType, tc.number)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)

			got, err = ParseXRTReadingValue(tc.valueType+"Array", []any{tc.number})
			require.NoError(t, err)
			assert.Equal(t, tc.want, reflect.ValueOf(got).Index(0).Interface())
		})
	}
}

// Every value type, valid and not; the 64-bit and smaller integer types have their own tests above.
func TestParseXRTReadingValue_Types(t *testing.T) {
	object := map[string]any{"a": float64(1)}
	for _, tc := range []struct {
		name      string
		valueType string
		reading   any
		want      any
		wantErr   string
	}{
		{"String", edgexCommon.ValueTypeString, "abc", "abc", ""},
		{"String invalid", edgexCommon.ValueTypeString, 1.5, nil, "invalid string value '1.5'"},
		{"Bool", edgexCommon.ValueTypeBool, true, true, ""},
		{"Bool invalid", edgexCommon.ValueTypeBool, "true", nil, "invalid bool value 'true'"},
		{"Uint8 invalid", edgexCommon.ValueTypeUint8, "1", nil, "invalid number '1'"},
		{"Float32", edgexCommon.ValueTypeFloat32, 1.5, float32(1.5), ""},
		{"Float32 invalid", edgexCommon.ValueTypeFloat32, "1.5", nil, "invalid number '1.5'"},
		{"Float64", edgexCommon.ValueTypeFloat64, 1.5, 1.5, ""},
		{"Float64 invalid", edgexCommon.ValueTypeFloat64, true, nil, "invalid number 'true'"},
		{"Binary", edgexCommon.ValueTypeBinary, "AQID", []byte{1, 2, 3}, ""},
		{"Binary not a string", edgexCommon.ValueTypeBinary, 1.5, nil, "invalid string value '1.5'"},
		{"Binary not base64", edgexCommon.ValueTypeBinary, "!!!", nil, "fail to decode the base64 string '!!!'"},
		{"BoolArray", edgexCommon.ValueTypeBoolArray, []any{true, false}, []bool{true, false}, ""},
		{"BoolArray invalid element", edgexCommon.ValueTypeBoolArray, []any{true, 1.5}, nil, "invalid bool value '1.5'"},
		{"StringArray", edgexCommon.ValueTypeStringArray, []any{"a"}, []string{"a"}, ""},
		{"StringArray invalid element", edgexCommon.ValueTypeStringArray, []any{true}, nil, "invalid string value 'true'"},
		{"StringArray not an array", edgexCommon.ValueTypeStringArray, "a", nil, "invalid array 'a'"},
		{"Uint8Array empty", edgexCommon.ValueTypeUint8Array, []any{}, []uint8{}, ""},
		{"Uint8Array in base64", edgexCommon.ValueTypeUint8Array, "AQID", []byte{1, 2, 3}, ""},
		{"Uint8Array neither", edgexCommon.ValueTypeUint8Array, "!!!", nil, "invalid array '!!!'"},
		{"Uint8Array invalid element", edgexCommon.ValueTypeUint8Array, []any{"1"}, nil, "invalid number '1'"},
		{"Float32Array", edgexCommon.ValueTypeFloat32Array, []any{1.5}, []float32{1.5}, ""},
		{"Float32Array invalid element", edgexCommon.ValueTypeFloat32Array, []any{"1.5"}, nil, "invalid number '1.5'"},
		{"Float64Array", edgexCommon.ValueTypeFloat64Array, []any{1.5}, []float64{1.5}, ""},
		{"Float64Array invalid element", edgexCommon.ValueTypeFloat64Array, []any{nil}, nil, "invalid number '<nil>'"},
		{"Object", edgexCommon.ValueTypeObject, object, object, ""},
		// As a device:put decodes it (see PutResourceRequest.UnmarshalJSON).
		{"Object number", edgexCommon.ValueTypeObject, json.Number("3"), float64(3), ""},
		{"Object array", edgexCommon.ValueTypeObject, []any{json.Number("1"), "x"}, []any{float64(1), "x"}, ""},
		{"ObjectArray", edgexCommon.ValueTypeObjectArray, []any{object}, []map[string]any{object}, ""},
		{"ObjectArray with null", edgexCommon.ValueTypeObjectArray, []any{object, nil}, []map[string]any{object, nil}, ""},
		{"ObjectArray invalid element", edgexCommon.ValueTypeObjectArray, []any{object, 1.5}, nil, "invalid object '1.5'"},
		{"ObjectArray not an array", edgexCommon.ValueTypeObjectArray, object, nil, "invalid array 'map[a:1]'"},
		{"unsupported type", "Bogus", 1.5, nil, "none supported value type 'Bogus'"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseXRTReadingValue(tc.valueType, tc.reading)
			assertParsed(t, tc.want, tc.wantErr, got, err)
		})
	}
}

func assertParsed(t *testing.T, want any, wantErr string, got any, err errors.EdgeX) {
	if wantErr != "" {
		require.Error(t, err)
		assert.Equal(t, errors.KindContractInvalid, errors.Kind(err))
		assert.Contains(t, err.Error(), wantErr)
		assert.Nil(t, got)
		return
	}
	require.NoError(t, err)
	assert.Equal(t, want, got)
}
