// Copyright (C) 2023-2024 IOTech Ltd

package xlsx

import (
	"fmt"
	"io"

	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

type mappingField struct {
	defaultValue string // the default value defined in the MappingTable sheet
	path         string // the path value defined in the MappingTable sheet
}

func ConvertDeviceXlsx(file io.Reader) (Converter[[]*edgexDtos.Device], errors.EdgeX) {
	deviceX, err := newDeviceXlsx(file)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindServerError, "failed to create deviceXlsx instance", err)
	}

	err = deviceX.ConvertToDTO()
	if err != nil {
		return nil, errors.NewCommonEdgeXWrapper(err)
	}

	return deviceX, nil
}

func ConvertDeviceProfileXlsx(file io.Reader) (Converter[*edgexDtos.DeviceProfile], error) {
	deviceProfileX, err := newDeviceProfileXlsx(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create deviceProfileXlsx instance: %w", err)
	}

	err = deviceProfileX.ConvertToDTO()
	if err != nil {
		return nil, err
	}

	return deviceProfileX, nil
}

// ConvertToXlsx converts the DTOs to the xlsx file and writes to io.Writer
func ConvertToXlsx[T AllowedDTOConverterTypes](fileReader io.Reader, w io.Writer, convertData T) errors.EdgeX {
	xlsxWriter, edgexErr := newXlsxWriter(convertData, fileReader)
	if edgexErr != nil {
		return edgexErr
	}
	defer func() { _ = xlsxWriter.closeXlsxFile() }()

	edgexErr = xlsxWriter.ConvertToXlsx()
	if edgexErr != nil {
		return edgexErr
	}

	edgexErr = xlsxWriter.Write(w)
	if edgexErr != nil {
		return edgexErr
	}
	return nil
}
