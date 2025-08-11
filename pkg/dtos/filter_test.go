// Copyright (C) 2025 IOTech Ltd

package dtos

import (
	"testing"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/models"

	"github.com/stretchr/testify/require"
)

var (
	testFilterId          = "431f0134-ae07-45ac-8577-1c1ac2b74fb3"
	testFilterType        = "IN"
	testDeviceName        = "^device+"
	testEventSourceName   = "^sensor+"
	testResourceName      = "Bool"
	testOnChangeThreshold = 5.0
	testFilterDTO         = Filter{
		Id:                testFilterId,
		Type:              testFilterType,
		DeviceName:        testDeviceName,
		EventSourceName:   testEventSourceName,
		ResourceName:      testResourceName,
		OnChange:          true,
		OnChangeThreshold: testOnChangeThreshold,
	}
	testFilterModel = models.Filter{
		Id:                testFilterId,
		Type:              testFilterType,
		DeviceName:        testDeviceName,
		EventSourceName:   testEventSourceName,
		ResourceName:      testResourceName,
		OnChange:          true,
		OnChangeThreshold: testOnChangeThreshold,
	}
)

func TestToFilterModel(t *testing.T) {
	actualFilterModel := ToFilterModel(testFilterDTO)
	require.Equal(t, testFilterModel, actualFilterModel)
}

func TestFromFilterModelToDTO(t *testing.T) {
	filterModel := testFilterModel
	filterModel.Created = createdTimestamp
	filterModel.Modified = modifiedTimestamp
	expectedFilterDTO := testFilterDTO
	expectedFilterDTO.Created = createdTimestamp
	expectedFilterDTO.Modified = modifiedTimestamp
	actualFilterDTO := FromFilterModelToDTO(filterModel)
	require.Equal(t, expectedFilterDTO, actualFilterDTO)
}
