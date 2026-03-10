// Copyright (C) 2026 IOTech Ltd

package dtos

import "github.com/edgexfoundry/go-mod-core-contracts/v4/models"

type TimeSeriesResource struct {
	ValueType  string  `json:"valueType"`
	Units      string  `json:"units"`
	TimeSeries [][]any `json:"timeSeries"`
	MediaType  string  `json:"mediaType,omitempty"`
}
type TimeSeriesResourceMap map[string]*TimeSeriesResource

func FromReadingModelsToTimeSeriesResourceMap(readings []models.Reading) TimeSeriesResourceMap {
	tsResources := make(TimeSeriesResourceMap)

	for _, reading := range readings {
		var name, valueType, units string
		var origin int64
		var value any

		mediaType := ""
		switch r := reading.(type) {
		case models.BinaryReading:
			name, valueType, units, origin, value, mediaType = r.ResourceName, r.ValueType, r.Units, r.Origin, r.BinaryValue, r.MediaType
		case models.ObjectReading:
			name, valueType, units, origin, value = r.ResourceName, r.ValueType, r.Units, r.Origin, r.ObjectValue
		case models.SimpleReading:
			name, valueType, units, origin, value = r.ResourceName, r.ValueType, r.Units, r.Origin, r.Value
		case models.NumericReading:
			name, valueType, units, origin, value = r.ResourceName, r.ValueType, r.Units, r.Origin, r.NumericValue
		case models.NullReading:
			name, valueType, units, origin, value = r.ResourceName, r.ValueType, r.Units, r.Origin, r.Value
		default:
			continue // Skip unknown types
		}

		resource, exists := tsResources[name]
		if !exists {
			resource = &TimeSeriesResource{
				ValueType:  valueType,
				Units:      units,
				MediaType:  mediaType,
				TimeSeries: [][]any{},
			}
			tsResources[name] = resource
		}

		resource.TimeSeries = append(resource.TimeSeries, []any{origin, value})
	}

	return tsResources
}
