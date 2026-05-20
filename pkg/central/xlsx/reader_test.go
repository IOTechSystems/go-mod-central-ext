// Copyright (C) 2023-2025 IOTech Ltd

package xlsx

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/xrtmodels"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"

	"github.com/stretchr/testify/require"
)

func Test_readStruct(t *testing.T) {
	testStr := "testString"
	testValidDevice := edgexDtos.Device{}
	deviceX, err := createDeviceXlsxInst()
	require.NoError(t, err)

	validMappings := deviceX.(*deviceXlsx).fieldMappings
	tests := []struct {
		name        string
		structPtr   *edgexDtos.Device
		headerRow   []string
		dataRow     []string
		expectError bool
	}{
		{"readStruct with invalid ptr", nil, nil, nil, true},
		{"readStruct with valid value type", &testValidDevice, []string{"Location"}, []string{"test"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.structPtr != nil {
				_, err = readStruct(tt.structPtr, tt.headerRow, tt.dataRow, validMappings)
			} else {
				_, err = readStruct(&testStr, tt.headerRow, tt.dataRow, validMappings)
			}
			if tt.expectError {
				require.Error(t, err, "Expected readStruct parse error not occurred")
			} else {
				require.NoError(t, err, "Unexpected readStruct parse error occurred")
				require.Equal(t, "test", testValidDevice.Location)
			}
		})
	}
}

func Test_getStructFieldByHeader(t *testing.T) {
	rowElement := reflect.New(reflect.TypeOf(edgexDtos.DeviceProfile{})).Elem()
	colNameWithoutSpace := "Manufacturer"
	colNameWithSpace := " " + colNameWithoutSpace
	headerCol := []string{"Name", colNameWithSpace}
	headerName, field := getStructFieldByHeader(&rowElement, 0, headerCol)
	require.Equal(t, "Name", headerName)
	require.Equal(t, reflect.String, field.Kind())

	headerName2, field2 := getStructFieldByHeader(&rowElement, 1, headerCol)
	require.Equal(t, colNameWithoutSpace, headerName2)
	require.Equal(t, reflect.String, field2.Kind())
}

func Test_setStdStructFieldValue(t *testing.T) {
	rowElement := reflect.New(reflect.TypeOf(edgexDtos.Device{})).Elem()
	labels := rowElement.FieldByName("Labels")
	tests := []struct {
		name        string
		cellValue   string
		field       reflect.Value
		expectError bool
	}{
		{"setStdStructFieldValue - fail to parse cell to bool field", "invalid", reflect.ValueOf(true), true},
		{"setStdStructFieldValue - success to parse cell to slice field", "a,b,c", labels, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setStdStructFieldValue(tt.cellValue, tt.field)
			if tt.expectError {
				require.Error(t, err, "Expected cell conversion error not occurred")
			} else {
				require.NoError(t, err, "Unexpected error occurred")

			}
		})
	}
}

func Test_setStdStructFieldValue_ptr_field(t *testing.T) {
	rowElement := reflect.New(reflect.TypeOf(edgexDtos.ResourceProperties{})).Elem()

	tests := []struct {
		name        string
		cellValue   string
		fieldName   string
		expectError bool
	}{
		{"setStdStructFieldValue_ptr_field - success to parse a cell to a pointer field", "0.9", "Minimum", false},
		{"setStdStructFieldValue_ptr_field - success to parse empty cell to a pointer field", "", "Scale", false},
		{"setStdStructFieldValue_ptr_field - fail to parse float number to a uint64 pointer field", "1.25", "Mask", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := rowElement.FieldByName(tt.fieldName)
			err := setStdStructFieldValue(tt.cellValue, field)
			if tt.expectError {
				require.Error(t, err, "Expected cell conversion error not occurred")
			} else {
				require.NoError(t, err, "Unexpected error occurred")
			}
		})
	}
}

func TestParseCellToField(t *testing.T) {

	tests := []struct {
		name        string
		cellValue   string
		kind        reflect.Kind
		expectValue any
		expectError bool
	}{
		{"success to parse a cell to float32 field", "0.9", reflect.Float32, float32(0.9), false},
		{"success to parse a cell to int64 field", "54321", reflect.Int64, int64(54321), false},
		{"success to parse a cell to uint32 field", "12345678", reflect.Uint32, uint32(12345678), false},
		{"fail to parse a cell to a unhandled type field", "12345678", reflect.Chan, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := parseCellToField(tt.cellValue, tt.kind)
			if tt.expectError {
				require.Error(t, err, "Expected cell conversion error not occurred")
			} else {
				require.Equal(t, tt.expectValue, value)
				require.NoError(t, err, "Unexpected error occurred")
			}
		})
	}
}

func Test_convertDeviceFields(t *testing.T) {
	extraPrtPropValue := "bar"
	headerCol := []string{"Name", common.ModbusAddress, common.ModbusBaudRate, "ProtocolName", mockTagsHeader, mockExtraPrtPropName}

	validDataRow := []string{"TestDevice", mockDeviceAddress, strconv.FormatInt(int64(mockDeviceBaudRate), 10), "", mockTags1, extraPrtPropValue}
	deviceX, err := createDeviceXlsxInst()
	require.NoError(t, err)

	validMappings := deviceX.(*deviceXlsx).fieldMappings

	tests := []struct {
		name          string
		dataRow       []string
		headerCol     []string
		fieldMappings map[string]mappingField
		expectError   bool
	}{
		{"Invalid convertDeviceFields - no fieldMappings", validDataRow, headerCol, nil, true},
		{"Valid convertDeviceFields", validDataRow, headerCol, validMappings, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			structPtr := edgexDtos.Device{}
			v := reflect.ValueOf(&structPtr)
			elementType := v.Elem().Type()
			element := reflect.New(elementType).Elem()

			err := convertDeviceFields(&element, tt.dataRow, tt.headerCol, tt.fieldMappings)
			if tt.expectError {
				require.Error(t, err, "Expected convertDeviceFields error not occurred")
			} else {
				require.NoError(t, err, "Unexpected convertDeviceFields error occurred")
				v.Elem().Set(element)
				require.Equal(t, modbusRTU, structPtr.Properties[common.ProtocolName])
				require.Equal(t, mockDeviceAddress, structPtr.Protocols[modbusRTU][common.ModbusAddress])
				require.Equal(t, int64(mockDeviceBaudRate), structPtr.Protocols[modbusRTU][common.ModbusBaudRate])
				require.Equal(t, extraPrtPropValue, structPtr.Protocols[mockExtraPropObj][mockExtraPrtPropName])
				require.Equal(t, mockTags1, structPtr.Tags[mockTagsHeader])
			}
		})
	}
}

func Test_convertAutoEventFields(t *testing.T) {
	rowElement := reflect.New(reflect.TypeOf(edgexDtos.AutoEvent{})).Elem()
	headerCol := []string{"Interval", "OnChange", "Reference Device Name"}
	invalidDataRow := []string{"3s", "invalid"}
	expectedDevice := common.TestDeviceName
	validDataRow := []string{"3s", "true", expectedDevice}
	deviceX, err := createDeviceXlsxInst()
	require.NoError(t, err)

	validMappings := deviceX.(*deviceXlsx).fieldMappings

	tests := []struct {
		name          string
		rowElement    *reflect.Value
		dataRow       []string
		headerCol     []string
		fieldMappings map[string]mappingField
		expectError   bool
	}{
		{"Invalid convertAutoEventFields - no fieldMappings", &rowElement, validDataRow, headerCol, nil, true},
		{"Invalid convertAutoEventFields - invalid OnChange cell", &rowElement, invalidDataRow, headerCol, validMappings, true},
		{"Valid convertAutoEventFields", &rowElement, validDataRow, headerCol, validMappings, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deviceNames, err := convertAutoEventFields(tt.rowElement, tt.dataRow, tt.headerCol, tt.fieldMappings)
			if tt.expectError {
				require.Error(t, err, "Expected convertAutoEventFields error not occurred")
			} else {
				require.NoError(t, err, "Unexpected convertAutoEventFields error occurred")
				require.Equal(t, 1, len(deviceNames))
				require.Equal(t, expectedDevice, deviceNames[0])
			}
		})
	}
}

func Test_convertScheduleFields(t *testing.T) {
	deviceX, err := createScheduleDeviceXlsxInst()
	require.NoError(t, err)
	validMappings := deviceX.(*deviceXlsx).fieldMappings

	t.Run("RefDeviceName as last header collects extra device names via fallback", func(t *testing.T) {
		rowElement := reflect.New(reflect.TypeOf(xrtmodels.Schedule{})).Elem()
		header := []string{"Interval", refDeviceName}
		// row has more cells than header: cells past last header fall back to refDeviceName
		row := []string{"1s", "Sensor0001", "Sensor0002"}

		deviceNames, err := convertScheduleFields(&rowElement, row, header, validMappings)
		require.NoError(t, err)
		require.Equal(t, []string{"Sensor0001", "Sensor0002"}, deviceNames)
	})

	t.Run("non-refDeviceName unmapped column is ignored", func(t *testing.T) {
		rowElement := reflect.New(reflect.TypeOf(xrtmodels.Schedule{})).Elem()
		header := []string{"Interval", refDeviceName, "Note"}
		row := []string{"1s", "Sensor0001", "irrelevant"}

		deviceNames, err := convertScheduleFields(&rowElement, row, header, validMappings)
		require.NoError(t, err)
		require.Equal(t, []string{"Sensor0001"}, deviceNames)
	})

	t.Run("nil fieldMappings returns error", func(t *testing.T) {
		rowElement := reflect.New(reflect.TypeOf(xrtmodels.Schedule{})).Elem()
		_, err := convertScheduleFields(&rowElement, []string{"1s", "Sensor0001"}, []string{"Interval", refDeviceName}, nil)
		require.Error(t, err)
	})
}

func Test_convertScheduleFields_MultiDeviceNamesAfterInsertedColumn(t *testing.T) {
	// simulate post-insertion state: parseSchedulesHeader appended "InsertedCol" so RefDeviceName is no longer last.
	// row still extends past header to provide a second device name; both must still be collected.
	deviceX, err := createScheduleDeviceXlsxInst()
	require.NoError(t, err)
	validMappings := deviceX.(*deviceXlsx).fieldMappings

	rowElement := reflect.New(reflect.TypeOf(xrtmodels.Schedule{})).Elem()
	header := []string{"Interval", refDeviceName, "InsertedCol"}
	row := []string{"1s", "Sensor0001", "default-value", "Sensor0002"}

	deviceNames, err := convertScheduleFields(&rowElement, row, header, validMappings)
	require.NoError(t, err)
	require.Equal(t, []string{"Sensor0001", "Sensor0002"}, deviceNames)
}

func Test_convertAutoEventFields_MultiDeviceNamesAfterInsertedColumn(t *testing.T) {
	deviceX, err := createDeviceXlsxInst()
	require.NoError(t, err)
	validMappings := deviceX.(*deviceXlsx).fieldMappings

	rowElement := reflect.New(reflect.TypeOf(edgexDtos.AutoEvent{})).Elem()
	header := []string{"Interval", "OnChange", "SourceName", refDeviceName, "InsertedCol"}
	row := []string{"1s", "true", "temperature", "Sensor0001", "default-value", "Sensor0002"}

	deviceNames, err := convertAutoEventFields(&rowElement, row, header, validMappings)
	require.NoError(t, err)
	require.Equal(t, []string{"Sensor0001", "Sensor0002"}, deviceNames)
}

func Test_convertAutoEventFields_IgnoresUnmappedColumns(t *testing.T) {
	deviceX, err := createDeviceXlsxInst()
	require.NoError(t, err)
	validMappings := deviceX.(*deviceXlsx).fieldMappings

	rowElement := reflect.New(reflect.TypeOf(edgexDtos.AutoEvent{})).Elem()
	header := []string{"Interval", "OnChange", "SourceName", refDeviceName, "Note"}
	row := []string{"1s", "true", "temperature", "Sensor0001", "irrelevant"}

	deviceNames, err := convertAutoEventFields(&rowElement, row, header, validMappings)
	require.NoError(t, err)
	require.Equal(t, []string{"Sensor0001"}, deviceNames)
}

func Test_convertDeviceCommandFields(t *testing.T) {
	rowElement := reflect.New(reflect.TypeOf(edgexDtos.DeviceCommand{})).Elem()
	headerCol := []string{"Name", "IsHidden"}
	invalidDataRow := []string{"testCommand", "invalid"}
	validDataRow := []string{"testCommand", "true"}

	tests := []struct {
		name        string
		rowElement  *reflect.Value
		dataRow     []string
		headerCol   []string
		expectError bool
	}{
		{"Invalid convertDeviceCommandFields - invalid IsHidden cell", &rowElement, invalidDataRow, headerCol, true},
		{"Valid convertDeviceCommandFields", &rowElement, validDataRow, headerCol, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := convertDeviceCommandFields(tt.rowElement, tt.dataRow, tt.headerCol)
			if tt.expectError {
				require.Error(t, err, "Expected convertAutoEventFields error not occurred")
			} else {
				require.NoError(t, err, "Unexpected convertAutoEventFields error occurred")
			}
		})
	}
}

func Test_convertResourcesFields(t *testing.T) {
	headerCol := []string{"Name", "IsHidden", "ValueType", "nodeAttribute"}
	invalidIsHiddenRow := []string{"testCommand", "invalid", "Int64"}
	validDataRow := []string{"testCommand", "true", "Int64", "value"}
	deviceX, err := createDeviceProfileXlsxInst()
	require.NoError(t, err)

	validMappings := deviceX.(*deviceProfileXlsx).fieldMappings
	tests := []struct {
		name          string
		dataRow       []string
		headerCol     []string
		fieldMappings map[string]mappingField
		expectError   bool
	}{
		{"Invalid convertResourcesFields - no fieldMappings", validDataRow, headerCol, nil, true},
		{"Invalid convertResourcesFields - invalid IsHidden cell", invalidIsHiddenRow, headerCol, validMappings, true},
		{"Valid convertResourcesFields", validDataRow, headerCol, validMappings, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			structPtr := edgexDtos.DeviceResource{}
			v := reflect.ValueOf(&structPtr)
			elementType := v.Elem().Type()
			element := reflect.New(elementType).Elem()

			err := convertResourcesFields(&element, tt.dataRow, tt.headerCol, tt.fieldMappings)
			v.Elem().Set(element)
			if tt.expectError {
				require.Error(t, err, "Expected convertResourcesFields error not occurred")
			} else {
				require.NoError(t, err, "Unexpected convertResourcesFields error occurred")
				require.Equal(t, validDataRow[0], structPtr.Name)
				require.Equal(t, validDataRow[1], strconv.FormatBool(structPtr.IsHidden))
				require.Equal(t, validDataRow[2], structPtr.Properties.ValueType)
				require.Equal(t, validDataRow[3], structPtr.Attributes[headerCol[3]])
			}
		})
	}
}

func Test_convertResourcesFields_Nested_Attributes(t *testing.T) {
	nestedAttrName1 := "dataTypeId.identifier"
	nestedAttrName2 := "dataTypeId.identifierType"

	headerCol := []string{"Name", nestedAttrName1, nestedAttrName2}
	dataRow := []string{"testCommand", "8", "NUMERIC"}
	deviceProfileX, err := createDeviceProfileXlsxInst()
	require.NoError(t, err)

	fieldMappings := deviceProfileX.(*deviceProfileXlsx).fieldMappings

	structPtr := edgexDtos.DeviceResource{}
	v := reflect.ValueOf(&structPtr)
	elementType := v.Elem().Type()
	element := reflect.New(elementType).Elem()

	err = convertResourcesFields(&element, dataRow, headerCol, fieldMappings)
	require.NoError(t, err)
	v.Elem().Set(element)

	require.Equal(t, dataRow[0], structPtr.Name)

	// check the converted nested attributes int64 value
	splitAttrNames := strings.Split(nestedAttrName1, mappingPathSeparator)
	if innerAttr, ok := structPtr.Attributes[splitAttrNames[0]].(map[string]any); ok {
		if attrVal, innerOk := innerAttr[splitAttrNames[1]].(int64); innerOk {
			require.Equal(t, dataRow[1], strconv.FormatInt(attrVal, 10))
		}
	}

	// check the converted nested attributes string value
	splitAttrNames = strings.Split(nestedAttrName2, mappingPathSeparator)
	if innerAttr, ok := structPtr.Attributes[splitAttrNames[0]].(map[string]any); ok {
		if attrVal, innerOk := innerAttr[splitAttrNames[1]]; innerOk {
			require.Equal(t, dataRow[2], attrVal)
		}
	}
}
