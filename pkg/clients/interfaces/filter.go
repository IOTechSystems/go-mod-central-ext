// Copyright (C) 2026 IOTech Ltd

package interfaces

import (
	"context"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/requests"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/responses"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

// FilterClient defines the interface for interactions with the filter API endpoints on the IOTech core-data service.
type FilterClient interface {
	// AllFilters returns all filters.
	// The result can be limited in a certain range by specifying the offset and limit parameters.
	// offset: The number of items to skip before starting to collect the result set. Default is 0.
	// limit: The number of items to return. Specifying -1 will return all remaining items after offset. The maximum will be the MaxResultCount as defined in the configuration of service. Default is 20.
	AllFilters(ctx context.Context, offset int, limit int) (responses.MultiFiltersResponse, errors.EdgeX)
	// FilterById returns a filter by id.
	FilterById(ctx context.Context, id string) (responses.FilterResponse, errors.EdgeX)
	// FiltersByDeviceName returns filters by device name.
	// The result can be limited in a certain range by specifying the offset and limit parameters.
	// offset: The number of items to skip before starting to collect the result set. Default is 0.
	// limit: The number of items to return. Specifying -1 will return all remaining items after offset. The maximum will be the MaxResultCount as defined in the configuration of service. Default is 20.
	FiltersByDeviceName(ctx context.Context, deviceName string, offset int, limit int) (responses.MultiFiltersResponse, errors.EdgeX)
	// Add adds new filters.
	Add(ctx context.Context, reqs []requests.FilterRequest) ([]common.BaseWithIdResponse, errors.EdgeX)
	// PutFilters replaces all filters with the given filters. This is destructive:
	// every existing filter is deleted first, then the supplied filters are added,
	// so any filter not included in reqs is removed.
	PutFilters(ctx context.Context, reqs []requests.FilterRequest) ([]common.BaseWithIdResponse, errors.EdgeX)
	// Update updates a filter.
	Update(ctx context.Context, req requests.FilterRequest) (common.BaseResponse, errors.EdgeX)
	// DeleteFilterById deletes a filter by id.
	DeleteFilterById(ctx context.Context, id string) (common.BaseResponse, errors.EdgeX)
	// DeleteFiltersByDeviceName deletes filters by device name.
	DeleteFiltersByDeviceName(ctx context.Context, deviceName string) (common.BaseResponse, errors.EdgeX)
	// DeleteFilters deletes all filters. This is destructive and removes every filter.
	DeleteFilters(ctx context.Context) (common.BaseResponse, errors.EdgeX)
}
