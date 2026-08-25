package meta

import (
	"context"
	"fmt"

	akalog "github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/hashicorp/terraform-plugin-framework/action"
)

// Action is the base struct for all actions.
type Action struct {
	Client edgegrid.Client
	meta   Meta
}

// Configure implements action.ActionWithConfigure.
func (a *Action) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		// ProviderData is nil when Configure is run first time as part of ValidateConfig
		// in framework provider.
		return
	}

	m, ok := req.ProviderData.(Meta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Action Configure Type",
			fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.",
				req.ProviderData),
		)
		return
	}

	a.Client = m.Client()
	a.meta = m
}

// Log creates a logger for the action from its provider meta.
func (a *Action) Log(args ...any) akalog.Interface {
	return a.meta.Log(args...)
}
