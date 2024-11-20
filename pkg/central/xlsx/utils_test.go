// Copyright (C) 2023-2024 IOTech Ltd

package xlsx

import (
	"fmt"
	"reflect"
	"testing"

	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"

	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

func Test_checkMappingObject_WithoutSheet(t *testing.T) {
	f := excelize.NewFile()
	defer f.Close()

	mockColCount := 2

	err := createMappingTableSheet(f)
	require.Error(t, err)

	err = checkMappingObject(f, devicesSheetName, &mockColCount, 1, "UNLOCKED", "AdminState", &[]string{"Name"})
	require.Error(t, err, "Expected no sheet defined error not occurred")
}

func Test_checkMappingObject_WithSheet(t *testing.T) {
	f := excelize.NewFile()
	defer f.Close()

	_, err := f.NewSheet(devicesSheetName)
	require.NoError(t, err)

	_, err = f.NewSheet(mappingTableSheetName)
	require.NoError(t, err)

	err = createMappingTableSheet(f)
	require.NoError(t, err)

	mockColCount := 2
	invalidMockColCount := -1
	mockHeader1 := []string{"Name"}
	mockHeader2 := mockHeader1

	tests := []struct {
		name         string
		colCount     *int
		rowCount     int
		defaultValue string
		objectField  string
		headerRow    *[]string
		expectError  bool
	}{
		{"Invalid - checkMappingObject with no header", &mockColCount, 2, "", "", nil, true},
		{"Invalid - checkMappingObject with invalid column count", &invalidMockColCount, 2, "UNLOCKED", "AdminState", &mockHeader1, true},
		{"Valid - checkMappingObject with header", &mockColCount, 2, "", "", &mockHeader1, false},
		{"Valid - checkMappingObject with header & objectField/defaultValue", &mockColCount, 2, "UNLOCKED", "AdminState", &mockHeader2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkMappingObject(f, devicesSheetName, tt.colCount, tt.rowCount, tt.defaultValue, tt.objectField, tt.headerRow)
			if tt.expectError {
				require.Error(t, err, "Expected checkMappingObject error not occurred")
			} else {
				require.NoError(t, err, "Unexpected checkMappingObject error occurred")
				if tt.defaultValue != "" && tt.objectField != "" {
					newHeaderCol := *tt.headerRow
					require.Equal(t, newHeaderCol[1], tt.objectField)
				}
			}
		})
	}
}

func Test_setMapToStructField(t *testing.T) {
	mockFieldName := "testField"
	mockFieldValue := int64(123)
	mockMap := map[string]any{mockFieldName: mockFieldValue}

	tests := []struct {
		name        string
		fieldName   string
		mapValue    map[string]any
		expectError bool
	}{
		{"Invalid - wrong field name", "INVALID", make(map[string]any), true},
		{"Invalid - non-map field", "Description", mockMap, true},
		{"Valid - set map field", properties, mockMap, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			structPtr := edgexDtos.Device{}
			v := reflect.ValueOf(&structPtr)
			elementType := v.Elem().Type()
			element := reflect.New(elementType).Elem()

			err := setMapToStructField(&element, tt.fieldName, tt.mapValue)
			if tt.expectError {
				require.Error(t, err, "Expected setMapToStructField error not occurred")
			} else {
				require.NoError(t, err, "Unexpected setMapToStructField error occurred")
				v.Elem().Set(element)
				require.Equal(t, int64(mockFieldValue), structPtr.Properties[mockFieldName])
			}
		})
	}
}

func Test_parseStringToActualType(t *testing.T) {
	mockFloatValue := 3.14
	mockIntValue := 6999
	mockBoolValue := true

	tests := []struct {
		name          string
		originString  string
		expectedValue any
	}{
		{"Parse string to float64", fmt.Sprintf("%f", mockFloatValue), mockFloatValue},
		{"Parse string to int64", fmt.Sprintf("%d", mockIntValue), int64(mockIntValue)},
		{"Parse boolean to int64", fmt.Sprintf("%t", mockBoolValue), mockBoolValue},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseStringToActualType(tt.originString)
			require.Equal(t, result, tt.expectedValue)
		})
	}
}
