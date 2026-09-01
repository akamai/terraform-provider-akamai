package clientlists

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/clientlists"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	clientListDataSource struct {
		meta.DataSource
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
		Key              types.String `tfsdk:"key"`
		Values           types.List   `tfsdk:"values"`
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
						Description: "Type of client list, which can be IP, GEO, ASN, TLS_FINGERPRINT, FILE_HASH, USER_ID, DOMAIN, or REQUEST_HEADER_NAME_VALUE.",
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
									Description: "Type of client list, which can be IP, GEO, ASN, TLS_FINGERPRINT, FILE_HASH, USER_ID, DOMAIN, or REQUEST_HEADER_NAME_VALUE.",
								},
								"value": schema.StringAttribute{
									Computed:    true,
									Description: "Value of the item (e.g. IP address, AS Number, GEO, domain, TLS fingerprint, file hash, user ID). Not applicable for REQUEST_HEADER_NAME_VALUE list type.",
								},
								"key": schema.StringAttribute{
									Computed:    true,
									Description: "Key of the item (e.g. request header name). Applicable only for REQUEST_HEADER_NAME_VALUE list type.",
								},
								"values": schema.ListAttribute{
									ElementType: types.StringType,
									Computed:    true,
									Description: "Values of the item (e.g. request header name values). Applicable only for REQUEST_HEADER_NAME_VALUE list type.",
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

func (d *clientListDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Client List data source")

	var data clientListDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	client := d.Client.GetClientLists()

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
		response.Diagnostics.Append(diags...)
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
	result := make([]clientListItemModel, 0, len(src))
	for _, item := range src {
		tags, diags := basetypes.NewListValueFrom(ctx, types.StringType, item.Tags)
		if diags.HasError() {
			return nil, diags
		}

		value, key, values, diags := getItemFields(ctx, item)
		if diags.HasError() {
			return nil, diags
		}

		result = append(result, clientListItemModel{
			Value:            value,
			Key:              key,
			Values:           values,
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
		})
	}
	return result, nil
}

func getItemFields(ctx context.Context, item clientlists.ListItemContent) (types.String, types.String, types.List, diag.Diagnostics) {
	if isKeyValuesItem(string(item.Type)) {
		return getKeyValuesItemFields(ctx, item)
	}
	return getValueItemFields(item)
}

func getKeyValuesItemFields(ctx context.Context, item clientlists.ListItemContent) (types.String, types.String, types.List, diag.Diagnostics) {
	key := stringValueOrNull(item.Key)
	if len(item.Values) == 0 {
		return types.StringNull(), key, types.ListNull(types.StringType), nil
	}
	values, diags := basetypes.NewListValueFrom(ctx, types.StringType, item.Values)
	return types.StringNull(), key, values, diags
}

func getValueItemFields(item clientlists.ListItemContent) (types.String, types.String, types.List, diag.Diagnostics) {
	return stringValueOrNull(calculateValue(item)), types.StringNull(), types.ListNull(types.StringType), nil
}

func stringValueOrNull(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
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
	case clientlists.DOMAIN:
		return "domainClientListItemsDS"
	case clientlists.RequestHeaderNameValue:
		return "requestHeaderNameValueClientListItemsDS"
	default:
		return "unknownClientListItemsDS" // fallback or handle error
	}
}
