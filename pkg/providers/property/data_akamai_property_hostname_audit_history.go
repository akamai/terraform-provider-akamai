package property

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &hostnameAuditHistoryDataSource{}
var _ datasource.DataSourceWithConfigure = &hostnameAuditHistoryDataSource{}

// NewHostnameAuditHistoryDataSource returns a new property hostname audit history data source.
func NewHostnameAuditHistoryDataSource() datasource.DataSource {
	return &hostnameAuditHistoryDataSource{}
}

type (
	// hostnameAuditHistoryDataSource defines the data source implementation for fetching property hostname audit history.
	hostnameAuditHistoryDataSource struct {
		meta.DataSource
	}

	hostnameAuditHistoryDataSourceModel struct {
		Hostname types.String                    `tfsdk:"hostname"`
		History  []hostnameAuditHistoryItemModel `tfsdk:"history"`
	}

	hostnameAuditHistoryItemModel struct {
		Action               types.String `tfsdk:"action"`
		CertProvisioningType types.String `tfsdk:"cert_provisioning_type"`
		CnameTo              types.String `tfsdk:"cname_to"`
		ContractID           types.String `tfsdk:"contract_id"`
		EdgeHostnameID       types.String `tfsdk:"edge_hostname_id"`
		GroupID              types.String `tfsdk:"group_id"`
		Network              types.String `tfsdk:"network"`
		PropertyID           types.String `tfsdk:"property_id"`
		Timestamp            types.String `tfsdk:"timestamp"`
		User                 types.String `tfsdk:"user"`
	}
)

// Metadata configures data source's meta information.
func (d *hostnameAuditHistoryDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_property_hostname_audit_history"
}

// Schema is used to define data source's terraform schema.
func (d *hostnameAuditHistoryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Property Hostname Audit History data source",
		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				Required:    true,
				Description: "The hostname to fetch audit history for. It is the cnameFrom for the hostname your end users see, indicated by the Host header in end user requests.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 253),
				},
			},
			"history": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Audit history of the hostname.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"action": schema.StringAttribute{
							Computed: true,
							MarkdownDescription: "The type of action performed to the property hostname. Possible values are: \n" +
								"* `ACTIVATE` - When the hostname is currently serving traffic.\n" +
								"* `DEACTIVATE` - When the hostname isn't serving traffic.\n" +
								"* `ADD` - When the user requested to add the hostname to a property.\n" +
								"* `REMOVE` - When the user requested to remove the hostname from a property.\n" +
								"* `MOVE` - When the hostname was moved from one property to another.\n" +
								"* `MODIFY` - When the user changed the edgeHostnameId or certProvisioningType values for an already-activated hostname.\n" +
								"* `ABORTED` - When the user request to cancel the hostname activation.\n" +
								"* `ERROR` - When the hostname activation failed.",
						},
						"cert_provisioning_type": schema.StringAttribute{
							Computed: true,
							MarkdownDescription: "The type of certificate used in the property hostname. Possible values are: \n" +
								"* `CPS_MANAGED` - For certificates you create with the Certificate Provisioning System API (CPS).\n" +
								"* `DEFAULT` - For Default Domain Validation (DV) certificates deployed automatically.\n" +
								"* `CCM` - For the third party certificates created with the Cloud Certificate Manager.",
						},
						"cname_to": schema.StringAttribute{
							Computed:    true,
							Description: "The edge hostname that the hostname points to.",
						},
						"contract_id": schema.StringAttribute{
							Computed:    true,
							Description: "Identifies the prevailing contract under which the data was requested.",
						},
						"edge_hostname_id": schema.StringAttribute{
							Computed:    true,
							Description: "Id of the edge hostname the hostname points to.",
						},
						"group_id": schema.StringAttribute{
							Computed:    true,
							Description: "Identifies the group under which the property is activated.",
						},
						"network": schema.StringAttribute{
							Computed: true,
							MarkdownDescription: "The network of activated hostnames. Possible values are: \n" +
								"* `STAGING` - Staging network.\n" +
								"* `PRODUCTION` - Production network.",
						},
						"property_id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier for the property.",
						},
						"timestamp": schema.StringAttribute{
							Computed:    true,
							Description: "Indicates when the action occurred.",
						},
						"user": schema.StringAttribute{
							Computed:    true,
							Description: "The user who initiated the action.",
						},
					},
				},
			},
		},
	}
}

func (d *hostnameAuditHistoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Hostname Audit History Data Source Read")

	var data hostnameAuditHistoryDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	historyResp, err := d.Client.GetPAPI().GetAuditHistory(ctx, papi.GetAuditHistoryRequest{
		Hostname: data.Hostname.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Read Hostname Audit History failed", err.Error())
		return
	}

	data.History = hostnameHistoryItemsToModels(historyResp.History.Items)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func hostnameHistoryItemsToModels(history []papi.HostnameHistoryItem) []hostnameAuditHistoryItemModel {
	var model []hostnameAuditHistoryItemModel
	for _, entry := range history {
		model = append(model, hostnameAuditHistoryItemModel{
			Action:               types.StringValue(entry.Action),
			CertProvisioningType: types.StringValue(entry.CertProvisioningType),
			CnameTo:              types.StringValue(entry.CnameTo),
			ContractID:           types.StringValue(entry.ContractID),
			EdgeHostnameID:       types.StringValue(entry.EdgeHostnameID),
			GroupID:              types.StringValue(entry.GroupID),
			Network:              types.StringValue(entry.Network),
			PropertyID:           types.StringValue(entry.PropertyID),
			Timestamp:            types.StringValue(entry.Timestamp),
			User:                 types.StringValue(entry.User),
		})
	}
	return model
}
