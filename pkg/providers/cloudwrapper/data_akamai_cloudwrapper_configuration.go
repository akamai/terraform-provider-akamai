package cloudwrapper

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudwrapper"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &configurationDataSource{}
	_ datasource.DataSourceWithConfigure = &configurationDataSource{}
)

type (
	configurationDataSource struct {
		client cloudwrapper.CloudWrapper
	}

	configurationDataSourceModel struct {
		ID                      types.Int64                  `tfsdk:"id"`
		CapacityAlertsThreshold types.Int64                  `tfsdk:"capacity_alerts_threshold"`
		Comments                types.String                 `tfsdk:"comments"`
		ConfigName              types.String                 `tfsdk:"config_name"`
		ContractID              types.String                 `tfsdk:"contract_id"`
		LastActivatedBy         types.String                 `tfsdk:"last_activated_by"`
		LastActivatedDate       types.String                 `tfsdk:"last_activated_date"`
		LastUpdatedBy           types.String                 `tfsdk:"last_updated_by"`
		LastUpdatedDate         types.String                 `tfsdk:"last_updated_date"`
		Locations               []configurationLocationModel `tfsdk:"locations"`
		MultiCDNSettings        *multiCDNSettingsModel       `tfsdk:"multi_cdn_settings"`
		NotificationEmails      types.Set                    `tfsdk:"notification_emails"`
		PropertyIDs             types.Set                    `tfsdk:"property_ids"`
		RetainIdleObjects       types.Bool                   `tfsdk:"retain_idle_objects"`
		Status                  types.String                 `tfsdk:"status"`
	}

	configurationLocationModel struct {
		Capacity      capacityModel `tfsdk:"capacity"`
		Comments      types.String  `tfsdk:"comments"`
		TrafficTypeID types.Int64   `tfsdk:"traffic_type_id"`
		MapName       types.String  `tfsdk:"map_name"`
	}

	boccModel struct {
		ConditionalSamplingFrequency types.String `tfsdk:"conditional_sampling_frequency"`
		Enabled                      types.Bool   `tfsdk:"enabled"`
		ForwardType                  types.String `tfsdk:"forward_type"`
		RequestType                  types.String `tfsdk:"request_type"`
		SamplingFrequency            types.String `tfsdk:"sampling_frequency"`
	}

	cdnsModel struct {
		CDNAuthKeys []cdnAuthKeyModel `tfsdk:"cdn_auth_keys"`
		CDNCode     types.String      `tfsdk:"cdn_code"`
		Enabled     types.Bool        `tfsdk:"enabled"`
		HTTPSOnly   types.Bool        `tfsdk:"https_only"`
		IPACLCIDRs  types.Set         `tfsdk:"ip_acl_cidrs"`
	}

	cdnAuthKeyModel struct {
		AuthKeyName types.String `tfsdk:"auth_key_name"`
		ExpiryDate  types.String `tfsdk:"expiry_date"`
		HeaderName  types.String `tfsdk:"header_name"`
		Secret      types.String `tfsdk:"secret"`
	}

	dataStreamsModel struct {
		DataStreamsIDs types.Set   `tfsdk:"data_stream_ids"`
		Enabled        types.Bool  `tfsdk:"enabled"`
		SamplingRate   types.Int64 `tfsdk:"sampling_rate"`
	}

	originModel struct {
		Hostname   types.String `tfsdk:"hostname"`
		OriginID   types.String `tfsdk:"origin_id"`
		PropertyID types.Int64  `tfsdk:"property_id"`
	}

	multiCDNSettingsModel struct {
		EnableSoftAlerts types.Bool       `tfsdk:"enable_soft_alerts"`
		BOCC             boccModel        `tfsdk:"bocc"`
		CDNs             []cdnsModel      `tfsdk:"cdns"`
		DataStreams      dataStreamsModel `tfsdk:"data_streams"`
		Origins          []originModel    `tfsdk:"origins"`
	}
)

// NewConfigurationDataSource returns a new configuration data source
func NewConfigurationDataSource() datasource.DataSource {
	return &configurationDataSource{}
}

// SetClient assigns given client to configuration data source
func (d *configurationDataSource) SetClient(client cloudwrapper.CloudWrapper) {
	d.client = client
}

// Metadata configures data source's meta information
func (d *configurationDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudwrapper_configuration"
}

// Configure configures data source at the beginning of the lifecycle
func (d *configurationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		// ProviderData is nil when Configure is run first time as part of ValidateDataSourceConfig in framework provider
		return
	}

	if d.client != nil {
		return
	}

	m, ok := req.ProviderData.(meta.Meta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
	d.client = cloudwrapper.Client(m.Session())
}

// Schema is used to define data source's terraform schema
func (d *configurationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "CloudWrapper configuration",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of a Cloud Wrapper configuration.",
			},
			"capacity_alerts_threshold": schema.Int64Attribute{
				Computed:    true,
				Description: "Represents the threshold for sending alerts.",
			},
			"comments": schema.StringAttribute{
				Computed:    true,
				Description: "Additional information provided by user which can help to differentiate or track changes of the configuration.",
			},
			"config_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the configuration.",
			},
			"contract_id": schema.StringAttribute{
				Computed:    true,
				Description: "Contract ID with Cloud Wrapper entitlement.",
			},
			"last_activated_by": schema.StringAttribute{
				Computed:    true,
				Description: "User to last activate the configuration.",
			},
			"last_activated_date": schema.StringAttribute{
				Computed:    true,
				Description: "ISO format date that represents when the configuration was last activated successfully.",
			},
			"last_updated_by": schema.StringAttribute{
				Computed:    true,
				Description: "User to last modify the configuration.",
			},
			"last_updated_date": schema.StringAttribute{
				Computed:    true,
				Description: "ISO format date that represents when the configuration was last edited.",
			},
			"notification_emails": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Email addresses to receive notifications.",
			},
			"property_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of properties belonging to media delivery products. Properties need to be unique across configurations.",
			},
			"retain_idle_objects": schema.BoolAttribute{
				Computed:    true,
				Description: "Retain idle objects beyond their max idle lifetime.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Current state of the provisioning of the configuration, either SAVED, IN_PROGRESS, ACTIVE, DELETE_IN_PROGRESS, or FAILED.",
			},
		},
		Blocks: map[string]schema.Block{
			"multi_cdn_settings": schema.SingleNestedBlock{
				Description: "Specify details about the Multi CDN settings.",
				Attributes: map[string]schema.Attribute{
					"enable_soft_alerts": schema.BoolAttribute{
						Computed:    true,
						Description: "Option to opt out of alerts based on soft limits of bandwidth usage.",
					},
				},
				Blocks: map[string]schema.Block{
					"bocc": schema.SingleNestedBlock{
						Description: "Specify diagnostic data beacons details.",
						Attributes: map[string]schema.Attribute{
							"conditional_sampling_frequency": schema.StringAttribute{
								Computed:    true,
								Description: "The sampling frequency of requests and forwards for EDGE, MIDGRESS, and ORIGIN beacons.",
							},
							"enabled": schema.BoolAttribute{
								Computed:    true,
								Description: "Enable diagnostic data beacons for consumption by the Broadcast Operations Control Center.",
							},
							"forward_type": schema.StringAttribute{
								Computed:    true,
								Description: "Select whether to beacon diagnostics data for internal ORIGIN_ONLY, MIDGRESS_ONLY, or both ORIGIN_AND_MIDGRESS forwards.",
							},
							"request_type": schema.StringAttribute{
								Computed:    true,
								Description: "Select whether to beacon diagnostics data for EDGE_ONLY or EDGE_AND_MIDGRESS requests.",
							},
							"sampling_frequency": schema.StringAttribute{
								Computed:    true,
								Description: "The sampling frequency of requests and forwards for EDGE, MIDGRESS, and ORIGIN beacons.",
							},
						},
					},
					"cdns": schema.SetNestedBlock{
						Description: "List of CDN added for the configuration.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"cdn_code": schema.StringAttribute{
									Computed:    true,
									Description: "Unique identifier for the CDN.",
								},
								"enabled": schema.BoolAttribute{
									Computed:    true,
									Description: "Enable CDN.",
								},
								"https_only": schema.BoolAttribute{
									Computed:    true,
									Description: "Specify whether CDN communication is HTTPS only.",
								},
								"ip_acl_cidrs": schema.SetAttribute{
									ElementType: types.StringType,
									Computed:    true,
									Description: "Configure an access control list using IP addresses in CIDR notation.",
								},
							},
							Blocks: map[string]schema.Block{
								"cdn_auth_keys": schema.SetNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"auth_key_name": schema.StringAttribute{
												Computed:    true,
												Description: "The name of the auth key.",
											},
											"expiry_date": schema.StringAttribute{
												Computed:    true,
												Description: "The expirty date of an auth key.",
											},
											"header_name": schema.StringAttribute{
												Computed:    true,
												Description: "The header name of an auth key.",
											},
											"secret": schema.StringAttribute{
												Computed:    true,
												Description: "The secret of an auth key.",
											},
										},
									},
									Description: "List of auth keys configured for the CDN.",
								},
							},
						},
					},
					"data_streams": schema.SingleNestedBlock{
						Description: "Specifies data streams details.",
						Attributes: map[string]schema.Attribute{
							"data_stream_ids": schema.SetAttribute{
								ElementType: types.Int64Type,
								Computed:    true,
								Description: "Unique identifiers of the Data Streams.",
							},
							"enabled": schema.BoolAttribute{
								Computed:    true,
								Description: "Enables DataStream reporting.",
							},
							"sampling_rate": schema.Int64Attribute{
								Computed:    true,
								Description: "Specifies the percentage of log data you want to collect for this configuration.",
							},
						},
					},
					"origins": schema.SetNestedBlock{
						Description: "List of origins corresponding to the properties selected in the configuration.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"hostname": schema.StringAttribute{
									Computed:    true,
									Description: "Origins hostname corresponding to the Akamai Delivery Property.",
								},
								"origin_id": schema.StringAttribute{
									Computed:    true,
									Description: "Origin identifier and will be used to generated Multi CDN host names.",
								},
								"property_id": schema.Int64Attribute{
									Computed:    true,
									Description: "Property ID of the property that origin belongs to.",
								},
							},
						},
					},
				},
			},
			"locations": schema.SetNestedBlock{
				Description: "List of all unused properties.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"comments": schema.StringAttribute{
							Computed:    true,
							Description: "Additional comments provided by user.",
						},
						"traffic_type_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Unique identifier for the location and traffic type combination.",
						},
						"map_name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the map.",
						},
						"capacity": schema.ObjectAttribute{
							AttributeTypes: map[string]attr.Type{
								"unit":  types.StringType,
								"value": types.Int64Type,
							},
							Computed:    true,
							Description: "The capacity assigned to this configuration's location.",
						},
					},
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state
func (d *configurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "CloudWrapper Configuration DataSource Read")

	var data configurationDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	cfg, err := d.client.GetConfiguration(ctx, cloudwrapper.GetConfigurationRequest{
		ConfigID: data.ID.ValueInt64(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Reading CloudWrapper Configuration", err.Error())
		return
	}

	if resp.Diagnostics.Append(data.populate(ctx, cfg)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(resp.State.Set(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
}

func (m *configurationDataSourceModel) populate(ctx context.Context, cfg *cloudwrapper.Configuration) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ContractID = types.StringValue(cfg.ContractID)
	m.ConfigName = types.StringValue(cfg.ConfigName)
	m.Comments = types.StringValue(cfg.Comments)
	m.RetainIdleObjects = types.BoolValue(cfg.RetainIdleObjects)
	m.LastUpdatedBy = types.StringValue(cfg.LastUpdatedBy)
	m.LastUpdatedDate = types.StringValue(cfg.LastUpdatedDate)
	m.Status = types.StringValue(string(cfg.Status))
	m.setLocations(cfg.Locations)

	if cfg.CapacityAlertsThreshold == nil {
		m.CapacityAlertsThreshold = types.Int64Null()
	} else {
		m.CapacityAlertsThreshold = types.Int64Value(int64(*cfg.CapacityAlertsThreshold))
	}
	if cfg.LastActivatedBy == nil {
		m.LastActivatedBy = types.StringNull()
	} else {
		m.LastActivatedBy = types.StringValue(*cfg.LastActivatedBy)
	}
	if cfg.LastActivatedDate == nil {
		m.LastActivatedDate = types.StringNull()
	} else {
		m.LastActivatedDate = types.StringValue(*cfg.LastActivatedDate)
	}

	if diags.Append(m.setPropertyIDs(ctx, cfg.PropertyIDs)...); diags.HasError() {
		return diags
	}
	if diags.Append(m.setNotificationEmails(ctx, cfg.NotificationEmails)...); diags.HasError() {
		return diags
	}

	if cfg.MultiCDNSettings != nil {
		diags.Append(m.setMultiCDNSettings(ctx, cfg.MultiCDNSettings)...)
	}

	return diags
}

func (m *configurationDataSourceModel) setLocations(locations []cloudwrapper.ConfigLocationResp) {
	for _, loc := range locations {
		m.Locations = append(m.Locations, configurationLocationModel{
			Capacity: capacityModel{
				Unit:  types.StringValue(string(loc.Capacity.Unit)),
				Value: types.Int64Value(loc.Capacity.Value),
			},
			Comments:      types.StringValue(loc.Comments),
			TrafficTypeID: types.Int64Value(int64(loc.TrafficTypeID)),
			MapName:       types.StringValue(loc.MapName),
		})
	}
}

func (m *configurationDataSourceModel) setPropertyIDs(ctx context.Context, propertyIDs []string) diag.Diagnostics {
	var diags diag.Diagnostics
	m.PropertyIDs, diags = types.SetValueFrom(ctx, types.StringType, propertyIDs)

	return diags
}

func (m *configurationDataSourceModel) setNotificationEmails(ctx context.Context, notificationEmails []string) diag.Diagnostics {
	var diags diag.Diagnostics
	m.NotificationEmails, diags = types.SetValueFrom(ctx, types.StringType, notificationEmails)

	return diags
}

func (m *configurationDataSourceModel) setMultiCDNSettings(ctx context.Context, multiCDNSettings *cloudwrapper.MultiCDNSettings) diag.Diagnostics {
	m.MultiCDNSettings = &multiCDNSettingsModel{
		EnableSoftAlerts: types.BoolValue(multiCDNSettings.EnableSoftAlerts),
	}
	m.setBOCC(multiCDNSettings.BOCC)
	m.setOrigins(multiCDNSettings.Origins)
	m.setCDNs(ctx, multiCDNSettings.CDNs)

	return m.setDataStreams(ctx, multiCDNSettings.DataStreams)
}

func (m *configurationDataSourceModel) setBOCC(bocc *cloudwrapper.BOCC) {
	m.MultiCDNSettings.BOCC = boccModel{
		ConditionalSamplingFrequency: types.StringValue(string(bocc.ConditionalSamplingFrequency)),
		Enabled:                      types.BoolValue(bocc.Enabled),
		ForwardType:                  types.StringValue(string(bocc.ForwardType)),
		RequestType:                  types.StringValue(string(bocc.RequestType)),
		SamplingFrequency:            types.StringValue(string(bocc.SamplingFrequency)),
	}
}

func (m *configurationDataSourceModel) setOrigins(origins []cloudwrapper.Origin) {
	for _, origin := range origins {
		m.MultiCDNSettings.Origins = append(m.MultiCDNSettings.Origins, originModel{
			Hostname:   types.StringValue(origin.Hostname),
			OriginID:   types.StringValue(origin.OriginID),
			PropertyID: types.Int64Value(int64(origin.PropertyID)),
		})
	}
}

func (m *configurationDataSourceModel) setCDNs(ctx context.Context, cdns []cloudwrapper.CDN) diag.Diagnostics {
	var diags diag.Diagnostics
	for _, cdn := range cdns {
		var ips types.Set
		ips, diags = types.SetValueFrom(ctx, types.StringType, cdn.IPACLCIDRs)
		if diags.HasError() {
			return diags
		}

		authKeys := make([]cdnAuthKeyModel, 0, len(cdn.CDNAuthKeys))
		for _, authKey := range cdn.CDNAuthKeys {
			authKeys = append(authKeys, cdnAuthKeyModel{
				AuthKeyName: types.StringValue(authKey.AuthKeyName),
				ExpiryDate:  types.StringValue(authKey.ExpiryDate),
				HeaderName:  types.StringValue(authKey.HeaderName),
				Secret:      types.StringValue(authKey.Secret),
			})
		}

		m.MultiCDNSettings.CDNs = append(m.MultiCDNSettings.CDNs, cdnsModel{
			CDNAuthKeys: authKeys,
			CDNCode:     types.StringValue(cdn.CDNCode),
			Enabled:     types.BoolValue(cdn.Enabled),
			HTTPSOnly:   types.BoolValue(cdn.HTTPSOnly),
			IPACLCIDRs:  ips,
		})
	}

	return diags
}

func (m *configurationDataSourceModel) setDataStreams(ctx context.Context, datastreams *cloudwrapper.DataStreams) diag.Diagnostics {
	dataStreamsIDs, diags := types.SetValueFrom(ctx, types.Int64Type, datastreams.DataStreamIDs)
	if diags.HasError() {
		return diags
	}
	m.MultiCDNSettings.DataStreams = dataStreamsModel{
		DataStreamsIDs: dataStreamsIDs,
		Enabled:        types.BoolValue(datastreams.Enabled),
	}
	if datastreams.SamplingRate != nil {
		m.MultiCDNSettings.DataStreams.SamplingRate = types.Int64Value(int64(*datastreams.SamplingRate))
	} else {
		m.MultiCDNSettings.DataStreams.SamplingRate = types.Int64Null()
	}

	return diags
}
