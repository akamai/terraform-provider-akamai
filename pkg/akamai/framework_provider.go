package akamai

import (
	"context"
	"os"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/cache"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/config"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/logger"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v4/version"
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/spf13/cast"
)

var _ provider.Provider = &Provider{}

// Provider is the implementation of akamai terraform provider which uses terraform-plugin-framework
type Provider struct{}

// ProviderModel represents the model of Provider configuration
type ProviderModel struct {
	Edgerc       types.String `tfsdk:"edgerc"`
	Section      types.String `tfsdk:"config_section"`
	Config       types.Set    `tfsdk:"config"`
	CacheEnabled types.Bool   `tfsdk:"cache_enabled"`
	RequestLimit types.Int64  `tfsdk:"request_limit"`
}

// NewFrameworkProvider returns a function returning Provider as provider.Provider
func NewFrameworkProvider() func() provider.Provider {
	return func() provider.Provider {
		return &Provider{}
	}
}

// Metadata configures provider's metadata
func (p *Provider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "akamai"
	resp.Version = version.ProviderVersion
}

// Schema sets provider's configuration schema
func (p *Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"edgerc": schema.StringAttribute{
				Optional: true,
			},
			"config_section": schema.StringAttribute{
				Description: "The section of the edgerc file to use for configuration",
				Optional:    true,
			},
			"cache_enabled": schema.BoolAttribute{
				Optional: true,
			},
			"request_limit": schema.Int64Attribute{
				Description: "The maximum number of API requests to be made per second (0 for no limit)",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"config": config.FrameworkOptions(),
		},
	}
}

// Configure configures provider context at the beginning of the lifecycle
// based on the values user specified in the provider configuration block
func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Edgerc.IsNull() {
		if v := os.Getenv("EDGERC"); v != "" {
			data.Edgerc = types.StringValue(v)
		}
	}

	if data.CacheEnabled.IsNull() {
		data.CacheEnabled = types.BoolValue(true)
	}

	if data.RequestLimit.IsNull() {
		v := os.Getenv("AKAMAI_REQUEST_LIMIT")
		if v == "" {
			data.RequestLimit = types.Int64Value(0)
		} else {
			reqLimit, err := strconv.Atoi(v)
			if err != nil {
				resp.Diagnostics.Append(diag.NewErrorDiagnostic("configuring context failed", err.Error()))
				return
			}
			data.RequestLimit = types.Int64Value(int64(reqLimit))
		}
	}

	var edgercConfig map[string]any
	resp.Diagnostics.Append(data.Config.ElementsAs(ctx, &edgercConfig, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cache.Enable(data.CacheEnabled.ValueBool())
	edgerc, err := newEdgegridConfig(data.Edgerc.ValueString(), data.Section.ValueString(), edgercConfig)
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("configuring context failed", err.Error()))
		return
	}

	opid := uuid.NewString()
	log := hclog.FromContext(ctx).With(
		"OperationID", opid,
	)
	logger := logger.FromHCLog(log)
	userAgent := userAgent(req.TerraformVersion)

	sess, err := session.New(
		session.WithSigner(edgerc),
		session.WithUserAgent(userAgent),
		session.WithLog(logger),
		session.WithHTTPTracing(cast.ToBool(os.Getenv("AKAMAI_HTTP_TRACE_ENABLED"))),
		session.WithRequestLimit(int(data.RequestLimit.ValueInt64())),
	)
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("configuring context failed", err.Error()))
		return
	}

	meta, err := meta.New(sess, log, opid)
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("configuring context failed", err.Error()))
		return
	}

	resp.DataSourceData = meta
	resp.ResourceData = meta
}

// Resources returns slice of fuctions used to instantiate resource implementations
func (p *Provider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

// DataSources returns slice of fuctions used to instantiate data source implementations
func (p *Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
