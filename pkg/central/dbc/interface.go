// Copyright (C) 2026 IOTech Ltd

package dbc

import (
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

type AllowedDTOTypes interface {
	[]edgexDtos.DeviceProfile | []edgexDtos.Device
}

type Converter[T AllowedDTOTypes] interface {
	// ConvertToDTO parses the DBC file content to DTOs
	ConvertToDTO() errors.EdgeX
	// GetDTOs returns the converted DTOs
	GetDTOs() T
	// GetValidateErrors returns the deviceName-validationError key-value map while parsing the DBC data to DTOs
	GetValidateErrors() map[string]error
}
