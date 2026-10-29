// Copyright (C) 2026 IOTech Ltd

package dtos

// AIPointsSummary lists the tags of one ontology that the site's points carry, one row per tag
type AIPointsSummary struct {
	Ontology     string         `json:"ontology"`
	TotalPoints  int            `json:"totalPoints"`
	TotalDevices int            `json:"totalDevices"`
	Tags         []AITagSummary `json:"tags"`
}

// AITagSummary is one tag of an AIPointsSummary; DeviceNames is set only when the site has at most 10 devices
type AITagSummary struct {
	Tag         string   `json:"tag"`
	Points      int      `json:"points"`
	Devices     int      `json:"devices"`
	DeviceNames []string `json:"deviceNames,omitempty"`
	Units       []string `json:"units"`
	Description string   `json:"description"`
}

// AIPointsPage is one page of the points carrying Tag in Ontology; NextOffset is nil on the last page
type AIPointsPage struct {
	Ontology   string    `json:"ontology"`
	Tag        string    `json:"tag"`
	Total      int       `json:"total"`
	Offset     int       `json:"offset"`
	NextOffset *int      `json:"nextOffset"`
	Points     []AIPoint `json:"points"`
}

// AIPoint is a device resource carrying a tag
type AIPoint struct {
	DeviceName   string `json:"deviceName"`
	ResourceName string `json:"resourceName"`
	Units        string `json:"units"`
	Description  string `json:"description"`
}
