package cloudwrapper

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/cloudwrapper"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/framework/replacer"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConfigurationResourceModel is a model for akamai_cloudwrapper_configuration resource
type ConfigurationResourceModel struct {
	ID                      types.Int64      `tfsdk:"id"`
	ConfigName              types.String     `tfsdk:"config_name"`
	ContractID              types.String     `tfsdk:"contract_id"`
	PropertyIDs             types.Set        `tfsdk:"property_ids"`
	Revision                types.String     `tfsdk:"revision"`
	Comments                types.String     `tfsdk:"comments"`
	RetainIdleObjects       types.Bool       `tfsdk:"retain_idle_objects"`
	NotificationEmails      types.Set        `tfsdk:"notification_emails"`
	CapacityAlertsThreshold types.Int64      `tfsdk:"capacity_alerts_threshold"`
	Locations               []ConfigLocation `tfsdk:"location"`
	Timeouts                timeouts.Value   `tfsdk:"timeouts"`
}

// ConfigLocation represents location item
type ConfigLocation struct {
	Comments      types.String     `tfsdk:"comments"`
	TrafficTypeID types.Int64      `tfsdk:"traffic_type_id"`
	Capacity      LocationCapacity `tfsdk:"capacity"`
}

// LocationCapacity holds capacity details for some location
type LocationCapacity struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func (m *ConfigurationResourceModel) hasUnknown() bool {
	for _, loc := range m.Locations {
		if loc.TrafficTypeID.IsUnknown() ||
			loc.Comments.IsUnknown() ||
			loc.Capacity.Unit.IsUnknown() ||
			loc.Capacity.Value.IsUnknown() {
			return true
		}
	}

	return m.ID.IsUnknown() ||
		m.ConfigName.IsUnknown() ||
		m.ContractID.IsUnknown() ||
		m.PropertyIDs.IsUnknown() ||
		m.Comments.IsUnknown() ||
		m.RetainIdleObjects.IsUnknown() ||
		m.NotificationEmails.IsUnknown() ||
		m.CapacityAlertsThreshold.IsUnknown()
}

func (m *ConfigurationResourceModel) buildCreateRequest(ctx context.Context) cloudwrapper.CreateConfigurationRequest {
	return cloudwrapper.CreateConfigurationRequest{
		Body: cloudwrapper.CreateConfigurationBody{
			CapacityAlertsThreshold: m.getCapacityAlertsThreshold(),
			Comments:                m.Comments.ValueString(),
			ContractID:              m.ContractID.ValueString(),
			Locations:               m.getLocationsReq(),
			MultiCDNSettings:        nil,
			ConfigName:              m.ConfigName.ValueString(),
			PropertyIDs:             m.getPropertyIDs(ctx),
			RetainIdleObjects:       m.RetainIdleObjects.ValueBool(),
			NotificationEmails:      m.getNotificationEmails(ctx),
		},
	}
}

func (m *ConfigurationResourceModel) buildUpdateRequest(ctx context.Context) cloudwrapper.UpdateConfigurationRequest {
	return cloudwrapper.UpdateConfigurationRequest{
		ConfigID: m.ID.ValueInt64(),
		Body: cloudwrapper.UpdateConfigurationBody{
			CapacityAlertsThreshold: m.getCapacityAlertsThreshold(),
			Comments:                m.Comments.ValueString(),
			Locations:               m.getLocationsReq(),
			MultiCDNSettings:        nil,
			NotificationEmails:      m.getNotificationEmails(ctx),
			PropertyIDs:             m.getPropertyIDs(ctx),
			RetainIdleObjects:       m.RetainIdleObjects.ValueBool(),
		},
	}
}

func (m *ConfigurationResourceModel) populateFrom(ctx context.Context, config *cloudwrapper.Configuration) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.Int64Value(config.ConfigID)
	m.ConfigName = types.StringValue(config.ConfigName)
	m.Comments = types.StringValue(config.Comments)
	m.RetainIdleObjects = types.BoolValue(config.RetainIdleObjects)

	m.setContractID(config.ContractID)
	m.setLocations(config.Locations)

	if config.CapacityAlertsThreshold == nil {
		m.CapacityAlertsThreshold = types.Int64Null()
	} else {
		m.CapacityAlertsThreshold = types.Int64Value(int64(*config.CapacityAlertsThreshold))
	}

	diags.Append(m.setPropertyIDs(ctx, config.PropertyIDs)...)
	if diags.HasError() {
		return diags
	}

	diags.Append(m.setNotificationEmails(ctx, config.NotificationEmails)...)
	if diags.HasError() {
		return diags
	}

	m.Revision = types.StringValue(calculateRevision(config))

	return diags
}

func (m *ConfigurationResourceModel) getCapacityAlertsThreshold() *int {
	var capacityAlertsThreshold *int
	if !m.CapacityAlertsThreshold.IsNull() {
		cat := int(m.CapacityAlertsThreshold.ValueInt64())
		capacityAlertsThreshold = &cat
	}
	return capacityAlertsThreshold
}

func (m *ConfigurationResourceModel) getPropertyIDs(ctx context.Context) []string {
	var propertyIDs []string
	m.PropertyIDs.ElementsAs(ctx, &propertyIDs, false)

	for i := range propertyIDs {
		propertyIDs[i] = strings.TrimPrefix(propertyIDs[i], "prp_")
	}

	return propertyIDs
}

func (m *ConfigurationResourceModel) setPropertyIDs(ctx context.Context, propIDs []string) diag.Diagnostics {
	var mProps []string
	if !m.PropertyIDs.IsNull() {
		diags := m.PropertyIDs.ElementsAs(ctx, &mProps, false)
		if diags.HasError() {
			return diags
		}
	}

	replaced := replacer.Replacer{
		Source:       propIDs,
		Replacements: mProps,
		EqFunc:       modifiers.EqualUpToPrefixFunc("prp_"),
	}.Replace()

	newPropIDs, diags := types.SetValueFrom(ctx, types.StringType, replaced)
	if diags.HasError() {
		return diags
	}

	m.PropertyIDs = newPropIDs

	return diags
}

func (m *ConfigurationResourceModel) setContractID(contract string) {
	if strings.TrimPrefix(contract, "ctr_") == strings.TrimPrefix(m.ContractID.ValueString(), "ctr_") {
		return
	}
	m.ContractID = types.StringValue(contract)
}

func (m *ConfigurationResourceModel) setNotificationEmails(ctx context.Context, emails []string) diag.Diagnostics {
	notificationEmails, diags := types.SetValueFrom(ctx, types.StringType, emails)
	if diags.HasError() {
		return diags
	}
	m.NotificationEmails = notificationEmails
	return diags
}

func (m *ConfigurationResourceModel) getNotificationEmails(ctx context.Context) []string {
	var emails []string
	m.NotificationEmails.ElementsAs(ctx, &emails, false)
	return emails
}

func (m *ConfigurationResourceModel) setLocations(locs []cloudwrapper.ConfigLocationResp) {
	m.Locations = make([]ConfigLocation, 0, len(locs))
	for _, loc := range locs {
		m.Locations = append(m.Locations, ConfigLocation{
			Comments:      types.StringValue(loc.Comments),
			TrafficTypeID: types.Int64Value(int64(loc.TrafficTypeID)),
			Capacity: LocationCapacity{
				Value: types.Int64Value(int64(loc.Capacity.Value)),
				Unit:  types.StringValue(string(loc.Capacity.Unit)),
			},
		})
	}
}

func (m *ConfigurationResourceModel) getLocationsReq() []cloudwrapper.ConfigLocationReq {
	locations := make([]cloudwrapper.ConfigLocationReq, 0, len(m.Locations))
	for _, loc := range m.Locations {
		locations = append(locations, cloudwrapper.ConfigLocationReq{
			Comments:      loc.Comments.ValueString(),
			TrafficTypeID: int(loc.TrafficTypeID.ValueInt64()),
			Capacity: cloudwrapper.Capacity{
				Value: loc.Capacity.Value.ValueInt64(),
				Unit:  cloudwrapper.Unit(loc.Capacity.Unit.ValueString()),
			},
		})
	}
	return locations
}

func (m *ConfigurationResourceModel) getLocationsResp() []cloudwrapper.ConfigLocationResp {
	locations := make([]cloudwrapper.ConfigLocationResp, 0, len(m.Locations))
	for _, loc := range m.Locations {
		locations = append(locations, cloudwrapper.ConfigLocationResp{
			Comments:      loc.Comments.ValueString(),
			TrafficTypeID: int(loc.TrafficTypeID.ValueInt64()),
			Capacity: cloudwrapper.Capacity{
				Value: loc.Capacity.Value.ValueInt64(),
				Unit:  cloudwrapper.Unit(loc.Capacity.Unit.ValueString()),
			},
		})
	}
	return locations
}

func (m *ConfigurationResourceModel) revision(ctx context.Context) string {
	return calculateRevision(&cloudwrapper.Configuration{
		CapacityAlertsThreshold: m.getCapacityAlertsThreshold(),
		Comments:                m.Comments.ValueString(),
		ContractID:              m.ContractID.ValueString(),
		ConfigID:                m.ID.ValueInt64(),
		Locations:               m.getLocationsResp(),
		ConfigName:              m.ConfigName.ValueString(),
		NotificationEmails:      m.getNotificationEmails(ctx),
		PropertyIDs:             m.getPropertyIDs(ctx),
		RetainIdleObjects:       m.RetainIdleObjects.ValueBool(),
	})
}

func calculateRevision(config *cloudwrapper.Configuration) string {
	sha := sha256.New()

	buffer := bytes.Buffer{}

	buffer.WriteString(config.Comments)

	sort.Strings(config.PropertyIDs)
	buffer.WriteString(strings.Join(config.PropertyIDs, ":"))

	buffer.WriteString(strconv.FormatBool(config.RetainIdleObjects))
	if config.CapacityAlertsThreshold != nil {
		buffer.WriteString(strconv.Itoa(*config.CapacityAlertsThreshold))
	}
	buffer.WriteString(strings.Join(config.NotificationEmails, ":"))

	sort.Slice(config.Locations, func(i, j int) bool {
		return config.Locations[i].TrafficTypeID < config.Locations[j].TrafficTypeID
	})
	for _, loc := range config.Locations {
		buffer.WriteString(loc.Comments)
		buffer.WriteString(strconv.FormatInt(loc.Capacity.Value, 10))
		buffer.WriteString(string(loc.Capacity.Unit))
	}

	_, err := buffer.WriteTo(sha)
	if err != nil {
		panic("calculate revision: writing to buffer failed")
	}

	return hex.EncodeToString(sha.Sum([]byte{}))[:20]
}
