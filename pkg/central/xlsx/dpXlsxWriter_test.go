//
// Copyright (C) 2024 IOTech Ltd
//

package xlsx

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"testing"

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"

	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

func initialXlsxFileReader() (*excelize.File, io.Reader, error) {
	f, err := createXlsxTemplateFile()
	if err != nil {
		return nil, nil, err
	}
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, nil, err
	}
	return f, buffer, nil
}

func Test_newDPXlsxWriter(t *testing.T) {
	f, buffer, err := initialXlsxFileReader()
	require.NoError(t, err)
	defer f.Close()

	xlsxWriter, err := newXlsxWriter(mockDeviceProfile, buffer)
	defer xlsxWriter.(*dpXlsxWriter).xlsFile.Close()

	require.NoError(t, err)
	require.Equal(t, mockProfileName, xlsxWriter.(*dpXlsxWriter).deviceProfile.Name)
}

func Test_ConvertToXlsx(t *testing.T) {
	f, buffer, err := initialXlsxFileReader()
	require.NoError(t, err)
	defer f.Close()

	xlsxWriter, err := newXlsxWriter(mockDeviceProfile, buffer)
	defer xlsxWriter.(*dpXlsxWriter).xlsFile.Close()
	require.NoError(t, err)

	err = xlsxWriter.ConvertToXlsx()
	require.NoError(t, err)
}

func Test_DPToXlsxWrite(t *testing.T) {
	f, buffer, err := initialXlsxFileReader()
	require.NoError(t, err)
	defer f.Close()

	xlsxWriter, err := newXlsxWriter(mockDeviceProfile, buffer)
	defer xlsxWriter.(*dpXlsxWriter).xlsFile.Close()
	require.NoError(t, err)

	var outputBuffer bytes.Buffer
	err = xlsxWriter.Write(&outputBuffer)
	require.NoError(t, err)
}

func Test_DPToXlsxCloseFile(t *testing.T) {
	f, buffer, err := initialXlsxFileReader()
	require.NoError(t, err)
	defer f.Close()

	xlsxWriter, err := newXlsxWriter(mockDeviceProfile, buffer)
	require.NoError(t, err)

	err = xlsxWriter.closeXlsxFile()
	require.NoError(t, err)
}

func Test_dpWriter_convertDeviceInfo(t *testing.T) {
	f, buffer, err := initialXlsxFileReader()
	require.NoError(t, err)
	defer f.Close()

	xlsxWriter, err := newXlsxWriter(mockDeviceProfile, buffer)
	defer xlsxWriter.(*dpXlsxWriter).xlsFile.Close()
	require.NoError(t, err)

	err = xlsxWriter.(*dpXlsxWriter).convertDeviceInfo()
	require.NoError(t, err)

	value, err := xlsxWriter.(*dpXlsxWriter).xlsFile.GetCellValue(deviceInfoSheetName, "B1")
	require.NoError(t, err)
	require.Equal(t, mockProfileName, value)
}

func Test_dpWriter_convertDeviceResources(t *testing.T) {
	f, buffer, err := initialXlsxFileReader()
	require.NoError(t, err)
	defer f.Close()

	mockAttrValue := "HOLDING_REGISTERS"
	minValue := float64(0)
	mockNestedIdAttrValue := 8
	mockResource := edgexDtos.DeviceResource{
		Description: "this is the mockRes1 resource",
		Name:        "mockRes1",
		IsHidden:    false,
		Properties: edgexDtos.ResourceProperties{
			ValueType: "Float32",
			ReadWrite: edgexCommon.ReadWrite_R,
			Minimum:   &minValue,
		},
		Attributes: map[string]any{"primaryTable": mockAttrValue, "dataTypeId": map[string]any{"identifier": mockNestedIdAttrValue}},
	}
	mockDeviceProfile.DeviceResources = []edgexDtos.DeviceResource{mockResource}
	xlsxWriter, err := newXlsxWriter(mockDeviceProfile, buffer)
	require.NoError(t, err)

	err = xlsxWriter.(*dpXlsxWriter).convertDeviceResources()
	require.NoError(t, err)

	fileReader := xlsxWriter.(*dpXlsxWriter).xlsFile
	defer xlsxWriter.(*dpXlsxWriter).xlsFile.Close()

	value, err := fileReader.GetCellValue(deviceResourceSheetName, "A2")
	require.NoError(t, err)
	require.Equal(t, mockResource.Name, value)

	value, err = fileReader.GetCellValue(deviceResourceSheetName, "B2")
	require.NoError(t, err)
	require.Equal(t, strconv.FormatBool(mockResource.IsHidden), value)

	value, err = fileReader.GetCellValue(deviceResourceSheetName, "C2")
	require.NoError(t, err)
	require.Equal(t, mockResource.Description, value)

	value, err = fileReader.GetCellValue(deviceResourceSheetName, "D2")
	require.NoError(t, err)
	require.Equal(t, mockResource.Properties.ValueType, value)

	value, err = fileReader.GetCellValue(deviceResourceSheetName, "E2")
	require.NoError(t, err)
	require.Equal(t, mockResource.Properties.ReadWrite, value)

	value, err = fileReader.GetCellValue(deviceResourceSheetName, "F2")
	require.NoError(t, err)
	require.Equal(t, mockAttrValue, value)

	value, err = fileReader.GetCellValue(deviceResourceSheetName, "G2")
	require.NoError(t, err)
	require.Equal(t, strconv.FormatFloat(*mockResource.Properties.Minimum, 'g', -1, 64), value)

	value, err = fileReader.GetCellValue(deviceResourceSheetName, "H2")
	require.NoError(t, err)
	require.Equal(t, strconv.FormatInt(int64(mockNestedIdAttrValue), 10), value)
}

func Test_dpWriter_convertDeviceCommand(t *testing.T) {
	f, buffer, err := initialXlsxFileReader()
	require.NoError(t, err)
	defer f.Close()

	mockDeviceCommand := edgexDtos.DeviceCommand{
		Name:     "mockCmd1",
		IsHidden: false,
	}
	testProfile := mockDeviceProfile
	testProfile.DeviceCommands = []edgexDtos.DeviceCommand{mockDeviceCommand}
	xlsxWriter, err := newXlsxWriter(testProfile, buffer)
	require.NoError(t, err)
	defer xlsxWriter.(*dpXlsxWriter).xlsFile.Close()

	err = xlsxWriter.(*dpXlsxWriter).convertDeviceCommand()
	require.NoError(t, err)

	value, err := xlsxWriter.(*dpXlsxWriter).xlsFile.GetCellValue(deviceCommandSheetName, "B1")
	require.NoError(t, err)
	require.Equal(t, mockDeviceCommand.Name, value)
}

func Test_dpWriter_setResourceNameCells(t *testing.T) {
	file, buffer, err := initialXlsxFileReader()
	require.NoError(t, err)
	defer file.Close()

	mockResOp1 := edgexDtos.ResourceOperation{
		DeviceResource: "res1",
	}
	mockResOp2 := edgexDtos.ResourceOperation{
		DeviceResource: "res2",
	}
	mockResOp3 := edgexDtos.ResourceOperation{
		DeviceResource: "res3",
	}
	mockResOPs := []edgexDtos.ResourceOperation{mockResOp1, mockResOp2, mockResOp3}

	xlsxWriter, err := newXlsxWriter(mockDeviceProfile, buffer)
	require.NoError(t, err)

	startRow := 0
	colNumber := 0
	err = xlsxWriter.(*dpXlsxWriter).setResourceNameCells(startRow, colNumber, mockResOPs)
	require.NoError(t, err)

	fileReader := xlsxWriter.(*dpXlsxWriter).xlsFile
	defer xlsxWriter.(*dpXlsxWriter).xlsFile.Close()

	value, err := fileReader.GetCellValue(deviceCommandSheetName, fmt.Sprintf("B%d", startRow+1))
	require.NoError(t, err)
	require.Equal(t, mockResOp1.DeviceResource, value)

	value, err = fileReader.GetCellValue(deviceCommandSheetName, fmt.Sprintf("B%d", startRow+2))
	require.NoError(t, err)
	require.Equal(t, mockResOp2.DeviceResource, value)

	value, err = fileReader.GetCellValue(deviceCommandSheetName, fmt.Sprintf("B%d", startRow+3))
	require.NoError(t, err)
	require.Equal(t, mockResOp3.DeviceResource, value)
}
