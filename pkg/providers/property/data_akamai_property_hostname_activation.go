package property

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &hostnameActivationDataSource{}
var _ datasource.DataSourceWithConfigure = &hostnameActivationDataSource{}

// NewHostnameActivationDataSource returns a new property hostname activation data source.
func NewHostnameActivationDataSource() datasource.DataSource {
	return &hostnameActivationDataSource{}
}

// hostnameActivationDataSource defines the data source implementation for fetching property hostname activation information.
type hostnameActivationDataSource struct {
	meta meta.Meta
}

type hostnameModel struct {
	EdgeHostnameID       types.String `tfsdk:"edge_hostname_id"`
	CnameFrom            types.String `tfsdk:"cname_from"`
	CnameTo              types.String `tfsdk:"cname_to"`
	CertProvisioningType types.String `tfsdk:"cert_provisioning_type"`
	Action               types.String `tfsdk:"action"`
}

type hostnameActivationDataSourceModel struct {
	PropertyID           types.String    `tfsdk:"property_id"`
	HostnameActivationID types.String    `tfsdk:"hostname_activation_id"`
	ContractID           types.String    `tfsdk:"contract_id"`
	GroupID              types.String    `tfsdk:"group_id"`
	IncludeHostnames     types.Bool      `tfsdk:"include_hostnames"`
	AccountID            types.String    `tfsdk:"account_id"`
	ActivationType       types.String    `tfsdk:"activation_type"`
	Network              types.String    `tfsdk:"network"`
	Note                 types.String    `tfsdk:"note"`
	NotifyEmails         types.List      `tfsdk:"notify_emails"`
	PropertyName         types.String    `tfsdk:"property_name"`
	Status               types.String    `tfsdk:"status"`
	SubmitDate           types.String    `tfsdk:"submit_date"`
	UpdateDate           types.String    `tfsdk:"update_date"`
	Hostnames            []hostnameModel `tfsdk:"hostnames"`
}

func (p *hostnameActivationDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_property_hostname_activation"
}

func (p *hostnameActivationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	p.meta = meta.Must(req.ProviderData)
}

func (p *hostnameActivationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Property hostname activation data source",
		Attributes: map[string]schema.Attribute{
			"property_id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier for the property.",
			},
			"hostname_activation_id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier for the hostname activation.",
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
			"include_hostnames": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether the response should include the property hostnames associated with an activation and the related certificate status on staging and production networks.",
			},
			"account_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifies the prevailing account under which you requested the data.",
			},
			"activation_type": schema.StringAttribute{
				Computed:    true,
				Description: "The activation type, either `ACTIVATE` or `DEACTIVATE`.",
			},
			"network": schema.StringAttribute{
				Computed:    true,
				Description: "The network of activation, either `STAGING` or `PRODUCTION`.",
			},
			"note": schema.StringAttribute{
				Computed:    true,
				Description: "Assigns a log message to the activation request.",
			},
			"notify_emails": schema.ListAttribute{
				Computed:    true,
				Description: "Email addresses to notify when the activation status changes.",
				ElementType: types.StringType,
			},
			"property_name": schema.StringAttribute{
				Computed:    true,
				Description: "A descriptive name for the property with the hostname bucket the activated property hostnames belong to.",
			},
			"status": schema.StringAttribute{
				Computed: true,
				Description: "The activation's status. `ACTIVE` if currently serving traffic. `INACTIVE` if another activation has superseded this one. " +
					"`PENDING` if not yet active. `ABORTED` if the client followed up with a `DELETE` request in time. " +
					"`FAILED` if the activation causes a range of edge network errors that may cause a fallback to the previous activation. " +
					"`PENDING_DEACTIVATION` or `DEACTIVATED` when the `activation_type` is `DEACTIVATE` to no longer serve traffic.",
			},
			"submit_date": schema.StringAttribute{
				Computed:    true,
				Description: "The ISO 8601 timestamp indicating when the activation was initiated.",
			},
			"update_date": schema.StringAttribute{
				Computed:    true,
				Description: "The ISO 8601 timestamp indicating when the `status` last changed.",
			},
			"hostnames": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The set of hostnames activated within a given property activation on the staging and production networks.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"edge_hostname_id": schema.StringAttribute{
							Computed:    true,
							Description: "Identifies the edge hostname you mapped your traffic to on the production network.",
						},
						"cname_from": schema.StringAttribute{
							Computed:    true,
							Description: "The hostname that your end users see, indicated by the Host header in end user requests.",
						},
						"cname_to": schema.StringAttribute{
							Computed:    true,
							Description: "The edge hostname you point the property hostname to so that you can start serving traffic through Akamai servers.",
						},
						"cert_provisioning_type": schema.StringAttribute{
							Computed: true,
							Description: "Indicates the certificate's provisioning type. Either `CPS_MANAGED` for the certificates you create with the Certificate Provisioning System (CPS) API, " +
								"or `DEFAULT` for the Domain Validation (DV) certificates created automatically. " +
								"Note that you can't specify the `DEFAULT` value if your property hostname uses the `akamaized.net` domain suffix.",
						},
						"action": schema.StringAttribute{
							Computed:    true,
							Description: "Specifies whether a given activation adds or removes a hostname item. Available options are `ADD` and `REMOVE`.",
						},
					},
				},
			},
		},
	}
}

func (p *hostnameActivationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Property HostnameActivationDataSource Read")
	var diags diag.Diagnostics
	var data hostnameActivationDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	client := Client(p.meta)
	activation, err := client.GetPropertyHostnameActivation(ctx, papi.GetPropertyHostnameActivationRequest{
		PropertyID:           data.PropertyID.ValueString(),
		HostnameActivationID: data.HostnameActivationID.ValueString(),
		ContractID:           data.ContractID.ValueString(),
		GroupID:              data.GroupID.ValueString(),
		IncludeHostnames:     data.IncludeHostnames.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError("fetching hostname activation failed", err.Error())
		return
	}

	if diags = data.assignActivationToModel(ctx, activation); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (m *hostnameActivationDataSourceModel) assignActivationToModel(ctx context.Context, activation *papi.GetPropertyHostnameActivationResponse) diag.Diagnostics {
	var diags diag.Diagnostics
	m.HostnameActivationID = types.StringValue(activation.HostnameActivation.HostnameActivationID)
	m.GroupID = types.StringValue(activation.GroupID)
	m.ContractID = types.StringValue(activation.ContractID)
	m.AccountID = types.StringValue(activation.AccountID)
	m.ActivationType = types.StringValue(activation.HostnameActivation.ActivationType)
	m.Network = types.StringValue(string(activation.HostnameActivation.Network))
	m.Note = types.StringValue(activation.HostnameActivation.Note)
	m.PropertyName = types.StringValue(activation.HostnameActivation.PropertyName)
	m.Status = types.StringValue(activation.HostnameActivation.Status)
	m.SubmitDate = types.StringValue(date.FormatRFC3339Nano(activation.HostnameActivation.SubmitDate))
	m.UpdateDate = types.StringValue(date.FormatRFC3339Nano(activation.HostnameActivation.UpdateDate))
	notifyEmails, diags := types.ListValueFrom(ctx, types.StringType, activation.HostnameActivation.NotifyEmails)
	if diags.HasError() {
		return diags
	}
	m.NotifyEmails = notifyEmails
	for _, h := range activation.HostnameActivation.Hostnames {
		var hostnameData hostnameModel
		hostnameData.EdgeHostnameID = types.StringValue(h.EdgeHostnameID)
		hostnameData.CnameTo = types.StringValue(h.CnameTo)
		hostnameData.CnameFrom = types.StringValue(h.CnameFrom)
		hostnameData.CertProvisioningType = types.StringValue(string(h.CertProvisioningType))
		hostnameData.Action = types.StringValue(h.Action)
		m.Hostnames = append(m.Hostnames, hostnameData)
	}
	return diags
}
