// Copyright (C) 2021-2024 IOTech Ltd

package xrtmodels

import (
	"encoding/json"
	"strings"

	"github.com/edgexfoundry/go-mod-central-ext/v4/pkg/common"
	"github.com/edgexfoundry/go-mod-central-ext/v4/pkg/v2models"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
	edgexModels "github.com/edgexfoundry/go-mod-core-contracts/v4/models"
)

type DeviceInfo struct {
	edgexDtos.Device
}

// ToEdgeXV2Device converts the XRT model to EdgeX v2 model
func ToEdgeXV2Device(device DeviceInfo, serviceName string) v2models.Device {
	protocols := make(map[string]v2models.ProtocolProperties)
	protocolName := ""
	for protocol, protocolProperties := range device.Protocols {
		protocols[protocol] = toEdgeXProperties(protocol, protocolProperties)
		protocolName = strings.ToLower(protocol)
	}
	return v2models.Device{
		Name:           device.Name,
		Description:    "",
		AdminState:     edgexModels.Unlocked,
		OperatingState: edgexModels.Up,
		ProtocolName:   protocolName,
		Protocols:      protocols,
		Labels:         nil,
		Location:       nil,
		ServiceName:    serviceName,
		ProfileName:    device.ProfileName,
		AutoEvents:     nil,
		Properties:     device.Properties,
	}
}

// ToEdgeXV3Device converts the XRT model to EdgeX v3 model
func ToEdgeXV3Device(device DeviceInfo, serviceName string) edgexDtos.Device {
	if device.Properties == nil {
		device.Properties = make(map[string]any)
	}
	for protocol := range device.Protocols {
		device.Properties[common.ProtocolName] = strings.ToLower(protocol)
	}
	return edgexDtos.Device{
		Name:           device.Name,
		Description:    "",
		AdminState:     edgexModels.Unlocked,
		OperatingState: edgexModels.Up,
		Protocols:      device.Protocols,
		Labels:         nil,
		Location:       nil,
		ServiceName:    serviceName,
		ProfileName:    device.ProfileName,
		AutoEvents:     nil,
		Properties:     device.Properties,
	}
}

// ToXrtDevice converts the EdgeX model to XRT model
func ToXrtDevice(device edgexDtos.Device) (deviceInfo DeviceInfo, edgexErr errors.EdgeX) {
	deviceData, err := json.Marshal(device)
	if err != nil {
		return deviceInfo, errors.NewCommonEdgeXWrapper(err)
	}
	err = json.Unmarshal(deviceData, &deviceInfo)
	if err != nil {
		return deviceInfo, errors.NewCommonEdgeXWrapper(err)
	}

	// Process the specified protocol for XRT
	for protocol := range deviceInfo.Protocols {
		switch protocol {
		case common.EtherNetIP:
			processEtherNetIP(deviceInfo.Protocols)
		}
	}

	return deviceInfo, nil
}

func processEtherNetIP(protocolProperties map[string]edgexDtos.ProtocolProperties) {
	// Combine ExplicitConnected, O2T, T2O and Key into EtherNet-IP
	if v, ok := protocolProperties[common.EtherNetIP]; ok {
		protocolProperties[common.EtherNetIPXRT] = v
		delete(protocolProperties, common.EtherNetIP)
	}
	if v, ok := protocolProperties[common.EtherNetIPExplicitConnected]; ok {
		protocolProperties[common.EtherNetIPXRT][common.EtherNetIPExplicitConnected] = v
		delete(protocolProperties, common.EtherNetIPExplicitConnected)
	}
	if v, ok := protocolProperties[common.EtherNetIPO2T]; ok {
		protocolProperties[common.EtherNetIPXRT][common.EtherNetIPO2T] = v
		delete(protocolProperties, common.EtherNetIPO2T)
	}
	if v, ok := protocolProperties[common.EtherNetIPT2O]; ok {
		protocolProperties[common.EtherNetIPXRT][common.EtherNetIPT2O] = v
		delete(protocolProperties, common.EtherNetIPT2O)
	}
	if v, ok := protocolProperties[common.EtherNetIPKey]; ok {
		protocolProperties[common.EtherNetIPXRT][common.EtherNetIPKey] = v
		delete(protocolProperties, common.EtherNetIPKey)
	}
}
