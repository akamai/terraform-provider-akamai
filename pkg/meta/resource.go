package meta

import (
	"context"
	"fmt"

	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Resource is the base struct for all resources.
type Resource struct {
	Client edgegrid.Client
}

// Configure implements resource.ResourceWithConfigure.
func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		// ProviderData is nil when Configure is run first time as part of ValidateResourceConfig
		// in framework provider.
		return
	}

	m, ok := req.ProviderData.(Meta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.",
				req.ProviderData),
		)
		return
	}

	r.Client = m.Client()
}
