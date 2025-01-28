package cloudwrapper

import (
	"context"
	"fmt"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cloudwrapper"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/collections"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type capacitiesDataSourceModel struct {
	ContractIDs types.List              `tfsdk:"contract_ids"`
	Capacities  []locationCapacityModel `tfsdk:"capacities"`
}

func (m capacitiesDataSourceModel) getContractIDs(ctx context.Context) ([]string, diag.Diagnostics) {
	var contracts []string
	diags := m.ContractIDs.ElementsAs(ctx, &contracts, false)
	if diags.HasError() {
		return nil, diags
	}

	collections.ForEachInSlice(contracts, func(s string) string {
		return strings.TrimPrefix(s, "ctr_")
	})

	return contracts, nil
}

type locationCapacityModel struct {
	LocationID         types.Int64   `tfsdk:"location_id"`
	LocationName       types.String  `tfsdk:"location_name"`
	ContractID         types.String  `tfsdk:"contract_id"`
	Type               types.String  `tfsdk:"type"`
	ApprovedCapacity   capacityModel `tfsdk:"approved"`
	AssignedCapacity   capacityModel `tfsdk:"assigned"`
	UnassignedCapacity capacityModel `tfsdk:"unassigned"`
}

type capacityModel struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func newLocationCapacityModel(capacity cloudwrapper.LocationCapacity) locationCapacityModel {
	return locationCapacityModel{
		LocationID:         types.Int64Value(int64(capacity.LocationID)),
		LocationName:       types.StringValue(capacity.LocationName),
		ContractID:         types.StringValue(capacity.ContractID),
		Type:               types.StringValue(string(capacity.Type)),
		ApprovedCapacity:   newCapacityModel(capacity.ApprovedCapacity),
		AssignedCapacity:   newCapacityModel(capacity.AssignedCapacity),
		UnassignedCapacity: newCapacityModel(capacity.UnassignedCapacity),
	}
}

func newCapacityModel(capacity cloudwrapper.Capacity) capacityModel {
	return capacityModel{
		Value: types.Int64Value(int64(capacity.Value)),
		Unit:  types.StringValue(string(capacity.Unit)),
	}
}

var (
	_ datasource.DataSource              = &capacitiesDataSource{}
	_ datasource.DataSourceWithConfigure = &capacitiesDataSource{}
)

type capacitiesDataSource struct {
	client cloudwrapper.CloudWrapper
}

// NewCapacitiesDataSource returns a new capacity data source
func NewCapacitiesDataSource() datasource.DataSource {
	return &capacitiesDataSource{}
}

func (d *capacitiesDataSource) setClient(client cloudwrapper.CloudWrapper) {
	d.client = client
}

// Metadata configures data source's meta information
func (d *capacitiesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.name()
}

func (d *capacitiesDataSource) name() string {
	return "akamai_cloudwrapper_capacities"
}

// Configure configures data source at the beginning of the lifecycle
func (d *capacitiesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if d.client != nil {
		return
	}

	m, ok := req.ProviderData.(meta.Meta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
	d.client = cloudwrapper.Client(m.Session())
}

// Schema is used to define data source's terraform schema
func (d *capacitiesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	capacityAttributes := map[string]attr.Type{
		"value": types.Int64Type,
		"unit":  types.StringType,
	}

	resp.Schema = schema.Schema{
		Description: "CloudWrapper capacities",
		Attributes: map[string]schema.Attribute{
			"contract_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of contract IDs with Cloud Wrapper entitlement.",
			},
		},
		Blocks: map[string]schema.Block{
			"capacities": schema.ListNestedBlock{
				Description: "List of all location capacities.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"location_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Unique identifier of the configured location.",
						},
						"location_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the location.",
						},
						"contract_id": schema.StringAttribute{
							Computed:    true,
							Description: "Contract ID having Cloud Wrapper entitlement.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of property this capacity is related to.",
						},
						"approved": schema.ObjectAttribute{
							Computed:       true,
							Description:    "Capacity allocated for the location.",
							AttributeTypes: capacityAttributes,
						},
						"assigned": schema.ObjectAttribute{
							Computed:       true,
							Description:    "Capacity already assigned to Cloud Wrapper configurations.",
							AttributeTypes: capacityAttributes,
						},
						"unassigned": schema.ObjectAttribute{
							Computed:       true,
							Description:    "Capacity value that can be assigned to Cloud Wrapper configurations.",
							AttributeTypes: capacityAttributes,
						},
					},
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state
func (d *capacitiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "CloudWrapper Capacities DataSource Read")

	var data capacitiesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	contractIDs, diags := data.getContractIDs(ctx)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	listCapacitiesResponse, err := d.client.ListCapacities(ctx, cloudwrapper.ListCapacitiesRequest{
		ContractIDs: contractIDs,
	})
	if err != nil {
		resp.Diagnostics.AddError("Reading CloudWrapper Capacities Failed", err.Error())
		return
	}

	for _, capacity := range listCapacitiesResponse.Capacities {
		data.Capacities = append(data.Capacities, newLocationCapacityModel(capacity))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
