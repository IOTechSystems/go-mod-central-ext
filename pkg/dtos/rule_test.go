// Copyright (C) 2023-2024 IOTech Ltd

package dtos

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRule(t *testing.T) {
	expectedName := "rule name"
	expectedRule := []byte("rule")
	actual := NewRule(expectedName, expectedRule)
	assert.Equal(t, expectedName, actual.Name)
	assert.Equal(t, expectedRule, actual.Rule)
}
