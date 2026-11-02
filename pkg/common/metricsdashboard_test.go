// Copyright (C) 2026 IOTech Ltd

package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetricsDashboardLastDuration(t *testing.T) {
	expected := map[string]time.Duration{
		"5m": 5 * time.Minute, "15m": 15 * time.Minute, "30m": 30 * time.Minute,
		"1h": time.Hour, "2h": 2 * time.Hour, "3h": 3 * time.Hour, "6h": 6 * time.Hour,
		"12h": 12 * time.Hour, "24h": 24 * time.Hour, "2d": 48 * time.Hour,
	}
	assert.Len(t, MetricsDashboardLastChoices, len(expected))
	for last, d := range expected {
		got, ok := MetricsDashboardLastDuration(last)
		assert.True(t, ok, last)
		assert.Equal(t, d, got, last)
	}
	for _, last := range []string{"7d", "", "1d"} {
		_, ok := MetricsDashboardLastDuration(last)
		assert.False(t, ok, last)
	}
	assert.Equal(t, MetricsDashboardLastChoices[0].Last, DefaultMetricsDashboardRangeLast)
}
