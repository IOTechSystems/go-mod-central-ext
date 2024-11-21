// Copyright (C) 2024 IOTech Ltd

package xlsx

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"

	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

func initialDevicesXlsxFileReader() (*excelize.File, io.Reader, error) {
	f, err := createDeviceXlsxTemplateFile()
	if err != nil {
		return nil, nil, err
	}
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, nil, err
	}
	return f, buffer, nil
}

func Test_newDevicesXlsxWriter(t *testing.T) {
	f, buffer, err := initialDevicesXlsxFileReader()
	require.NoError(t, err)
	defer f.Close()

	xlsxWriter, err := newXlsxWriter(mockDevices, buffer)
	defer xlsxWriter.(*devicesXlsxWriter).xlsFile.Close()

	require.NoError(t, err)
	require.Equal(t, mockDevices[0].Name, xlsxWriter.(*devicesXlsxWriter).devices[0].Name)
}

func Test_ConvertDevicesToXlsx(t *testing.T) {
	f, buffer, err := initialDevicesXlsxFileReader()
	require.NoError(t, err)
	defer f.Close()

	xlsxWriter, err := newXlsxWriter(mockDevices, buffer)
	defer xlsxWriter.(*devicesXlsxWriter).xlsFile.Close()
	require.NoError(t, err)

	err = xlsxWriter.ConvertToXlsx()
	require.NoError(t, err)
}

func Test_DevicesToXlsxWrite(t *testing.T) {
	f, buffer, err := initialDevicesXlsxFileReader()
	require.NoError(t, err)
	defer f.Close()

	xlsxWriter, err := newXlsxWriter(mockDevices, buffer)
	defer xlsxWriter.(*devicesXlsxWriter).xlsFile.Close()
	require.NoError(t, err)

	var outputBuffer bytes.Buffer
	err = xlsxWriter.Write(&outputBuffer)
	require.NoError(t, err)
}

func Test_DevicesToXlsxCloseFile(t *testing.T) {
	f, buffer, err := initialDevicesXlsxFileReader()
	require.NoError(t, err)
	defer f.Close()

	xlsxWriter, err := newXlsxWriter(mockDevices, buffer)
	require.NoError(t, err)

	err = xlsxWriter.closeXlsxFile()
	require.NoError(t, err)
}

func Test_deviceWriter_convertDevices(t *testing.T) {
	f, buffer, err := initialDevicesXlsxFileReader()
	require.NoError(t, err)
	defer f.Close()

	mockDevice1 := edgexDtos.Device{
		Name:           "mock-device-1",
		Description:    "test-rtu-device",
		AdminState:     "UNLOCKED",
		OperatingState: "UP",
		ServiceName:    "device-modbus",
		ProfileName:    "rtu-profile",
		Labels:         []string{"modbus-rtu-labels1", "modbus-rtu-labels2"},
		Protocols: map[string]edgexDtos.ProtocolProperties{"modbus-rtu": map[string]any{common.ModbusAddress: "/dev/virtualport", common.ModbusBaudRate: 19200, common.ModbusDataBits: 8,
			common.ModbusParity: 0, common.ModbusStopBits: 1, common.ModbusUnitID: 247,
		}},
		Tags:       map[string]any{"MachineType": "chip"},
		Properties: map[string]any{"IOTech_ProtocolName": "modbus-rtu"},
	}
	testDevices := []edgexDtos.Device{mockDevice1}
	xlsxWriter, err := newXlsxWriter(testDevices, buffer)
	require.NoError(t, err)

	err = xlsxWriter.(*devicesXlsxWriter).convertDevices()
	require.NoError(t, err)

	fileReader := xlsxWriter.(*devicesXlsxWriter).xlsFile
	defer xlsxWriter.(*devicesXlsxWriter).xlsFile.Close()

	value, err := fileReader.GetCellValue(devicesSheetName, "A2")
	require.NoError(t, err)
	require.Equal(t, mockDevice1.Name, value)

	value, err = fileReader.GetCellValue(devicesSheetName, "B2")
	require.NoError(t, err)
	require.Equal(t, mockDevice1.Description, value)

	value, err = fileReader.GetCellValue(devicesSheetName, "C2")
	require.NoError(t, err)
	require.Equal(t, mockDevice1.ServiceName, value)

	value, err = fileReader.GetCellValue(devicesSheetName, "D2")
	require.NoError(t, err)
	require.Equal(t, mockDevice1.Properties["IOTech_ProtocolName"], value)

	value, err = fileReader.GetCellValue(devicesSheetName, "E2")
	require.NoError(t, err)
	require.Equal(t, strings.Join(mockDevice1.Labels, ","), value)

	value, err = fileReader.GetCellValue(devicesSheetName, "F2")
	require.NoError(t, err)
	require.Equal(t, mockDevice1.AdminState, value)

	value, err = fileReader.GetCellValue(devicesSheetName, "G2")
	require.NoError(t, err)
	require.Equal(t, mockDevice1.Protocols["modbus-rtu"][common.ModbusAddress], value)

	value, err = fileReader.GetCellValue(devicesSheetName, "H2")
	require.NoError(t, err)
	require.Equal(t, strconv.FormatInt(int64(mockDevice1.Protocols["modbus-rtu"][common.ModbusBaudRate].(int)), 10), value)

	value, err = fileReader.GetCellValue(devicesSheetName, "I2")
	require.NoError(t, err)
	require.Equal(t, strconv.FormatInt(int64(mockDevice1.Protocols["modbus-rtu"][common.ModbusDataBits].(int)), 10), value)

	value, err = fileReader.GetCellValue(devicesSheetName, "J2")
	require.NoError(t, err)
	require.Equal(t, strconv.FormatInt(int64(mockDevice1.Protocols["modbus-rtu"][common.ModbusParity].(int)), 10), value)

	value, err = fileReader.GetCellValue(devicesSheetName, "K2")
	require.NoError(t, err)
	require.Equal(t, strconv.FormatInt(int64(mockDevice1.Protocols["modbus-rtu"][common.ModbusStopBits].(int)), 10), value)

	value, err = fileReader.GetCellValue(devicesSheetName, "L2")
	require.NoError(t, err)
	require.Equal(t, strconv.FormatInt(int64(mockDevice1.Protocols["modbus-rtu"][common.ModbusUnitID].(int)), 10), value)

	value, err = fileReader.GetCellValue(devicesSheetName, "M2")
	require.NoError(t, err)
	require.Equal(t, mockDevice1.ProfileName, value)

	value, err = fileReader.GetCellValue(devicesSheetName, "N2")
	require.NoError(t, err)
	require.Equal(t, mockDevice1.Tags[mockTagsHeader], value)
}

func Test_deviceWriter_convertAutoEvents(t *testing.T) {
	f, buffer, err := initialDevicesXlsxFileReader()
	require.NoError(t, err)
	defer f.Close()

	mockAutoEvent := edgexDtos.AutoEvent{
		Interval:   "5s",
		OnChange:   false,
		SourceName: "refSource",
	}
	testdevice := mockDevice
	testdevice.AutoEvents = []edgexDtos.AutoEvent{mockAutoEvent}

	xlsxWriter, err := newXlsxWriter([]edgexDtos.Device{testdevice}, buffer)
	require.NoError(t, err)

	err = xlsxWriter.(*devicesXlsxWriter).convertAutoEvents()
	require.NoError(t, err)

	fileReader := xlsxWriter.(*devicesXlsxWriter).xlsFile
	defer xlsxWriter.(*devicesXlsxWriter).xlsFile.Close()

	value, err := fileReader.GetCellValue(autoEventsSheetName, "A2")
	require.NoError(t, err)
	require.Equal(t, mockAutoEvent.Interval, value)
}

func Test_getNestedMapValue(t *testing.T) {
	mockFieldNames := []string{"foo", "bar"}
	expectedValue := "baz"
	mockMap := map[string]any{"foo": map[string]any{"bar": expectedValue}}
	result, err := getNestedMapValue(mockFieldNames, mockMap)
	require.NoError(t, err)
	require.Equal(t, expectedValue, result)
}
