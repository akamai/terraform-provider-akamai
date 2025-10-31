package cloudcertificates

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ccm"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf/validators"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &hostnameBindingsDataSource{}
	_ datasource.DataSourceWithConfigure = &hostnameBindingsDataSource{}
)

type (
	hostnameBindingsDataSource struct {
		meta meta.Meta
	}

	hostnameBindingsDataSourceModel struct {
		ContractID     types.String   `tfsdk:"contract_id"`
		GroupID        types.String   `tfsdk:"group_id"`
		Domain         types.String   `tfsdk:"domain"`
		Network        types.String   `tfsdk:"network"`
		ExpiringInDays types.Int64    `tfsdk:"expiring_in_days"`
		Bindings       []bindingModel `tfsdk:"bindings"`
	}

	bindingModel struct {
		CertificateID types.String `tfsdk:"certificate_id"`
		Hostname      types.String `tfsdk:"hostname"`
		Network       types.String `tfsdk:"network"`
		ResourceType  types.String `tfsdk:"resource_type"`
	}
)

var pageSize int64 = 100

// NewCloudCertificatesHostnameBindingsDataSource returns a new cloud certificates hostname bindings data source.
func NewCloudCertificatesHostnameBindingsDataSource() datasource.DataSource {
	return &hostnameBindingsDataSource{}
}

// Metadata configures data source's meta information.
func (d *hostnameBindingsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_hostname_bindings"
}

// Configure configures data source at the beginning of the lifecycle.
func (d *hostnameBindingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Schema is used to define data source's terraform schema.
func (d *hostnameBindingsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve details for hostname bindings for all certificates you have access to.",
		Attributes: map[string]schema.Attribute{
			"contract_id": schema.StringAttribute{
				Description: "Filter results to certificates created under the specified contract.",
				Optional:    true,
				Validators:  []validator.String{validators.NotEmptyString()},
			},
			"group_id": schema.StringAttribute{
				Description: "Filter results to certificates created under the specified group.",
				Optional:    true,
				Validators:  []validator.String{validators.NotEmptyString()},
			},
			"domain": schema.StringAttribute{
				Description: "Filter results to certificates associated with the specified domain. " +
					"Returns SANs and Subject CN for certificates matching the specified domain. " +
					"Matches are case-insensitive, and support wildcards. " +
					"For example, domain=example.com returns certificates with *.akamai.com or akamai.com in the SAN list or Subject CN.",
				Optional:   true,
				Validators: []validator.String{validators.NotEmptyString()},
			},
			"network": schema.StringAttribute{
				Description: "Filter results to certificates in the specified network, either 'STAGING' or 'PRODUCTION'.",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("STAGING", "PRODUCTION")},
			},
			"expiring_in_days": schema.Int64Attribute{
				Description: "Filter results to certificates expiring in the specified number of days from the request date. " +
					"For example, a value of 5 returns certificates that expire in the next five days or less. " +
					"A value of 0 returns only expired certificates.",
				Optional:   true,
				Validators: []validator.Int64{int64validator.AtLeast(0)},
			},
			"bindings": schema.ListNestedAttribute{
				Description: "Certificate bindings filtered per request.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"certificate_id": schema.StringAttribute{
							Description: "Unique identifier for the third-party certificate.",
							Computed:    true,
						},
						"hostname": schema.StringAttribute{
							Description: "Hostname on the Akamai CDN the certificate applies to.",
							Computed:    true,
						},
						"network": schema.StringAttribute{
							Description: "The deployment network, either 'STAGING' or 'PRODUCTION', on which the certificate is active for a property version.",
							Computed:    true,
						},
						"resource_type": schema.StringAttribute{
							Description: "Resource type this binding applies to. Currently, only 'CDN_HOSTNAME' is available.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state.
func (d *hostnameBindingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Cloud Certificates Hostname Bindings DataSource Read")

	var data hostnameBindingsDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	client := Client(d.meta)

	var bindings []ccm.CertificateBinding
	for page := int64(1); ; page++ {
		res, err := client.ListBindings(ctx, ccm.ListBindingsRequest{
			ContractID:     data.ContractID.ValueString(),
			GroupID:        data.GroupID.ValueString(),
			Domain:         data.Domain.ValueString(),
			Network:        ccm.Network(data.Network.ValueString()),
			ExpiringInDays: data.ExpiringInDays.ValueInt64Pointer(),
			Page:           page,
			PageSize:       pageSize,
		})
		if err != nil {
			resp.Diagnostics.AddError("List bindings failed", err.Error())
			return
		}
		bindings = append(bindings, res.Bindings...)
		if res.Links.Next == nil {
			break
		}
	}

	data.setData(bindings)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (m *hostnameBindingsDataSourceModel) setData(bindings []ccm.CertificateBinding) {
	m.Bindings = make([]bindingModel, len(bindings))
	for i, binding := range bindings {
		m.Bindings[i] = bindingModel{
			CertificateID: types.StringValue(binding.CertificateID.String()),
			Hostname:      types.StringValue(binding.Hostname),
			Network:       types.StringValue(binding.Network),
			ResourceType:  types.StringValue(binding.ResourceType),
		}
	}
}
