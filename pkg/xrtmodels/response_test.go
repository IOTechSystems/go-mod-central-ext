// Copyright (C) 2023-2026 IOTech Ltd

package xrtmodels

import (
	"encoding/json"
	"testing"

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
