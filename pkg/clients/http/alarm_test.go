// Copyright (C) 2026 IOTech Ltd

package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	pkgCommon "github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
)

// TestAlarmClient_queryAll_RequestsAllItems verifies that the "list all" path
// (here exercised via AllTemplates) asks support-alarm for every item by sending
// offset=0&limit=-1, rather than a capped page size.
func TestAlarmClient_queryAll_RequestsAllItems(t *testing.T) {
	var gotOffset, gotLimit string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotOffset = r.URL.Query().Get(pkgCommon.Offset)
		gotLimit = r.URL.Query().Get(pkgCommon.Limit)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{pkgCommon.AlarmJsonKeyTemplates: []any{}})
	}))
	defer ts.Close()

	client := NewAlarmClient(ts.URL, NewNullAuthenticationInjector(), false)
	_, err := client.AllTemplates(context.Background())
	require.NoError(t, err)

	require.Equal(t, "0", gotOffset)
	require.Equal(t, "-1", gotLimit)
}
