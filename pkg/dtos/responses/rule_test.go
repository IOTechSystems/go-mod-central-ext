// Copyright (C) 2023-2025 IOTech Ltd

package responses

import (
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMultiRulesResponse(t *testing.T) {
	expectedRules := []dtos.Rule{
		{Name: "rule1", Rule: []byte("rule1")},
		{Name: "rule2", Rule: []byte("rule2")},
	}
	expectedTotalCount := int64(len(expectedRules))
	actual := NewMultiRulesResponse(common.ExpectedRequestId, common.ExpectedMessage, common.ExpectedStatusCode, expectedTotalCount, expectedRules)
	assert.Equal(t, common.ExpectedRequestId, actual.RequestId)
	assert.Equal(t, common.ExpectedStatusCode, actual.StatusCode)
	assert.Equal(t, common.ExpectedMessage, actual.Message)
	assert.Equal(t, expectedTotalCount, actual.TotalCount)
	assert.Equal(t, expectedRules, actual.Rules)
}
