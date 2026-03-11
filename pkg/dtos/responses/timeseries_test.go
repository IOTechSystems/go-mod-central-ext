// Copyright (C) 2026 IOTech Ltd

package responses

import (
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"

	"github.com/stretchr/testify/assert"
)

func TestNewTimeSeriesResponse(t *testing.T) {
	TestResourceNames := []string{"resource01", "resource02"}
	TestValueType := "int8"
	TestUnits := "kw"
	TestTimeSeries := [][]any{{int64(0), 0}, {int64(1), 1}}
	TestTimeSeriesResource := dtos.TimeSeriesResource{
		ValueType:  TestValueType,
		Units:      TestUnits,
		TimeSeries: TestTimeSeries,
	}

	TestTSResourceMap_1 := dtos.TimeSeriesResourceMap{
		TestResourceNames[0]: &TestTimeSeriesResource,
	}

	TestTSResourceMap_2 := dtos.TimeSeriesResourceMap{
		TestResourceNames[0]: &TestTimeSeriesResource,
		TestResourceNames[1]: &TestTimeSeriesResource,
	}

	tests := []struct {
		name           string
		inputResources dtos.TimeSeriesResourceMap
		expected       map[string]interface{}
	}{
		{
			name:           "Empty resource map",
			inputResources: dtos.TimeSeriesResourceMap{},
			expected: TimeSeriesResponse{
				"apiVersion": common.NewVersionable().ApiVersion,
			},
		},
		{
			name:           "Single resource in map",
			inputResources: TestTSResourceMap_1,
			expected: TimeSeriesResponse{
				"apiVersion":         common.NewVersionable().ApiVersion,
				TestResourceNames[0]: &TestTimeSeriesResource,
			},
		},
		{
			name:           "Multiple resources in map",
			inputResources: TestTSResourceMap_2,
			expected: TimeSeriesResponse{
				"apiVersion":         common.NewVersionable().ApiVersion,
				TestResourceNames[0]: &TestTimeSeriesResource,
				TestResourceNames[1]: &TestTimeSeriesResource,
			},
		},
		{
			name:           "Nil input map",
			inputResources: nil,
			expected: TimeSeriesResponse{
				"apiVersion": common.NewVersionable().ApiVersion,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := NewTimeSeriesResponse(tt.inputResources)
			assert.IsType(t, TimeSeriesResponse{}, response)

			for key, value := range tt.expected {
				val, found := response[key]
				assert.True(t, found, "Key %s should be present in the response", key)
				assert.Equal(t, value, val, "Unexpected value for key %s", key)
			}
		})
	}
}
