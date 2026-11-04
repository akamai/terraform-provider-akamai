package meta

import (
	"context"
	"fmt"

	"github.com/akamai/terraform-provider-akamai/v12/internal/edgegrid"
	"github.com/hashicorp/terraform-plugin-framework/action"
)

// Action is the base struct for all actions.
type Action struct {
	Client edgegrid.Client
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
}
