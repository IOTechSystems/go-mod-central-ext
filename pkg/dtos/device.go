// Copyright (C) 2020-2025 IOTech Ltd

package dtos

import (
	"strings"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	edgexModels "github.com/edgexfoundry/go-mod-core-contracts/v4/models"

	"github.com/spf13/cast"
)

// ToDeviceModel transforms the Device DTO to the Device Model
func ToDeviceModel(dto edgexDtos.Device) edgexModels.Device {
	d := edgexDtos.ToDeviceModel(dto)
	// Central
	if protocolName, ok := dto.Properties[common.ProtocolName]; ok {
		d.Properties[common.ProtocolName] = strings.ToLower(cast.ToString(protocolName))
	}
	return d
}
