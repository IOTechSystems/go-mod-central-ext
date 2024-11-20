// Copyright (C) 2024 IOTech Ltd

package responses

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/edgexfoundry/go-mod-central-ext/v4/pkg/common"
)

func TestNewTokenResponse(t *testing.T) {
	expectedToken := "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzEzODQxNDksImlzcyI6IklPVGVjaCIsImxjciI6MTczMTM4NTMxOSwidG9rZW5faWQiOiI3YWI4YTQwMS04ODEyLTRkMTgtOTgyMS05NTA3OGRiOGI2MDcifQ.arKj3pfSXRm2wH5chVaSMTBUA-cgSu_0CW2AbvXVEsiSIbB_KOt9p3pt2V1WWml2Tzvk7m_tLo-W_1HJVhuiCA" // nolint:gosec
	actual := NewTokenResponse(common.ExpectedRequestId, common.ExpectedMessage, common.ExpectedStatusCode, expectedToken)

	require.Equal(t, common.ExpectedRequestId, actual.RequestId)
	require.Equal(t, common.ExpectedStatusCode, actual.StatusCode)
	require.Equal(t, common.ExpectedMessage, actual.Message)
	require.Equal(t, expectedToken, actual.JWT)
}
