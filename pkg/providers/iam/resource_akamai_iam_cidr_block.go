package iam

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &cidrBlockResource{}
	_ resource.ResourceWithImportState = &cidrBlockResource{}
)

type cidrBlockResource struct {
	meta meta.Meta
}

// NewCIDRBlockResource returns new akamai_iam_cidr_block resource.
func NewCIDRBlockResource() resource.Resource { return &cidrBlockResource{} }

type cidrBlockResourceModel struct {
	CIDR         types.String `tfsdk:"cidr_block"`
	Enabled      types.Bool   `tfsdk:"enabled"`
	Comments     types.String `tfsdk:"comments"`
	Actions      types.Object `tfsdk:"actions"`
	CIDRBlockID  types.Int64  `tfsdk:"cidr_block_id"`
	CreatedBy    types.String `tfsdk:"created_by"`
	CreatedDate  types.String `tfsdk:"created_date"`
	ModifiedBy   types.String `tfsdk:"modified_by"`
	ModifiedDate types.String `tfsdk:"modified_date"`
}

func (r *cidrBlockResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_iam_cidr_block"
}

func (r *cidrBlockResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cidr_block": schema.StringAttribute{
				Required:    true,
				Description: "The value of an IP address or IP address range.",
			},
			"enabled": schema.BoolAttribute{
				Required:    true,
				Description: "Enables the CIDR block on the account.",
			},
			"comments": schema.StringAttribute{
				Optional:    true,
				Description: "Descriptive label you provide for the CIDR block.",
			},
			"actions": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Specifies activities available for the CIDR block.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"delete": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can delete this CIDR block. You can't delete a CIDR block from an IP address not on the allowlist, or if the CIDR block is the only one on the allowlist.",
					},
					"edit": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can edit this CIDR block. You can't edit CIDR block from an IP address not on the allowlist, or if the CIDR block is the only one on the allowlist.",
					},
				},
			},
			"cidr_block_id": schema.Int64Attribute{
				Computed:    true,
				Description: "Unique identifier for each CIDR block.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who created the CIDR block.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_date": schema.StringAttribute{
				Computed:    true,
				Description: "ISO 8601 timestamp indicating when the CIDR block was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who last edited the CIDR block.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified_date": schema.StringAttribute{
				Computed:    true,
				Description: "ISO 8601 timestamp indicating when the CIDR block was last modified.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *cidrBlockResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		// ProviderData is nil when Configure is run first time as part of ValidateDataSourceConfig in framework provider
		return
	}

	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unexpected Resource Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
			)
		}
	}()

	r.meta = meta.Must(req.ProviderData)
}

func (r *cidrBlockResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating CIDR Block resource")
	var plan cidrBlockResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.create(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("create cidr block failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *cidrBlockResource) create(ctx context.Context, plan *cidrBlockResourceModel) error {
	client := inst.Client(r.meta)

	resp, err := client.CreateCIDRBlock(ctx, iam.CreateCIDRBlockRequest{
		CIDRBlock: plan.CIDR.ValueString(),
		Comments:  plan.Comments.ValueStringPointer(),
		Enabled:   plan.Enabled.ValueBool(),
	})
	if err != nil {
		return err
	}

	// Use GetCIDRBlock to fetch current actions, as they are not present in the response from CreateCIDRBlock
	cidr, err := client.GetCIDRBlock(ctx, iam.GetCIDRBlockRequest{
		CIDRBlockID: resp.CIDRBlockID,
		Actions:     true,
	})
	if err != nil {
		return err
	}

	plan.setData(cidr)

	return nil
}

func (r *cidrBlockResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading CIDR Block Resource")
	var state cidrBlockResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.read(ctx, &state); err != nil {
		resp.Diagnostics.AddError("read cidr block error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *cidrBlockResource) read(ctx context.Context, data *cidrBlockResourceModel) error {
	client := inst.Client(r.meta)

	cidr, err := client.GetCIDRBlock(ctx, iam.GetCIDRBlockRequest{
		CIDRBlockID: data.CIDRBlockID.ValueInt64(),
		Actions:     true,
	})
	if err != nil {
		return err
	}

	data.setData(cidr)

	return nil
}

func (r *cidrBlockResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating CIDR Block Resource")
	var plan cidrBlockResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.update(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("update cidr block failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *cidrBlockResource) update(ctx context.Context, plan *cidrBlockResourceModel) error {
	client := inst.Client(r.meta)

	_, err := client.UpdateCIDRBlock(ctx, iam.UpdateCIDRBlockRequest{
		CIDRBlockID: plan.CIDRBlockID.ValueInt64(),
		Body: iam.UpdateCIDRBlockRequestBody{
			CIDRBlock: plan.CIDR.ValueString(),
			Comments:  plan.Comments.ValueStringPointer(),
			Enabled:   plan.Enabled.ValueBool(),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *cidrBlockResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting CIDR Block Resource")

	var state *cidrBlockResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := inst.Client(r.meta)

	if err := client.DeleteCIDRBlock(ctx, iam.DeleteCIDRBlockRequest{
		CIDRBlockID: state.CIDRBlockID.ValueInt64(),
	}); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("delete cidr block %d failed", state.CIDRBlockID), err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *cidrBlockResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing CIDR Block Resource")

	cidrBlockID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("could not convert import ID to int", err.Error())
		return
	}

	data := &cidrBlockResourceModel{}

	// in import, we only need to set cidr block ID to allow read function to fill other attributes
	data.CIDRBlockID = types.Int64Value(cidrBlockID)
	// we also need to satisfy framework with a correct value for actions object
	data.Actions = types.ObjectNull(map[string]attr.Type{
		"edit":   types.BoolType,
		"delete": types.BoolType,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *cidrBlockResourceModel) setData(resp *iam.GetCIDRBlockResponse) {
	if resp.Actions != nil {
		m.Actions = getActionFor(resp.Actions.Edit, resp.Actions.Delete)
	} else {
		m.Actions = types.ObjectNull(map[string]attr.Type{
			"edit":   types.BoolType,
			"delete": types.BoolType,
		})
	}
	m.CIDRBlockID = types.Int64Value(resp.CIDRBlockID)
	m.ModifiedBy = types.StringValue(resp.ModifiedBy)
	m.ModifiedDate = types.StringValue(resp.ModifiedDate.Format(time.RFC3339Nano))
	m.CreatedBy = types.StringValue(resp.CreatedBy)
	m.CreatedDate = types.StringValue(resp.CreatedDate.Format(time.RFC3339Nano))
	m.CIDR = types.StringValue(resp.CIDRBlock)
	m.Comments = types.StringPointerValue(resp.Comments)
	m.Enabled = types.BoolValue(resp.Enabled)
}

func getActionFor(editAction, deleteAction bool) basetypes.ObjectValue {
	return types.ObjectValueMust(
		map[string]attr.Type{
			"edit":   types.BoolType,
			"delete": types.BoolType,
		},
		map[string]attr.Value{
			"edit":   basetypes.NewBoolValue(editAction),
			"delete": basetypes.NewBoolValue(deleteAction),
		},
	)
}
