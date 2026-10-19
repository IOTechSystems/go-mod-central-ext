// Copyright (C) 2026 IOTech Ltd

package requests

import (
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
)

// UpdateTagsRequest represents the request DTO for invoking PATCH /api/v3/deviceprofile/name/:name/resourcetag/result
type UpdateTagsRequest struct {
	dtoCommon.BaseRequest `json:",inline"`
	UpdatedTags           []dtos.UpdatePredictedResult `json:"updatedTags"`
}

type ResourceTagJobRequest struct {
	dtoCommon.BaseRequest `json:",inline"`
	ProfileNames          []string `json:"profileNames" validate:"required"`
}
