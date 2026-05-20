// Copyright (C) 2023-2026 IOTech Ltd

package xlsx

import (
	"io"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/xrtmodels"
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

// DeviceScheduleReader exposes Schedules parsed from a device xlsx, looked up by device name.
// A Converter[[]*edgexDtos.Device] returned from ConvertDeviceXlsx implements this interface;
// callers should type-assert to access it.
//
// Example:
//
//	deviceXlsx, edgexErr := xlsx.ConvertDeviceXlsx(f) // ConvertDeviceXlsx already runs ConvertToDTO internally
//	if edgexErr != nil { /* handle */ }
//
//	scheduleReader := deviceXlsx.(xlsx.DeviceScheduleReader)
//	for _, device := range deviceXlsx.GetDTOs() {
//	    schedules := scheduleReader.GetSchedulesByDeviceName(device.Name)
//	    // send schedules to XRT, etc.
//	}
type DeviceScheduleReader interface {
	GetSchedulesByDeviceName(name string) []xrtmodels.Schedule
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
