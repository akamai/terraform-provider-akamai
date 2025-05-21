package property

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	AssetID           types.String `tfsdk:"asset_id"`
	GroupID           types.String `tfsdk:"group_id"`
	ContractID        types.String `tfsdk:"contract_id"`
	ProductID         types.String `tfsdk:"product_id"`
	UseHostnameBucket types.Bool   `tfsdk:"use_hostname_bucket"`
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
			"use_hostname_bucket": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				Description: "Specifies whether hostname bucket is used with this property. " +
					"It allows you to add or remove property hostnames without incrementing property versions.",
				PlanModifiers: []planmodifier.Bool{
					modifiers.PreventBoolUpdate(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the Property",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"asset_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the property in the Identity and Access Management API.",
				PlanModifiers: []planmodifier.String{
					modifiers.StringUseStateIf(modifiers.EqualUpToPrefixFunc("aid_")),
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
	propertyID, err := createProperty(ctx, client, papi.CreatePropertyRequest{
		ContractID: contractID,
		GroupID:    groupID,
		Property: papi.PropertyCreate{
			ProductID:         productID,
			PropertyName:      data.Name.ValueString(),
			RuleFormat:        "",
			UseHostnameBucket: data.UseHostnameBucket.ValueBool(),
		},
	})
	if err != nil {
		err = interpretCreatePropertyErrorFramework(ctx, err, client, groupID, contractID, productID)
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), "")
			return
		}
	}

	data.ID = types.StringValue(propertyID)

	prop, err := fetchLatestProperty(ctx, client, propertyID, groupID, contractID)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	data.AssetID = types.StringValue(prop.AssetID)
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
	prop, err := fetchLatestProperty(ctx, client, propertyID, groupID, contractID)
	if errors.Is(err, papi.ErrNotFound) {
		tflog.Warn(ctx, fmt.Sprintf("property %q removed on server. Removing from local state", propertyID))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
	}

	data.AssetID = types.StringValue(prop.AssetID)
	// Protection against drift: if group id was changed from outside terraform,
	// store its current value
	data.GroupID = types.StringValue(prop.GroupID)
	useHostnameBucket := prop.PropertyType != nil && *prop.PropertyType == "HOSTNAME_BUCKET"
	data.UseHostnameBucket = types.BoolValue(useHostnameBucket)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update supports change for the following attributes:
// - `group_id` using a dedicated endpoint from the IAM API,
// - `name`, which results in resource replacement.
// Trying to update `contract_id` or `product_id` will result in an error.
func (r *BootstrapResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state BootstrapResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupsDiffer, err := areGroupIDsDifferent(state.GroupID.ValueString(), plan.GroupID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update Resource",
			fmt.Sprintf("An error occurred while parsing the group ids: %s, %s. Error: %s",
				state.GroupID.ValueString(), plan.GroupID.ValueString(), err.Error()))
		return
	}

	if groupsDiffer {
		hlp := helper{Client(r.meta), IAMClient(r.meta)}
		key := papiKey{
			propertyID: state.ID.ValueString(),
			groupID:    state.GroupID.ValueString(),
			contractID: state.ContractID.ValueString(),
		}
		err := hlp.moveProperty(ctx, key, state.AssetID.ValueString(), plan.GroupID.ValueString())

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Update Resource",
				"An error occurred while moving the property. Error: "+err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
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
