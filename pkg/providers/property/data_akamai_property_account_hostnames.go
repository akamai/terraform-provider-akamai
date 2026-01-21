package property

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &accountHostnamesDataSource{}
	_ datasource.DataSourceWithConfigure = &accountHostnamesDataSource{}
)

type (
	accountHostnamesDataSource struct {
		meta.DataSource
	}

	accountHostnamesDataSourceModel struct {
		ContractID    types.String           `tfsdk:"contract_id"`
		GroupID       types.String           `tfsdk:"group_id"`
		Hostname      types.String           `tfsdk:"hostname"`
		CnameTo       types.String           `tfsdk:"cname_to"`
		Network       types.String           `tfsdk:"network"`
		Sort          types.String           `tfsdk:"sort"`
		AccountID     types.String           `tfsdk:"account_id"`
		AvailableSort types.List             `tfsdk:"available_sort"`
		CurrentSort   types.String           `tfsdk:"current_sort"`
		DefaultSort   types.String           `tfsdk:"default_sort"`
		Hostnames     []accountHostnameModel `tfsdk:"hostnames"`
	}

	accountHostnameModel struct {
		CnameFrom                types.String `tfsdk:"cname_from"`
		ContractID               types.String `tfsdk:"contract_id"`
		GroupID                  types.String `tfsdk:"group_id"`
		LatestVersion            types.Int64  `tfsdk:"latest_version"`
		ProductionCertType       types.String `tfsdk:"production_cert_type"`
		ProductionCnameTo        types.String `tfsdk:"production_cname_to"`
		ProductionCnameType      types.String `tfsdk:"production_cname_type"`
		ProductionEdgeHostnameID types.String `tfsdk:"production_edge_hostname_id"`
		ProductionProductID      types.String `tfsdk:"production_product_id"`
		PropertyID               types.String `tfsdk:"property_id"`
		PropertyName             types.String `tfsdk:"property_name"`
		PropertyType             types.String `tfsdk:"property_type"`
		StagingCertType          types.String `tfsdk:"staging_cert_type"`
		StagingCnameTo           types.String `tfsdk:"staging_cname_to"`
		StagingCnameType         types.String `tfsdk:"staging_cname_type"`
		StagingEdgeHostnameID    types.String `tfsdk:"staging_edge_hostname_id"`
		StagingProductID         types.String `tfsdk:"staging_product_id"`
	}
)

// NewAccountHostnamesDataSource returns a new accountHostnamesDataSource.
func NewAccountHostnamesDataSource() datasource.DataSource {
	return &accountHostnamesDataSource{}
}

// Metadata configures data source's meta information.
func (d *accountHostnamesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_property_account_hostnames"
}

// Schema is used to define data source's terraform schema.
func (d *accountHostnamesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List active property hostnames for an account.",
		Attributes: map[string]schema.Attribute{
			"contract_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Unique identifier for the contract.",
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Unique identifier for the group.",
			},
			"hostname": schema.StringAttribute{
				Optional:    true,
				Description: "Filter the results by `cnameFrom`. Supports wildcard matches with *.",
			},
			"cname_to": schema.StringAttribute{
				Optional:    true,
				Description: "Filter the results by edge hostname. Supports wildcard matches with *.",
			},
			"network": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Network of activated hostnames, either `STAGING` or `PRODUCTION`.",
				Validators: []validator.String{
					stringvalidator.OneOf("STAGING", "PRODUCTION"),
				},
			},
			"sort": schema.StringAttribute{
				Optional: true,
				Description: "Sort the results based on the cnameFrom value, either `hostname:a` for ascending or " +
					"`hostname:d` for descending order. The default is `hostname:a`.",
				Validators: []validator.String{
					stringvalidator.OneOf("hostname:a", "hostname:d"),
				},
			},
			"account_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifies the prevailing account under which you requested the data.",
			},
			"available_sort": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of available sort options.",
			},
			"current_sort": schema.StringAttribute{
				Computed:    true,
				Description: "The current sort order applied to the results.",
			},
			"default_sort": schema.StringAttribute{
				Computed:    true,
				Description: "The default sort order for the results.",
			},
			"hostnames": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of active property hostnames.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cname_from": schema.StringAttribute{
							Computed:    true,
							Description: "The hostname that your end users see, indicated by the Host header in end user requests.",
						},
						"contract_id": schema.StringAttribute{
							Computed:    true,
							Description: "Identifies the prevailing contract under which you requested the data.",
						},
						"group_id": schema.StringAttribute{
							Computed:    true,
							Description: "Identifies the prevailing group under which you requested the data.",
						},
						"latest_version": schema.Int64Attribute{
							Computed:    true,
							Description: "Specifies the most recent version of the property.",
						},
						"production_cert_type": schema.StringAttribute{
							Computed: true,
							Description: "Indicates the certificate's provisioning type. " +
								"Either `CPS_MANAGED` for the certificates created with the Certificate Provisioning System (CPS) API, " +
								"`CCM` for the certificates created with the Cloud Certificate Manager (CCM) API, " +
								"or `DEFAULT` for the Domain Validation (DV) certificates created automatically. " +
								"Note that you can't specify the `DEFAULT` value if your property hostname uses the `akamaized.net` domain suffix.",
						},
						"production_cname_to": schema.StringAttribute{
							Computed: true,
							Description: "The edge hostname you point the property hostname to so that you can start serving traffic through Akamai servers. " +
								"This member corresponds to the edge hostname object's edgeHostnameDomain member.",
						},
						"production_cname_type": schema.StringAttribute{
							Computed:    true,
							Description: "Indicates the type of CNAME you used in the production network, either `EDGE_HOSTNAME` or `CUSTOM`.",
						},
						"production_edge_hostname_id": schema.StringAttribute{
							Computed:    true,
							Description: "Identifies the edge hostname you mapped your traffic to on the production network.",
						},
						"production_product_id": schema.StringAttribute{
							Computed:    true,
							Description: "Identifies the product association on the network.",
						},
						"property_id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier for the property.",
						},
						"property_name": schema.StringAttribute{
							Computed:    true,
							Description: "A unique, descriptive name for the property.",
						},
						"property_type": schema.StringAttribute{
							Computed: true,
							Description: "Specifies the type of the property. Either `TRADITIONAL` for properties where you pair property hostnames " +
								"with the property version, or `HOSTNAME_BUCKET` where you manage property hostnames independently of the property version.",
						},
						"staging_cert_type": schema.StringAttribute{
							Computed: true,
							Description: "Indicates the certificate's provisioning type. " +
								"Either `CPS_MANAGED` for the certificates created with the Certificate Provisioning System (CPS) API, " +
								"`CCM` for the certificates created with the Cloud Certificate Manager (CCM) API, " +
								"or `DEFAULT` for the Domain Validation (DV) certificates created automatically. " +
								"Note that you can't specify the `DEFAULT` value if your property hostname uses the `akamaized.net` domain suffix.",
						},
						"staging_cname_to": schema.StringAttribute{
							Computed: true,
							Description: "The edge hostname you point the property hostname to so that you can start serving traffic through Akamai servers. " +
								"This member corresponds to the edge hostname object's edgeHostnameDomain member.",
						},
						"staging_cname_type": schema.StringAttribute{
							Computed:    true,
							Description: "Indicates the type of CNAME you used in the staging network, either `EDGE_HOSTNAME` or `CUSTOM`.",
						},
						"staging_edge_hostname_id": schema.StringAttribute{
							Computed:    true,
							Description: "Identifies the edge hostname you mapped your traffic to on the staging network.",
						},
						"staging_product_id": schema.StringAttribute{
							Computed:    true,
							Description: "Identifies the product association on the network.",
						},
					},
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state.
func (d *accountHostnamesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Account Hostnames DataSource Read")

	var data accountHostnamesDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	hostnamesResp, err := listAccountHostnames(ctx, d.Client.GetPAPI(), data)
	if err != nil {
		resp.Diagnostics.AddError("Read Property Account Hostnames failed", err.Error())
		return
	}

	resp.Diagnostics.Append(data.convertResponseToModel(ctx, hostnamesResp)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func listAccountHostnames(ctx context.Context, client papi.PAPI, data accountHostnamesDataSourceModel) (*papi.ListActiveAccountHostnamesResponse, error) {
	req := papi.ListActiveAccountHostnamesRequest{
		ContractID: data.ContractID.ValueString(),
		GroupID:    data.GroupID.ValueString(),
		Hostname:   data.Hostname.ValueString(),
		CnameTo:    data.CnameTo.ValueString(),
		Network:    papi.ActivationNetwork(data.Network.ValueString()),
		Sort:       papi.SortOrder(data.Sort.ValueString()),
	}

	resp, err := client.ListActiveAccountHostnames(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *accountHostnamesDataSourceModel) convertResponseToModel(ctx context.Context, response *papi.ListActiveAccountHostnamesResponse) diag.Diagnostics {
	m.AccountID = types.StringValue(response.AccountID)
	m.CurrentSort = types.StringValue(response.CurrentSort)
	m.DefaultSort = types.StringValue(response.DefaultSort)

	availableSortList, diags := types.ListValueFrom(ctx, types.StringType, response.AvailableSort)
	if diags.HasError() {
		return diags
	}
	m.AvailableSort = availableSortList

	m.Hostnames = make([]accountHostnameModel, 0, len(response.Hostnames.Items))
	for _, item := range response.Hostnames.Items {
		m.Hostnames = append(m.Hostnames, accountHostnameModel{
			CnameFrom:                types.StringValue(item.CnameFrom),
			ContractID:               types.StringValue(item.ContractID),
			GroupID:                  types.StringValue(item.GroupID),
			LatestVersion:            types.Int64Value(int64(item.LatestVersion)),
			ProductionCertType:       types.StringPointerValue(item.ProductionCertType),
			ProductionCnameTo:        types.StringPointerValue(item.ProductionCnameTo),
			ProductionCnameType:      types.StringPointerValue(item.ProductionCnameType),
			ProductionEdgeHostnameID: types.StringPointerValue(item.ProductionEdgeHostnameID),
			ProductionProductID:      types.StringPointerValue(item.ProductionProductID),
			PropertyID:               types.StringValue(item.PropertyID),
			PropertyName:             types.StringValue(item.PropertyName),
			PropertyType:             types.StringValue(item.PropertyType),
			StagingCertType:          types.StringPointerValue(item.StagingCertType),
			StagingCnameTo:           types.StringPointerValue(item.StagingCnameTo),
			StagingCnameType:         types.StringPointerValue(item.StagingCnameType),
			StagingEdgeHostnameID:    types.StringPointerValue(item.StagingEdgeHostnameID),
			StagingProductID:         types.StringPointerValue(item.StagingProductID),
		})
	}

	return nil
}
