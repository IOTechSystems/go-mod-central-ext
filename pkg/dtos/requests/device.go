// Copyright (C) 2020-2024 IOTech Ltd

package requests

import (
	"github.com/edgexfoundry/go-mod-central-ext/v4/pkg/dtos"
	dtoRequests "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/requests"
	edgexModels "github.com/edgexfoundry/go-mod-core-contracts/v4/models"
)

// AddDeviceReqToDeviceModels transforms the AddDeviceRequest DTO array to the Device model array
func AddDeviceReqToDeviceModels(addRequests []dtoRequests.AddDeviceRequest) (Devices []edgexModels.Device) {
	for _, req := range addRequests {
		d := dtos.ToDeviceModel(req.Device)
		Devices = append(Devices, d)
	}
	return Devices
}
