package clientlists

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/clientlists"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/hash"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	clientListsDataSource struct {
		meta meta.Meta
	}

	// clientListsDataSourceModel describes the data source data model for ClientListsDataSource.
	clientListsDataSourceModel struct {
		ID          types.String      `tfsdk:"id"`
		Name        types.String      `tfsdk:"name"`
		TypeFilters []types.String    `tfsdk:"type"`
		ListIDs     []types.String    `tfsdk:"list_ids"`
		Lists       []clientListModel `tfsdk:"lists"`
		JSON        types.String      `tfsdk:"json"`
		OutputText  types.String      `tfsdk:"output_text"`
	}

	clientListModel struct {
		Name                       types.String   `tfsdk:"name"`
		Type                       types.String   `tfsdk:"type"`
		Notes                      types.String   `tfsdk:"notes"`
		Tags                       []types.String `tfsdk:"tags"`
		ListID                     types.String   `tfsdk:"list_id"`
		Version                    types.Int64    `tfsdk:"version"`
		ItemsCount                 types.Int64    `tfsdk:"items_count"`
		CreateDate                 types.String   `tfsdk:"create_date"`
		CreatedBy                  types.String   `tfsdk:"created_by"`
		UpdateDate                 types.String   `tfsdk:"update_date"`
		UpdatedBy                  types.String   `tfsdk:"updated_by"`
		ProductionActivationStatus types.String   `tfsdk:"production_activation_status"`
		StagingActivationStatus    types.String   `tfsdk:"staging_activation_status"`
		ListType                   types.String   `tfsdk:"list_type"`
		Shared                     types.Bool     `tfsdk:"shared"`
		ReadOnly                   types.Bool     `tfsdk:"read_only"`
		Deprecated                 types.Bool     `tfsdk:"deprecated"`
	}
)

var (
	_ datasource.DataSource              = &clientListsDataSource{}
	_ datasource.DataSourceWithConfigure = &clientListsDataSource{}
)

// NewClientListsDataSource returns a new client lists data source
func NewClientListsDataSource() datasource.DataSource { return &clientListsDataSource{} }

// Metadata configures data source's meta information
func (d *clientListsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "akamai_clientlist_lists"
}

// Schema is used to define data source's terraform schema
func (d *clientListsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Client lists data source.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Filter client lists by name",
			},
			"type": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Filter client lists by type. Valid values: IP, GEO, ASN, TLS_FINGERPRINT, FILE_HASH.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf(getValidListTypes()...)),
				},
			},
			"list_ids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A set of client list ids.",
			},
			"lists": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A set of client lists.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the client list",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the client list",
						},
						"notes": schema.StringAttribute{
							Computed:    true,
							Description: "The client list notes",
						},
						"tags": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "The client list tags",
						},
						"list_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the client list.",
						},
						"version": schema.Int64Attribute{
							Computed:    true,
							Description: "The current version of the client list.",
						},
						"items_count": schema.Int64Attribute{
							Computed:    true,
							Description: "The number of items that a client list contains.",
						},
						"create_date": schema.StringAttribute{
							Computed:    true,
							Description: "The client list creation date.",
						},
						"created_by": schema.StringAttribute{
							Computed:    true,
							Description: "The username of the user who created the client list.",
						},
						"update_date": schema.StringAttribute{
							Computed:    true,
							Description: "The date of last update.",
						},
						"updated_by": schema.StringAttribute{
							Computed:    true,
							Description: "The username of the user that updated the client list last.",
						},
						"production_activation_status": schema.StringAttribute{
							Computed:    true,
							Description: "The activation status in production environment.",
						},
						"staging_activation_status": schema.StringAttribute{
							Computed:    true,
							Description: "The activation status in staging environment.",
						},
						"list_type": schema.StringAttribute{
							Computed:    true,
							Description: "The client list type.",
						},
						"shared": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the client list is shared.",
						},
						"read_only": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the client is editable for the authenticated user.",
						},
						"deprecated": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the client list was removed.",
						},
					},
				},
			},
			"json": schema.StringAttribute{
				Computed:    true,
				Description: "JSON representation of the client lists.",
			},
			"output_text": schema.StringAttribute{
				Computed:    true,
				Description: "Tabular representation of the client lists.",
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the data source",
				Computed:            true,
			},
		},
	}
}

func (d *clientListsDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "Configuring Client Lists data source")

	if request.ProviderData == nil {
		return
	}

	meta, ok := request.ProviderData.(meta.Meta)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)
		return
	}

	d.meta = meta
}

func (d *clientListsDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Client Lists data source")

	var data clientListsDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	client := inst.Client(d.meta)
	logger := d.meta.Log("CLIENTLIST", "dataSourceClientListRead")

	name := data.Name.ValueString()

	var listTypes []clientlists.ClientListType
	if len(data.TypeFilters) > 0 {
		for _, t := range data.TypeFilters {
			listTypes = append(listTypes, clientlists.ClientListType(t.ValueString()))
		}
	}

	lists, err := client.GetClientLists(ctx, clientlists.GetClientListsRequest{
		Name: name,
		Type: listTypes,
	})
	if err != nil {
		logger.Errorf("calling 'GetClientLists': %s", err.Error())
		response.Diagnostics.AddError("get client lists error", err.Error())
		return
	}

	listIDs := make([]types.String, 0, len(lists.Content))
	clientLists := make([]clientListModel, 0, len(lists.Content))
	for _, cl := range lists.Content {
		tags := make([]types.String, 0, len(cl.Tags))
		for _, tag := range cl.Tags {
			tags = append(tags, types.StringValue(tag))
		}

		clientList := clientListModel{
			Name:                       types.StringValue(cl.Name),
			Type:                       types.StringValue(string(cl.Type)),
			Notes:                      types.StringValue(cl.Notes),
			Tags:                       tags,
			ListID:                     types.StringValue(cl.ListID),
			Version:                    types.Int64Value(cl.Version),
			ItemsCount:                 types.Int64Value(cl.ItemsCount),
			CreateDate:                 types.StringValue(cl.CreateDate),
			CreatedBy:                  types.StringValue(cl.CreatedBy),
			UpdateDate:                 types.StringValue(cl.UpdateDate),
			UpdatedBy:                  types.StringValue(cl.UpdatedBy),
			ProductionActivationStatus: types.StringValue(cl.ProductionActivationStatus),
			StagingActivationStatus:    types.StringValue(cl.StagingActivationStatus),
			ListType:                   types.StringValue(cl.ListType),
			Shared:                     types.BoolValue(cl.Shared),
			ReadOnly:                   types.BoolValue(cl.ReadOnly),
			Deprecated:                 types.BoolValue(cl.Deprecated),
		}
		clientLists = append(clientLists, clientList)
		listIDs = append(listIDs, types.StringValue(cl.ListID))
	}
	data.Lists = clientLists
	data.ListIDs = listIDs

	jsonBody, err := json.MarshalIndent(lists.Content, "", "  ")
	if err != nil {
		response.Diagnostics.AddError("Error marshaling JSON", err.Error())
		return
	}
	data.JSON = types.StringValue(string(jsonBody))

	ots := OutputTemplates{}
	InitTemplates(ots)
	outputText, err := RenderTemplates(ots, "clientListsDS", lists)
	if err != nil {
		response.Diagnostics.AddError("Error rendering output text", err.Error())
		return
	}
	data.OutputText = types.StringValue(outputText)
	data.ID = types.StringValue(hash.GetSHAString(string(jsonBody)))
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func getValidListTypes() []string {
	return []string{
		string(clientlists.IP),
		string(clientlists.GEO),
		string(clientlists.ASN),
		string(clientlists.TLSFingerprint),
		string(clientlists.FileHash),
		string(clientlists.USER),
	}
}
