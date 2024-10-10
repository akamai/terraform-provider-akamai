package testprovider

import (
	"context"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &TestResource{}
	_ resource.ResourceWithImportState = &TestResource{}
)

// TestResource represents akamai_test to be used in internal testing
type TestResource struct{}

// NewTestResource returns new akamai test resource to be used in internal testing
func NewTestResource() resource.Resource {
	return &TestResource{}
}

// testResourceModel is a model for akamai_test_resource resource
type testResourceModel struct {
	ID     types.Int64  `tfsdk:"id"`
	Input  types.String `tfsdk:"input"`
	Output types.String `tfsdk:"output"`
}

// Metadata implements resource.Resource.
func (r *TestResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_test"
}

// Schema implements resource.Resource.
func (r *TestResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"input": schema.StringAttribute{
				Required:    true,
				Description: "Required attribute, its value is copied to 'output' attribute.",
			},
			"output": schema.StringAttribute{
				Computed:    true,
				Description: "Read-only attribute, its value comes from 'input' attribute.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "Resource's unique and randomly generated identifier.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *TestResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating Akamai Test Resource")

	var data *testResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Output = data.Input
	data.ID = types.Int64Value(int64(uuid.New().ID()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read implements resource.Resource.
func (r *TestResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading Akamai Test Resource")

	var data *testResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update implements resource.Resource.
func (r *TestResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating Akamai Test Resource")

	var data *testResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Output = data.Input

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete implements resource.Resource.
func (r *TestResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting Akamai Test Resource")
	resp.State.RemoveResource(ctx)
}

// ImportState implements resource.ResourceWithImportState.
func (r *TestResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing Configuration Resource")

	var data = &testResourceModel{}

	id := strings.Split(req.ID, ",")
	if len(id) != 2 {
		resp.Diagnostics.AddError("incorrect import ID", "import ID should have format: <id>,<input>")
		return
	}

	idInt, err := strconv.ParseInt(id[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("parsing error", err.Error())
		return
	}
	data.ID = types.Int64Value(idInt)
	data.Input = types.StringValue(id[1])
	data.Output = types.StringValue(id[1])

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
