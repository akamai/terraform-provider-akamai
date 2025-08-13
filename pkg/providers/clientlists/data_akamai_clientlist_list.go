package clientlists

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/clientlists"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	clientListDataSource struct {
		meta meta.Meta
	}

	// clientListDataSourceModel describes the data source data model for ClientListDataSource.
	clientListDataSourceModel struct {
		ListID     types.String `tfsdk:"list_id"`
		List       *listModel   `tfsdk:"list"`
		JSON       types.String `tfsdk:"json"`
		OutputText types.String `tfsdk:"output_text"`
	}

	listModel struct {
		clientListModel
		Items []clientListItemModel `tfsdk:"items"`
	}

	clientListItemModel struct {
		Value            types.String `tfsdk:"value"`
		Tags             types.List   `tfsdk:"tags"`
		Description      types.String `tfsdk:"description"`
		ExpirationDate   types.String `tfsdk:"expiration_date"`
		CreateDate       types.String `tfsdk:"create_date"`
		CreatedBy        types.String `tfsdk:"created_by"`
		CreatedVersion   types.Int64  `tfsdk:"created_version"`
		ProductionStatus types.String `tfsdk:"production_activation_status"`
		StagingStatus    types.String `tfsdk:"staging_activation_status"`
		Type             types.String `tfsdk:"type"`
		UpdateDate       types.String `tfsdk:"update_date"`
		UpdatedBy        types.String `tfsdk:"updated_by"`
	}
)

var (
	_ datasource.DataSource              = &clientListDataSource{}
	_ datasource.DataSourceWithConfigure = &clientListDataSource{}
)

// NewClientListDataSource returns a new client list data source
func NewClientListDataSource() datasource.DataSource { return &clientListDataSource{} }

// Metadata configures data source's meta information
func (d *clientListDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "akamai_clientlist_list"
}

// Schema is used to define data source's terraform schema
func (d *clientListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Client lists data source.",
		Attributes: map[string]schema.Attribute{
			"list_id": schema.StringAttribute{
				Required:    true,
				Description: "A client list id.",
			},
			"list": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "A client list.",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "The name of the client list.",
					},
					"type": schema.StringAttribute{
						Computed:    true,
						Description: "The type of the client list.",
					},
					"notes": schema.StringAttribute{
						Computed:    true,
						Description: "The client list notes.",
					},
					"tags": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "The client list tags.",
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
					"items": schema.ListNestedAttribute{
						Computed:    true,
						Description: "A set of client list values.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"create_date": schema.StringAttribute{
									Computed:    true,
									Description: "The client list item creation date.",
								},
								"created_by": schema.StringAttribute{
									Computed:    true,
									Description: "The username of the person who created the client list item.",
								},
								"created_version": schema.Int64Attribute{
									Computed:    true,
									Description: "The version of the client list when item was created.",
								},
								"update_date": schema.StringAttribute{
									Computed:    true,
									Description: "The date of last update.",
								},
								"updated_by": schema.StringAttribute{
									Computed:    true,
									Description: "The username of the person that updated the client list item last.",
								},
								"description": schema.StringAttribute{
									Optional:    true,
									Description: "The description of the client list item.",
								},
								"expiration_date": schema.StringAttribute{
									Computed:    true,
									Description: "The client list item expiration date.",
								},
								"production_activation_status": schema.StringAttribute{
									Computed:    true,
									Description: "The client list activation status in production environment.",
								},
								"staging_activation_status": schema.StringAttribute{
									Computed:    true,
									Description: "The client list activation status in staging environment.",
								},
								"type": schema.StringAttribute{
									Computed:    true,
									Description: "Type of client list, which can be IP, GEO, ASN, TLS_FINGERPRINT, FILE_HASH, or USER.",
								},
								"value": schema.StringAttribute{
									Computed:    true,
									Description: "Value of the item, which is either an IP address, an Autonomous System Number (ASN), a Geo location, a TLS fingerprint, a file hash, or User ID.",
								},
								"tags": schema.ListAttribute{
									ElementType: types.StringType,
									Computed:    true,
									Description: "A list of tags associated with the client list item.",
								},
							},
						},
					},
				},
			},
			"json": schema.StringAttribute{
				Computed:    true,
				Description: "JSON-formatted information about the client list.",
			},
			"output_text": schema.StringAttribute{
				Computed:    true,
				Description: "Tabular representation of the client lists.",
			},
		},
	}
}

func (d *clientListDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "Configuring Client List data source")

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

func (d *clientListDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Client List data source")

	var data clientListDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	client := inst.Client(d.meta)

	getClientListReq := clientlists.GetClientListRequest{
		ListID:       data.ListID.ValueString(),
		IncludeItems: true,
	}

	cl, err := client.GetClientList(ctx, getClientListReq)
	if err != nil {
		tflog.Error(ctx, "calling 'getClientList' failed", map[string]any{
			"error": err.Error(),
		})
		response.Diagnostics.AddError("get client list error", err.Error())
		return
	}

	if cl.Type == clientlists.USER {
		getClientListItems := clientlists.GetClientListItemsRequest{
			ListID: cl.ListID,
		}
		items, err := client.GetClientListItems(ctx, getClientListItems)
		if err != nil {
			tflog.Error(ctx, "calling 'getClientListItems' failed", map[string]any{
				"error": err.Error(),
			})
			response.Diagnostics.AddError("get client list items error", err.Error())
			return
		} else if len(items.Items) > 0 {
			cl.Items = processListItemContent(items.Items)
		}
	}

	tags := make([]types.String, 0, len(cl.Tags))
	for _, tag := range cl.Tags {
		tags = append(tags, types.StringValue(tag))
	}

	clientList := listModel{
		clientListModel: clientListModel{
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
		},
	}

	items, diags := convertListItemContentModel(ctx, cl.Items)
	if diags.HasError() {
		response.Diagnostics.AddError("Error converting list items to model", err.Error())
		return
	}
	clientList.Items = items

	jsonBody, err := json.MarshalIndent(cl, "", "  ")
	if err != nil {
		response.Diagnostics.AddError("Error marshaling JSON", err.Error())
		return
	}
	data.JSON = types.StringValue(string(jsonBody))

	ots := OutputTemplates{}
	InitTemplates(ots)
	outputTextList, err := RenderTemplates(ots, "clientListDS", []clientlists.GetClientListResponse{*cl})
	if err != nil {
		response.Diagnostics.AddError("Error rendering output text", err.Error())
		return
	}

	clientListItemsTemplateName := getClientListItemsTemplateName(cl.Type)
	outputTextItems, err := RenderTemplates(ots, clientListItemsTemplateName, cl)
	if err != nil {
		response.Diagnostics.AddError("Error rendering output text", err.Error())
		return
	}

	data.OutputText = types.StringValue(outputTextList + outputTextItems)
	data.List = &clientList
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func convertListItemContentModel(ctx context.Context, src []clientlists.ListItemContent) ([]clientListItemModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	result := make([]clientListItemModel, 0, len(src))
	for _, item := range src {
		tags, diagnostics := basetypes.NewListValueFrom(ctx, types.StringType, item.Tags)
		if diagnostics.HasError() {
			diags.Append(diagnostics...)
			return nil, diags
		}
		itemModel := clientListItemModel{
			Value:            types.StringValue(calculateValue(item)),
			Tags:             tags,
			Description:      types.StringValue(item.Description),
			ExpirationDate:   types.StringValue(item.ExpirationDate),
			CreateDate:       types.StringValue(item.CreateDate),
			CreatedBy:        types.StringValue(item.CreatedBy),
			CreatedVersion:   types.Int64Value(item.CreatedVersion),
			ProductionStatus: types.StringValue(item.ProductionStatus),
			StagingStatus:    types.StringValue(item.StagingStatus),
			Type:             types.StringValue(string(item.Type)),
			UpdateDate:       types.StringValue(item.UpdateDate),
			UpdatedBy:        types.StringValue(item.UpdatedBy),
		}
		result = append(result, itemModel)
	}
	return result, diags
}

func processListItemContent(src []clientlists.ListItemContent) []clientlists.ListItemContent {
	result := make([]clientlists.ListItemContent, 0, len(src))
	for _, item := range src {
		itemModel := clientlists.ListItemContent{
			Value:            calculateValue(item),
			Tags:             item.Tags,
			Description:      item.Description,
			ExpirationDate:   item.ExpirationDate,
			CreateDate:       item.CreateDate,
			CreatedBy:        item.CreatedBy,
			CreatedVersion:   item.CreatedVersion,
			ProductionStatus: item.ProductionStatus,
			StagingStatus:    item.StagingStatus,
			Type:             item.Type,
			UpdateDate:       item.UpdateDate,
			UpdatedBy:        item.UpdatedBy,
		}
		result = append(result, itemModel)
	}
	return result
}

func calculateValue(item clientlists.ListItemContent) string {
	if item.Type == clientlists.USER && item.Username != "" {
		return fmt.Sprintf("%s (%s)", item.Value, item.Username)
	}
	return item.Value
}

func getClientListItemsTemplateName(listType clientlists.ClientListType) string {
	switch listType {
	case clientlists.USER:
		return "userClientListItemsDS"
	case clientlists.IP:
		return "ipClientListItemsDS"
	case clientlists.ASN:
		return "asnClientListItemsDS"
	case clientlists.GEO:
		return "geoClientListItemsDS"
	case clientlists.TLSFingerprint:
		return "tlsFingerprintClientListItemsDS"
	case clientlists.FileHash:
		return "fileHashClientListItemsDS"
	default:
		return "unknownClientListItemsDS" // fallback or handle error
	}
}
