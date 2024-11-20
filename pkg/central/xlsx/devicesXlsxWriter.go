// Copyright (C) 2024 IOTech Ltd

package xlsx

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"

	"github.com/xuri/excelize/v2"
)

// devicesXlsxWriter stores the worksheets processed result and the converted []Device DTO
type devicesXlsxWriter struct {
	baseXlsx
	devices []edgexDtos.Device
}

// ConvertToXlsx converts the []Device DTO into xlsx file
func (deviceWriter *devicesXlsxWriter) ConvertToXlsx() errors.EdgeX {
	// convert to the Devices sheet
	edgexErr := deviceWriter.convertDevices()
	if edgexErr != nil {
		return errors.NewCommonEdgeXWrapper(edgexErr)
	}

	// convert to the AutoEvents sheet
	edgexErr = deviceWriter.convertAutoEvents()
	if edgexErr != nil {
		return errors.NewCommonEdgeXWrapper(edgexErr)
	}

	return nil
}

// Write writes the xlsx file content to io.Writer
func (deviceWriter *devicesXlsxWriter) Write(w io.Writer) errors.EdgeX {
	// write the file to io.Writer
	err := deviceWriter.xlsFile.Write(w)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "failed to write xlsx file to io.Writer", err)
	}
	return nil
}

// closeXlsxFile closes the xlsx file reader
func (deviceWriter *devicesXlsxWriter) closeXlsxFile() errors.EdgeX {
	err := deviceWriter.xlsFile.Close()
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "failed to close xlsx file", err)
	}
	return nil
}

// convertDevices converts the []Device DTO into the Devices worksheet
func (deviceWriter *devicesXlsxWriter) convertDevices() errors.EdgeX {
	f := deviceWriter.xlsFile
	rows, err := f.GetRows(devicesSheetName)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to retrieve all rows from %s worksheet", devicesSheetName), err)
	}

	if len(rows) == 0 {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("no header row defined in %s worksheet", devicesSheetName), nil)
	}

	headerRow := rows[0]
	structPtr := &edgexDtos.Device{}
	v := reflect.ValueOf(structPtr)
	devices := deviceWriter.devices
	for deviceIndex, device := range devices {
	OUTER:
		for colIndex, headerCell := range headerRow {
			if headerCell == "" {
				continue
			}

			var cell any
			field := v.Elem().FieldByName(headerCell)
			if field.Kind() != reflect.Invalid {
				// header matches the Device field name (one of the Name, Description, Labels, AdminState, OperatingState
				// ServiceName or ProfileName field)
				switch strings.ToLower(headerCell) {
				case edgexCommon.Name:
					cell = device.Name
				case strings.ToLower(description):
					cell = device.Description
				case strings.ToLower(edgexCommon.Labels):
					cell = strings.Join(device.Labels, ",")
				case strings.ToLower(adminState):
					cell = device.AdminState
				case strings.ToLower(operatingState):
					cell = device.OperatingState
				case strings.ToLower(edgexCommon.ServiceName):
					cell = device.ServiceName
				case strings.ToLower(edgexCommon.ProfileName):
					cell = device.ProfileName
				default:
					continue
				}
			} else {
				for objectField, mapping := range deviceWriter.fieldMappings {
					// if header matches the MappingTable Object field from worksheet
					// check the Device Protocols/Properties/Tags map field from DTO and get the value
					if headerCell == objectField {
						mappingPath := strings.Split(mapping.path, mappingPathSeparator)
						mappingPathLength := len(mappingPath)
						if mappingPathLength < 2 {
							// invalid path defined in the MappingTable sheet, at least 1 dot needs to exist, e.g., properties.IOTech_ProtocolName
							continue OUTER
						}

						switch strings.ToLower(mappingPath[0]) {
						case strings.ToLower(protocols):
							if mappingPathLength < 3 {
								// invalid protocols path defined in the mapping table, at least 2 dots needs to exist
								// e.g., protocols.modbus-rtu.Address
								continue OUTER
							}

							if topLevelPrtProp, ok := device.Protocols[mappingPath[1]]; ok {
								cell, err = getNestedMapValue(mappingPath[2:], topLevelPrtProp)
								if err != nil {
									return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to get '%s' field from Protocols map in %s worksheet", headerCell, devicesSheetName), err)
								}
							} else {
								continue OUTER
							}
						case strings.ToLower(properties):
							cell, err = getNestedMapValue(mappingPath[1:], device.Properties)
							if err != nil {
								return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to get '%s' field from Properties map in %s worksheet", headerCell, devicesSheetName), err)
							}
						case strings.ToLower(tags):
							cell, err = getNestedMapValue(mappingPath[1:], device.Tags)
							if err != nil {
								return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to get '%s' field from Tags map in %s worksheet", headerCell, devicesSheetName), err)
							}
						default:
							continue
						}
					}
				}
			}

			if cell == "" {
				continue
			}

			// get the current column name of the cell will be set
			columnName, err := excelize.ColumnNumberToName(colIndex + 1)
			if err != nil {
				return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to convert column number %d to name all rows from %s worksheet", colIndex+1, devicesSheetName), err)
			}

			err = f.SetCellValue(devicesSheetName, fmt.Sprintf("%s%d", columnName, deviceIndex+2), cell)
			if err != nil {
				return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to set cell value in the '%s' sheet", devicesSheetName), err)
			}
		}
	}
	return nil
}

// convertAutoEvents converts the []AutoEvent DTO into the AutoEvents worksheet
func (deviceWriter *devicesXlsxWriter) convertAutoEvents() errors.EdgeX {
	f := deviceWriter.xlsFile
	rows, err := f.GetRows(autoEventsSheetName)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to retrieve all rows from %s worksheet", autoEventsSheetName), err)
	}

	if len(rows) == 0 {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("no header row defined in %s worksheet", autoEventsSheetName), nil)
	}

	headerRow := rows[0]
	structPtr := &edgexDtos.AutoEvent{}
	v := reflect.ValueOf(structPtr)
	totalAutoEventCount := 0

	for _, device := range deviceWriter.devices {
		for _, autoEvent := range device.AutoEvents {
			for colIndex, headerCell := range headerRow {
				if headerCell == "" {
					continue
				}

				var cell any
				field := v.Elem().FieldByName(headerCell)
				if field.Kind() != reflect.Invalid {
					// header matches the Device field name (one of the Interval, Description, Labels, AdminState, OperatingState
					// ServiceName or ProfileName field)
					switch strings.ToLower(headerCell) {
					case edgexCommon.Interval:
						cell = autoEvent.Interval
					case strings.ToLower(onChange):
						cell = autoEvent.OnChange
					case strings.ToLower(edgexCommon.SourceName):
						cell = autoEvent.SourceName
					default:
						continue
					}
				} else {
					// set device name to cell if header is "Reference Device Name"
					if strings.EqualFold(headerCell, refDeviceName) {
						cell = device.Name
					}
				}

				if cell == "" {
					continue
				}

				// get the current column name of the cell will be set
				columnName, err := excelize.ColumnNumberToName(colIndex + 1)
				if err != nil {
					return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to convert column number %d to name all rows from %s worksheet", colIndex+1, autoEventsSheetName), err)
				}

				err = f.SetCellValue(autoEventsSheetName, fmt.Sprintf("%s%d", columnName, totalAutoEventCount+2), cell)
				if err != nil {
					return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to set cell value in the '%s' sheet", autoEventsSheetName), err)
				}
			}

			totalAutoEventCount++
		}
	}
	return nil
}

// getNestedMapValue get the value of map from the passed MappingTable Path in the worksheet
// e.g., modbus-rtu.Address returns the nested 'Address' value from device protocol properties map
func getNestedMapValue(fieldNames []string, topLevelMap map[string]any) (any, errors.EdgeX) {
	var innerValue any
	fieldNameLength := len(fieldNames)

	for i := 0; i < fieldNameLength; i++ {
		if childLevelMap, ok := topLevelMap[fieldNames[i]].(map[string]any); ok {
			topLevelMap = childLevelMap
		} else {
			if i != fieldNameLength-1 {
				return innerValue, errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to get the inner value from map based on the MappingTable path %s", strings.Join(fieldNames, ".")), nil)
			}
			innerValue = topLevelMap[fieldNames[i]]
		}
	}
	return innerValue, nil
}
