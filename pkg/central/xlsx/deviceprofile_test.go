// Copyright (C) 2023-2024 IOTech Ltd

package xlsx

import (
	"errors"
	"testing"

	"github.com/edgexfoundry/go-mod-central-ext/v4/pkg/common"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"

	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

const (
	mockProfileName1 = "Sensor30001Profile"
	mockManufacturer = "IOTech"
)

var (
	validResourceRow    = []any{"IP_Curing_time_St_1", "true", "St_1", "Int16", "W", "INPUT_REGISTERS"}
	validResourceHeader = []any{"Name", "IsHidden", "Description", "ValueType", "ReadWrite", "primaryTable"}
)

func Test_newDeviceProfileXlsx(t *testing.T) {
	f := excelize.NewFile()
	defer f.Close()

	_, err := f.NewSheet(mappingTableSheetName)
	require.NoError(t, err)
	err = createProfileMappingTableSheet(f)
	require.NoError(t, err)
	buffer, err := f.WriteToBuffer()
	require.NoError(t, err)

	deviceXls, err := newDeviceProfileXlsx(buffer)
	require.NoError(t, err)
	require.NotEmpty(t, deviceXls)
}

func Test_newDeviceProfileXlsx_NoMappingTable(t *testing.T) {
	f, err := mockExcelFile([]string{deviceInfoSheetName})
	require.NoError(t, err)
	buffer, err := f.WriteToBuffer()
	require.NoError(t, err)

	_, err = newDeviceProfileXlsx(buffer)
	require.Error(t, err, "Expected MappingTable worksheet not defined error not occurred")
}

func createProfileMappingTableSheet(f *excelize.File) error {
	sw, err := f.NewStreamWriter(mappingTableSheetName)
	if err != nil {
		return err
	}

	err = sw.SetRow("A1",
		[]any{
			"Object", "Path", "Default Value",
		})
	if err != nil {
		return err
	}

	err = sw.SetRow("A2",
		[]any{
			"ValueType", "deviceResources[].properties.valueType", "String",
		})
	if err != nil {
		return err
	}

	err = sw.SetRow("A3",
		[]any{
			"ReadWrite", "deviceResources[].properties.readWrite", "R",
		})
	if err != nil {
		return err
	}

	err = sw.SetRow("A4",
		[]any{
			"nodeAttribute", "deviceResources[].attributes.nodeAttribute",
		})
	if err != nil {
		return err
	}

	err = sw.SetRow("A5",
		[]any{
			"deviceResources[].attributes.dataTypeId.identifier",
		})
	if err != nil {
		return err
	}

	err = sw.SetRow("A6",
		[]any{
			"deviceResources[].attributes.dataTypeId.identifierType",
		})
	if err != nil {
		return err
	}

	err = sw.Flush()
	if err != nil {
		return err
	}

	return nil
}

func createDeviceProfileXlsxInst() (Converter[*edgexDtos.DeviceProfile], error) {
	f, err := mockExcelFile([]string{deviceInfoSheetName, mappingTableSheetName, deviceResourceSheetName})
	if err != nil {
		return nil, err
	}

	err = createMappingTableSheet(f)
	if err != nil {
		return nil, err
	}

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	deviceProfileXls, err := newDeviceProfileXlsx(buffer)
	if err != nil {
		return nil, err
	}
	return deviceProfileXls, err
}

func Test_DeviceProfile_convertToDTO_InvalidSheets(t *testing.T) {
	f, err := mockExcelFile([]string{deviceInfoSheetName, mappingTableSheetName})
	require.NoError(t, err)
	err = createProfileMappingTableSheet(f)
	require.NoError(t, err)

	buffer, err := f.WriteToBuffer()
	require.NoError(t, err)
	dpX, err := newDeviceProfileXlsx(buffer)
	require.NoError(t, err)
	defer dpX.(*deviceProfileXlsx).xlsFile.Close()

	err = dpX.ConvertToDTO()
	require.Error(t, err, "Expected required worksheet not defined error not occurred")
}

func Test_deviceProfileXlsx_convertToDTO(t *testing.T) {
	dpX, err := createDeviceProfileXlsxInst()
	require.NoError(t, err)
	xlsFile := dpX.(*deviceProfileXlsx).xlsFile
	defer xlsFile.Close()

	sw, err := xlsFile.NewStreamWriter(deviceInfoSheetName)
	require.NoError(t, err)
	err = sw.SetRow("A1", []any{"Name", mockProfileName1})
	require.NoError(t, err)
	err = sw.SetRow("A2", []any{"Manufacturer", mockManufacturer})
	require.NoError(t, err)
	err = sw.Flush()
	require.NoError(t, err)

	sw, err = xlsFile.NewStreamWriter(deviceResourceSheetName)
	require.NoError(t, err)
	err = sw.SetRow("A1", validResourceHeader)
	require.NoError(t, err)
	err = sw.SetRow("A2", validResourceRow)
	require.NoError(t, err)

	err = sw.Flush()
	require.NoError(t, err)
	require.NotEmpty(t, dpX)

	err = dpX.ConvertToDTO()
	require.NoError(t, err)

	deviceProfile := dpX.GetDTOs()
	require.Equal(t, mockProfileName1, deviceProfile.Name)
	require.Equal(t, 1, len(deviceProfile.DeviceResources))
	require.Equal(t, "IP_Curing_time_St_1", deviceProfile.DeviceResources[0].Name)
}

func Test_convertDeviceInfo(t *testing.T) {
	dpX, err := createDeviceProfileXlsxInst()
	require.NoError(t, err)
	xlsFile := dpX.(*deviceProfileXlsx).xlsFile
	defer xlsFile.Close()

	convertedProfile := &edgexDtos.DeviceProfile{}
	err = dpX.(*deviceProfileXlsx).convertDeviceInfo(convertedProfile)
	require.Error(t, err, "expected \"at least 2 columns needs to be defined\" not occurred")

	err = xlsFile.SetSheetRow(deviceInfoSheetName, "A1", &[]any{"Name", mockProfileName1})
	require.NoError(t, err)
	err = xlsFile.SetSheetRow(deviceInfoSheetName, "A2", &[]any{"Manufacturer", mockManufacturer})
	require.NoError(t, err)

	err = dpX.(*deviceProfileXlsx).convertDeviceInfo(convertedProfile)
	require.NoError(t, err, "Unexpected error occurred")

	require.Equal(t, mockProfileName1, convertedProfile.Name)
	require.Equal(t, mockManufacturer, convertedProfile.Manufacturer)
}

func Test_convertDeviceResources(t *testing.T) {
	invalidIsHiddenRow := append([]any(nil), validResourceRow...)
	invalidIsHiddenRow[1] = "invalid"

	invalidReadWriteRow := append([]any(nil), validResourceRow...)
	invalidReadWriteRow[4] = "invalid"

	tests := []struct {
		name                string
		dataRow             []any
		expectError         bool
		expectValidateError bool
	}{
		{"convertDeviceResources with row count less than 2", nil, true, false},
		{"convertDeviceResources - success", validResourceRow, false, false},
		{"convertDeviceResources - invalid IsHidden", invalidIsHiddenRow, true, false},
		{"convertDeviceResources - invalid ReadWrite", invalidReadWriteRow, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dpX, err := createDeviceProfileXlsxInst()
			require.NoError(t, err)
			xlsFile := dpX.(*deviceProfileXlsx).xlsFile
			defer xlsFile.Close()

			convertedProfile := &edgexDtos.DeviceProfile{}
			err = xlsFile.SetSheetRow(deviceResourceSheetName, "A1", &validResourceHeader)
			require.NoError(t, err)
			dataRow := tt.dataRow
			err = xlsFile.SetSheetRow(deviceResourceSheetName, "A2", &dataRow)
			require.NoError(t, err)
			err = dpX.(*deviceProfileXlsx).convertDeviceResources(convertedProfile)
			if tt.expectError {
				require.Error(t, err, "Expected convertDeviceResources error not generated")
			} else {
				require.NoError(t, err)
				if tt.expectValidateError {
					require.NotNil(t, dpX.GetValidateErrors(), "Expected convertDeviceResources validation error not generated")
				} else {
					require.Equal(t, emptyValidateErr, dpX.GetValidateErrors(), "Unexpected convertDeviceResources validation error")
					require.Equal(t, 1, len(convertedProfile.DeviceResources))
					require.Equal(t, tt.dataRow[0], convertedProfile.DeviceResources[0].Name)
				}
			}
		})
	}
}

func Test_parseDeviceResourceHeader_Fail_WithoutDeviceResourceSheet(t *testing.T) {
	f, err := mockExcelFile([]string{deviceInfoSheetName, mappingTableSheetName})
	require.NoError(t, err)
	err = createProfileMappingTableSheet(f)
	require.NoError(t, err)

	buffer, err := f.WriteToBuffer()
	require.NoError(t, err)
	dpX, err := newDeviceProfileXlsx(buffer)
	require.NoError(t, err)
	defer dpX.(*deviceProfileXlsx).xlsFile.Close()

	err = dpX.(*deviceProfileXlsx).parseDeviceResourceHeader(&[]string{"Name"}, 1)
	require.Error(t, err, "Expected parseDeviceResourceHeader error not occurred")
}

func Test_parseDeviceResourceHeader_Success_WithDeviceResourceSheet(t *testing.T) {
	dpX, err := createDeviceProfileXlsxInst()
	require.NoError(t, err)
	defer dpX.(*deviceProfileXlsx).xlsFile.Close()

	err = dpX.(*deviceProfileXlsx).parseDeviceResourceHeader(&[]string{"Name"}, 1)
	require.NoError(t, err, "Unexpected parseDeviceResourceHeader error occurred")
}

func Test_DeviceProfile_convertDeviceCommands_NoDeviceCommandSheet(t *testing.T) {
	dpX, err := createDeviceProfileXlsxInst()
	require.NoError(t, err)
	defer dpX.(*deviceProfileXlsx).xlsFile.Close()

	convertedProfile := &edgexDtos.DeviceProfile{}
	err = dpX.(*deviceProfileXlsx).convertDeviceCommands(convertedProfile)
	require.Error(t, err, "Expected no DeviceCommand sheet error not occurred")
}

func Test_DeviceProfile_convertDeviceCommands(t *testing.T) {
	validDeviceCommandHeader := []any{"Name", "IsHidden", "ReadWrite", "ResourceName"}
	validDeviceCommandCol := []any{"Curing_time", "FALSE", "R", "IP_Curing_time_St_4"}
	invalidIsHiddenCol := append([]any(nil), validDeviceCommandCol...)
	invalidIsHiddenCol[1] = "invalid"

	invalidReadWriteCol := append([]any(nil), validDeviceCommandCol...)
	invalidReadWriteCol[2] = "invalid"

	tests := []struct {
		name                string
		dataRow             []any
		expectError         bool
		expectValidateError bool
	}{
		{"convertDeviceCommands with row count less than 2", []any{}, false, false},
		{"convertDeviceCommands - success", validDeviceCommandCol, false, false},
		{"convertDeviceCommands - invalid IsHidden", invalidIsHiddenCol, true, false},
		{"convertDeviceCommands - invalid ReadWrite", invalidReadWriteCol, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dpX, err := createDeviceProfileXlsxInst()
			require.NoError(t, err)
			xlsFile := dpX.(*deviceProfileXlsx).xlsFile
			defer xlsFile.Close()

			_, err = xlsFile.NewSheet(deviceCommandSheetName)
			require.NoError(t, err)
			convertedProfile := &edgexDtos.DeviceProfile{}
			err = xlsFile.SetSheetCol(deviceCommandSheetName, "A1", &validDeviceCommandHeader)
			require.NoError(t, err)
			dataRow := tt.dataRow
			err = xlsFile.SetSheetCol(deviceCommandSheetName, "B1", &dataRow)
			require.NoError(t, err)
			err = dpX.(*deviceProfileXlsx).convertDeviceCommands(convertedProfile)
			if tt.expectError {
				require.Error(t, err, "Expected convertDeviceCommands error not generated")
			} else {
				require.NoError(t, err)
				if tt.expectValidateError {
					require.NotNil(t, dpX.GetValidateErrors(), "Expected convertDeviceCommands validation error not generated")
				} else {
					require.Equal(t, emptyValidateErr, dpX.GetValidateErrors(), "Unexpected convertDeviceCommands validation error")
					if len(convertedProfile.DeviceCommands) > 0 {
						require.Equal(t, 1, len(convertedProfile.DeviceCommands))
						require.Equal(t, tt.dataRow[0], convertedProfile.DeviceCommands[0].Name)
					}
				}
			}
		})
	}
}

func Test_deviceProfileXlsx_GetDTOs(t *testing.T) {
	dpX, err := createDeviceProfileXlsxInst()
	require.NoError(t, err)
	defer dpX.(*deviceProfileXlsx).xlsFile.Close()

	deviceProfileDTO := dpX.GetDTOs()
	require.Nil(t, deviceProfileDTO)

	deviceProfileName := common.TestProfileName
	mockDeviceProfile := edgexDtos.DeviceProfile{
		DeviceProfileBasicInfo: edgexDtos.DeviceProfileBasicInfo{Name: deviceProfileName},
	}
	dpX.(*deviceProfileXlsx).deviceProfile = &mockDeviceProfile

	deviceProfile := dpX.GetDTOs()
	require.Equal(t, deviceProfileName, deviceProfile.Name)
}

func Test_deviceProfileXlsx_GetValidateErrors(t *testing.T) {
	mockDeviceProfileName := common.TestProfileName
	dpX, err := createDeviceProfileXlsxInst()
	require.NoError(t, err)
	defer dpX.(*deviceProfileXlsx).xlsFile.Close()

	validateErrs := dpX.GetValidateErrors()
	require.Equal(t, validateErrs, emptyValidateErr)

	errMsg := "test error"
	mockError := errors.New(errMsg)
	dpX.(*deviceProfileXlsx).validateErrors[mockDeviceProfileName] = mockError

	validateErrs = dpX.GetValidateErrors()
	require.NotNil(t, validateErrs)
	if actualErr, ok := validateErrs[mockDeviceProfileName]; ok {
		require.EqualError(t, actualErr, errMsg)
	} else {
		require.Fail(t, "Expected device profile validation error not found")
	}
}
