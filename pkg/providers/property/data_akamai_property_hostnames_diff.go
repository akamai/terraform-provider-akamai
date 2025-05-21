package property

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &hostnamesDiffDataSource{}
var _ datasource.DataSourceWithConfigure = &hostnamesDiffDataSource{}

// NewHostnamesDiffDataSource returns a new property hostnames diff data source.
func NewHostnamesDiffDataSource() datasource.DataSource {
	return &hostnamesDiffDataSource{}
}

// hostnamesDiffDataSource defines the data source implementation for fetching property hostnames diff information.
type hostnamesDiffDataSource struct {
	meta meta.Meta
}

// hostnamesDiffDataSourceModel describes the data source data model for PropertyHostnamesDiffDataSource.
type hostnamesDiffDataSourceModel struct {
	PropertyID types.String    `tfsdk:"property_id"`
	ContractID types.String    `tfsdk:"contract_id"`
	GroupID    types.String    `tfsdk:"group_id"`
	AccountID  types.String    `tfsdk:"account_id"`
	Hostnames  []hostnamesDiff `tfsdk:"hostnames"`
}

type hostnamesDiff struct {
	CNameFrom                      types.String `tfsdk:"cname_from"`
	StagingCNameType               types.String `tfsdk:"staging_cname_type"`
	ProductionCNameType            types.String `tfsdk:"production_cname_type"`
	StagingEdgeHostnameID          types.String `tfsdk:"staging_edge_hostname_id"`
	ProductionEdgeHostnameID       types.String `tfsdk:"production_edge_hostname_id"`
	StagingCertProvisioningType    types.String `tfsdk:"staging_cert_provisioning_type"`
	ProductionCertProvisioningType types.String `tfsdk:"production_cert_provisioning_type"`
	StagingCNameTo                 types.String `tfsdk:"staging_cname_to"`
	ProductionCNameTo              types.String `tfsdk:"production_cname_to"`
}

// Metadata configures data source's meta information.
func (d *hostnamesDiffDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_property_hostnames_diff"
}

// Schema is used to define data source's terraform schema.
func (d *hostnamesDiffDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Property hostnames diff data source",
		Attributes: map[string]schema.Attribute{
			"property_id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier for the property.",
			},
			"contract_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier for the contract.",
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier for the group.",
			},
			"account_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifies the prevailing account under which you requested the data.",
			},
			"hostnames": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cname_from": schema.StringAttribute{
							Computed:    true,
							Description: "The hostname that your end users see, indicated by the Host header in end user requests.",
						},
						"staging_cname_type": schema.StringAttribute{
							Computed:    true,
							Description: "A hostname's CNAME type. Supports only the `EDGE_HOSTNAME` value.",
						},
						"production_cname_type": schema.StringAttribute{
							Computed:    true,
							Description: "A hostname's CNAME type. Supports only the `EDGE_HOSTNAME` value.",
						},
						"staging_edge_hostname_id": schema.StringAttribute{
							Computed:    true,
							Description: "The unique identifier for the edge hostname.",
						},
						"production_edge_hostname_id": schema.StringAttribute{
							Computed:    true,
							Description: "The unique identifier for the edge hostname.",
						},
						"staging_cert_provisioning_type": schema.StringAttribute{
							Computed:    true,
							Description: "Indicates the certificate's provisioning type. Either `CPS_MANAGED` for the certificates you create with the Certificate Provisioning System (CPS) API, or `DEFAULT` for the Domain Validation (DV) certificates created automatically. Note that you can't specify the `DEFAULT` value if your property hostname uses the `akamaized.net` domain suffix.",
						},
						"production_cert_provisioning_type": schema.StringAttribute{
							Computed:    true,
							Description: "Indicates the certificate's provisioning type. Either `CPS_MANAGED` for the certificates you create with the Certificate Provisioning System (CPS) API, or `DEFAULT` for the Domain Validation (DV) certificates created automatically. Note that you can't specify the `DEFAULT` value if your property hostname uses the `akamaized.net` domain suffix.",
						},
						"staging_cname_to": schema.StringAttribute{
							Computed:    true,
							Description: "The edge hostname you point the property hostname to so that you can start serving traffic through Akamai servers. This member corresponds to the edge hostname object's `edgeHostnameDomain` member.",
						},
						"production_cname_to": schema.StringAttribute{
							Computed:    true,
							Description: "The edge hostname you point the property hostname to so that you can start serving traffic through Akamai servers. This member corresponds to the edge hostname object's `edgeHostnameDomain` member.",
						},
					},
				},
			},
		},
	}
}

// Configure  configures data source at the beginning of the lifecycle.
func (d *hostnamesDiffDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		// ProviderData is nil when Configure is run first time as part of ValidateDataSourceConfig in framework provider
		return
	}

	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unexpected Data Source Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
			)
		}
	}()

	d.meta = meta.Must(req.ProviderData)
}

func getHostnamesDiff(ctx context.Context, client papi.PAPI, contractID, groupID, propertyID string) (*papi.GetActivePropertyHostnamesDiffResponse, error) {
	pageSize, offset := 999, 0
	response := &papi.GetActivePropertyHostnamesDiffResponse{}

	for {
		act, err := client.GetActivePropertyHostnamesDiff(ctx, papi.GetActivePropertyHostnamesDiffRequest{
			PropertyID: propertyID,
			Offset:     offset,
			Limit:      pageSize,
			ContractID: contractID,
			GroupID:    groupID,
		})
		if err != nil {
			return nil, err
		}

		response.AccountID = act.AccountID
		response.ContractID = act.ContractID
		response.GroupID = act.GroupID
		response.PropertyID = act.PropertyID
		response.Hostnames.Items = append(response.Hostnames.Items, act.Hostnames.Items...)

		offset += pageSize
		if offset >= act.Hostnames.TotalItems {
			break
		}
	}

	return response, nil

}

// Read is called when the provider must read data source values in order to update state.
func (d *hostnamesDiffDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "PropertyHostnamesDiffDataSource Read")

	var data hostnamesDiffDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	client := Client(d.meta)
	activations, err := getHostnamesDiff(ctx, client, data.ContractID.ValueString(), data.GroupID.ValueString(), data.PropertyID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("fetching property hostnames diff failed", err.Error())
		return
	}

	data.AccountID = types.StringValue(activations.AccountID)
	data.ContractID = types.StringValue(activations.ContractID)
	data.GroupID = types.StringValue(activations.GroupID)

	for _, item := range activations.Hostnames.Items {
		a := hostnamesDiff{
			CNameFrom:                      types.StringValue(item.CnameFrom),
			StagingCNameType:               types.StringValue(string(item.StagingCnameType)),
			ProductionCNameType:            types.StringValue(string(item.ProductionCnameType)),
			StagingEdgeHostnameID:          types.StringValue(item.StagingEdgeHostnameID),
			ProductionEdgeHostnameID:       types.StringValue(item.ProductionEdgeHostnameID),
			StagingCertProvisioningType:    types.StringValue(string(item.StagingCertProvisioningType)),
			ProductionCertProvisioningType: types.StringValue(string(item.ProductionCertProvisioningType)),
			StagingCNameTo:                 types.StringValue(item.StagingCnameTo),
			ProductionCNameTo:              types.StringValue(item.ProductionCnameTo),
		}
		data.Hostnames = append(data.Hostnames, a)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
