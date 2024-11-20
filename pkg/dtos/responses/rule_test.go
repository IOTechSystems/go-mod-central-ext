// Copyright (C) 2023-2024 IOTech Ltd

package responses

import (
	"github.com/edgexfoundry/go-mod-central-ext/v4/pkg/common"
	"github.com/edgexfoundry/go-mod-central-ext/v4/pkg/dtos"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMultiRulesResponse(t *testing.T) {
	expectedRules := []dtos.Rule{
		{Name: "rule1", Rule: []byte("rule1")},
		{Name: "rule2", Rule: []byte("rule2")},
	}
	expectedTotalCount := uint32(len(expectedRules))
	actual := NewMultiRulesResponse(common.ExpectedRequestId, common.ExpectedMessage, common.ExpectedStatusCode, expectedTotalCount, expectedRules)
	assert.Equal(t, common.ExpectedRequestId, actual.RequestId)
	assert.Equal(t, common.ExpectedStatusCode, actual.StatusCode)
	assert.Equal(t, common.ExpectedMessage, actual.Message)
	assert.Equal(t, expectedTotalCount, actual.TotalCount)
	assert.Equal(t, expectedRules, actual.Rules)
}
