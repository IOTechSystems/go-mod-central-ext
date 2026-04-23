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
