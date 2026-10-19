// Copyright (C) 2026 IOTech Ltd

package responses

import (
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
)

type ResourceTagJobResponse struct {
	dtoCommon.BaseResponse `json:",inline"`
	ResourceTagJob         dtos.ResourceTagBasicInfo `json:"resourceTagJob"`
}

func NewResourceTagJobBasicInfoResponse(requestId string, message string, statusCode int, resourceTagJob dtos.ResourceTagBasicInfo) ResourceTagJobResponse {
	return ResourceTagJobResponse{
		BaseResponse:   dtoCommon.NewBaseResponse(requestId, message, statusCode),
		ResourceTagJob: resourceTagJob,
	}
}

type ResourceTagJobResultResponse struct {
	dtoCommon.BaseResponse `json:",inline"`
	Ontology               string                          `json:"ontology"`
	ResourceTagJobResult   dtos.ResourceTagResultMapFormat `json:"resourceTagJobResult"`
}

func NewResourceTagJobResultResponse(requestId string, message string, statusCode int, ontology string, resourceTagJobResult dtos.ResourceTagResultMapFormat) ResourceTagJobResultResponse {
	return ResourceTagJobResultResponse{
		BaseResponse:         dtoCommon.NewBaseResponse(requestId, message, statusCode),
		Ontology:             ontology,
		ResourceTagJobResult: resourceTagJobResult,
	}
}
