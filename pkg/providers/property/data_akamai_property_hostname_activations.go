package property

import (
	"context"
	"fmt"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &hostnameActivationsDataSource{}
var _ datasource.DataSourceWithConfigure = &hostnameActivationsDataSource{}

// NewHostnameActivationsDataSource returns a new property hostname activations data source
func NewHostnameActivationsDataSource() datasource.DataSource {
	return &hostnameActivationsDataSource{}
}

// hostnameActivationsDataSource defines the data source implementation for fetching property hostname activations information.
type hostnameActivationsDataSource struct {
	meta meta.Meta
}

// hostnameActivationsDataSourceModel describes the data source data model for PropertyHostnameActivationsDataSource.
type hostnameActivationsDataSourceModel struct {
	PropertyID          types.String         `tfsdk:"property_id"`
	Network             types.String         `tfsdk:"network"`
	ContractID          types.String         `tfsdk:"contract_id"`
	GroupID             types.String         `tfsdk:"group_id"`
	AccountID           types.String         `tfsdk:"account_id"`
	PropertyName        types.String         `tfsdk:"property_name"`
	HostnameActivations []hostnameActivation `tfsdk:"hostname_activations"`
}

type hostnameActivation struct {
	ActivationType       types.String   `tfsdk:"activation_type"`
	Network              types.String   `tfsdk:"network"`
	Note                 types.String   `tfsdk:"note"`
	NotifyEmails         []types.String `tfsdk:"notify_emails"`
	Status               types.String   `tfsdk:"status"`
	SubmitDate           types.String   `tfsdk:"submit_date"`
	UpdateDate           types.String   `tfsdk:"update_date"`
	HostnameActivationID types.String   `tfsdk:"hostname_activation_id"`
}

// Metadata configures data source's meta information.
func (d *hostnameActivationsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_property_hostname_activations"
}

// Schema is used to define data source's terraform schema.
func (d *hostnameActivationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Property hostname activations data source",
		Attributes: map[string]schema.Attribute{
			"property_id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier for the property.",
			},
			"network": schema.StringAttribute{
				Optional:    true,
				Description: "The network of activation, either `STAGING`, `PRODUCTION`, or none.",
				Validators:  []validator.String{stringvalidator.OneOfCaseInsensitive("production", "staging", "")},
			},
			"contract_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier for the contract.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
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
			"property_name": schema.StringAttribute{
				Computed:    true,
				Description: "A descriptive name for the property with the hostname bucket the activated property hostnames belong to.",
			},
			"hostname_activations": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"activation_type": schema.StringAttribute{
							Computed:    true,
							Description: "The activation type, either `ACTIVATE` or `DEACTIVATE`.",
						},
						"network": schema.StringAttribute{
							Computed:    true,
							Description: "The network of activation, either `STAGING` or `PRODUCTION`.`",
						},
						"note": schema.StringAttribute{
							Computed:    true,
							Description: "Assigns a log message to the activation request.",
						},
						"notify_emails": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Email addresses to notify when the activation status changes.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "The activation's status. `ACTIVE` if currently serving traffic. `INACTIVE` if another activation has superseded this one. `PENDING` if not yet active. `ABORTED` if the client followed up with a `DELETE` request in time. `FAILED` if the activation causes a range of edge network errors that may cause a fallback to the previous activation. `PENDING_DEACTIVATION` or `DEACTIVATED` when the `activation_type` is `DEACTIVATE` to no longer serve traffic.",
						},
						"submit_date": schema.StringAttribute{
							Computed:    true,
							Description: "The timestamp indicating when the activation was initiated.",
						},
						"update_date": schema.StringAttribute{
							Computed:    true,
							Description: "The ISO 8601 timestamp indicating when the status last changed.",
						},
						"hostname_activation_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ISO 8601 timestamp property hostname activation's unique identifier.",
						},
					},
				},
			},
		},
	}
}

// Configure configures data source at the beginning of the lifecycle.
func (d *hostnameActivationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func getAllActivations(ctx context.Context, client papi.PAPI, contractID, groupID, propertyID, network string) (*papi.ListPropertyHostnameActivationsResponse, error) {
	pageSize, offset := 999, 0
	response := &papi.ListPropertyHostnameActivationsResponse{}
	for {
		act, err := client.ListPropertyHostnameActivations(ctx, papi.ListPropertyHostnameActivationsRequest{
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
		for _, item := range act.HostnameActivations.Items {
			if network != "" && !strings.EqualFold(network, string(item.Network)) {
				continue
			}
			response.HostnameActivations.Items = append(response.HostnameActivations.Items, item)
		}

		offset += pageSize
		if offset >= act.HostnameActivations.TotalItems {
			break
		}
	}

	return response, nil

}

// Read is called when the provider must read data source values in order to update state.
func (d *hostnameActivationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Property HostnameActivationsDataSource Read")

	var data hostnameActivationsDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	client := Client(d.meta)
	activations, err := getAllActivations(ctx, client, data.ContractID.ValueString(), data.GroupID.ValueString(), data.PropertyID.ValueString(), data.Network.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("fetching property hostname activations failed", err.Error())
		return
	}

	data.AccountID = types.StringValue(activations.AccountID)
	data.ContractID = types.StringValue(activations.ContractID)
	data.GroupID = types.StringValue(activations.GroupID)
	if len(activations.HostnameActivations.Items) > 0 {
		data.PropertyName = types.StringValue(activations.HostnameActivations.Items[0].PropertyName)
	} else {
		data.PropertyName = types.StringNull()
	}
	for _, item := range activations.HostnameActivations.Items {
		emails := make([]types.String, len(item.NotifyEmails))
		for i, email := range item.NotifyEmails {
			emails[i] = types.StringValue(email)
		}
		a := hostnameActivation{
			ActivationType:       types.StringValue(item.ActivationType),
			Network:              types.StringValue(string(item.Network)),
			Note:                 types.StringValue(item.Note),
			NotifyEmails:         emails,
			Status:               types.StringValue(item.Status),
			SubmitDate:           types.StringValue(date.FormatRFC3339Nano(item.SubmitDate)),
			UpdateDate:           types.StringValue(date.FormatRFC3339Nano(item.UpdateDate)),
			HostnameActivationID: types.StringValue(item.HostnameActivationID),
		}
		data.HostnameActivations = append(data.HostnameActivations, a)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
