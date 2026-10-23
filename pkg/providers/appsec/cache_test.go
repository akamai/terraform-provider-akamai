package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	akalog "github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	akameta "github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/stretchr/testify/require"
)

// mockMeta implements akameta.Meta for unit tests. It wires a TestClient so that
// functions calling meta.Client().GetAPPSEC() receive the mock APPSEC client.
type mockMeta struct {
	client *edgegrid.TestClient
}

func (m *mockMeta) Log(_ ...interface{}) akalog.Interface { return akalog.NOPLogger() }
func (m *mockMeta) OperationID() string                   { return "test-op" }
func (m *mockMeta) Session() session.Session              { return nil }
func (m *mockMeta) Client() edgegrid.Client               { return m.client }

func newTestMeta(client *edgegrid.TestClient) akameta.Meta {
	return &mockMeta{client: client}
}

// buildRatePoliciesResponse creates a GetRatePoliciesResponse from a JSON string.
func buildRatePoliciesResponse(t *testing.T, raw string) *appsec.GetRatePoliciesResponse {
	t.Helper()
	out := &appsec.GetRatePoliciesResponse{}
	require.NoError(t, json.Unmarshal([]byte(raw), out))
	return out
}

func buildMatchTargetsResponse(t *testing.T, raw string) *appsec.GetMatchTargetsResponse {
	t.Helper()
	out := &appsec.GetMatchTargetsResponse{}
	require.NoError(t, json.Unmarshal([]byte(raw), out))
	return out
}

func buildReputationProfilesResponse(t *testing.T, raw string) *appsec.GetReputationProfilesResponse {
	t.Helper()
	out := &appsec.GetReputationProfilesResponse{}
	require.NoError(t, json.Unmarshal([]byte(raw), out))
	return out
}
