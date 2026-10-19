// Copyright (C) 2026 IOTech Ltd

package dtos_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/requests"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos/responses"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/common"
)

func TestResourceTagJSON(t *testing.T) {
	basicInfo := dtos.ResourceTagBasicInfo{
		ProfileName:   "p1",
		Status:        "COMPLETED",
		Ontology:      "haystack",
		FailedMessage: "none",
		ResourceCountInfo: dtos.ResourceCountInfo{
			TotalResourceCount:     3,
			ProcessedResourceCount: 2,
		},
		DBTimestamp: edgexDtos.DBTimestamp{Created: 1, Modified: 2},
		Synced:      3,
	}
	update := dtos.UpdatePredictedResult{ResourceName: "r1", Tag: "t1", Skipped: true, AssetType: "AHU"}
	result := dtos.ResourceTagResultMapFormat{{"ResourceName": "r1", "Tag": "t1"}}

	tests := []struct {
		name     string
		value    any
		expected string
	}{
		{"ResourceCountInfo", basicInfo.ResourceCountInfo,
			`{"totalResourceCount":3,"processedResourceCount":2}`},
		{"ResourceTagBasicInfo", basicInfo,
			`{"profileName":"p1","status":"COMPLETED","ontology":"haystack","failedMessage":"none","totalResourceCount":3,"processedResourceCount":2,"created":1,"modified":2,"synced":3}`},
		{"ResourceTagBasicInfo omitempty", dtos.ResourceTagBasicInfo{ProfileName: "p1"},
			`{"profileName":"p1"}`},
		{"UpdatePredictedResult", update,
			`{"resourceName":"r1","tag":"t1","skipped":true,"assetType":"AHU"}`},
		{"ResourceTagResultMapFormat", result,
			`[{"ResourceName":"r1","Tag":"t1"}]`},
		{"UpdateTagsRequest", requests.UpdateTagsRequest{BaseRequest: dtoCommon.BaseRequest{Versionable: dtoCommon.NewVersionable()}, UpdatedTags: []dtos.UpdatePredictedResult{update}},
			`{"apiVersion":"v3","requestId":"","updatedTags":[{"resourceName":"r1","tag":"t1","skipped":true,"assetType":"AHU"}]}`},
		{"ResourceTagJobRequest", requests.ResourceTagJobRequest{BaseRequest: dtoCommon.BaseRequest{Versionable: dtoCommon.NewVersionable()}, ProfileNames: []string{"p1", "p2"}},
			`{"apiVersion":"v3","requestId":"","profileNames":["p1","p2"]}`},
		{"ResourceTagJobResponse", responses.NewResourceTagJobBasicInfoResponse("", "", 200, dtos.ResourceTagBasicInfo{ProfileName: "p1", Status: "RUNNING"}),
			`{"apiVersion":"v3","statusCode":200,"resourceTagJob":{"profileName":"p1","status":"RUNNING"}}`},
		{"ResourceTagJobResultResponse", responses.NewResourceTagJobResultResponse("", "", 200, "haystack", result),
			`{"apiVersion":"v3","statusCode":200,"ontology":"haystack","resourceTagJobResult":[{"ResourceName":"r1","Tag":"t1"}]}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.value)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(b))
		})
	}
}
