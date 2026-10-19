// Copyright (C) 2026 IOTech Ltd

package requests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/dtos"
	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
)

func TestUpdatePredictedResultValidate(t *testing.T) {
	tests := []struct {
		name         string
		resourceName string
		expectError  bool
	}{
		{"valid", "r1", false},
		{"empty resourceName", "", true},
		{"blank resourceName", " ", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := edgexCommon.Validate(dtos.UpdatePredictedResult{ResourceName: tt.resourceName, Tag: "t1"})
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
