// Copyright (C) 2025 IOTech Ltd

package dtos

import (
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/models"

	"github.com/stretchr/testify/require"
)

var (
	testSourceName       = "testSource"
	testDuration         = "10h"
	testId               = "c3a535a6-cea4-43c5-aeda-14fcb1243e9d"
	testRetentionPolicy1 = RetentionPolicy{
		Id:         testId,
		DeviceName: testDeviceName,
		Duration:   testDuration,
	}
	testRetentionPolicy2 = RetentionPolicy{
		Id:         testId,
		SourceName: testSourceName,
		Duration:   testDuration,
	}
	testRetentionPolicyModel1 = models.RetentionPolicy{
		Id:         testId,
		DeviceName: testDeviceName,
		Duration:   testDuration,
	}
	testRetentionPolicyModel2 = models.RetentionPolicy{
		Id:         testId,
		SourceName: testSourceName,
		Duration:   testDuration,
	}
)

func TestToRetentionPolicyModel(t *testing.T) {
	actualPolicyModel := ToRetentionPolicyModel(testRetentionPolicy1)

	require.Equal(t, testRetentionPolicyModel1, actualPolicyModel)
}

func TestFromRetentionPolicyModelToDTO(t *testing.T) {
	policyModel := testRetentionPolicyModel2
	policyModel.Created = createdTimestamp
	policyModel.Modified = modifiedTimestamp

	expectedPolicyDTO := testRetentionPolicy2
	expectedPolicyDTO.Created = createdTimestamp
	expectedPolicyDTO.Modified = modifiedTimestamp

	actualPolicyDTO := FromRetentionPolicyModelToDTO(policyModel)
	require.Equal(t, expectedPolicyDTO, actualPolicyDTO)
}
