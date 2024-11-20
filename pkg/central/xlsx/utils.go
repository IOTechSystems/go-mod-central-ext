// Copyright (C) 2023-2024 IOTech Ltd

package xlsx

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"

	"github.com/xuri/excelize/v2"
)

// checkMappingObject checks if the object field from MappingTable is defined in the provided worksheet
// if not found, adds the new object column with defaultValue to the provided sheetName
func checkMappingObject(xlsFile *excelize.File, sheetName string, totalColCount *int, totalRowCount int,
	defaultValue, objectField string, header *[]string) errors.EdgeX {
	if header == nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "header cannot be nil", nil)
	}

	found := false
	for _, headerCell := range *header {
		if headerCell == objectField {
			// check if the mapping object field is defined in the Devices sheet header column
			found = true
			break
		}
	}

	if found || defaultValue == "" {
		return nil
	}

	// if DefaultValue column of the missing field in the MappingTable is not empty
	// Append a new column before colName for this missing field with DefaultValue set on each row
	colName, err := excelize.ColumnNumberToName(*totalColCount + 1)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "failed to covert column number to name", err)
	}

	err = xlsFile.InsertCols(sheetName, colName, 1)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to insert empty column to %s", sheetName), err)
	}

	// create the new column values to append to the worksheet for the missing column
	newColValues := make([]any, totalRowCount)
	for i := range newColValues {
		if i == 0 {
			newColValues[i] = objectField
		} else {
			newColValues[i] = defaultValue
		}
	}

	err = xlsFile.SetSheetCol(sheetName, colName+"1", &newColValues)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to set new column to %s in %s", colName, sheetName), err)
	}

	*header = append(*header, objectField)
	*totalColCount++
	return nil
}

// checkRequiredSheets examines if all the required sheets are defined in the xlsx
func checkRequiredSheets(allSheetNames, requiredSheets []string) errors.EdgeX {
	for _, requiredSheet := range requiredSheets {
		if !slices.Contains(allSheetNames, requiredSheet) {
			return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("%s worksheet not found in the file", requiredSheet), nil)
		}
	}
	return nil
}

// setMapToStructField sets the map value to the given field of the struct pointer
func setMapToStructField(rowElement *reflect.Value, fieldName string, mapValue map[string]any) errors.EdgeX {
	propsField := rowElement.FieldByName(fieldName)
	if propsField.Kind() == reflect.Invalid {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to find %s field in struct", fieldName), nil)
	}
	if propsField.Kind() != reflect.Map {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to set map to non-map '%s' field in struct", fieldName), nil)
	}
	propsField.Set(reflect.ValueOf(mapValue))
	return nil
}

// parseStringToActualType parses the string value to the actual type (int64, float64 or boolean) if the value can be converted
func parseStringToActualType(strValue string) any {
	var convertedValue any

	if intValue, err := strconv.ParseInt(strValue, 10, 64); err == nil {
		convertedValue = intValue
	} else if floatValue, err := strconv.ParseFloat(strValue, 64); err == nil {
		convertedValue = floatValue
	} else if boolValue, err := strconv.ParseBool(strValue); err == nil {
		convertedValue = boolValue
	} else {
		convertedValue = strValue
	}
	return convertedValue
}
