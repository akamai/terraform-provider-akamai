package cloudcertificates

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &bindingsDataSource{}
	_ datasource.DataSourceWithConfigure = &bindingsDataSource{}
)

type (
	bindingsDataSource struct {
		meta.DataSource
	}

	bindingsDataSourceModel struct {
		LineageID types.Int64    `tfsdk:"lineage_id"`
		Network   types.String   `tfsdk:"network"`
		SortOrder types.String   `tfsdk:"sort_order"`
		Bindings  []bindingModel `tfsdk:"bindings"`
	}

	bindingModel struct {
		Active   types.Bool     `tfsdk:"active"`
		Hostname types.String   `tfsdk:"hostname"`
		Networks []types.String `tfsdk:"networks"`
	}
)

// NewBindingsDataSource returns a new bindings data source.
func NewBindingsDataSource() datasource.DataSource {
	return &bindingsDataSource{}
}

// Metadata configures the data source's type name.
func (d *bindingsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_bindings"
}

// Schema defines the Terraform schema for the bindings data source.
func (d *bindingsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for retrieving the hostnames bound to a certificate lineage.",
		Attributes: map[string]schema.Attribute{
			"lineage_id": schema.Int64Attribute{
				Required:    true,
				Description: "The unique numeric identifier of the certificate lineage.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"network": schema.StringAttribute{
				Optional:    true,
				Description: "Filter bindings by network. When not specified, bindings on both networks are returned.",
				Validators: []validator.String{stringvalidator.OneOf(
					string(cloudcertificates.TargetNetworkStaging),
					string(cloudcertificates.TargetNetworkProduction),
				)},
			},
			"sort_order": schema.StringAttribute{
				Optional: true,
				Description: "The order in which to sort the returned bindings, oldest first (`ASC`) or newest " +
					"first (`DESC`). Both options reflect the chronological order of the bindings. The default is `ASC`.",
				Validators: []validator.String{stringvalidator.OneOf(
					string(cloudcertificates.SortOrderAscending),
					string(cloudcertificates.SortOrderDescending),
				)},
			},
			"bindings": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Hostnames bound to the certificate lineage.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"active": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the binding is currently active.",
						},
						"hostname": schema.StringAttribute{
							Computed:    true,
							Description: "The bound hostname.",
						},
						"networks": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "The networks to which the hostname is bound, e.g., `STAGING`, `PRODUCTION`, or both.",
						},
					},
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state.
func (d *bindingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Cloud Certificates Bindings DataSource Read")

	var data bindingsDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	listReq := cloudcertificates.ListLineageBindingsRequest{
		LineageID: data.LineageID.ValueInt64(),
		Network:   cloudcertificates.TargetNetwork(data.Network.ValueString()),
		Sort:      cloudcertificates.SortOrder(data.SortOrder.ValueString()),
		PageSize:  cloudcertificates.MaxListLineageBindingsPageSize,
	}

	bindings, err := d.listAllBindings(ctx, listReq)
	if err != nil {
		if errors.Is(err, cloudcertificates.ErrLineageNotFound) {
			resp.Diagnostics.AddError("Lineage Not Found", fmt.Sprintf("No certificate lineage found with ID %d", data.LineageID.ValueInt64()))
			return
		}
		resp.Diagnostics.AddError("Failed to retrieve lineage bindings", fmt.Sprintf(
			"Unable to retrieve bindings for certificate lineage ID %d (network %q, sort order %q): %v. "+
				"Verify that the lineage and filters are valid, that your credentials can access it, and retry.",
			data.LineageID.ValueInt64(), data.Network.ValueString(), data.SortOrder.ValueString(), err))
		return
	}

	data.Bindings = bindingsToModels(bindings)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// listAllBindings follows the API's keyset pagination cursor, querying the maximum page size internally, until
// the full list of hostname bindings for the lineage has been collected.
func (d *bindingsDataSource) listAllBindings(ctx context.Context, listReq cloudcertificates.ListLineageBindingsRequest) ([]cloudcertificates.LineageBinding, error) {
	client := d.Client.GetCloudCertificates()

	var bindings []cloudcertificates.LineageBinding
	for {
		resp, err := client.ListLineageBindings(ctx, listReq)
		if err != nil {
			return nil, err
		}

		bindings = append(bindings, resp.Bindings...)

		nextCursor := resp.NextCursor
		if nextCursor == nil || len(resp.Bindings) == 0 || *nextCursor == "" {
			break
		}
		if *nextCursor == listReq.After {
			return nil, fmt.Errorf("lineage bindings pagination returned a non-advancing cursor")
		}
		listReq.After = *nextCursor
	}

	return bindings, nil
}

func bindingsToModels(bindings []cloudcertificates.LineageBinding) []bindingModel {
	models := make([]bindingModel, len(bindings))
	for i, b := range bindings {
		networks := make([]types.String, len(b.Networks))
		for j, n := range b.Networks {
			networks[j] = types.StringValue(n)
		}
		models[i] = bindingModel{
			Active:   types.BoolValue(b.Active),
			Hostname: types.StringValue(b.Hostname),
			Networks: networks,
		}
	}
	return models
}
