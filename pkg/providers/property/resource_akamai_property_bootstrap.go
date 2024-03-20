package property

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &BootstrapResource{}
	_ resource.ResourceWithImportState = &BootstrapResource{}
	_ resource.ResourceWithConfigure   = &BootstrapResource{}
)

// BootstrapResource represents akamai_property_bootstrap resource
type BootstrapResource struct {
	meta meta.Meta
}

// BootstrapResourceModel is a model for akamai_property_bootstrap resource
type BootstrapResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	GroupID    types.String `tfsdk:"group_id"`
	ContractID types.String `tfsdk:"contract_id"`
	ProductID  types.String `tfsdk:"product_id"`
}

// NewBootstrapResource returns new property bootstrap resource
func NewBootstrapResource() resource.Resource {
	return &BootstrapResource{}
}

// Metadata implements resource.Resource.
func (r *BootstrapResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_property_bootstrap"
}

// Schema implements resource's Schema
func (r *BootstrapResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Name to give to the Property (must be unique)",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(85),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[A-Za-z0-9.\-_]+$`),
						"a name must only contain letters, numbers, and these characters: . _ -"),
				},
			},
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: "Group ID to be assigned to the Property",
				PlanModifiers: []planmodifier.String{
					modifiers.StringUseStateIf(modifiers.EqualUpToPrefixFunc("grp_")),
					modifiers.PreventStringUpdate(),
				},
			},
			"contract_id": schema.StringAttribute{
				Required:    true,
				Description: "Contract ID to be assigned to the Property",
				PlanModifiers: []planmodifier.String{
					modifiers.StringUseStateIf(modifiers.EqualUpToPrefixFunc("ctr_")),
					modifiers.PreventStringUpdate(),
				},
			},
			"product_id": schema.StringAttribute{
				Required:    true,
				Description: "Product ID to be assigned to the Property",
				PlanModifiers: []planmodifier.String{
					modifiers.StringUseStateIf(modifiers.EqualUpToPrefixFunc("prd_")),
					modifiers.PreventStringUpdate(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the Property",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure implements resource.ResourceWithConfigure.
func (r *BootstrapResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create implements resource's Create method
func (r *BootstrapResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating Bootstrap Resource")

	var data *BootstrapResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	contractID := str.AddPrefix(data.ContractID.ValueString(), "ctr_")
	groupID := str.AddPrefix(data.GroupID.ValueString(), "grp_")
	productID := str.AddPrefix(data.ProductID.ValueString(), "prd_")

	client := Client(r.meta)
	propertyID, err := createProperty(ctx, client, data.Name.ValueString(), groupID, contractID, productID, "")
	if err != nil {
		err = interpretCreatePropertyErrorFramework(ctx, err, client, groupID, contractID, productID)
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), "")
			return
		}
	}

	data.ID = types.StringValue(propertyID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func interpretCreatePropertyErrorFramework(ctx context.Context, err error, client papi.PAPI, groupID string, contractID string, productID string) error {
	if errors.Is(err, papi.ErrNotFound) {
		if _, err = getGroup(ctx, client, groupID); err != nil {
			if errors.Is(err, ErrGroupNotFound) {
				return fmt.Errorf("%v: %s", ErrGroupNotFound, groupID)
			}
			return err
		}
		if _, err = getContract(ctx, client, contractID); err != nil {
			if errors.Is(err, ErrContractNotFound) {
				return fmt.Errorf("%v: %s", ErrContractNotFound, contractID)
			}
			return err
		}
		if _, err = getProduct(ctx, client, productID, contractID); err != nil {
			if errors.Is(err, ErrProductNotFound) {
				return fmt.Errorf("%v: %s", ErrProductNotFound, productID)
			}
			return err
		}
	}
	return err
}

// Read implements resource's Read method
func (r *BootstrapResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading Bootstrap Resource")

	var data *BootstrapResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	propertyID := data.ID.ValueString()
	contractID := str.AddPrefix(data.ContractID.ValueString(), "ctr_")
	groupID := str.AddPrefix(data.GroupID.ValueString(), "grp_")

	client := Client(r.meta)
	_, err := fetchLatestProperty(ctx, client, propertyID, groupID, contractID)
	if errors.Is(err, papi.ErrNotFound) {
		tflog.Warn(ctx, fmt.Sprintf("property %q removed on server. Removing from local state", propertyID))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
	}
}

// Update of group, contract, product is noop, it will return an error before invoking Update. Updating name will result in resource replacement
func (r *BootstrapResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

// Delete implements resource's Delete method
func (r *BootstrapResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting Property Bootstrap")

	var data *BootstrapResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	propertyID := data.ID.ValueString()
	contractID := str.AddPrefix(data.ContractID.ValueString(), "ctr_")
	groupID := str.AddPrefix(data.GroupID.ValueString(), "grp_")

	client := Client(r.meta)
	if err := removeProperty(ctx, client, propertyID, groupID, contractID); err != nil {
		resp.Diagnostics.AddError("removeProperty:", err.Error())
	}
}

// ImportState implements resource's ImportState method
func (r *BootstrapResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing Property Bootstrap")

	var propertyID, groupID, contractID string
	parts := strings.Split(req.ID, ",")

	switch len(parts) {
	case 3:
		propertyID = str.AddPrefix(parts[0], "prp_")
		contractID = str.AddPrefix(parts[1], "ctr_")
		groupID = str.AddPrefix(parts[2], "grp_")
	case 2:
		resp.Diagnostics.AddError("missing group id or contract id", "")
		return
	case 1:
		propertyID = str.AddPrefix(parts[0], "prp_")
	default:
		resp.Diagnostics.AddError(fmt.Sprintf("invalid property identifier: %s", req.ID), "")
		return
	}

	client := Client(r.meta)
	property, err := fetchLatestProperty(ctx, client, propertyID, groupID, contractID)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	res, err := fetchPropertyVersion(ctx, client, property.PropertyID, property.GroupID, property.ContractID, property.LatestVersion)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	data := BootstrapResourceModel{
		ProductID:  types.StringValue(res.Version.ProductID),
		ContractID: types.StringValue(property.ContractID),
		GroupID:    types.StringValue(property.GroupID),
		Name:       types.StringValue(property.PropertyName),
		ID:         types.StringValue(propertyID),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
