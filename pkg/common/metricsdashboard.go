// Copyright (C) 2026 IOTech Ltd

package common

import "time"

// MetricsDashboardLastChoice is one relative time range a metrics dashboard may show, such as the last 5 minutes
type MetricsDashboardLastChoice struct {
	Last     string
	Duration time.Duration
}

// MetricsDashboardLastChoices lists the relative ranges of a metrics dashboard, in the order the Data Dashboard offers them
var MetricsDashboardLastChoices = []MetricsDashboardLastChoice{
	{"5m", 5 * time.Minute},
	{"15m", 15 * time.Minute},
	{"30m", 30 * time.Minute},
	{"1h", time.Hour},
	{"2h", 2 * time.Hour},
	{"3h", 3 * time.Hour},
	{"6h", 6 * time.Hour},
	{"12h", 12 * time.Hour},
	{"24h", 24 * time.Hour},
	{"2d", 48 * time.Hour},
}

// DefaultMetricsDashboardRangeLast is the range of a metrics dashboard created without one
const DefaultMetricsDashboardRangeLast = "5m"

// MetricsDashboardLastDuration returns the duration of a choice in MetricsDashboardLastChoices, and false for any other value
func MetricsDashboardLastDuration(last string) (time.Duration, bool) {
	for _, c := range MetricsDashboardLastChoices {
		if c.Last == last {
			return c.Duration, true
		}
	}
	return 0, false
}
