// Copyright (C) 2023-2024 IOTech Ltd

package xlsx

import (
	"bytes"
	"testing"

	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"

	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

var (
	mockProfileName   = "test"
	mockDeviceProfile = edgexDtos.DeviceProfile{
		DeviceProfileBasicInfo: edgexDtos.DeviceProfileBasicInfo{
			Name: mockProfileName,
		},
	}

	mockDevice = edgexDtos.Device{
		Name: "test-device",
	}
	mockDevices = []edgexDtos.Device{mockDevice}
)

func initialXlsxFile(sheetNames []string) (*excelize.File, error) {
	f, err := mockExcelFile(sheetNames)
	if err != nil {
		return nil, err
	}

	err = createMappingTableSheet(f)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func Test_ConvertDeviceXlsx_WithoutDeviceSheet(t *testing.T) {
	f, err := initialXlsxFile([]string{mappingTableSheetName})
	require.NoError(t, err)
	defer f.Close()

	buffer, err := f.WriteToBuffer()
	require.NoError(t, err)

	_, err = ConvertDeviceXlsx(buffer)
	require.Error(t, err, "Expected no Devices sheet error not occurred")
}

func Test_ConvertDeviceXlsx_WithDeviceSheet(t *testing.T) {
	f, err := initialXlsxFile([]string{mappingTableSheetName, devicesSheetName})
	require.NoError(t, err)
	defer f.Close()

	sw, err := f.NewStreamWriter(devicesSheetName)
	require.NoError(t, err)
	err = sw.SetRow("A1", validDeviceHeader)
	require.NoError(t, err)
	err = sw.SetRow("A2",
		[]any{
			"Sensor30001", "test-rtu-device 30001", "device-modbus", "modbus-rtu", "modbus-rtu-labels1,modbus-rtu-labels2", "LOCKED", "/dev/virtualport", "19200", "8", "O", "1", "247", "rtu-profile",
		})
	require.NoError(t, err)
	err = sw.Flush()
	require.NoError(t, err)

	buffer, err := f.WriteToBuffer()
	require.NoError(t, err)
	_, err = ConvertDeviceXlsx(buffer)
	require.NoError(t, err, "Unexpected ConvertXlsx error occurred")
}

func Test_ConvertDeviceProfileXlsx_WithoutDeviceInfoSheet(t *testing.T) {
	f, err := initialXlsxFile([]string{mappingTableSheetName})
	require.NoError(t, err)
	defer f.Close()

	buffer, err := f.WriteToBuffer()
	require.NoError(t, err)

	_, err = ConvertDeviceProfileXlsx(buffer)
	require.Error(t, err, "Expected no Devices sheet error not occurred")
}

func Test_ConvertDeviceProfileXlsx_WithDeviceInfoSheet(t *testing.T) {
	f, err := initialXlsxFile([]string{mappingTableSheetName, deviceInfoSheetName, deviceResourceSheetName})
	require.NoError(t, err)
	defer f.Close()

	err = createProfileMappingTableSheet(f)
	require.NoError(t, err)
	sw, err := f.NewStreamWriter(deviceInfoSheetName)
	require.NoError(t, err)
	err = sw.SetRow("A1", []any{"Name", mockProfileName1})
	require.NoError(t, err)
	err = sw.SetRow("A2", []any{"Manufacturer", mockManufacturer})
	require.NoError(t, err)
	err = sw.Flush()
	require.NoError(t, err)

	sw, err = f.NewStreamWriter(deviceResourceSheetName)
	require.NoError(t, err)
	err = sw.SetRow("A1", validResourceHeader)
	require.NoError(t, err)
	err = sw.SetRow("A2", validResourceRow)
	require.NoError(t, err)
	err = sw.Flush()
	require.NoError(t, err)

	buffer, err := f.WriteToBuffer()
	require.NoError(t, err)
	_, err = ConvertDeviceProfileXlsx(buffer)
	require.NoError(t, err, "Unexpected ConvertDeviceProfileXlsx error occurred")
}

func Test_ConvertToDeviceProfileXlsx(t *testing.T) {
	f, err := createXlsxTemplateFile()
	require.NoError(t, err)
	defer f.Close()

	buffer, err := f.WriteToBuffer()
	require.NoError(t, err)

	var outputBuffer bytes.Buffer
	edgexErr := ConvertToXlsx(buffer, &outputBuffer, mockDeviceProfile)
	require.NoError(t, edgexErr)
}

func Test_ConvertToDeviceProfileXlsx_InvalidTemplate(t *testing.T) {
	f, err := initialXlsxFile([]string{mappingTableSheetName, deviceInfoSheetName, deviceResourceSheetName})
	require.NoError(t, err)
	defer f.Close()

	buffer, err := f.WriteToBuffer()
	require.NoError(t, err)

	var outputBuffer bytes.Buffer
	mockDeviceProfile := edgexDtos.DeviceProfile{
		DeviceProfileBasicInfo: edgexDtos.DeviceProfileBasicInfo{
			Name: "test",
		},
	}
	edgexErr := ConvertToXlsx(buffer, &outputBuffer, mockDeviceProfile)
	require.Error(t, edgexErr)
}

func Test_ConvertToDevicesXlsx(t *testing.T) {
	f, err := createDeviceXlsxTemplateFile()
	require.NoError(t, err)
	defer f.Close()

	buffer, err := f.WriteToBuffer()
	require.NoError(t, err)

	var outputBuffer bytes.Buffer
	edgexErr := ConvertToXlsx(buffer, &outputBuffer, mockDevices)
	require.NoError(t, edgexErr)
}

func Test_ConvertToDevicesXlsx_InvalidTemplate(t *testing.T) {
	f, err := initialXlsxFile([]string{mappingTableSheetName, devicesSheetName, autoEventsSheetName})
	require.NoError(t, err)
	defer f.Close()

	buffer, err := f.WriteToBuffer()
	require.NoError(t, err)

	var outputBuffer bytes.Buffer
	edgexErr := ConvertToXlsx(buffer, &outputBuffer, mockDevices)
	require.Error(t, edgexErr)
}

func createXlsxTemplateFile() (*excelize.File, error) {
	f, err := initialXlsxFile([]string{mappingTableSheetName, deviceInfoSheetName, deviceResourceSheetName, deviceCommandSheetName})
	if err != nil {
		return nil, err
	}
	sw, err := f.NewStreamWriter(deviceInfoSheetName)
	if err != nil {
		return nil, err
	}
	err = sw.SetRow("A1", []any{"Name"})
	if err != nil {
		return nil, err
	}
	err = sw.Flush()
	if err != nil {
		return nil, err
	}

	resourceHeader := append(validResourceHeader, "Minimum", "dataTypeId.identifier")
	sw, err = f.NewStreamWriter(deviceResourceSheetName)
	if err != nil {
		return nil, err
	}
	err = sw.SetRow("A1", resourceHeader)
	if err != nil {
		return nil, err
	}
	err = sw.Flush()
	if err != nil {
		return nil, err
	}

	sw, err = f.NewStreamWriter(deviceCommandSheetName)
	if err != nil {
		return nil, err
	}
	err = sw.SetRow("A1", []any{"Name"})
	if err != nil {
		return nil, err
	}
	err = sw.Flush()
	if err != nil {
		return nil, err
	}
	return f, nil
}

func createDeviceXlsxTemplateFile() (*excelize.File, error) {
	f, err := initialXlsxFile([]string{mappingTableSheetName, devicesSheetName, autoEventsSheetName})
	if err != nil {
		return nil, err
	}

	err = createMappingTableSheet(f)
	if err != nil {
		return nil, err
	}

	sw, err := f.NewStreamWriter(devicesSheetName)
	if err != nil {
		return nil, err
	}
	err = sw.SetRow("A1", validDeviceHeader)
	if err != nil {
		return nil, err
	}
	err = sw.Flush()
	if err != nil {
		return nil, err
	}

	autoEventHeader := []any{"Interval", "OnChange", "SourceName", "Reference Device Name"}
	sw, err = f.NewStreamWriter(autoEventsSheetName)
	if err != nil {
		return nil, err
	}
	err = sw.SetRow("A1", autoEventHeader)
	if err != nil {
		return nil, err
	}
	err = sw.Flush()
	if err != nil {
		return nil, err
	}

	return f, nil
}
