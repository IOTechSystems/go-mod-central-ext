// Copyright (C) 2023-2024 IOTech Ltd

package xlsx

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

func Test_convertMappingTable_WithoutMappingSheet(t *testing.T) {
	f := excelize.NewFile()
	defer f.Close()

	_, err := convertMappingTable(f)
	require.Error(t, err, "Expected convertMappingTable error not occurred")
}

func Test_convertMappingTable_WithMappingSheet(t *testing.T) {
	f := excelize.NewFile()
	defer f.Close()

	_, err := f.NewSheet(mappingTableSheetName)
	require.NoError(t, err)

	mockObj := "mockObj"
	mockPath := "mockPath"
	mockDefault := "mockDefault"

	tests := []struct {
		name        string
		headerRow   []any
		dataRow     []any
		expectError bool
	}{
		{"convertMappingTable with no data row", nil, nil, true},
		{"convertMappingTable with invalid header", []any{"a", "b", "c"}, []any{"xxx"}, true},
		{"convertMappingTable with valid header", []any{"Object", "Path", "Default Value"}, []any{mockObj, mockPath, mockDefault}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.headerRow != nil {
				headerRow := tt.headerRow
				err = f.SetSheetRow(mappingTableSheetName, "A1", &headerRow)
				require.NoError(t, err)
			}
			if tt.dataRow != nil {
				dataRow := tt.dataRow
				err = f.SetSheetRow(mappingTableSheetName, "A2", &dataRow)
				require.NoError(t, err)
			}
			mappingResult, err := convertMappingTable(f)
			if tt.expectError {
				require.Error(t, err, "Expected convertMappingTable error not occurred")
			} else {
				require.NoError(t, err, "Unexpected convertMappingTable error occurred")
				require.Equal(t, mockPath, mappingResult[mockObj].path)
				require.Equal(t, mockDefault, mappingResult[mockObj].defaultValue)
			}
		})
	}
}
