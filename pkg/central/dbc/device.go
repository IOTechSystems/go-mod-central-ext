// Copyright (C) 2026 IOTech Ltd

package dbc

import (
	"strconv"

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
	edgexModels "github.com/edgexfoundry/go-mod-core-contracts/v4/models"
)

// deviceDBC stores the compiled DBC database, conversion arguments, and the converted Device DTOs
type deviceDBC struct {
	compileResult  *CompileResult
	args           map[string]string
	devices        []edgexDtos.Device
	validateErrors map[string]error
}

func newDeviceDBC(data []byte, args map[string]string) (Converter[[]edgexDtos.Device], errors.EdgeX) {
	compileResult, err := Compile("", data)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindServerError, "failed to compile DBC file", err)
	}

	return &deviceDBC{
		compileResult:  compileResult,
		args:           args,
		devices:        make([]edgexDtos.Device, 0),
		validateErrors: make(map[string]error, len(compileResult.Database.Messages)),
	}, nil
}

// ConvertToDTO parses the DBC messages and converts them to Device DTOs
func (dDBC *deviceDBC) ConvertToDTO() errors.EdgeX {
	for _, m := range dDBC.compileResult.Database.Messages {
		deviceDTO := edgexDtos.Device{
			Name:           m.Name,
			Description:    m.Description,
			AdminState:     edgexModels.Unlocked,
			OperatingState: edgexModels.Up,
			ProfileName:    m.Name,
			ServiceName:    dDBC.args[ServiceName],
			Protocols: map[string]edgexDtos.ProtocolProperties{
				Canbus: {
					NetType:  dDBC.args[NetType],
					CommType: dDBC.args[CommType],
					Network:  dDBC.args[Network],
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
		if dDBC.args[NetType] == NetTypeEthernet {
			deviceDTO.Protocols[Canbus][Port] = dDBC.args[Port]
		}

		validateErr := edgexCommon.Validate(deviceDTO)
		if validateErr != nil {
			dDBC.validateErrors[deviceDTO.Name] = validateErr
		} else {
			dDBC.devices = append(dDBC.devices, deviceDTO)
		}
	}
	return nil
}

func (dDBC *deviceDBC) GetDTOs() []edgexDtos.Device {
	return dDBC.devices
}

func (dDBC *deviceDBC) GetValidateErrors() map[string]error {
	return dDBC.validateErrors
}
