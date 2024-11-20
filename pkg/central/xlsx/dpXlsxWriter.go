//
// Copyright (C) 2024 IOTech Ltd
//

package xlsx

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"

	"github.com/xuri/excelize/v2"
)

// deviceProfileXlsx stores the worksheets processed result and the converted DeviceProfile DTO
type dpXlsxWriter struct {
	baseXlsx
	deviceProfile edgexDtos.DeviceProfile
}

func newXlsxWriter[T AllowedDTOConverterTypes](s T, xlsxReader io.Reader) (DTOConverter[T], errors.EdgeX) {
	// file io.Reader should be closed from the caller in ConvertToDeviceProfileXlsx method
	f, err := excelize.OpenReader(xlsxReader)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindServerError, "failed to open xlsx template file from io.Reader", err)
	}

	switch any(s).(type) {
	case []edgexDtos.Device:
		fieldMappings, edgexErr := convertMappingTable(f)
		if edgexErr != nil {
			return nil, errors.NewCommonEdgeXWrapper(edgexErr)
		}
		return &devicesXlsxWriter{
			baseXlsx: baseXlsx{
				xlsFile:       f,
				fieldMappings: fieldMappings,
			},
			devices: any(s).([]edgexDtos.Device),
		}, nil
	case edgexDtos.DeviceProfile:
		return &dpXlsxWriter{
			baseXlsx: baseXlsx{
				xlsFile: f,
			},
			deviceProfile: any(s).(edgexDtos.DeviceProfile),
		}, nil
	default:
		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, "unknown DTO type for xlsx writer", nil)
	}
}

// ConvertToXlsx converts the DeviceProfile DTO into xlsx file
func (dpWriter *dpXlsxWriter) ConvertToXlsx() errors.EdgeX {
	// convert to the DeviceInfo sheet
	edgexErr := dpWriter.convertDeviceInfo()
	if edgexErr != nil {
		return errors.NewCommonEdgeXWrapper(edgexErr)
	}

	// convert to the DeviceResource sheet
	if len(dpWriter.deviceProfile.DeviceResources) > 0 {
		edgexErr = dpWriter.convertDeviceResources()
		if edgexErr != nil {
			return errors.NewCommonEdgeXWrapper(edgexErr)
		}
	}

	// convert to the DeviceCommand sheet
	if len(dpWriter.deviceProfile.DeviceCommands) > 0 {
		edgexErr = dpWriter.convertDeviceCommand()
		if edgexErr != nil {
			return errors.NewCommonEdgeXWrapper(edgexErr)
		}
	}

	return nil
}
func (dpWriter *dpXlsxWriter) Write(w io.Writer) errors.EdgeX {
	// write the file to io.Writer
	err := dpWriter.xlsFile.Write(w)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "failed to write xlsx file to io.Writer", err)
	}
	return nil
}

func (dpWriter *dpXlsxWriter) closeXlsxFile() errors.EdgeX {
	err := dpWriter.xlsFile.Close()
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "failed to close xlsx file", err)
	}
	return nil
}

// convertDeviceInfo converts the DeviceProfile DTO into the DeviceInfo worksheet
func (dpWriter *dpXlsxWriter) convertDeviceInfo() errors.EdgeX {
	f := dpWriter.xlsFile
	cols, err := f.GetCols(deviceInfoSheetName)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to retrieve all columns from %s worksheet", deviceInfoSheetName), err)
	}

	if len(cols) == 0 {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("no header column defined in %s worksheet", deviceInfoSheetName), nil)
	}

	profile := dpWriter.deviceProfile
	for i, cell := range cols[0] {
		var value string
		switch strings.ToLower(cell) {
		case strings.ToLower(apiVersion):
			value = edgexCommon.ApiVersion
		case edgexCommon.Name:
			value = profile.Name
		case edgexCommon.Manufacturer:
			value = profile.Manufacturer
		case edgexCommon.Model:
			value = profile.Model
		case strings.ToLower(description):
			value = profile.Description
		case edgexCommon.Labels:
			value = strings.Join(profile.Labels, ", ")
		default:
			// unknown header
			continue
		}
		err = f.SetCellValue(deviceInfoSheetName, fmt.Sprintf("B%d", i+1), value)
		if err != nil {
			return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to set cell value in the '%s' sheet", deviceInfoSheetName), err)
		}
	}

	return nil
}

// convertDeviceResources converts the []DeviceResource DTO into the DeviceResource worksheet
func (dpWriter *dpXlsxWriter) convertDeviceResources() errors.EdgeX {
	f := dpWriter.xlsFile
	rows, err := f.GetRows(deviceResourceSheetName)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to retrieve all rows from %s worksheet", deviceResourceSheetName), err)
	}

	if len(rows) == 0 {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("no header row defined in %s worksheet", deviceResourceSheetName), nil)
	}

	headerRow := rows[0]
	structPtr := &edgexDtos.DeviceResource{}
	v := reflect.ValueOf(structPtr)
	resources := dpWriter.deviceProfile.DeviceResources
	for resIndex, res := range resources {
		for colIndex, headerCell := range headerRow {
			if headerCell == "" {
				continue
			}

			var cell any
			field := v.Elem().FieldByName(headerCell)
			if field.Kind() != reflect.Invalid {
				// header matches the DeviceResource field name (one of the Name, Description or IsHidden field name)
				switch strings.ToLower(headerCell) {
				case edgexCommon.Name:
					cell = res.Name
				case strings.ToLower(description):
					cell = res.Description
				case strings.ToLower(isHidden):
					cell = strconv.FormatBool(res.IsHidden)
				default:
					continue
				}
			} else {
				//  header matches the ResourceProperties field name (one of ValueType, ReadWrite)
				resPropField := v.Elem().FieldByName(properties).FieldByName(headerCell)
				if resPropField.Kind() != reflect.Invalid {
					switch strings.ToLower(headerCell) {
					case strings.ToLower(edgexCommon.ValueType):
						cell = res.Properties.ValueType
					case strings.ToLower(readWrite):
						cell = res.Properties.ReadWrite
					case strings.ToLower(units):
						cell = res.Properties.Units
					case strings.ToLower(minimum):
						if res.Properties.Minimum != nil {
							cell = *res.Properties.Minimum
						}
					case strings.ToLower(maximum):
						if res.Properties.Maximum != nil {
							cell = *res.Properties.Maximum
						}
					case strings.ToLower(defaultValue):
						cell = res.Properties.DefaultValue
					case strings.ToLower(mask):
						if res.Properties.Mask != nil {
							cell = *res.Properties.Mask
						}
					case strings.ToLower(shift):
						if res.Properties.Shift != nil {
							cell = *res.Properties.Shift
						}
					case strings.ToLower(scale):
						if res.Properties.Scale != nil {
							cell = *res.Properties.Scale
						}
					case strings.ToLower(edgexCommon.Offset):
						if res.Properties.Offset != nil {
							cell = *res.Properties.Offset
						}
					case strings.ToLower(base):
						if res.Properties.Base != nil {
							cell = *res.Properties.Base
						}
					case strings.ToLower(assertion):
						cell = res.Properties.Assertion
					case strings.ToLower(mediaType):
						cell = res.Properties.MediaType
					default:
						continue
					}
				} else {
					attrNames := strings.Split(headerCell, mappingPathSeparator)
					attrNameLength := len(attrNames)
					if attrNameLength >= 1 {
						if topLevelAttr, ok := res.Attributes[attrNames[0]]; ok {
							if attrNameLength == 1 {
								// header is only 1-level attribute name, e.g., primaryTable
								cell = topLevelAttr
							} else {
								//handle the multi-level Attributes map, e.g., dataTypeId.identifier attribute in opc-ua
								for i := 1; i < attrNameLength; i++ {
									if topAttrMap, ok := topLevelAttr.(map[string]any); ok {
										topLevelAttr = topAttrMap[attrNames[i]]
									} else {
										return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to convert device resource attribute into column '%s' from %s worksheet", headerCell, deviceResourceSheetName), err)
									}
								}
								cell = topLevelAttr
							}
						}
					} else if tag, ok := res.Tags[headerCell]; ok {
						// header belongs to the Tags map field
						cell = tag
					}
				}
			}

			if cell == "" {
				continue
			}

			// get the current column name of the cell will be set
			columnName, err := excelize.ColumnNumberToName(colIndex + 1)
			if err != nil {
				return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to convert column number %d to name all rows from %s worksheet", colIndex+1, deviceResourceSheetName), err)
			}

			err = f.SetCellValue(deviceResourceSheetName, fmt.Sprintf("%s%d", columnName, resIndex+2), cell)
			if err != nil {
				return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to set cell value in the '%s' sheet", deviceResourceSheetName), err)
			}
		}
	}
	return nil
}

// convertDeviceCommand converts the []DeviceCommand DTO into the DeviceCommand worksheet
func (dpWriter *dpXlsxWriter) convertDeviceCommand() errors.EdgeX {
	f := dpWriter.xlsFile
	cols, err := f.GetCols(deviceCommandSheetName)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to retrieve all cols from %s worksheet", deviceCommandSheetName), err)
	}

	if len(cols) == 0 {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("no header column defined in %s worksheet", deviceCommandSheetName), nil)
	}

	headerCol := cols[0]
	commands := dpWriter.deviceProfile.DeviceCommands
OUTER:
	for cmdIndex, cmd := range commands {
		var cell any

		for colIndex, headerCell := range headerCol {
			if headerCell == "" {
				continue
			}
			//  check if header matches the DeviceCommand field name (one of IsHidden, ReadWrite or ResourceOperation)
			switch strings.ToLower(headerCell) {
			case edgexCommon.Name:
				cell = cmd.Name
			case strings.ToLower(isHidden):
				cell = cmd.IsHidden
			case strings.ToLower(readWrite):
				cell = cmd.ReadWrite
			case strings.ToLower(edgexCommon.ResourceName):
			case strings.ToLower(resourceOperation):
				edgexErr := dpWriter.setResourceNameCells(colIndex, cmdIndex, cmd.ResourceOperations)
				if err != nil {
					return errors.NewCommonEdgeXWrapper(edgexErr)
				}
				// ResourceName should be the last row in header column
				continue OUTER
			default:
				continue
			}

			// get the current column name of the cell will be set
			columnName, err := excelize.ColumnNumberToName(cmdIndex + 2)
			if err != nil {
				return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to convert column number %d to name all rows from %s worksheet", cmdIndex+2, deviceCommandSheetName), err)
			}
			err = f.SetCellValue(deviceCommandSheetName, fmt.Sprintf("%s%d", columnName, colIndex+1), cell)
			if err != nil {
				return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to set cell value in the '%s' sheet", deviceCommandSheetName), err)
			}
		}
	}
	return nil
}

// setResourceNameCells converts the []ResourceOperation.DeviceResource values into multiple cells (e.g., B4,B5,B6 will be the xlsx cell number) in DeviceCommand worksheet
func (dpWriter *dpXlsxWriter) setResourceNameCells(startRow int, colNumber int, resOps []edgexDtos.ResourceOperation) errors.EdgeX {
	// get the current column name of the ResourceName cell will be set
	columnName, err := excelize.ColumnNumberToName(colNumber + 2)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to convert column number %d to name all rows from %s worksheet", colNumber+2, deviceCommandSheetName), err)
	}

	rowNum := startRow + 1
	for i, op := range resOps {
		err = dpWriter.xlsFile.SetCellValue(deviceCommandSheetName, fmt.Sprintf("%s%d", columnName, rowNum+i), op.DeviceResource)
		if err != nil {
			return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to set cell value '%s' to ResourceName header in the '%s' sheet", op.DeviceResource, deviceCommandSheetName), err)
		}
	}
	return nil
}
