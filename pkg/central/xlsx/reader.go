// Copyright (C) 2023-2025 IOTech Ltd

package xlsx

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

func readStruct(structPtr any, headerCol []string, row []string, mapppingTable map[string]mappingField) (any, errors.EdgeX) {
	var extraReturnedCols any
	v := reflect.ValueOf(structPtr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, errors.NewCommonEdgeX(errors.KindServerError, "the structPtr argument should be a pointer of struct", nil)
	}

	elementType := v.Elem().Type()
	rowElement := reflect.New(elementType).Elem()

	var err errors.EdgeX
	switch elementType {
	case reflect.TypeOf(edgexDtos.DeviceProfile{}):
		err = convertDTOStdTypeFields(&rowElement, row, headerCol, mapppingTable)
	case reflect.TypeOf(edgexDtos.AutoEvent{}):
		extraReturnedCols, err = convertAutoEventFields(&rowElement, row, headerCol, mapppingTable)
	case reflect.TypeOf(edgexDtos.Device{}):
		err = convertDeviceFields(&rowElement, row, headerCol, mapppingTable)
	case reflect.TypeOf(edgexDtos.DeviceCommand{}):
		err = convertDeviceCommandFields(&rowElement, row, headerCol)
	case reflect.TypeOf(edgexDtos.DeviceResource{}):
		err = convertResourcesFields(&rowElement, row, headerCol, mapppingTable)
	default:
		// skip the processing of the not found field name
		err = errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("unknown converted DTO type '%T'", elementType), nil)
	}
	if err != nil {
		return nil, err
	}

	v.Elem().Set(rowElement)
	return extraReturnedCols, nil
}

// getStructFieldByHeader returns the passed structEle struct field by headerName
func getStructFieldByHeader(structEle *reflect.Value, colIndex int, headerCol []string) (string, reflect.Value) {
	var headerName string
	headerLastIndex := len(headerCol) - 1
	// check if row length is larger than the header
	if colIndex > headerLastIndex {
		headerName = strings.TrimSpace(headerCol[headerLastIndex])
	} else {
		headerName = strings.TrimSpace(headerCol[colIndex])
	}
	field := structEle.FieldByName(headerName)
	return headerName, field
}

// setStdStructFieldValue set the struct field with Go standard types to the xlsx cell value
func setStdStructFieldValue(originValue string, field reflect.Value) errors.EdgeX {
	var fieldValue any
	var err errors.EdgeX

	// Handle the struct pointer field
	if field.Kind() == reflect.Ptr {
		if originValue == "" {
			return nil
		}

		fieldValue, err = parseCellToField(originValue, field.Type().Elem().Kind())
		if err != nil {
			return errors.NewCommonEdgeXWrapper(err)
		}

		// Create a new pointer field based on the primitive data type
		ptrValue := reflect.New(field.Type().Elem())
		// Set value to the pointer
		ptrValue.Elem().Set(reflect.ValueOf(fieldValue))
		// Set the struct field to the pointer value
		field.Set(ptrValue)
	} else {
		fieldValue, err = parseCellToField(originValue, field.Kind())
		if err != nil {
			return errors.NewCommonEdgeXWrapper(err)
		}
		field.Set(reflect.ValueOf(fieldValue))
	}

	return nil
}

// parseCellToField parses the xlsx cell string to the Go primitive type value based on the data type
func parseCellToField(originValue string, kind reflect.Kind) (any, errors.EdgeX) {
	var fieldValue any

	switch kind {
	case reflect.String:
		fieldValue = originValue
	case reflect.Slice:
		values := strings.Split(originValue, edgexCommon.CommaSeparator)
		fieldValue = values
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(originValue)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("failed to parse originValue '%v' to bool type", originValue), err)
		}
		fieldValue = boolValue
	case reflect.Int32:
		intValue, err := strconv.ParseInt(originValue, 10, 32)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("failed to parse originValue '%v' to Int32 type", originValue), err)
		}
		fieldValue = int32(intValue)
	case reflect.Int64:
		int64Value, err := strconv.ParseInt(originValue, 10, 64)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("failed to parse originValue '%v' to Int64 type", originValue), err)
		}
		fieldValue = int64Value
	case reflect.Float32:
		floatValue, err := strconv.ParseFloat(originValue, 32)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("failed to parse originValue '%v' to Float32 type", originValue), err)
		}
		fieldValue = float32(floatValue)
	case reflect.Float64:
		floatValue, err := strconv.ParseFloat(originValue, 64)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("failed to parse originValue '%v' to Float64 type", originValue), err)
		}
		fieldValue = floatValue
	case reflect.Uint32:
		uintValue, err := strconv.ParseUint(originValue, 10, 32)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("failed to parse originValue '%v' to Uint32 type", originValue), err)
		}
		fieldValue = uint32(uintValue)
	case reflect.Uint64:
		uintValue, err := strconv.ParseUint(originValue, 10, 64)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("failed to parse originValue '%v' to Uint64 type", originValue), err)
		}
		fieldValue = uintValue
	case reflect.Interface:
		fieldValue = originValue
	default:
		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("failed to parse originValue '%v' to %s type", originValue, kind), nil)
	}

	return fieldValue, nil
}

// convertDTOStdTypeFields unmarshalls the xlsx cells into the standard type fields of the DTO struct
func convertDTOStdTypeFields(rowElement *reflect.Value, xlsxRow []string, headerCol []string, fieldMappings map[string]mappingField) errors.EdgeX {
	for colIndex, cell := range xlsxRow {
		headerName, field := getStructFieldByHeader(rowElement, colIndex, headerCol)
		fieldValue := strings.TrimSpace(cell)

		if field.Kind() != reflect.Invalid {
			if fieldValue == "" {
				// set the struct field value to 'default value' defined in mapping Table if not empty
				if mapping, ok := fieldMappings[headerName]; ok && mapping.defaultValue != "" {
					fieldValue = mapping.defaultValue
				}
			}

			err := setStdStructFieldValue(fieldValue, field)
			if err != nil {
				return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("error occurred on '%s' column", headerName), err)
			}
		} else {
			// field not found in the DTO struct, skip this column
			continue
		}
	}
	return nil
}

// convertDeviceFields convert the xlsx row to the Device DTO
func convertDeviceFields(rowElement *reflect.Value, xlsxRow []string, headerCol []string, fieldMappings map[string]mappingField) errors.EdgeX {
	if fieldMappings == nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "fieldMappings not defined while converting device fields", nil)
	}

	prtPropMap := make(map[string]edgexDtos.ProtocolProperties)
	propertiesMap := make(map[string]any)
	tagsMap := make(map[string]any)

	for colIndex, cell := range xlsxRow {
		headerName, field := getStructFieldByHeader(rowElement, colIndex, headerCol)
		fieldValue := strings.TrimSpace(cell)
		if fieldValue == "" {
			// set fieldValue to 'default value' defined in mapping Table if not empty
			if mapping, ok := fieldMappings[headerName]; ok && mapping.defaultValue != "" {
				fieldValue = mapping.defaultValue
			}
		}

		if field.Kind() != reflect.Invalid {
			// header matches the Device DTO field name (one of the Name, Description, AdminState, OperatingState, etc)
			err := setStdStructFieldValue(fieldValue, field)
			if err != nil {
				return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("error occurred on '%s' column", headerName), err)
			}
		} else {
			// header not belongs to the above fields with standard types
			// map the cell to the Protocols or Tags field
			if fieldValue != "" {
				// get the Path defined in the MappingTable
				if mapping, ok := fieldMappings[headerName]; ok && mapping.path != "" {
					splitPaths := strings.SplitN(mapping.path, mappingPathSeparator, 2)
					fieldPrefix := strings.TrimSpace(splitPaths[0])

					fieldName := headerName
					if len(splitPaths) > 1 {
						fieldName = strings.TrimSpace(splitPaths[1])
					}

					convertedValue := parseStringToActualType(fieldValue)

					switch fieldPrefix {
					case strings.ToLower(protocols):
						// get the Device protocols field
						prtMapField := rowElement.FieldByName(protocols)
						if prtMapField.Len() > 0 {
							// convert the prtMapField reflect.Value to map[string]ProtocolProperties
							prtPropMap, ok = prtMapField.Interface().(map[string]edgexDtos.ProtocolProperties)
							if !ok {
								return errors.NewCommonEdgeX(errors.KindServerError, "failed to convert Device Protocols field to map[string]ProtocolProperties data type", nil)
							}
						}

						// to handle the nested ProtocolProperties name
						// split the ProtocolProperties name using the "." separator into array
						prtPropNames := strings.Split(fieldName, mappingPathSeparator)
						lastPropNameIdx := len(prtPropNames) - 1

						var innerPrtProp edgexDtos.ProtocolProperties
						for i, propName := range prtPropNames {
							if i == lastPropNameIdx {
								// the last part of ProtocolProperties property name
								innerPrtProp[propName] = convertedValue
							} else {
								if i == 0 {
									if _, ok := prtPropMap[propName]; !ok {
										prtPropMap[propName] = make(edgexDtos.ProtocolProperties)
									}
									// assign prtPropMap[propName] to innerPrtProp
									innerPrtProp = prtPropMap[propName]
								} else {
									if _, ok := innerPrtProp[propName]; !ok {
										// initialize a new ProtocolProperties map for inner node
										innerPrtProp[propName] = make(edgexDtos.ProtocolProperties)
									}
									innerPrtProp, ok = innerPrtProp[propName].(edgexDtos.ProtocolProperties)
									if !ok {
										return errors.NewCommonEdgeX(errors.KindServerError,
											fmt.Sprintf("failed to convert property '%s' from '%s' path to ProtocolProperties type", propName, mapping.path), nil)
									}
								}
							}
						}
						prtMapField.Set(reflect.ValueOf(prtPropMap))
					case strings.ToLower(properties):
						// set the cell to Protocols map
						propertiesMap[fieldName] = convertedValue
					case strings.ToLower(tags):
						// set the cell to Tags map
						tagsMap[headerName] = convertedValue
					default:
						// unknown column header
						continue
					}
				}
			}
		}
	}

	// set Properties field to the Device DTO struct
	if len(propertiesMap) > 0 {
		err := setMapToStructField(rowElement, properties, propertiesMap)
		if err != nil {
			return errors.NewCommonEdgeXWrapper(err)
		}
	}
	// set Tags field to the Device DTO struct
	if len(tagsMap) > 0 {
		err := setMapToStructField(rowElement, tags, tagsMap)
		if err != nil {
			return errors.NewCommonEdgeXWrapper(err)
		}
	}

	return nil
}

// convertAutoEventFields convert the xlsx row to the AutoEvent DTO
func convertAutoEventFields(rowElement *reflect.Value, xlsxRow []string, headerCol []string, fieldMappings map[string]mappingField) ([]string, errors.EdgeX) {
	if fieldMappings == nil {
		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, "fieldMappings not defined while converting AutoEvent fields", nil)
	}
	var deviceNames []string

	for colIndex, cell := range xlsxRow {
		headerName, field := getStructFieldByHeader(rowElement, colIndex, headerCol)
		fieldValue := strings.TrimSpace(cell)
		if fieldValue == "" {
			// set fieldValue to 'default value' defined in mapping Table if not empty
			if mapping, ok := fieldMappings[headerName]; ok && mapping.defaultValue != "" {
				fieldValue = mapping.defaultValue
			}
		}

		if field.Kind() != reflect.Invalid {
			// header matches the AutoEvent DTO field name (one of the Interval, OnChange, SourceName field)
			err := setStdStructFieldValue(fieldValue, field)
			if err != nil {
				return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("error occurred on '%s' column", headerName), err)
			}
		} else {
			// the cell belongs to the "Reference Device Name" column, append it to deviceNames
			if fieldValue != "" {
				deviceNames = append(deviceNames, fieldValue)
			}
		}
	}

	return deviceNames, nil
}

// convertDeviceCommandFields convert the xlsx row to the DeviceCommand DTO
func convertDeviceCommandFields(rowElement *reflect.Value, xlsxCol []string, headerCol []string) errors.EdgeX {
	var resOpSlice []edgexDtos.ResourceOperation
	for colIndex, cell := range xlsxCol {
		// skip the empty cell, all the cell should have value in DeviceCommand sheet
		if cell == "" {
			continue
		}

		headerName, field := getStructFieldByHeader(rowElement, colIndex, headerCol)
		cell = strings.TrimSpace(cell)

		if field.Kind() != reflect.Invalid {
			// header matches the DeviceCommand field name (one of the Name, IsHidden or ReadWrite field name)
			err := setStdStructFieldValue(cell, field)
			if err != nil {
				return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("error occurred on '%s' row", headerName), err)
			}
		} else {
			// parse the rest ResourceName columns in the xlsx row and convert to the ResourceOperation DTO
			resOp := edgexDtos.ResourceOperation{
				DeviceResource: cell,
			}
			resOpSlice = append(resOpSlice, resOp)
		}
	}

	if len(resOpSlice) > 0 {
		// set resOpSlice to the ResourceOperations field of DeviceCommand struct
		resOpField := rowElement.FieldByName(resourceOperations)
		if resOpField.Kind() == reflect.Invalid {
			return errors.NewCommonEdgeX(errors.KindServerError, "failed to find ResourceOperations field in DeviceCommand DTO", nil)
		}
		resOpField.Set(reflect.ValueOf(resOpSlice))
	}
	return nil
}

// convertResourcesFields convert the xlsx row to the DeviceResource DTO
func convertResourcesFields(rowElement *reflect.Value, xlsxRow []string, headerCol []string, fieldMappings map[string]mappingField) errors.EdgeX {
	if fieldMappings == nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "fieldMappings not defined while converting DeviceResource fields", nil)
	}

	for colIndex, cell := range xlsxRow {
		headerName, field := getStructFieldByHeader(rowElement, colIndex, headerCol)
		fieldValue := strings.TrimSpace(cell)
		if fieldValue == "" {
			// set fieldValue to 'default value' defined in mapping Table if not empty
			if mapping, ok := fieldMappings[headerName]; ok && mapping.defaultValue != "" {
				fieldValue = mapping.defaultValue
			}
		}

		if field.Kind() != reflect.Invalid {
			// header matches the DeviceResource field name (one of the Name, Description or IsHidden field name)
			err := setStdStructFieldValue(fieldValue, field)
			if err != nil {
				return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("error occurred on '%s' column", headerName), err)
			}
		} else {
			resPropField := rowElement.FieldByName(properties).FieldByName(headerName)
			if resPropField.Kind() != reflect.Invalid {
				// header matches the ResourceProperties DTO field name (one of the ValueType, ReadWrite, Units, etc)
				err := setStdStructFieldValue(fieldValue, resPropField)
				if err != nil {
					return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("error occurred on '%s' column", headerName), err)
				}
			} else {
				// set the cell to Attributes map if header not belongs to Properties field
				if fieldValue != "" {
					// check if the header defined in the mapping table first, and if the path contains "attributes"
					// if not, skip this column and move to the next
					if fieldMapping, ok := fieldMappings[headerName]; ok {
						if !strings.Contains(strings.ToLower(fieldMapping.path), strings.ToLower(attributes)) {
							continue
						}
					}

					var attrMap map[string]any
					attrMapField := rowElement.FieldByName(attributes)
					if attrMapField.Len() == 0 {
						// initialize the Attributes map
						attrMap = make(map[string]any)
					} else {
						attrMap = attrMapField.Interface().(map[string]any)
					}

					// parse the attribute value to the actual data type other than string if needed
					attrValue := parseStringToActualType(fieldValue)

					// to handle the nested attribute name, split the attribute name using the "." separator into array
					attrNames := strings.Split(headerName, mappingPathSeparator)
					attrNameLength := len(attrNames)
					currentAttrMap := attrMap

					for i, attrName := range attrNames {
						if i == attrNameLength-1 {
							// the last part of attribute name
							currentAttrMap[attrName] = attrValue
						} else {
							if _, ok := currentAttrMap[attrName]; !ok {
								currentAttrMap[attrName] = make(map[string]any)
							}
							if innerMap, ok := currentAttrMap[attrName].(map[string]any); ok {
								// set the current attribute map to the inner attribute map
								currentAttrMap = innerMap
							} else {
								return errors.NewCommonEdgeX(errors.KindContractInvalid,
									fmt.Sprintf("error occurred while converting the nested attribute of '%s' column", headerName), nil)
							}
						}
					}

					// set the attrMap back to the attrMapField
					attrMapField.Set(reflect.ValueOf(attrMap))
				}
			}
		}
	}

	return nil
}
