// Copyright (C) 2026 IOTech Ltd

package dtos

import (
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
)

// ResourceTagBasicInfo is the status of a central-ui resource tag job, without its result
type ResourceTagBasicInfo struct {
	ProfileName   string `json:"profileName" validate:"required,edgex-dto-none-empty-string"`
	Status        string `json:"status,omitempty"`
	Ontology      string `json:"ontology,omitempty"`
	FailedMessage string `json:"failedMessage,omitempty"`
	ResourceCountInfo
	edgexDtos.DBTimestamp
	Synced int64 `json:"synced,omitempty"`
}

type ResourceCountInfo struct {
	TotalResourceCount     int64 `json:"totalResourceCount,omitempty"`
	ProcessedResourceCount int64 `json:"processedResourceCount,omitempty"`
}

// ResourceTagResultMapFormat is a resource tag job result, one map per resource. Its keys are the common.ResourceTagHeader*
// constants plus the job's tag metadata headers, which vary by ontology; the OptionalFields value is a nested map.
type ResourceTagResultMapFormat []map[string]any

// UpdatePredictedResult is the edit of one resource in a resource tag job result
type UpdatePredictedResult struct {
	ResourceName string `json:"resourceName" validate:"required,edgex-dto-none-empty-string"`
	Tag          string `json:"tag"`
	Skipped      bool   `json:"skipped"`
	AssetType    string `json:"assetType"`
}
