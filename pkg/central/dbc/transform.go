// Copyright (C) 2023-2024 IOTech Ltd

package dbc

import (
	"fmt"
	"math"
	"strconv"

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	edgexModels "github.com/edgexfoundry/go-mod-core-contracts/v4/models"

	"go.einride.tech/can/pkg/descriptor"
)

func valueType(s *descriptor.Signal) string {
	_, offsetFrac := math.Modf(s.Offset)
	_, scaleFrac := math.Modf(s.Scale)
	if offsetFrac != 0 || scaleFrac != 0 {
		return edgexCommon.ValueTypeFloat64
	}
	if s.IsSigned {
		return edgexCommon.ValueTypeInt64
	} else {
		return edgexCommon.ValueTypeUint64
	}
}

func ConvertDBCtoProfile(data []byte) (profileDTOs []edgexDtos.DeviceProfile, err error, validateErrors map[string]error) {
	compileResult, err := Compile("", data)
	if err != nil {
		return
	}

	validateErrors = make(map[string]error, len(compileResult.Database.Messages))

	for _, m := range compileResult.Database.Messages {
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
			validateErrors[profileDto.Name] = validateErr
		} else {
			profileDTOs = append(profileDTOs, profileDto)
		}
	}
	return
}

func ConvertDBCtoDevice(data []byte, args map[string]string) (deviceDTOs []edgexDtos.Device, err error, validateErrors map[string]error) {
	compileResult, err := Compile("", data)
	if err != nil {
		return
	}

	validateErrors = make(map[string]error, len(compileResult.Database.Messages))

	for _, m := range compileResult.Database.Messages {
		deviceDTO := edgexDtos.Device{
			Name:           m.Name,
			Description:    m.Description,
			AdminState:     edgexModels.Unlocked,
			OperatingState: edgexModels.Up,
			ProfileName:    m.Name,
			ServiceName:    args[ServiceName],
			Protocols: map[string]edgexDtos.ProtocolProperties{
				Canbus: {
					NetType:  args[NetType],
					CommType: args[CommType],
					Network:  args[Network],
					Standard: J1939,
					ID:       getOriginalCanId(m.ID),
					DataSize: strconv.Itoa(int(m.Length)),
					Sender:   m.SenderNode,
				},
			},
			Tags: map[string]any{
				PGN: getPGN(m.ID),
			},
		}
		if args[NetType] == NetTypeEthernet {
			deviceDTO.Protocols[Canbus][Port] = args[Port]
		}

		validateErr := edgexCommon.Validate(deviceDTO)
		if validateErr != nil {
			validateErrors[deviceDTO.Name] = validateErr
		} else {
			deviceDTOs = append(deviceDTOs, deviceDTO)
		}
	}
	return
}

func getOriginalCanId(canID uint32) string {
	id := canID | messageIDExtendedFlag
	return strconv.FormatUint(uint64(id), 10)
}

func getPGN(canID uint32) string {
	// J1939 PGN bit start from 9, length is 18
	pgn := (canID >> j1939PGNOffset) & j1939PGNMask
	return fmt.Sprintf("%X", pgn)
}
