// Copyright (C) 2026 IOTech Ltd

package responses

import (
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
)

type TimeSeriesResponse map[string]any

func NewTimeSeriesResponse(tsResources dtos.TimeSeriesResourceMap) TimeSeriesResponse {
	resp := TimeSeriesResponse{
		"apiVersion": edgexCommon.NewVersionable().ApiVersion,
	}
	for resourceName, tsResource := range tsResources {
		resp[resourceName] = tsResource
	}
	return resp
}
