// Copyright (C) 2025 IOTech Ltd

package dtos

import (
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/models"

	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
)

type RetentionPolicy struct {
	edgexDtos.DBTimestamp
	Id         string `json:"id"`
	DeviceName string `json:"deviceName" validate:"required_without=SourceName"`
	SourceName string `json:"sourceName" validate:"required_without=DeviceName"`
	Duration   string `json:"duration" validate:"required,edgex-dto-duration"`
}

// ToRetentionPolicyModel transforms the RetentionPolicy DTO to the RetentionPolicy Model
func ToRetentionPolicyModel(policy RetentionPolicy) models.RetentionPolicy {
	return models.RetentionPolicy{
		Id:         policy.Id,
		DeviceName: policy.DeviceName,
		SourceName: policy.SourceName,
		Duration:   policy.Duration,
	}
}

// FromRetentionPolicyModelToDTO transforms the RetentionPolicy Model to the RetentionPolicy DTO
func FromRetentionPolicyModelToDTO(policy models.RetentionPolicy) RetentionPolicy {
	return RetentionPolicy{
		DBTimestamp: edgexDtos.DBTimestamp(policy.DBTimestamp),
		Id:          policy.Id,
		DeviceName:  policy.DeviceName,
		SourceName:  policy.SourceName,
		Duration:    policy.Duration,
	}
}
