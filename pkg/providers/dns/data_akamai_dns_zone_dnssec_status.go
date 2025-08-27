package dns

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/dns"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	zoneDNSSecStatusDataSource struct {
		meta meta.Meta
	}

	zoneDNSSecStatusDataSourceModel struct {
		Zone           types.String     `tfsdk:"zone"`
		CurrentRecords *securityRecords `tfsdk:"current_records"`
		NewRecords     *securityRecords `tfsdk:"new_records"`
		Alerts         types.Set        `tfsdk:"alerts"`
	}

	// securityRecords represents a set of DNSSEC records for a DNS zone
	securityRecords struct {
		DNSKeyRecord     types.String `tfsdk:"dnskey_record"`
		DSRecord         types.String `tfsdk:"ds_record"`
		ExpectedTTL      types.Int64  `tfsdk:"expected_ttl"`
		LastModifiedDate types.String `tfsdk:"last_modified_date"`
	}
)

var (
	_ datasource.DataSource              = &zoneDNSSecStatusDataSource{}
	_ datasource.DataSourceWithConfigure = &zoneDNSSecStatusDataSource{}
)

// NewZoneDNSSecStatusDataSource returns a new single zone's DNSSEC status data source
func NewZoneDNSSecStatusDataSource() datasource.DataSource { return &zoneDNSSecStatusDataSource{} }

func (d *zoneDNSSecStatusDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "akamai_zone_dnssec_status"
}

func (d *zoneDNSSecStatusDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *zoneDNSSecStatusDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Single zone's DNSSEC status data source.",
		Attributes: map[string]schema.Attribute{
			"zone": schema.StringAttribute{
				Required:    true,
				Description: "The name of the zone.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"current_records": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The currently active set of generated DNSSEC records.",
				Attributes: map[string]schema.Attribute{
					"dnskey_record": schema.StringAttribute{
						Computed:    true,
						Description: "The generated DNSKEY record for this zone.",
					},
					"ds_record": schema.StringAttribute{
						Computed:    true,
						Description: "The generated DS record for this zone.",
					},
					"expected_ttl": schema.Int64Attribute{
						Computed:    true,
						Description: "The TTL on the NS record for this zone. This should match the TTL on the DS or DNSKEY record.",
					},
					"last_modified_date": schema.StringAttribute{
						Computed:    true,
						Description: "The ISO 8601 timestamp on which these records were generated.",
					},
				},
			},
			"new_records": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The newly generated set of DNSSEC records, if one exists.",
				Attributes: map[string]schema.Attribute{
					"dnskey_record": schema.StringAttribute{
						Computed:    true,
						Description: "The generated DNSKEY record for this zone.",
					},
					"ds_record": schema.StringAttribute{
						Computed:    true,
						Description: "The generated DS record for this zone.",
					},
					"expected_ttl": schema.Int64Attribute{
						Computed:    true,
						Description: "The TTL on the NS record for this zone. This should match the TTL on the DS or DNSKEY record.",
					},
					"last_modified_date": schema.StringAttribute{
						Computed:    true,
						Description: "The ISO 8601 timestamp on which these records were generated.",
					},
				},
			},
			"alerts": schema.SetAttribute{
				Computed:    true,
				Description: "A set of existing problems with the current DNSSEC configuration.",
				ElementType: types.StringType,
			},
		},
	}
}

func (d *zoneDNSSecStatusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "DNS ZoneDNSSecStatus DataSource Read")

	var data zoneDNSSecStatusDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := inst.Client(d.meta)
	zoneName := data.Zone.ValueString()
	zonesDNSSecStatus, err := client.GetZonesDNSSecStatus(ctx, dns.GetZonesDNSSecStatusRequest{
		Zones: []string{zoneName},
	})
	if err != nil {
		resp.Diagnostics.AddError("fetching DNS ZoneDNSSecStatus failed: ", err.Error())
		return
	}
	// No status object is returned by Edge DNS if the zone has DNSSEC disabled. For a single-zone request
	// this means an empty response list.
	if len(zonesDNSSecStatus.DNSSecStatuses) == 0 {
		resp.Diagnostics.AddError(fmt.Sprintf("no DNSSEC status for zone: %s", zoneName),
			"make sure that zone has DNSSEC enabled")
		return
	}
	if len(zonesDNSSecStatus.DNSSecStatuses) > 1 {
		resp.Diagnostics.AddWarning(fmt.Sprintf("unexpected multiple DNSSEC statuses for zone: %s", zoneName),
			fmt.Sprintf("%+v", zonesDNSSecStatus.DNSSecStatuses))
	}

	diags := data.setAttributes(ctx, zonesDNSSecStatus.DNSSecStatuses[0])
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *zoneDNSSecStatusDataSourceModel) setAttributes(ctx context.Context, secStatus dns.SecStatus) diag.Diagnostics {
	m.CurrentRecords = ptr.To(newSecurityRecords(secStatus.CurrentRecords))
	if secStatus.NewRecords != nil {
		m.NewRecords = ptr.To(newSecurityRecords(*secStatus.NewRecords))
	}
	alerts, diags := types.SetValueFrom(ctx, types.StringType, secStatus.Alerts)
	if diags.HasError() {
		return diags
	}
	m.Alerts = alerts
	return nil
}

func newSecurityRecords(records dns.SecRecords) securityRecords {
	return securityRecords{
		DNSKeyRecord:     types.StringValue(records.DNSKeyRecord),
		DSRecord:         types.StringValue(records.DSRecord),
		ExpectedTTL:      types.Int64Value(records.ExpectedTTL),
		LastModifiedDate: types.StringValue(date.FormatRFC3339(records.LastModifiedDate)),
	}
}
