// Copyright (C) 2020-2024 IOTech Ltd

package requests

import (
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/requests"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/models"

	"github.com/stretchr/testify/assert"
)

var testProtocolName = "Modbus-IP"
var testProtocols = map[string]dtos.ProtocolProperties{
	"modbus-ip": {
		"Address": "localhost",
		"Port":    "1502",
		"UnitID":  "1",
	},
}

var testAddDevice = requests.AddDeviceRequest{
	BaseRequest: dtoCommon.BaseRequest{
		RequestId:   common.ExampleUUID,
		Versionable: dtoCommon.NewVersionable(),
	},
	Device: dtos.Device{
		Name:           common.TestDeviceName,
		ServiceName:    common.TestServiceName,
		ProfileName:    common.TestProfileName,
		AdminState:     models.Locked,
		OperatingState: models.Up,
		Tags:           map[string]interface{}{"1": common.TestTag1, "2": common.TestTag2},
		Protocols:      testProtocols,
		Properties: map[string]any{
			common.ProtocolName: testProtocolName,
		},
	},
}

func Test_AddDeviceReqToDeviceModels(t *testing.T) {
	requests := []requests.AddDeviceRequest{testAddDevice}
	expectedDeviceModel := []models.Device{
		{
			Name:           common.TestDeviceName,
			ServiceName:    common.TestServiceName,
			ProfileName:    common.TestProfileName,
			AdminState:     models.Locked,
			OperatingState: models.Up,
			Tags:           map[string]any{"1": common.TestTag1, "2": common.TestTag2},
			Protocols: map[string]models.ProtocolProperties{
				"modbus-ip": {
					"Address": "localhost",
					"Port":    "1502",
					"UnitID":  "1",
				},
			},
			Properties: map[string]any{
				common.ProtocolName: "modbus-ip",
			},
			AutoEvents: []models.AutoEvent{},
		},
	}
	resultModels := AddDeviceReqToDeviceModels(requests)
	assert.Equal(t, expectedDeviceModel, resultModels, "AddDeviceReqToDeviceModels did not result in expected Device model.")
}
