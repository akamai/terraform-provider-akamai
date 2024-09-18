package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &CIDRBlocksResource{}
	_ resource.ResourceWithImportState = &CIDRBlocksResource{}
)

// CIDRBlocksResource represents akamai_iam_cidr_blocks resource
type CIDRBlocksResource struct {
	meta meta.Meta
}

// NewCIDRBlocksResource returns new akamai_iam_cidr_blocks resource
func NewCIDRBlocksResource() resource.Resource { return &CIDRBlocksResource{} }

// CIDRBlocksResourceModel represents model of akamai_iam_cidr_blocks resource
type CIDRBlocksResourceModel struct {
	CIDRBlocks []CIDRBlock `tfsdk:"cidr_blocks"`
}

// CIDRBlock represents a single cidr block
type CIDRBlock struct {
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

// Metadata implements resource.Resource.
func (r *CIDRBlocksResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_iam_cidr_blocks"
}

// Schema implements resource.Resource.
func (r *CIDRBlocksResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cidr_blocks": schema.ListNestedAttribute{
				Optional:    true,
				Description: "List of CIDR blocks.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cidr_block": schema.StringAttribute{
							Required:    true,
							Description: "The value of an IP address or IP address range.",
						},
						"enabled": schema.BoolAttribute{
							Required:    true,
							Description: "Enables the IP allowlist on the account.",
						},
						"comments": schema.StringAttribute{
							Optional:    true,
							Description: "Descriptive label you provide for the CIDR block.",
						},
						"actions": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Specifies activities available for the CIDR block.",
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
							Computed:      true,
							Description:   "Unique identifier for each CIDR block.",
							PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
						},
						"created_by": schema.StringAttribute{
							Computed:    true,
							Description: "The user who created the CIDR block.",
						},
						"created_date": schema.StringAttribute{
							Computed:    true,
							Description: "ISO 8601 timestamp indicating when the CIDR block was created.",
						},
						"modified_by": schema.StringAttribute{
							Computed:    true,
							Description: "The user who last edited the CIDR block.",
						},
						"modified_date": schema.StringAttribute{
							Computed:    true,
							Description: "ISO 8601 timestamp indicating when the CIDR block was last modified.",
						},
					},
				},
			},
		},
	}
}

// Configure implements resource.ResourceWithConfigure.
func (r *CIDRBlocksResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create implements resource.Resource.
func (r *CIDRBlocksResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating CIDR Blocks resource")
	var plan *CIDRBlocksResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.create(ctx, plan); err != nil {
		resp.Diagnostics.AddError("create cidr block failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CIDRBlocksResource) create(ctx context.Context, plan *CIDRBlocksResourceModel) error {
	client := inst.Client(r.meta)

	for i, c := range plan.CIDRBlocks {
		cidr, err := client.CreateCIDRBlock(ctx, c.buildCreateCIDRBlockRequest())
		if err != nil {
			return err
		}
		plan.setFromCreateCIDRBlock(cidr, i)
	}

	return nil
}

// Read implements resource.Resource.
func (r *CIDRBlocksResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading CIDR Blocks Resource")
	var state *CIDRBlocksResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.read(ctx, state); err != nil {
		resp.Diagnostics.AddError("read cidr blocks error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CIDRBlocksResource) read(ctx context.Context, data *CIDRBlocksResourceModel) error {
	client := inst.Client(r.meta)

	res, err := client.ListCIDRBlocks(ctx, iam.ListCIDRBlocksRequest{
		Actions: true,
	})
	if err != nil {
		return err
	}
	data.setFromListCIDRBlocks(res)

	return nil
}

// Update implements resource.Resource.
func (r *CIDRBlocksResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating CIDR Blocks Resource")
	var data *CIDRBlocksResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.update(ctx, data); err != nil {
		resp.Diagnostics.AddError("update cidr block failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CIDRBlocksResource) update(ctx context.Context, data *CIDRBlocksResourceModel) error {
	client := inst.Client(r.meta)

	for i, cidr := range data.CIDRBlocks {
		resp, err := client.UpdateCIDRBlock(ctx, cidr.buildUpdateRequest())
		if err != nil {
			return err
		}
		data.setFromUpdateCIDRBlock(resp, i)
	}

	return nil
}

// Delete implements resource.Resource.
func (r *CIDRBlocksResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting CIDR Blocks Resource")
	var state *CIDRBlocksResourceModel
	client := inst.Client(r.meta)

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, cidr := range state.CIDRBlocks {
		if err := client.DeleteCIDRBlock(ctx, iam.DeleteCIDRBlockRequest{
			CIDRBlockID: cidr.CIDRBlockID.ValueInt64(),
		}); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("delete cidr block %d failed", cidr.CIDRBlockID), err.Error())
			return
		}
	}
	resp.State.RemoveResource(ctx)
}

// ImportState implements resource.ResourceWithImportState.
func (r *CIDRBlocksResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing CIDR Blocks Resource")

	data := &CIDRBlocksResourceModel{}
	client := inst.Client(r.meta)

	res, err := client.ListCIDRBlocks(ctx, iam.ListCIDRBlocksRequest{
		Actions: true,
	})
	if err != nil {
		resp.Diagnostics.AddError("import cidr blocks error", err.Error())
		return
	}

	data.setFromImportCIDRBlocks(res)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m CIDRBlock) buildCreateCIDRBlockRequest() iam.CreateCIDRBlockRequest {
	return iam.CreateCIDRBlockRequest{
		CIDRBlock: m.CIDR.ValueString(),
		Comments:  m.Comments.ValueStringPointer(),
		Enabled:   m.Enabled.ValueBool(),
	}
}

func (m CIDRBlock) buildUpdateRequest() iam.UpdateCIDRBlockRequest {
	return iam.UpdateCIDRBlockRequest{
		CIDRBlockID: m.CIDRBlockID.ValueInt64(),
		Body: iam.UpdateCIDRBlockRequestBody{
			CIDRBlock: m.CIDR.ValueString(),
			Comments:  m.Comments.ValueStringPointer(),
			Enabled:   m.Enabled.ValueBool(),
		},
	}
}

func (c *CIDRBlocksResourceModel) setFromCreateCIDRBlock(cidr *iam.CreateCIDRBlockResponse, i int) {
	if cidr.Actions != nil {
		c.CIDRBlocks[i].Actions = getActionFor(cidr.Actions.Edit, cidr.Actions.Delete)
	} else {
		c.CIDRBlocks[i].Actions = types.ObjectNull(map[string]attr.Type{
			"edit":   types.BoolType,
			"delete": types.BoolType,
		})
	}
	c.CIDRBlocks[i].CIDR = types.StringValue(cidr.CIDRBlock)
	if cidr.Comments != nil {
		c.CIDRBlocks[i].Comments = types.StringValue(*cidr.Comments)
	} else {
		c.CIDRBlocks[i].Comments = types.StringNull()
	}
	c.CIDRBlocks[i].CIDRBlockID = types.Int64Value(cidr.CIDRBlockID)
	c.CIDRBlocks[i].ModifiedBy = types.StringValue(cidr.ModifiedBy)
	c.CIDRBlocks[i].ModifiedDate = types.StringValue(cidr.ModifiedDate.Format(time.RFC3339Nano))
	c.CIDRBlocks[i].CreatedBy = types.StringValue(cidr.CreatedBy)
	c.CIDRBlocks[i].CreatedDate = types.StringValue(cidr.CreatedDate.Format(time.RFC3339Nano))
}

func (c *CIDRBlocksResourceModel) setFromListCIDRBlocks(resp iam.ListCIDRBlocksResponse) {
	for _, cidr := range c.CIDRBlocks {
		for _, r := range resp {
			if cidr.CIDRBlockID.ValueInt64() == r.CIDRBlockID {
				cidr.CIDR = types.StringValue(r.CIDRBlock)
				if r.Actions != nil {
					cidr.Actions = getActionFor(r.Actions.Edit, r.Actions.Delete)
				} else {
					cidr.Actions = types.ObjectNull(map[string]attr.Type{
						"edit":   types.BoolType,
						"delete": types.BoolType,
					})
				}
				if r.Comments != nil {
					cidr.Comments = types.StringValue(*r.Comments)
				} else {
					cidr.Comments = types.StringNull()
				}
				cidr.Enabled = types.BoolValue(r.Enabled)
				cidr.ModifiedBy = types.StringValue(r.ModifiedBy)
				cidr.ModifiedDate = types.StringValue(r.ModifiedDate.Format(time.RFC3339Nano))
				cidr.CreatedBy = types.StringValue(r.CreatedBy)
				cidr.CreatedDate = types.StringValue(r.CreatedDate.Format(time.RFC3339Nano))
				break
			}
		}
	}
}

func (c *CIDRBlocksResourceModel) setFromUpdateCIDRBlock(resp *iam.UpdateCIDRBlockResponse, i int) {
	if resp.Actions != nil {
		c.CIDRBlocks[i].Actions = getActionFor(resp.Actions.Edit, resp.Actions.Delete)
	} else {
		c.CIDRBlocks[i].Actions = types.ObjectNull(map[string]attr.Type{
			"edit":   types.BoolType,
			"delete": types.BoolType,
		})
	}
	if resp.Comments != nil {
		c.CIDRBlocks[i].Comments = types.StringValue(*resp.Comments)
	} else {
		c.CIDRBlocks[i].Comments = types.StringNull()
	}
	c.CIDRBlocks[i].CIDR = types.StringValue(resp.CIDRBlock)
	c.CIDRBlocks[i].CIDRBlockID = types.Int64Value(resp.CIDRBlockID)
	c.CIDRBlocks[i].Enabled = types.BoolValue(resp.Enabled)
	c.CIDRBlocks[i].CreatedBy = types.StringValue(resp.CreatedBy)
	c.CIDRBlocks[i].CreatedDate = types.StringValue(resp.CreatedDate.Format(time.RFC3339Nano))
	c.CIDRBlocks[i].ModifiedBy = types.StringValue(resp.ModifiedBy)
	c.CIDRBlocks[i].ModifiedDate = types.StringValue(resp.ModifiedDate.Format(time.RFC3339Nano))
}

func (c *CIDRBlocksResourceModel) setFromImportCIDRBlocks(resp iam.ListCIDRBlocksResponse) {
	var action types.Object

	for _, cidr := range resp {
		if cidr.Actions != nil {
			action = getActionFor(cidr.Actions.Edit, cidr.Actions.Delete)
		} else {
			action = types.ObjectNull(map[string]attr.Type{
				"edit":   types.BoolType,
				"delete": types.BoolType,
			})
		}
		block := CIDRBlock{
			CIDR:         types.StringValue(cidr.CIDRBlock),
			Enabled:      types.BoolValue(cidr.Enabled),
			Comments:     types.StringPointerValue(cidr.Comments),
			Actions:      action,
			CIDRBlockID:  types.Int64Value(cidr.CIDRBlockID),
			CreatedBy:    types.StringValue(cidr.CreatedBy),
			CreatedDate:  types.StringValue(cidr.CreatedDate.Format(time.RFC3339Nano)),
			ModifiedBy:   types.StringValue(cidr.ModifiedBy),
			ModifiedDate: types.StringValue(cidr.ModifiedDate.Format(time.RFC3339Nano)),
		}
		c.CIDRBlocks = append(c.CIDRBlocks, block)
	}
}

func getActionFor(editActions, deleteActions bool) basetypes.ObjectValue {
	return types.ObjectValueMust(
		map[string]attr.Type{
			"edit":   types.BoolType,
			"delete": types.BoolType,
		},
		map[string]attr.Value{
			"edit":   basetypes.NewBoolValue(editActions),
			"delete": basetypes.NewBoolValue(deleteActions),
		},
	)
}
