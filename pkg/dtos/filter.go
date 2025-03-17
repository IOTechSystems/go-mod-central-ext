// Copyright (C) 2025 IOTech Ltd

package dtos

import (
	"strings"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/models"
)

type Filter struct {
	Id              string `json:"id"`
	Type            string `json:"type" validate:"required,oneof=IN OUT"`
	DeviceName      string `json:"deviceName" validate:"required,edgex-dto-none-empty-string"`
	EventSourceName string `json:"eventSourceName"`
	ResourceName    string `json:"resourceName"`
}

// ToFilterModel transforms the Filter DTO to the Filter Model
func ToFilterModel(filter Filter) models.Filter {
	return models.Filter{
		Id:              filter.Id,
		Type:            strings.ToUpper(filter.Type),
		DeviceName:      filter.DeviceName,
		EventSourceName: filter.EventSourceName,
		ResourceName:    filter.ResourceName,
	}
}

type UpdateFilter struct {
	Id              string `json:"id" validate:"required,edgex-dto-none-empty-string"`
	Type            string `json:"type" validate:"required,edgex-dto-none-empty-string"`
	DeviceName      string `json:"deviceName" validate:"required,edgex-dto-none-empty-string"`
	EventSourceName string `json:"eventSourceName"`
	ResourceName    string `json:"resourceName"`
}

// FromFilterModelToDTO transforms the Filter Model to the Filter DTO
func FromFilterModelToDTO(filter models.Filter) Filter {
	return Filter{
		Id:              filter.Id,
		Type:            filter.Type,
		DeviceName:      filter.DeviceName,
		EventSourceName: filter.EventSourceName,
		ResourceName:    filter.ResourceName,
	}
}
