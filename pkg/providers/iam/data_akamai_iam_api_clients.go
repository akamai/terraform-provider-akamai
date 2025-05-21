package iam

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &apiClientsDataSource{}
	_ datasource.DataSourceWithConfigure = &apiClientsDataSource{}
)

type (
	apiClientsDataSource struct {
		meta meta.Meta
	}

	apiClientActions struct {
		Delete        types.Bool `tfsdk:"delete"`
		DeactivateAll types.Bool `tfsdk:"deactivate_all"`
		Edit          types.Bool `tfsdk:"edit"`
		Lock          types.Bool `tfsdk:"lock"`
		Transfer      types.Bool `tfsdk:"transfer"`
		Unlock        types.Bool `tfsdk:"unlock"`
	}

	apiClientModel struct {
		AccessToken             types.String      `tfsdk:"access_token"`
		Actions                 *apiClientActions `tfsdk:"actions"`
		ActiveCredentialCount   types.Int64       `tfsdk:"active_credential_count"`
		AllowAccountSwitch      types.Bool        `tfsdk:"allow_account_switch"`
		AuthorizedUsers         []types.String    `tfsdk:"authorized_users"`
		CanAutoCreateCredential types.Bool        `tfsdk:"can_auto_create_credential"`
		ClientDescription       types.String      `tfsdk:"client_description"`
		ClientID                types.String      `tfsdk:"client_id"`
		ClientName              types.String      `tfsdk:"client_name"`
		ClientType              types.String      `tfsdk:"client_type"`
		CreatedBy               types.String      `tfsdk:"created_by"`
		CreatedDate             types.String      `tfsdk:"created_date"`
		IsLocked                types.Bool        `tfsdk:"is_locked"`
		NotificationEmails      []types.String    `tfsdk:"notification_emails"`
		ServiceConsumerToken    types.String      `tfsdk:"service_consumer_token"`
	}

	apiClientsSourceModel struct {
		APIClients []apiClientModel `tfsdk:"api_clients"`
	}
)

// NewAPIClientsDataSource returns a new iam API clients data source.
func NewAPIClientsDataSource() datasource.DataSource {
	return &apiClientsDataSource{}
}

func (d *apiClientsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_iam_api_clients"
}

func (d *apiClientsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unexpected Data Source Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.",
					req.ProviderData))
		}
	}()
	d.meta = meta.Must(req.ProviderData)
}

func (d *apiClientsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Identity and Access Management API clients.",
		Attributes: map[string]schema.Attribute{
			"api_clients": schema.ListNestedAttribute{
				Description: "List of API clients",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"access_token": schema.StringAttribute{
							Computed:    true,
							Description: "The part of the client secret that identifies your API client and lets you access applications and resources.",
						},
						"actions": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Specifies activities available for the API client.",
							Attributes: map[string]schema.Attribute{
								"delete": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether you can remove the API client.",
								},
								"deactivate_all": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether you can deactivate the API client's credentials.",
								},
								"edit": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether you can update the API client.",
								},
								"lock": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether you can lock the API client.",
								},
								"transfer": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether you can transfer the API client to a new owner.",
								},
								"unlock": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether you can unlock the API client.",
								},
							},
						},
						"active_credential_count": schema.Int64Attribute{
							Computed:    true,
							Description: "The number of credentials active for the API client.",
						},
						"allow_account_switch": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the API client can manage more than one account.",
						},
						"authorized_users": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "The API client's valid users.",
						},
						"can_auto_create_credential": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the API client can create a credential for a new API client.",
						},
						"client_description": schema.StringAttribute{
							Computed:    true,
							Description: "A human-readable description of the API client.",
						},
						"client_id": schema.StringAttribute{
							Computed:    true,
							Description: "A unique identifier for the API client.",
						},
						"client_name": schema.StringAttribute{
							Computed:    true,
							Description: "A human-readable name for the API client.",
						},
						"client_type": schema.StringAttribute{
							Computed:    true,
							Description: "Specifies the API client's ownership and credential management.",
						},
						"created_by": schema.StringAttribute{
							Computed:    true,
							Description: "The user who created the API client.",
						},
						"created_date": schema.StringAttribute{
							Computed:    true,
							Description: "The ISO 8601 timestamp indicating when the API client was created.",
						},
						"is_locked": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the API client is locked.",
						},
						"notification_emails": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "Email addresses to notify users when credentials expire.",
						},
						"service_consumer_token": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier for the service hostname.",
						},
					},
				},
			},
		},
	}
}

func (d *apiClientsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "IAM API Clients DataSource Read")

	client := inst.Client(d.meta)

	apiClients, err := client.ListAPIClients(ctx, iam.ListAPIClientsRequest{
		Actions: true,
	})
	if err != nil {
		resp.Diagnostics.AddError("IAM list API Clients failed", err.Error())
		return
	}

	var data apiClientsSourceModel
	data.read(apiClients)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *apiClientsSourceModel) read(apiClients iam.ListAPIClientsResponse) {
	for _, apiClient := range apiClients {
		authorizedUsers := make([]types.String, 0, len(apiClient.AuthorizedUsers))
		for _, authorizedUser := range apiClient.AuthorizedUsers {
			authorizedUsers = append(authorizedUsers, types.StringValue(authorizedUser))
		}

		notificationEmails := make([]types.String, 0, len(apiClient.NotificationEmails))
		for _, notificationEmail := range apiClient.NotificationEmails {
			notificationEmails = append(notificationEmails, types.StringValue(notificationEmail))
		}

		client := apiClientModel{
			AccessToken:             types.StringValue(apiClient.AccessToken),
			ActiveCredentialCount:   types.Int64Value(apiClient.ActiveCredentialCount),
			AllowAccountSwitch:      types.BoolValue(apiClient.AllowAccountSwitch),
			AuthorizedUsers:         authorizedUsers,
			CanAutoCreateCredential: types.BoolValue(apiClient.CanAutoCreateCredential),
			ClientDescription:       types.StringValue(apiClient.ClientDescription),
			ClientID:                types.StringValue(apiClient.ClientID),
			ClientName:              types.StringValue(apiClient.ClientName),
			ClientType:              types.StringValue(string(apiClient.ClientType)),
			CreatedBy:               types.StringValue(apiClient.CreatedBy),
			CreatedDate:             types.StringValue(date.FormatRFC3339Nano(apiClient.CreatedDate)),
			IsLocked:                types.BoolValue(apiClient.IsLocked),
			NotificationEmails:      notificationEmails,
			ServiceConsumerToken:    types.StringValue(apiClient.ServiceConsumerToken),
		}
		if apiClient.Actions != nil {
			client.Actions = &apiClientActions{
				Delete:        types.BoolValue(apiClient.Actions.Delete),
				DeactivateAll: types.BoolValue(apiClient.Actions.DeactivateAll),
				Edit:          types.BoolValue(apiClient.Actions.Edit),
				Lock:          types.BoolValue(apiClient.Actions.Lock),
				Transfer:      types.BoolValue(apiClient.Actions.Transfer),
				Unlock:        types.BoolValue(apiClient.Actions.Unlock),
			}
		}

		d.APIClients = append(d.APIClients, client)
	}
}
