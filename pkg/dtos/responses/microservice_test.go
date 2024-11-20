//
// Copyright (C) 2022-2024 IOTech Ltd
//

package responses

import (
	"testing"

	"github.com/edgexfoundry/go-mod-central-ext/v4/pkg/dtos"

	"github.com/stretchr/testify/assert"
)

func TestNewMultiMicroServicesResponse(t *testing.T) {
	expectedRequestId := "123456"
	expectedStatusCode := 200
	expectedMessage := "unit test message"
	expectedMicroServices := []dtos.MicroService{
		{
			ID:    "7102e1adab50bfc9894e88d0afe903c0f0cf69d5305cbac86b7f3c6ee28afde4",
			Names: []string{"/core-data"},
			Image: "iotechsys/dev-edgecentral-core-data:2.2.dev",
		},
		{
			ID:    "bbf4d4be01004080043a29185c8e5d01a3973f082fa069898c2432b0f44996ff",
			Names: []string{"/core-command"},
			Image: "iotechsys/dev-edgecentral-core-command:2.2.dev",
		},
	}
	actual := NewMultiMicroServicesResponse(expectedRequestId, expectedMessage, expectedStatusCode, expectedMicroServices)

	assert.Equal(t, expectedRequestId, actual.RequestId)
	assert.Equal(t, expectedStatusCode, actual.StatusCode)
	assert.Equal(t, expectedMessage, actual.Message)
	assert.Equal(t, expectedMicroServices, actual.MicroServices)
}
