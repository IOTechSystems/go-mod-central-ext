// Copyright (C) 2025 IOTech Ltd

package responses

import (
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"

	"github.com/stretchr/testify/require"
)

var (
	testDeviceName       = "testDevice"
	testSourceName       = "testSource"
	testDuration         = "10h"
	testRetentionPolicy1 = dtos.RetentionPolicy{
		DeviceName: testDeviceName,
		Duration:   testDuration,
	}
	testRetentionPolicy2 = dtos.RetentionPolicy{
		SourceName: testSourceName,
		Duration:   testDuration,
	}
)

func TestNewRetentionPolicyResponse(t *testing.T) {
	expectedPolicy := testRetentionPolicy1
	actual := NewRetentionPolicyResponse(common.ExpectedRequestId, common.ExpectedMessage, common.ExpectedStatusCode, expectedPolicy)

	require.Equal(t, common.ExpectedRequestId, actual.RequestId)
	require.Equal(t, common.ExpectedStatusCode, actual.StatusCode)
	require.Equal(t, common.ExpectedMessage, actual.Message)
	require.Equal(t, expectedPolicy, actual.RetentionPolicy)
}

func TestNewMultiRetentionPolicyResponse(t *testing.T) {
	expectedPolicies := []dtos.RetentionPolicy{testRetentionPolicy1, testRetentionPolicy2}
	expectedTotalCount := uint32(2)
	actual := NewMultiRetentionPolicyResponse(common.ExpectedRequestId, common.ExpectedMessage, common.ExpectedStatusCode, expectedTotalCount, expectedPolicies)

	require.Equal(t, common.ExpectedRequestId, actual.RequestId)
	require.Equal(t, common.ExpectedStatusCode, actual.StatusCode)
	require.Equal(t, common.ExpectedMessage, actual.Message)
	require.Equal(t, expectedTotalCount, actual.TotalCount)
	require.Equal(t, expectedPolicies, actual.RetentionPolicies)
}
