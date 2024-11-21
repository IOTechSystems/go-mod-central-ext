// Copyright (C) 2023-2024 IOTech Ltd

package xlsx

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"

	"github.com/xuri/excelize/v2"
)

// requiredSheets defines the required worksheet names in the xlsx file
var requiredSheets = []string{devicesSheetName}

// baseXlsx stores the basic worksheets information and the parsed validateErrors
type baseXlsx struct {
	xlsFile        *excelize.File
	fieldMappings  map[string]mappingField // fieldMappings defines all the device fields with default values defined in the xlsx
	validateErrors map[string]error
}

// deviceXlsx stores the worksheets processed result and the converted Device DTOs
type deviceXlsx struct {
	baseXlsx
	devices []*edgexDtos.Device
}

func newDeviceXlsx(file io.Reader) (Converter[[]*edgexDtos.Device], errors.EdgeX) {
	// file io.Reader should be closed from the caller in another module
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, errors.NewCommonEdgeXWrapper(err)
	}

	fieldMappings, edgexErr := convertMappingTable(f)
	if edgexErr != nil {
		return nil, errors.NewCommonEdgeXWrapper(edgexErr)
	}
	return &deviceXlsx{
		baseXlsx: baseXlsx{
			xlsFile:        f,
			fieldMappings:  fieldMappings,
			validateErrors: make(map[string]error),
		},
	}, nil
}

// ConvertToDTO parses the Devices sheet and convert the rows to Device DTOs
func (deviceXlsx *deviceXlsx) ConvertToDTO() errors.EdgeX {
	allSheetNames := deviceXlsx.xlsFile.GetSheetList()

	edgexErr := checkRequiredSheets(allSheetNames, requiredSheets)
	if edgexErr != nil {
		return errors.NewCommonEdgeXWrapper(edgexErr)
	}

	var header []string
	xlsFile := deviceXlsx.xlsFile
	protocol := deviceXlsx.fieldMappings[protocolName].defaultValue

	rows, err := xlsFile.GetRows(devicesSheetName)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to retrieve all rows from %s worksheet", devicesSheetName), err)
	}

	// checks at least 2 rows exists in the Devices sheet (1 header and 1 data row)
	// and parses the header row
	if len(rows) >= 2 {
		header = rows[0]
		edgexErr = deviceXlsx.parseDevicesHeader(&header, len(rows))
		if edgexErr != nil {
			return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("failed to parse the header row from %s worksheet", devicesSheetName), err)
		}
	} else {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("at least 2 rows need to be defined in %s worksheet", devicesSheetName), nil)
	}

	// retrieve all rows again as new columns might be added while the Header row
	rows, err = xlsFile.GetRows(devicesSheetName)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to retrieve all rows from %s worksheet after inserting misshing columns", devicesSheetName), err)
	}

	// parse the device data rows
	for rowIndex, row := range rows {
		if rowIndex == 0 {
			continue
		}

		convertedDevice := edgexDtos.Device{Properties: map[string]any{common.ProtocolName: protocol}}
		_, err = readStruct(&convertedDevice, header, row, deviceXlsx.fieldMappings)
		if err != nil {
			return errors.NewCommonEdgeX(errors.KindContractInvalid, "failed to unmarshal an xlsx row into Device DTO", err)
		}

		for k, v := range convertedDevice.Protocols {
			err = toXrtProperties(k, v)
			if err != nil {
				return errors.NewCommonEdgeX(errors.Kind(err), "failed to convert EdgeX protocol properties to device connector protocol properties", err)
			}
		}

		// validate the device DTO
		err = edgexCommon.Validate(convertedDevice)
		if err != nil {
			deviceXlsx.validateErrors[convertedDevice.Name] = err
		} else {
			deviceXlsx.devices = append(deviceXlsx.devices, &convertedDevice)
		}
	}

	if slices.Contains(allSheetNames, autoEventsSheetName) {
		err = deviceXlsx.convertAutoEvents()
		if err != nil {
			return errors.NewCommonEdgeXWrapper(err)
		}
	}

	return nil
}

func (deviceXlsx *deviceXlsx) parseDevicesHeader(header *[]string, rowCount int) errors.EdgeX {
	var err error
	// get the column count of the header row to see if any Object field from MappingTable sheet is not defined
	colCount := len(*header)

	for objectField, mapping := range deviceXlsx.fieldMappings {
		if startsWithAutoEvents(mapping.path) {
			// if the mapping path starts with autoEvents, skip the check of the Devices sheet header column
			continue
		}

		// check if the mapping object is defined in the Devices sheet if the defaultValue is not empty
		// if not, insert the mapping object as a new column in the Devices sheet with defaultValue set in each data row
		if mapping.defaultValue != "" {
			err = checkMappingObject(deviceXlsx.xlsFile, devicesSheetName, &colCount, rowCount, mapping.defaultValue, objectField, header)
			if err != nil {
				return errors.NewCommonEdgeX(errors.KindContractInvalid, "failed to check mapping object", err)
			}
		}
	}

	return nil
}

// convertAutoEvents parses the AutoEvents sheet and convert the rows to AutoEvent DTOs
func (deviceXlsx *deviceXlsx) convertAutoEvents() errors.EdgeX {
	var header []string
	xlsFile := deviceXlsx.xlsFile

	rows, err := xlsFile.GetRows(autoEventsSheetName)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to retrieve all rows from %s worksheet", autoEventsSheetName), err)
	}

	// checks at least 2 rows exists in the AutoEvents sheet (1 header and 1 data row)
	// and parses the header row
	if len(rows) >= 2 {
		header = rows[0]
		// parse the header row
		// get the column count of the header row to see if any Object field from MappingTable sheet is not defined
		colCount := len(header)

		// AutoEvents sheet should at least define 2 columns in the header row (SourceName and Reference Device Name)
		if colCount < 2 {
			err = deviceXlsx.parseAutoEventsHeader(header, len(rows))
			if err != nil {
				return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to parse the header row from %s worksheet", autoEventsSheetName), err)
			}
		}
	} else {
		return nil
	}

	rows, err = xlsFile.GetRows(autoEventsSheetName)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to retrieve all rows from %s worksheet", autoEventsSheetName), err)
	}

OUTER:
	// parse the device data rows
	for rowIndex, row := range rows {
		if rowIndex == 0 {
			continue
		}

		autoEvent := edgexDtos.AutoEvent{}
		deviceNameResult, edgexErr := readStruct(&autoEvent, header, row, deviceXlsx.fieldMappings)
		if edgexErr != nil {
			return errors.NewCommonEdgeX(errors.KindContractInvalid, "failed to unmarshal an excel row into AutoEvent DTO", err)
		}

		deviceNames, ok := deviceNameResult.([]string)
		if !ok {
			return errors.NewCommonEdgeX(errors.KindContractInvalid,
				fmt.Sprintf("failed to obtain the 'Reference Device Name' cell of the xlsx row from %s worksheet", autoEventsSheetName), nil)
		}

		// validate the AutoEvent DTO
		err = edgexCommon.Validate(autoEvent)
		if err != nil {
			for _, deviceName := range deviceNames {
				// find the matched device DTO index equals to the "Reference Device Name" on the AutoEvents row
				idx := slices.IndexFunc(deviceXlsx.devices, func(d *edgexDtos.Device) bool { return d.Name == deviceName })
				if idx > -1 {
					// delete the device element in deviceXlsx.devices slice if the referenced AutoEvent failed validation
					deviceXlsx.devices = slices.Delete(deviceXlsx.devices, idx, idx+1)
					deviceXlsx.validateErrors[deviceName] = err
				}
			}
			continue OUTER
		}

		for _, deviceName := range deviceNames {
			for _, device := range deviceXlsx.devices {
				if deviceName == device.Name {
					device.AutoEvents = append(device.AutoEvents, autoEvent)
				}
			}
		}
	}

	return nil
}

func (deviceXlsx *deviceXlsx) parseAutoEventsHeader(header []string, rowCount int) errors.EdgeX {
	var err error
	colCount := len(header)
	newColCount := &colCount

	for objectField, mapping := range deviceXlsx.fieldMappings {
		if !startsWithAutoEvents(mapping.path) {
			// if the mapping path doesn't start with autoEvents, skip the check of the AutoEvents sheet header column
			continue
		}

		// check if the mapping object is defined in the AutoEvents sheet if the defaultValue is not empty
		// if not, insert the mapping object as a new column in the Devices sheet with defaultValue set in each data row
		err = checkMappingObject(deviceXlsx.xlsFile, autoEventsSheetName, newColCount, rowCount, mapping.defaultValue, objectField, &header)
		if err != nil {
			return errors.NewCommonEdgeX(errors.KindContractInvalid, "failed to check mapping object", err)
		}
	}

	return nil
}

// startsWithAutoEvents checks if the path name defined in MappingTable sheet starts with autoEvents
func startsWithAutoEvents(path string) bool {
	return strings.HasPrefix(strings.ToLower(path), strings.ToLower(autoEvents))
}

func (deviceXlsx *deviceXlsx) GetDTOs() []*edgexDtos.Device {
	return deviceXlsx.devices
}

func (deviceXlsx *deviceXlsx) GetValidateErrors() map[string]error {
	return deviceXlsx.validateErrors
}
