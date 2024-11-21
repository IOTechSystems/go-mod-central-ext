// Copyright (C) 2024 IOTech Ltd

package xlsx

import (
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToXrtProperties(t *testing.T) {
	tests := []struct {
		protocol   string
		properties map[string]interface{}
		expected   map[string]interface{}
	}{
		{
			protocol: common.BacnetIP,
			properties: map[string]interface{}{
				common.BacnetDeviceInstance: "1234",
			},
			expected: map[string]interface{}{
				common.BacnetDeviceInstance: 1234,
			},
		},
		{
			protocol: common.BacnetIP,
			properties: map[string]interface{}{
				common.BacnetDeviceInstance: "4194302",
			},
			expected: map[string]interface{}{
				common.BacnetDeviceInstance: 4194302,
			},
		},
		{
			protocol: common.BacnetIP,
			properties: map[string]interface{}{
				common.BacnetAddress: "192.168.60.123",
				common.BacnetPort:    "47809",
			},
			expected: map[string]interface{}{
				common.BacnetAddress: "192.168.60.123",
				common.BacnetPort:    47809,
			},
		},
		{
			protocol: common.Opcua,
			properties: map[string]interface{}{
				common.OpcuaRequestedSessionTimeout:    "1200000",
				common.OpcuaBrowseDepth:                "0",
				common.OpcuaBrowsePublishInterval:      "5.2",
				common.OpcuaConnectionReadingPostDelay: "0",
				common.OpcuaIDType:                     "1",
				common.OpcuaReadBatchSize:              "15",
				common.OpcuaWriteBatchSize:             "10",
				common.OpcuaNodesPerBrowse:             "5",
				common.OpcuaSessionKeepAliveInterval:   "1000",
			},
			expected: map[string]interface{}{
				common.OpcuaRequestedSessionTimeout:    1200000,
				common.OpcuaBrowseDepth:                0,
				common.OpcuaBrowsePublishInterval:      5.2,
				common.OpcuaConnectionReadingPostDelay: 0,
				common.OpcuaIDType:                     "1",
				common.OpcuaReadBatchSize:              15,
				common.OpcuaWriteBatchSize:             10,
				common.OpcuaNodesPerBrowse:             5,
				common.OpcuaSessionKeepAliveInterval:   1000.0,
			},
		},
		{
			protocol: common.ModbusTcp,
			properties: map[string]interface{}{
				common.ModbusUnitID:                    "1",
				common.ModbusPort:                      "1234",
				common.ModbusReadMaxHoldingRegisters:   "125",
				common.ModbusReadMaxInputRegisters:     "125",
				common.ModbusReadMaxBitsCoils:          "2000",
				common.ModbusReadMaxBitsDiscreteInputs: "2000",
				common.ModbusWriteMaxHoldingRegisters:  "123",
				common.ModbusWriteMaxBitsCoils:         "1968",
			},
			expected: map[string]interface{}{
				common.ModbusUnitID:                    1,
				common.ModbusPort:                      1234,
				common.ModbusReadMaxHoldingRegisters:   125,
				common.ModbusReadMaxInputRegisters:     125,
				common.ModbusReadMaxBitsCoils:          2000,
				common.ModbusReadMaxBitsDiscreteInputs: 2000,
				common.ModbusWriteMaxHoldingRegisters:  123,
				common.ModbusWriteMaxBitsCoils:         1968,
			},
		},
		{
			protocol: common.ModbusRtu,
			properties: map[string]interface{}{
				common.ModbusUnitID:                    "1",
				common.ModbusBaudRate:                  "0",
				common.ModbusDataBits:                  "5",
				common.ModbusStopBits:                  "1",
				common.ModbusReadMaxHoldingRegisters:   "125",
				common.ModbusReadMaxInputRegisters:     "125",
				common.ModbusReadMaxBitsCoils:          "2000",
				common.ModbusReadMaxBitsDiscreteInputs: "2000",
				common.ModbusWriteMaxHoldingRegisters:  "123",
				common.ModbusWriteMaxBitsCoils:         "1968",
			},
			expected: map[string]interface{}{
				common.ModbusUnitID:                    1,
				common.ModbusBaudRate:                  0,
				common.ModbusDataBits:                  5,
				common.ModbusStopBits:                  1,
				common.ModbusReadMaxHoldingRegisters:   125,
				common.ModbusReadMaxInputRegisters:     125,
				common.ModbusReadMaxBitsCoils:          2000,
				common.ModbusReadMaxBitsDiscreteInputs: 2000,
				common.ModbusWriteMaxHoldingRegisters:  123,
				common.ModbusWriteMaxBitsCoils:         1968,
			},
		},
		{
			protocol: common.EtherNetIPExplicitConnected,
			properties: map[string]interface{}{
				common.EtherNetIPDeviceResource: "VendorID",
				common.EtherNetIPRPI:            "3000",
				common.EtherNetIPSaveValue:      "true",
			},
			expected: map[string]interface{}{
				common.EtherNetIPDeviceResource: "VendorID",
				common.EtherNetIPRPI:            3000,
				common.EtherNetIPSaveValue:      true,
			},
		},
		{
			protocol: common.EtherNetIPO2T,
			properties: map[string]interface{}{
				common.EtherNetIPConnectionType: "p2p",
				common.EtherNetIPRPI:            "10",
				common.EtherNetIPPriority:       "low",
				common.EtherNetIPOwnership:      "exclusive",
			},
			expected: map[string]interface{}{
				common.EtherNetIPConnectionType: "p2p",
				common.EtherNetIPRPI:            10,
				common.EtherNetIPPriority:       "low",
				common.EtherNetIPOwnership:      "exclusive",
			},
		},
		{
			protocol: common.EtherNetIPT2O,
			properties: map[string]interface{}{
				common.EtherNetIPConnectionType: "p2p",
				common.EtherNetIPRPI:            "10",
				common.EtherNetIPPriority:       "low",
				common.EtherNetIPOwnership:      "exclusive",
			},
			expected: map[string]interface{}{
				common.EtherNetIPConnectionType: "p2p",
				common.EtherNetIPRPI:            10,
				common.EtherNetIPPriority:       "low",
				common.EtherNetIPOwnership:      "exclusive",
			},
		},
		{
			protocol: common.EtherNetIPKey,
			properties: map[string]interface{}{
				common.EtherNetIPMethod:        "exact",
				common.EtherNetIPVendorID:      "10",
				common.EtherNetIPDeviceType:    "72",
				common.EtherNetIPProductCode:   "50",
				common.EtherNetIPMajorRevision: "12",
				common.EtherNetIPMinorRevision: "2",
			},
			expected: map[string]interface{}{
				common.EtherNetIPMethod:        "exact",
				common.EtherNetIPVendorID:      10,
				common.EtherNetIPDeviceType:    72,
				common.EtherNetIPProductCode:   50,
				common.EtherNetIPMajorRevision: 12,
				common.EtherNetIPMinorRevision: 2,
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.protocol, func(t *testing.T) {
			err := toXrtProperties(testCase.protocol, testCase.properties)
			require.NoError(t, err)
			assert.EqualValues(t, testCase.expected, testCase.properties)
		})
	}
}

func TestToXrType_Invalid(t *testing.T) {
	tests := []struct {
		protocol   string
		properties map[string]interface{}
	}{
		{
			protocol: common.ModbusTcp,
			properties: map[string]interface{}{
				common.ModbusPort: "test",
			},
		},
		{
			protocol: common.Opcua,
			properties: map[string]interface{}{
				common.OpcuaBrowsePublishInterval: "test",
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.protocol, func(t *testing.T) {
			err := toXrtProperties(testCase.protocol, testCase.properties)
			require.Error(t, err)
		})
	}
}
