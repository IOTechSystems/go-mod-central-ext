// Copyright (C) 2026 IOTech Ltd

package xrtmodels

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchedule_TagsRoundTrip(t *testing.T) {
	s := Schedule{
		Name:     "test-schedule",
		Device:   "test-device",
		Resource: []string{"temperature"},
		Interval: 5000000,
		OnChange: true,
		Tags:     map[string]any{"unit": "celsius", "priority": float64(1)},
	}

	data, err := json.Marshal(s)
	require.NoError(t, err)

	var decoded Schedule
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, s.Tags, decoded.Tags)
}

func TestSchedule_TagsOmittedWhenNil(t *testing.T) {
	s := Schedule{Name: "s", Device: "d", Resource: []string{"r"}}

	data, err := json.Marshal(s)
	require.NoError(t, err)

	assert.NotContains(t, string(data), `"tags"`)
}
