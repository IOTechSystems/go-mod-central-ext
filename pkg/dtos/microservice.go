//
// Copyright (C) 2022-2024 IOTech Ltd
//

package dtos

// Port and its properties are defined in the edgex-go-private sys-mgmt API specification
type Port struct {
	// Host IP address that the container's port is mapped to
	IP string `json:"ip,omitempty"`

	// Port on the container
	PrivatePort uint16 `json:"privatePort"`

	// Port exposed on the host
	PublicPort uint16 `json:"publicPort,omitempty"`

	// The type of the mapped port
	Type string `json:"type"`
}

// MicroService and its properties are defined in the edgex-go-private sys-mgmt API specification
type MicroService struct {
	ID        string            `json:"id"`
	Names     []string          `json:"names"`
	Image     string            `json:"image"`
	Created   int64             `json:"created"`
	Ports     []Port            `json:"ports"`
	Labels    map[string]string `json:"labels"`
	State     string            `json:"state"`
	IPAddress string            `json:"ipAddress"`
}
