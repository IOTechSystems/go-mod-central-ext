// Copyright (C) 2023-2025 IOTech Ltd

package xlsx

import (
	"fmt"
	"io"
	"slices"

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"

	"github.com/xuri/excelize/v2"
)

const (
	validateErrResourcePrefix = "deviceResource_"
	validateErrCommandPrefix  = "deviceCommand_"
	validateErrProfilePrefix  = "deviceProfile_"
)

// requiredProfileSheets defines the required worksheet names in the xlsx file
var requiredProfileSheets = []string{deviceInfoSheetName, deviceResourceSheetName}

// deviceProfileXlsx stores the worksheets processed result and the converted DeviceProfile DTO
type deviceProfileXlsx struct {
	baseXlsx
	deviceProfile *edgexDtos.DeviceProfile
}

func newDeviceProfileXlsx(file io.Reader) (Converter[*edgexDtos.DeviceProfile], errors.EdgeX) {
	// file io.Reader should be closed from the caller in another module
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, errors.NewCommonEdgeXWrapper(err)
	}

	fieldMappings, edgexErr := convertMappingTable(f)
	if edgexErr != nil {
		return nil, errors.NewCommonEdgeXWrapper(edgexErr)
	}
	return &deviceProfileXlsx{
		baseXlsx: baseXlsx{
			xlsFile:        f,
			fieldMappings:  fieldMappings,
			validateErrors: make(map[string]error),
		},
	}, nil
}

// ConvertToDTO parses the DeviceInfo/DeviceResource/DeviceCommand sheets and convert the rows to DeviceProfile DTO
func (dpXlsx *deviceProfileXlsx) ConvertToDTO() errors.EdgeX {
	allSheetNames := dpXlsx.xlsFile.GetSheetList()

	edgexErr := checkRequiredSheets(allSheetNames, requiredProfileSheets)
	if edgexErr != nil {
		return errors.NewCommonEdgeXWrapper(edgexErr)
	}

	convertedProfile := &edgexDtos.DeviceProfile{}
	edgexErr = dpXlsx.convertDeviceInfo(convertedProfile)
	if edgexErr != nil {
		return errors.NewCommonEdgeXWrapper(edgexErr)
	}

	// parse the DeviceResource sheet
	edgexErr = dpXlsx.convertDeviceResources(convertedProfile)
	if edgexErr != nil {
		return errors.NewCommonEdgeXWrapper(edgexErr)
	}

	if slices.Contains(allSheetNames, deviceCommandSheetName) {
		// parse the DeviceCommand sheet
		edgexErr = dpXlsx.convertDeviceCommands(convertedProfile)
		if edgexErr != nil {
			return errors.NewCommonEdgeXWrapper(edgexErr)
		}
	}

	// validate the device profile DTO
	err := convertedProfile.Validate()
	if err != nil {
		dpXlsx.validateErrors[validateErrProfilePrefix+convertedProfile.Name] = err
	} else if dpXlsx.validateErrors != nil {
		dpXlsx.deviceProfile = convertedProfile
	}
	return nil
}

// convertDeviceInfo parses the DeviceInfo sheet and convert the rows to DeviceProfile DTO
func (dpXlsx *deviceProfileXlsx) convertDeviceInfo(convertedProfile *edgexDtos.DeviceProfile) errors.EdgeX {
	var header []string
	cols, err := dpXlsx.xlsFile.GetCols(deviceInfoSheetName)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to retrieve all columns from %s worksheet", deviceInfoSheetName), err)
	}

	// checks at least 2 columns exists in the DeviceInfo sheet (1 header and 1 data column)
	// and parses the header column
	if len(cols) >= 2 {
		header = cols[0]
	} else {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("at least 2 columns need to be defined in %s worksheet", deviceInfoSheetName), nil)
	}

	// parse the DeviceInfo data column
	_, err = readStruct(convertedProfile, header, cols[1], dpXlsx.fieldMappings)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "failed to unmarshal an xlsx column into DeviceProfile DTO", err)
	}
	return nil
}

// convertDeviceResources parses the DeviceResource sheet and convert the rows to DeviceResource DTOs
func (dpXlsx *deviceProfileXlsx) convertDeviceResources(convertedProfile *edgexDtos.DeviceProfile) errors.EdgeX {
	var header []string
	rows, err := dpXlsx.xlsFile.GetRows(deviceResourceSheetName)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to retrieve all rows from %s worksheet", deviceResourceSheetName), err)
	}

	// checks at least 2 rows exists in the DeviceResource sheet (1 header and 1 data row)
	// and parses the header row
	if len(rows) >= 2 {
		header = rows[0]
		err = dpXlsx.parseDeviceResourceHeader(&header, len(rows))
		if err != nil {
			return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("failed to parse the header row from %s worksheet", deviceResourceSheetName), err)
		}
	} else {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("at least 2 rows need to be defined in %s worksheet", deviceResourceSheetName), nil)
	}

	// retrieve all rows again as new columns might be added while the Header row
	rows, err = dpXlsx.xlsFile.GetRows(deviceResourceSheetName)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to retrieve all rows from %s worksheet after inserting misshing columns", deviceResourceSheetName), err)
	}

	// parse the device resource data rows
	for rowIndex, row := range rows {
		if rowIndex == 0 {
			continue
		}

		convertedDR := edgexDtos.DeviceResource{}
		_, err = readStruct(&convertedDR, header, row, dpXlsx.fieldMappings)
		if err != nil {
			return errors.NewCommonEdgeX(errors.KindContractInvalid, "failed to unmarshal an xlsx row into DeviceResource DTO", err)
		}

		// validate the DeviceResource DTO
		err = convertedDR.Validate()
		if err != nil {
			dpXlsx.validateErrors[validateErrResourcePrefix+convertedDR.Name] = err
		} else {
			convertedProfile.DeviceResources = append(convertedProfile.DeviceResources, convertedDR)
		}
	}
	return nil
}

func (dpXlsx *deviceProfileXlsx) parseDeviceResourceHeader(header *[]string, rowCount int) error {
	var err error
	// get the column count of the header row to see if any Object field from MappingTable sheet is not defined
	colCount := len(*header)

	for objectField, mapping := range dpXlsx.fieldMappings {
		// check if the mapping object is defined in the DeviceResource sheet if the defaultValue is not empty
		// if not, insert the mapping object as a new column in the DeviceResource sheet with defaultValue set in each data row
		if mapping.defaultValue != "" {
			err = checkMappingObject(dpXlsx.xlsFile, deviceResourceSheetName, &colCount, rowCount, mapping.defaultValue, objectField, header)
			if err != nil {
				return fmt.Errorf("failed to check mapping object: %w", err)
			}
		}
	}

	return nil
}

// convertDeviceCommands parses the DeviceCommand sheet and convert the rows to DeviceCommand DTOs
func (dpXlsx *deviceProfileXlsx) convertDeviceCommands(convertedProfile *edgexDtos.DeviceProfile) errors.EdgeX {
	var header []string
	cols, err := dpXlsx.xlsFile.GetCols(deviceCommandSheetName)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to retrieve all columns from %s worksheet", deviceCommandSheetName), err)
	}

	// checks at least 2 columns exists in the DeviceCommand sheet (1 header and 1 data column)
	// and parses the header column
	if len(cols) >= 2 {
		header = cols[0]
	} else {
		// not enough column defined in the DeviceCommand sheet, skip the following code and return
		return nil
	}

	for colIndex, col := range cols {
		// skip the header column
		if colIndex == 0 {
			continue
		}

		// check if the column has any non-empty cell
		// if yes, convert the xlsx column to DeviceCommand DTO
		nonEmptyCol := false
		for _, rowCell := range col {
			if rowCell != "" {
				nonEmptyCol = true
				break
			}
		}

		if nonEmptyCol {
			// parse the DeviceCommand data columns
			convertedDC := edgexDtos.DeviceCommand{}
			_, err = readStruct(&convertedDC, header, col, dpXlsx.fieldMappings)
			if err != nil {
				return errors.NewCommonEdgeX(errors.KindContractInvalid, "failed to unmarshal an xlsx column into DeviceCommand DTO", err)
			}

			// validate the DeviceCommand DTO
			err = edgexCommon.Validate(convertedDC)
			if err != nil {
				dpXlsx.validateErrors[validateErrCommandPrefix+convertedDC.Name] = err
			} else {
				convertedProfile.DeviceCommands = append(convertedProfile.DeviceCommands, convertedDC)
			}
		}
	}
	return nil
}

func (dpXlsx *deviceProfileXlsx) GetDTOs() *edgexDtos.DeviceProfile {
	return dpXlsx.deviceProfile
}

func (dpXlsx *deviceProfileXlsx) GetValidateErrors() map[string]error {
	return dpXlsx.validateErrors
}
