// Copyright (C) 2021-2024 IOTech Ltd

package interfaces

import (
	"context"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/requests"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

// SystemManagementClient defines the interface for interactions with the API endpoint on the IOTech sys-mgmt service.
type SystemManagementClient interface {
	// GetHealth obtain health information of services via registry by their name
	GetHealth(ctx context.Context, services []string) ([]common.BaseWithServiceNameResponse, errors.EdgeX)
	// GetConfig obtain configuration from services by their name
	GetConfig(ctx context.Context, services []string) ([]common.BaseWithConfigResponse, errors.EdgeX)
	// DoOperation issue a start, stop, restart action to the targeted services
	DoOperation(ctx context.Context, reqs []requests.OperationRequest) ([]common.BaseResponse, errors.EdgeX)
}
