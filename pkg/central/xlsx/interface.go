// Copyright (C) 2023-2024 IOTech Ltd

package xlsx

import (
	"io"

	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

type AllowedDTOTypes interface {
	*edgexDtos.DeviceProfile | []*edgexDtos.Device
}

type Converter[T AllowedDTOTypes] interface {
	// ConvertToDTO parses the xlsx file content to DTOs
	ConvertToDTO() errors.EdgeX
	// GetDTOs returns the coverted DTOs
	GetDTOs() T
	// GetValidateErrors returns the deviceName-validationError key-value map while parsing the excel data rows to DTOs
	GetValidateErrors() map[string]error
}

type AllowedDTOConverterTypes interface {
	edgexDtos.DeviceProfile | []edgexDtos.Device
}

type DTOConverter[T AllowedDTOConverterTypes] interface {
	// ConvertToXlsx parses the DTOs to xlsx file content
	ConvertToXlsx() errors.EdgeX
	// Write writes xlsx file content to io.Writer
	Write(io.Writer) errors.EdgeX
	// closeXlsxFile closes the xlsx file reader
	closeXlsxFile() errors.EdgeX
}
