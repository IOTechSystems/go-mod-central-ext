// Copyright (C) 2023-2026 IOTech Ltd

package xrtmodels

import (
	"encoding/json"
	"math"
	"testing"

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXrtErrorCode(t *testing.T) {
	assert.Equal(t, XrtErrorCode(errors.NewCommonEdgeX(errors.KindEntityDoesNotExist, "", nil)), XrtSdkStatusNotFound)
	assert.Equal(t, XrtErrorCode(errors.NewCommonEdgeX(errors.KindNotImplemented, "", nil)), XrtSdkStatusNotSupported)
	assert.Equal(t, XrtErrorCode(errors.NewCommonEdgeX(errors.KindInvalidId, "", nil)), XrtSdkStatusInvalidOperation)
	assert.Equal(t, XrtErrorCode(errors.NewCommonEdgeX(errors.KindDuplicateName, "", nil)), XrtSdkStatusAlreadyExists)
	assert.Equal(t, XrtErrorCode(errors.NewCommonEdgeX(errors.KindServerError, "", nil)), XrtSdkStatusServerError)
}

func TestScheduleReadResponse_Unmarshal(t *testing.T) {
	raw := `{
		"client": "c1",
		"request_id": "abc",
		"type": "xrt.response:1.0",
		"result": {
			"status": 0,
			"schedule": {"name":"s1","device":"d1","resource":["r1"],"on_change":false,"bounds":{},"publish":false,"units":false}
		}
	}`

	var resp ScheduleReadResponse
	err := json.Unmarshal([]byte(raw), &resp)
	require.NoError(t, err)
	assert.Equal(t, "s1", resp.Result.Schedule.Name)
	assert.Equal(t, "d1", resp.Result.Schedule.Device)
}

// Batch add and delete report per-item status while the envelope stays successful, so a
// failed item must always carry an error — treating one as success would report a failure
// as done.
func TestNewBatchItemResult(t *testing.T) {
	t.Run("every non-zero status carries an error", func(t *testing.T) {
		// Includes codes BaseResult has no specific mapping for.
		for _, status := range []int{1, 3, 6, 7, 10, 13, 500, 9999} {
			result := NewBatchItemResult("dev-1", BaseResult{Status: status, ErrorMessage: "something went wrong"})
			assert.True(t, result.Failed(), "status %d", status)
			assert.Error(t, result.Err, "status %d", status)
		}
	})

	t.Run("a zero status carries no error", func(t *testing.T) {
		result := NewBatchItemResult("dev-1", BaseResult{})
		assert.False(t, result.Failed())
		assert.NoError(t, result.Err)
		assert.Equal(t, "dev-1", result.Name)
	})
}

func TestReading_UnmarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		name string
		json string
		want any
	}{
		{"Int64 at its minimum", `{"value": -9223372036854775808, "type": "Int64"}`, int64(math.MinInt64)},
		{"Int64 at its maximum", `{"value": 9223372036854775807, "type": "Int64"}`, int64(math.MaxInt64)},
		{"Int64 just beyond 2^53", `{"value": 9007199254740993, "type": "Int64"}`, int64(1<<53 + 1)},
		{"Uint64 at its maximum", `{"value": 18446744073709551615, "type": "Uint64"}`, uint64(math.MaxUint64)},
		{"type before value, in another case", `{"type": "uint64", "value": 7}`, uint64(7)},
		{"Int64Array", `{"value": [-1, 9223372036854775807], "type": "Int64Array"}`, []any{int64(-1), int64(math.MaxInt64)}},
		{"Uint64Array", `{"value": [18244850924796673191], "type": "Uint64Array"}`, []any{uint64(18244850924796673191)}},
		// Exponent and fraction forms (see integerDigits).
		{"Int64 in exponent form", `{"value": 9007199254740993e0, "type": "Int64"}`, int64(1<<53 + 1)},
		{"Int64 in fraction form", `{"value": 9007199254740993.000, "type": "Int64"}`, int64(1<<53 + 1)},
		{"Int64 with a negative exponent", `{"value": -90071992547409930E-1, "type": "Int64"}`, -int64(1<<53 + 1)},
		{"Int64 with leading zeros in the mantissa", `{"value": 0.0012e+4, "type": "Int64"}`, int64(12)},
		{"Int64 at its minimum with a leading zero", `{"value": -0.09223372036854775808e20, "type": "Int64"}`, int64(math.MinInt64)},
		{"Int64 with zeros from the exponent", `{"value": 9007199254740993e2, "type": "Int64"}`, int64(900719925474099300)},
		{"Int64 zero with a huge exponent", `{"value": 0e999999999, "type": "Int64"}`, int64(0)},
		{"Uint64 at its maximum in exponent form", `{"value": 1.8446744073709551615e19, "type": "Uint64"}`, uint64(math.MaxUint64)},
		{"Uint64 at its maximum with a leading zero", `{"value": 0.18446744073709551615e20, "type": "Uint64"}`, uint64(math.MaxUint64)},
		{"Uint64Array in exponent form", `{"value": [18244850924796673191e0], "type": "Uint64Array"}`, []any{uint64(18244850924796673191)}},
		{"Int64 not an integer in exponent form", `{"value": 15e-1, "type": "Int64"}`, 1.5},
		{"Uint64 beyond 20 digits", `{"value": 1e20, "type": "Uint64"}`, 1e20},
		{"Int64 just beyond its maximum", `{"value": 9223372036854775808, "type": "Int64"}`, float64(1 << 63)},
		{"Int64 zero with an exponent beyond int", `{"value": 0e99999999999999999999, "type": "Int64"}`, int64(0)},
		{"Int64 not an integer", `{"value": 1.5, "type": "Int64"}`, 1.5},
		{"Uint64 negative", `{"value": -1, "type": "Uint64"}`, float64(-1)},
		{"Int64 null", `{"value": null, "type": "Int64"}`, nil},
		{"Int32 as before", `{"value": 7, "type": "Int32"}`, float64(7)},
		{"Float64 as before", `{"value": 1.5, "type": "Float64"}`, 1.5},
		{"Object as before", `{"value": {"n": 1}, "type": "Object"}`, map[string]any{"n": float64(1)}},
		{"no value", `{"type": "Int64"}`, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var r Reading
			require.NoError(t, json.Unmarshal([]byte(tc.json), &r))
			assert.Equal(t, tc.want, r.Value)
		})
	}
}

// Decoding into a filled Reading, as into any struct, keeps the fields the JSON leaves out.
func TestReading_UnmarshalJSON_Filled(t *testing.T) {
	for _, tc := range []struct {
		json string
		want Reading
	}{
		{`null`, Reading{Value: int64(5), Type: "Int64", Origin: 7}},
		{`{"origin": 8}`, Reading{Value: int64(5), Type: "Int64", Origin: 8}},
	} {
		r := Reading{Value: int64(5), Type: "Int64", Origin: 7}
		require.NoError(t, json.Unmarshal([]byte(tc.json), &r), tc.json)
		assert.Equal(t, tc.want, r, tc.json)
	}
}

func TestIntegerDigits_Bound(t *testing.T) {
	digits, ok := integerDigits("1e19")
	assert.True(t, ok)
	assert.Equal(t, "10000000000000000000", digits)

	for _, n := range []string{"1e20", "1e999999999"} {
		_, ok := integerDigits(n)
		assert.False(t, ok, n)
	}
}

// The other fields decode as before.
func TestReading_UnmarshalJSON_OtherFields(t *testing.T) {
	var r Reading
	require.NoError(t, json.Unmarshal([]byte(`{"value": 1, "type": "Int64", "origin": 1700000000000000000, "tags": {"k": "v"}}`), &r))

	assert.Equal(t, Reading{Value: int64(1), Type: "Int64", Origin: 1700000000000000000, Tags: map[string]any{"k": "v"}}, r)
}

func TestReading_UnmarshalJSON_Invalid(t *testing.T) {
	for _, data := range []string{
		`{"value": [1,, "type": "Int64"}`,
		// Beyond float64's range, as it failed before the 64-bit types decoded exactly.
		`{"value": 1e400, "type": "Int64"}`,
		`{"value": [1, 1e400], "type": "Uint64Array"}`,
	} {
		var r Reading
		assert.Error(t, json.Unmarshal([]byte(data), &r), data)
	}
}

// An XRT reply decoded through MultiResourcesResponse reaches the EdgeX event exactly, and an ObjectArray encodes
// as it did when it was []any.
func TestToEdgeXV2EventDTO_FromJSON(t *testing.T) {
	var res MultiResourcesResponse
	require.NoError(t, json.Unmarshal([]byte(`{"result": {"device": "d", "profile": "p", "readings": {
		"count": {"value": 18446744073709551615, "type": "Uint64"},
		"offset": {"value": 9223372036854775807, "type": "Int64"},
		"objs": {"value": [{"a": 1, "b": "x"}, {}, null], "type": "ObjectArray"}}}}`), &res))

	event, err := ToEdgeXV2EventDTO(res.Result)

	require.NoError(t, err)
	values := map[string]string{}
	for _, reading := range event.Readings {
		values[reading.ResourceName] = reading.Value
		if reading.ValueType == edgexCommon.ValueTypeObjectArray {
			encoded, jsonErr := json.Marshal(reading.ObjectValue)
			require.NoError(t, jsonErr)
			values[reading.ResourceName] = string(encoded)
		}
	}
	assert.Equal(t, map[string]string{"count": "18446744073709551615", "offset": "9223372036854775807",
		"objs": `[{"a":1,"b":"x"},{},null]`}, values)
}
