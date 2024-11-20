// Copyright (C) 2024 IOTech Ltd

package dtos

import (
	"testing"

	"github.com/edgexfoundry/go-mod-central-ext/v4/pkg/models"

	"github.com/stretchr/testify/require"
)

var (
	mockRole    = "mockAdmin"
	mockPath    = "*/core-command/*"
	mockMethods = []string{"GET", "PUT"}
	mockEffect  = "allow"

	mockAccessPolicyDTO = AccessPolicy{
		Path:        mockPath,
		HttpMethods: mockMethods,
		Effect:      mockEffect,
	}
	mockAccessPolicyModel = models.AccessPolicy{
		Path:        mockPath,
		HttpMethods: mockMethods,
		Effect:      mockEffect,
	}
	mockRolePolicyDTO = RolePolicy{
		Role:           mockRole,
		AccessPolicies: []AccessPolicy{mockAccessPolicyDTO},
	}
	mockRolePolicyModel = models.RolePolicy{
		Role:           mockRole,
		AccessPolicies: []models.AccessPolicy{mockAccessPolicyModel},
	}
)

func TestToRolePolicyModel(t *testing.T) {
	actualRolePolicyModel := ToRolePolicyModel(mockRolePolicyDTO)

	require.Equal(t, mockRolePolicyModel, actualRolePolicyModel)
}

func TestFromRolePolicyModelToDTO(t *testing.T) {
	actualRolePolicyDTO := FromRolePolicyModelToDTO(mockRolePolicyModel)
	require.Equal(t, mockRolePolicyDTO, actualRolePolicyDTO)
}

func TestFromRolePolicyModelsToDTOs(t *testing.T) {
	actualRolePolicyDTOs := FromRolePolicyModelsToDTOs([]models.RolePolicy{mockRolePolicyModel})
	require.Equal(t, []RolePolicy{mockRolePolicyDTO}, actualRolePolicyDTOs)
}

func TestToAccessPolicyModel(t *testing.T) {
	actualAccessPolicyModel := ToAccessPolicyModel(mockAccessPolicyDTO)

	require.Equal(t, mockAccessPolicyModel, actualAccessPolicyModel)
}

func TestFromAccessPolicyModelToDTO(t *testing.T) {
	actualAccessPolicyDTO := FromAccessPolicyModelToDTO(mockAccessPolicyModel)
	require.Equal(t, mockAccessPolicyDTO, actualAccessPolicyDTO)
}

func TestFromAccessPolicyModelsToDTOs(t *testing.T) {
	actualAccessPolicyDTOs := FromAccessPolicyModelsToDTOs([]models.AccessPolicy{mockAccessPolicyModel})
	require.Equal(t, []AccessPolicy{mockAccessPolicyDTO}, actualAccessPolicyDTOs)
}
