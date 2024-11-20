// Copyright (C) 2023-2024 IOTech Ltd

package xrtmodels

import (
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
	"github.com/stretchr/testify/assert"
)

func TestXrtErrorCode(t *testing.T) {
	assert.Equal(t, XrtErrorCode(errors.NewCommonEdgeX(errors.KindEntityDoesNotExist, "", nil)), XrtSdkStatusNotFound)
	assert.Equal(t, XrtErrorCode(errors.NewCommonEdgeX(errors.KindNotImplemented, "", nil)), XrtSdkStatusNotSupported)
	assert.Equal(t, XrtErrorCode(errors.NewCommonEdgeX(errors.KindInvalidId, "", nil)), XrtSdkStatusInvalidOperation)
	assert.Equal(t, XrtErrorCode(errors.NewCommonEdgeX(errors.KindDuplicateName, "", nil)), XrtSdkStatusAlreadyExists)
	assert.Equal(t, XrtErrorCode(errors.NewCommonEdgeX(errors.KindServerError, "", nil)), XrtSdkStatusServerError)
}
