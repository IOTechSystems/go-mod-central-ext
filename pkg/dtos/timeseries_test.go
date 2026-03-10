// Copyright (C) 2026 IOTech Ltd

package dtos

import (
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/models"

	"github.com/stretchr/testify/assert"
)

func TestFromReadingModelsToTimeSeriesResourceMap(t *testing.T) {
	TestSimple := "Simple"
	TestUnit := "C"
	TestValueString1 := "test_1"
	TestValueString2 := "test_2"
	TestNumeric := "Numeric"
	TestValueNumeric := int16(25)
	TestImage := "Image"
	TestValueImage := []byte{0xFF, 0xD8}
	TestMediaType := "image/jpeg"
	TestOrigin := int64(1740992400000)
	NewTestOrigin := int64(1740992400001)

	readings := []models.Reading{
		models.SimpleReading{
			BaseReading: models.BaseReading{ResourceName: TestSimple, ValueType: common.ValueTypeString, Units: TestUnit, Origin: TestOrigin},
			Value:       TestValueString1,
		},
		models.SimpleReading{
			BaseReading: models.BaseReading{ResourceName: TestSimple, ValueType: common.ValueTypeString, Units: TestUnit, Origin: NewTestOrigin},
			Value:       TestValueString2,
		},
		models.NumericReading{
			BaseReading:  models.BaseReading{ResourceName: TestNumeric, ValueType: common.ValueTypeInt16, Units: TestUnit, Origin: TestOrigin},
			NumericValue: TestValueNumeric,
		},
		models.BinaryReading{
			BaseReading: models.BaseReading{ResourceName: TestImage, ValueType: common.ValueTypeBinary, Units: "", Origin: TestOrigin},
			BinaryValue: TestValueImage,
			MediaType:   TestMediaType,
		},
	}

	result := FromReadingModelsToTimeSeriesResourceMap(readings)

	// Assert TestSimple group
	assert.Contains(t, result, TestSimple)
	assert.Equal(t, common.ValueTypeString, result[TestSimple].ValueType)
	assert.Equal(t, TestUnit, result[TestSimple].Units)
	assert.Empty(t, result[TestSimple].MediaType)
	assert.Equal(t, 2, len(result[TestSimple].TimeSeries))
	assert.Equal(t, []any{TestOrigin, TestValueString1}, result[TestSimple].TimeSeries[0])
	assert.Equal(t, []any{NewTestOrigin, TestValueString2}, result[TestSimple].TimeSeries[1])

	// Assert TestNumeric group
	assert.Contains(t, result, TestNumeric)
	assert.Equal(t, common.ValueTypeInt16, result[TestNumeric].ValueType)
	assert.Equal(t, TestUnit, result[TestSimple].Units)
	assert.Empty(t, result[TestSimple].MediaType)
	assert.Equal(t, 1, len(result[TestNumeric].TimeSeries))
	assert.Equal(t, []any{TestOrigin, TestValueNumeric}, result[TestNumeric].TimeSeries[0])

	// Assert TestImage group
	assert.Contains(t, result, TestImage)
	assert.Equal(t, common.ValueTypeBinary, result[TestImage].ValueType)
	assert.Equal(t, TestMediaType, result[TestImage].MediaType)
	assert.Contains(t, TestUnit, "")
	assert.Equal(t, 1, len(result[TestImage].TimeSeries))
	assert.Equal(t, []any{TestOrigin, []byte{0xFF, 0xD8}}, result[TestImage].TimeSeries[0])

}

func TestFromReadingModelsToTimeSeriesResourceMap_Empty(t *testing.T) {
	result := FromReadingModelsToTimeSeriesResourceMap([]models.Reading{})
	assert.Empty(t, result)
}
