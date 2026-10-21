package testutils

import (
	"context"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
)

// ConfigureTestAction configures the given action (already constructed, with any action-specific fields set) with a
// fresh mock edgegrid client via the real Configure flow, and returns that client for mock expectation setup.
func ConfigureTestAction(t *testing.T, a action.ActionWithConfigure) *edgegrid.TestClient {
	t.Helper()

	client := edgegrid.NewTestClient()
	m := meta.NewTestMeta(client)

	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), action.ConfigureRequest{ProviderData: m}, configureResp)
	assert.False(t, configureResp.Diagnostics.HasError())

	return client
}

// BuildActionConfig derives a's schema object type and returns a tfsdk.Config populated with attrs, filling any
// attribute omitted from attrs with null.
func BuildActionConfig(ctx context.Context, t *testing.T, a action.Action, attrs map[string]tftypes.Value) tfsdk.Config {
	t.Helper()

	schemaResp := &action.SchemaResponse{}
	a.Schema(ctx, action.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())
	objectType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)

	values := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, attrType := range objectType.AttributeTypes {
		if v, ok := attrs[name]; ok {
			values[name] = v
			continue
		}
		values[name] = tftypes.NewValue(attrType, nil)
	}

	return tfsdk.Config{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objectType, values)}
}

// InvokeTestAction builds a's config from attrs (see BuildActionConfig) and invokes it under ctx, returning the
// response alongside every progress message sent during the invocation.
func InvokeTestAction(ctx context.Context, t *testing.T, a action.Action, attrs map[string]tftypes.Value) (*action.InvokeResponse, []string) {
	t.Helper()

	progress := []string{}
	resp := &action.InvokeResponse{
		SendProgress: func(e action.InvokeProgressEvent) {
			progress = append(progress, e.Message)
		},
	}
	a.Invoke(ctx, action.InvokeRequest{Config: BuildActionConfig(ctx, t, a, attrs)}, resp)
	return resp, progress
}

// AssertActionProgress asserts that the collected progress messages equal expected exactly, in order.
func AssertActionProgress(t *testing.T, expected, progress []string) {
	t.Helper()

	if len(expected) == 0 {
		assert.Empty(t, progress)
		return
	}
	assert.Equal(t, expected, progress)
}

// Diagnostic pairs a diagnostic's summary and detail for comparison in tests.
type Diagnostic struct {
	Summary string
	Detail  string
}

// Diagnostics builds a []Diagnostic from alternating summary, detail arguments, e.g.
// testutils.Diagnostics("Lineage Not Found", "No certificate lineage found with ID 999999.").
func Diagnostics(summaryDetailPairs ...string) []Diagnostic {
	if len(summaryDetailPairs)%2 != 0 {
		panic("testutils.Diagnostics: expected an even number of arguments (summary, detail, summary, detail, ...)")
	}
	diags := make([]Diagnostic, 0, len(summaryDetailPairs)/2)
	for i := 0; i < len(summaryDetailPairs); i += 2 {
		diags = append(diags, Diagnostic{Summary: summaryDetailPairs[i], Detail: summaryDetailPairs[i+1]})
	}
	return diags
}

// AssertActionDiagnostics asserts that diags contains exactly the given error and warning Diagnostics, regardless
// of order.
func AssertActionDiagnostics(t *testing.T, diags diag.Diagnostics, expectedErrors, expectedWarnings []Diagnostic) {
	t.Helper()

	var gotErrors, gotWarnings []Diagnostic
	for _, d := range diags {
		item := Diagnostic{Summary: d.Summary(), Detail: d.Detail()}
		switch d.Severity() {
		case diag.SeverityError:
			gotErrors = append(gotErrors, item)
		case diag.SeverityWarning:
			gotWarnings = append(gotWarnings, item)
		}
	}
	assert.ElementsMatch(t, expectedErrors, gotErrors, "unexpected errors")
	assert.ElementsMatch(t, expectedWarnings, gotWarnings, "unexpected warnings")
}
