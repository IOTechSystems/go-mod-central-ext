// Copyright (C) 2026 IOTech Ltd

package dbc

import (
	"strconv"

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

// deviceProfileDBC stores the compiled DBC database and the converted DeviceProfile DTOs
type deviceProfileDBC struct {
	compileResult  *CompileResult
	deviceProfiles []edgexDtos.DeviceProfile
	validateErrors map[string]error
}

func newDeviceProfileDBC(data []byte) (Converter[[]edgexDtos.DeviceProfile], errors.EdgeX) {
	compileResult, err := Compile("", data)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindServerError, "failed to compile DBC file", err)
	}

	return &deviceProfileDBC{
		compileResult:  compileResult,
		deviceProfiles: make([]edgexDtos.DeviceProfile, 0),
		validateErrors: make(map[string]error, len(compileResult.Database.Messages)),
	}, nil
}

// ConvertToDTO parses the DBC messages and converts them to DeviceProfile DTOs
func (dpDBC *deviceProfileDBC) ConvertToDTO() errors.EdgeX {
	for _, m := range dpDBC.compileResult.Database.Messages {
		var deviceResources []edgexDtos.DeviceResource
		var deviceCommands []edgexDtos.DeviceCommand
		var profileDto edgexDtos.DeviceProfile

		for _, s := range m.Signals {
			deviceResource := edgexDtos.DeviceResource{
				Name:        s.Name,
				Description: s.Description,
				Properties: edgexDtos.ResourceProperties{
					ValueType:    valueType(s),
					ReadWrite:    edgexCommon.ReadWrite_R,
					Units:        s.Unit,
					Minimum:      &s.Min,
					Maximum:      &s.Max,
					Scale:        &s.Scale,
					Offset:       &s.Offset,
					DefaultValue: strconv.FormatInt(int64(s.DefaultValue), 10),
				},
				Attributes: map[string]interface{}{
					BitStart:      s.Start,
					BitLen:        s.Length,
					LittleEndian:  !s.IsBigEndian,
					ReceiverNames: s.ReceiverNodes,
					MuxSignal:     s.IsMultiplexer,
					IsSigned:      s.IsSigned,
				},
			}
			if s.IsMultiplexed {
				deviceResource.Attributes[MuxNum] = s.MultiplexerValue
			}
			if len(s.ValueDescriptions) > 0 {
				var deviceCommand edgexDtos.DeviceCommand
				deviceCommand.Name = s.Name
				deviceCommand.ReadWrite = edgexCommon.ReadWrite_R
				mappings := make(map[string]string, len(s.ValueDescriptions))
				for _, valueDescription := range s.ValueDescriptions {
					mappings[strconv.FormatInt(valueDescription.Value, 10)] = valueDescription.Description
				}
				deviceCommand.ResourceOperations = []edgexDtos.ResourceOperation{
					{
						DeviceResource: s.Name,
						DefaultValue:   strconv.FormatInt(int64(s.DefaultValue), 10),
						Mappings:       mappings,
					},
				}
				deviceCommands = append(deviceCommands, deviceCommand)
			}
			deviceResources = append(deviceResources, deviceResource)
		}

		profileDto.Name = m.Name
		profileDto.Description = m.Description
		profileDto.DeviceResources = deviceResources
		profileDto.DeviceCommands = deviceCommands

		if validateErr := edgexCommon.Validate(profileDto); validateErr != nil {
			dpDBC.validateErrors[profileDto.Name] = validateErr
		} else {
			dpDBC.deviceProfiles = append(dpDBC.deviceProfiles, profileDto)
		}
	}
	return nil
}

func (dpDBC *deviceProfileDBC) GetDTOs() []edgexDtos.DeviceProfile {
	return dpDBC.deviceProfiles
}

func (dpDBC *deviceProfileDBC) GetValidateErrors() map[string]error {
	return dpDBC.validateErrors
}
