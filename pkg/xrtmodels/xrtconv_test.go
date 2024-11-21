// Copyright (C) 2022-2024 IOTech Ltd

package xrtmodels

import (
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	edgexModels "github.com/edgexfoundry/go-mod-core-contracts/v4/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToXrtDevice(t *testing.T) {
	deviceName := "test-device"
	profileName := "test-profile"
	serviceName := "device-bacnet-ip"
	device := edgexDtos.Device{
		Name: deviceName,
		Protocols: map[string]edgexDtos.ProtocolProperties{
			common.BacnetIP: {
				common.BacnetDeviceInstance: 1234,
			},
		},
		ProfileName:    profileName,
		ServiceName:    serviceName,
		AdminState:     edgexModels.Unlocked,
		OperatingState: edgexModels.Up,
		Properties: map[string]any{
			common.ProtocolName: common.BacnetIP,
		},
	}
	xrtDevice, err := ToXrtDevice(device)
	require.NoError(t, err)

	assert.Equal(t, deviceName, xrtDevice.Name)
	assert.Equal(t, common.BacnetIP, xrtDevice.Properties[common.ProtocolName])
	assert.Equal(t, float64(1234), xrtDevice.Protocols[common.BacnetIP][common.BacnetDeviceInstance])
	assert.Equal(t, profileName, xrtDevice.ProfileName)
	assert.Equal(t, serviceName, xrtDevice.ServiceName)
	assert.Equal(t, edgexModels.Unlocked, xrtDevice.AdminState)
	assert.Equal(t, edgexModels.Up, xrtDevice.OperatingState)
}

func TestToEdgeXProperties(t *testing.T) {
	tests := []struct {
		protocol   string
		properties map[string]interface{}
		expected   map[string]string
	}{
		{
			protocol: common.BacnetIP,
			properties: map[string]interface{}{
				common.BacnetDeviceInstance: float64(1234),
			},
			expected: map[string]string{
				common.BacnetDeviceInstance: "1234",
			},
		},
		{
			protocol: common.BacnetIP,
			properties: map[string]interface{}{
				common.BacnetDeviceInstance: float64(4194302),
			},
			expected: map[string]string{
				common.BacnetDeviceInstance: "4194302",
			},
		},
		{
			protocol: common.Opcua,
			properties: map[string]interface{}{
				common.OpcuaRequestedSessionTimeout:    float64(1200000),
				common.OpcuaBrowseDepth:                float64(0),
				common.OpcuaBrowsePublishInterval:      5.2,
				common.OpcuaConnectionReadingPostDelay: float64(0),
				common.OpcuaIDType:                     float64(1),
			},
			expected: map[string]string{
				common.OpcuaRequestedSessionTimeout:    "1200000",
				common.OpcuaBrowseDepth:                "0",
				common.OpcuaBrowsePublishInterval:      "5.2",
				common.OpcuaConnectionReadingPostDelay: "0",
				common.OpcuaIDType:                     "1",
			},
		},
		{
			protocol: common.ModbusTcp,
			properties: map[string]interface{}{
				common.ModbusUnitID: float64(1),
				common.ModbusPort:   float64(1234),
			},
			expected: map[string]string{
				common.ModbusUnitID: "1",
				common.ModbusPort:   "1234",
			},
		},
		{
			protocol: common.ModbusRtu,
			properties: map[string]interface{}{
				common.ModbusUnitID:   float64(1),
				common.ModbusBaudRate: float64(0),
				common.ModbusDataBits: float64(5),
				common.ModbusStopBits: float64(1),
			},
			expected: map[string]string{
				common.ModbusUnitID:   "1",
				common.ModbusBaudRate: "0",
				common.ModbusDataBits: "5",
				common.ModbusStopBits: "1",
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.protocol, func(t *testing.T) {
			result := toEdgeXProperties(testCase.protocol, testCase.properties)
			assert.EqualValues(t, testCase.expected, result)
		})
	}
}
