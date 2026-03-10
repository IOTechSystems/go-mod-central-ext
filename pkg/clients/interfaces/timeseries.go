// Copyright (C) 2026 IOTech Ltd

package interfaces

import (
	"context"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/responses"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

// TimeSeriesClient defines the interface for interactions with the timeseries API endpoints on the IOTech core-data service.
type TimeSeriesClient interface {
	// TimeSeriesByDeviceNameAndResourceNameAndTimeRange returns time series data by device name, resource name, and specified time range. Time series are sorted in descending order of origin time.
	// start, end: Unix timestamp, indicating the date/time range
	TimeSeriesByDeviceNameAndResourceNameAndTimeRange(ctx context.Context, deviceName, resourceName string, start, end int64) (responses.TimeSeriesResponse, errors.EdgeX)
	// TimeSeriesByDeviceNameAndMultiResourceNamesAndTimeRange returns time series by device name, multiple resource names and specified time range. Time series are sorted in descending order of origin time.
	// If none of resourceNames is specified, return all time series under the specified deviceName and within the specified time range
	// start, end: Unix timestamp, indicating the date/time range
	TimeSeriesByDeviceNameAndMultiResourceNamesAndTimeRange(ctx context.Context, deviceName string, resourceNames []string, start, end int64) (responses.TimeSeriesResponse, errors.EdgeX)
}
