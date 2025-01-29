package iam

import (
	"context"
	"fmt"

	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &ipAllowlistResource{}
	_ resource.ResourceWithConfigure   = &ipAllowlistResource{}
	_ resource.ResourceWithImportState = &ipAllowlistResource{}
)

type ipAllowlistResource struct {
	meta meta.Meta
}

// NewIPAllowlistResource returns new akamai_iam_ip_allowlist resource.
func NewIPAllowlistResource() resource.Resource {
	return &ipAllowlistResource{}
}

type ipAllowlistModel struct {
	Enable types.Bool `tfsdk:"enable"`
}

func (r *ipAllowlistResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"unexpected resource configure type",
				fmt.Sprintf("expected meta.Meta, got: %T. please report this issue to the provider developers.", req.ProviderData),
			)
		}
	}()

	r.meta = meta.Must(req.ProviderData)
}

func (r *ipAllowlistResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_iam_ip_allowlist"
}

func (r *ipAllowlistResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"enable": schema.BoolAttribute{
				Required:    true,
				Description: "Whether to enable or disable the allowlist.",
			},
		},
	}
}

func (r *ipAllowlistResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating IP Allowlist resource")
	var plan *ipAllowlistModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Enable.ValueBool() {
		resp.Diagnostics.Append(r.enableIPAllowlist(ctx)...)
	} else {
		resp.Diagnostics.Append(r.disableIPAllowlist(ctx)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ipAllowlistResource) enableIPAllowlist(ctx context.Context) diag.Diagnostics {
	var diags diag.Diagnostics
	client := inst.Client(r.meta)
	status, err := client.GetIPAllowlistStatus(ctx)
	if err != nil {
		diags.AddError("cannot fetch IP Allowlist status", err.Error())
		return diags
	}
	if !status.Enabled {
		err = client.EnableIPAllowlist(ctx)
		if err != nil {
			diags.AddError("enable ip allowlist fail", err.Error())
			return diags
		}
	}
	return nil
}

func (r *ipAllowlistResource) disableIPAllowlist(ctx context.Context) diag.Diagnostics {
	var diags diag.Diagnostics
	client := inst.Client(r.meta)
	status, err := client.GetIPAllowlistStatus(ctx)
	if err != nil {
		diags.AddError("cannot fetch IP Allowlist status", err.Error())
		return diags
	}
	if status.Enabled {
		err = client.DisableIPAllowlist(ctx)
		if err != nil {
			diags.AddError("disable IP allowlist fail", err.Error())
			return diags
		}
	}
	return nil
}

func (r *ipAllowlistResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading IP Allowlist Resource")
	var oldState *ipAllowlistModel
	client := inst.Client(r.meta)
	resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	status, err := client.GetIPAllowlistStatus(ctx)
	if err != nil {
		resp.Diagnostics.AddError("cannot fetch IP Allowlist status", err.Error())
		return
	}

	oldState.Enable = types.BoolValue(status.Enabled)
	resp.Diagnostics.Append(resp.State.Set(ctx, &oldState)...)
}

func (r *ipAllowlistResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating IP Allowlist Resource")
	var diags diag.Diagnostics
	var plan *ipAllowlistModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if plan.Enable.ValueBool() {
		diags.Append(r.enableIPAllowlist(ctx)...)
	} else {
		diags.Append(r.disableIPAllowlist(ctx)...)
	}
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ipAllowlistResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan *ipAllowlistModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Enable.ValueBool() {
		resp.Diagnostics.Append(r.disableIPAllowlist(ctx)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// It can be only removed from state
	resp.State.RemoveResource(ctx)
}

func (r *ipAllowlistResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing IP Allowlist Resource")
	client := inst.Client(r.meta)
	status, err := client.GetIPAllowlistStatus(ctx)
	if err != nil {
		resp.Diagnostics.AddError("cannot fetch IP Allowlist status", err.Error())
		return
	}
	var state = &ipAllowlistModel{}
	state.Enable = types.BoolValue(status.Enabled)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
