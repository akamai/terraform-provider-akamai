package akamai

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/tf/validators"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/subprovider"
	"github.com/akamai/terraform-provider-akamai/v7/version"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &Provider{}

// Provider is the implementation of akamai terraform provider which uses terraform-plugin-framework
type Provider struct {
	subproviders []subprovider.Subprovider
}

// ProviderModel represents the model of Provider configuration
type ProviderModel struct {
	EdgercPath    types.String `tfsdk:"edgerc"`
	EdgercSection types.String `tfsdk:"config_section"`
	EdgercConfig  types.Set    `tfsdk:"config"`
	CacheEnabled  types.Bool   `tfsdk:"cache_enabled"`
	RequestLimit  types.Int64  `tfsdk:"request_limit"`
	RetryMax      types.Int64  `tfsdk:"retry_max"`
	RetryWaitMin  types.Int64  `tfsdk:"retry_wait_min"`
	RetryWaitMax  types.Int64  `tfsdk:"retry_wait_max"`
	RetryDisabled types.Bool   `tfsdk:"retry_disabled"`
}

// ConfigModel represents the model of edgegrid configuration block
type ConfigModel struct {
	Host         types.String `tfsdk:"host"`
	AccessToken  types.String `tfsdk:"access_token"`
	ClientToken  types.String `tfsdk:"client_token"`
	ClientSecret types.String `tfsdk:"client_secret"`
	MaxBody      types.Int64  `tfsdk:"max_body"`
	AccountKey   types.String `tfsdk:"account_key"`
}

// NewFrameworkProvider returns a function returning Provider as provider.Provider
func NewFrameworkProvider(subproviders ...subprovider.Subprovider) func() provider.Provider {
	return func() provider.Provider {
		return &Provider{
			subproviders: subproviders,
		}
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
			"retry_max": schema.Int64Attribute{
				Description: "The maximum number retires of API requests, default 10",
				Optional:    true,
			},
			"retry_wait_min": schema.Int64Attribute{
				Description: "The minimum wait time in seconds between API requests retries, default is 1 sec",
				Optional:    true,
			},
			"retry_wait_max": schema.Int64Attribute{
				Description: "The maximum wait time in seconds between API requests retries, default is 30 sec",
				Optional:    true,
			},
			"retry_disabled": schema.BoolAttribute{
				Description: "Should the retries of API requests be disabled, default false",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"config": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							Required:   true,
							Validators: []validator.String{validators.NotEmptyString()},
						},
						"access_token": schema.StringAttribute{
							Required:   true,
							Validators: []validator.String{validators.NotEmptyString()},
						},
						"client_token": schema.StringAttribute{
							Required:   true,
							Validators: []validator.String{validators.NotEmptyString()},
						},
						"client_secret": schema.StringAttribute{
							Required:   true,
							Validators: []validator.String{validators.NotEmptyString()},
						},
						"max_body": schema.Int64Attribute{
							Optional: true,
						},
						"account_key": schema.StringAttribute{
							Optional: true,
						},
					},
				},
			},
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

	if data.EdgercPath.IsNull() {
		if v := os.Getenv("EDGERC"); v != "" {
			data.EdgercPath = types.StringValue(v)
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

	var edgegridConfigBearer configBearer
	if !data.EdgercConfig.IsNull() {
		var configModels []ConfigModel
		resp.Diagnostics.Append(data.EdgercConfig.ElementsAs(ctx, &configModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		configModel := configModels[0]
		edgegridConfigBearer = configBearer{
			accessToken:  configModel.AccessToken.ValueString(),
			accountKey:   configModel.AccountKey.ValueString(),
			clientSecret: configModel.ClientSecret.ValueString(),
			clientToken:  configModel.ClientToken.ValueString(),
			host:         configModel.Host.ValueString(),
			maxBody:      int(configModel.MaxBody.ValueInt64()),
		}

	}

	edgegridConfig, err := newEdgegridConfig(data.EdgercPath.ValueString(), data.EdgercSection.ValueString(), edgegridConfigBearer)
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("configuring context failed", err.Error()))
		return
	}

	requestLimit, err := getFrameworkConfigInt(data.RequestLimit, "AKAMAI_REQUEST_LIMIT")
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("configuring context failed", err.Error()))
		return
	}

	retryMax, err := getFrameworkConfigInt(data.RetryMax, "AKAMAI_RETRY_MAX")
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("configuring context failed", err.Error()))
		return
	}

	retryWaitMin, err := getFrameworkConfigInt(data.RetryWaitMin, "AKAMAI_RETRY_WAIT_MIN")
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("configuring context failed", err.Error()))
		return
	}

	retryWaitMax, err := getFrameworkConfigInt(data.RetryWaitMax, "AKAMAI_RETRY_WAIT_MAX")
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("configuring context failed", err.Error()))
		return
	}

	retryDisabled, err := getFrameworkConfigBool(data.RetryDisabled, "AKAMAI_RETRY_DISABLED")
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("configuring context failed", err.Error()))
		return
	}

	meta, err := configureContext(contextConfig{
		edgegridConfig: edgegridConfig,
		userAgent:      userAgent(req.TerraformVersion),
		ctx:            ctx,
		requestLimit:   requestLimit,
		enableCache:    data.CacheEnabled.ValueBool(),
		retryMax:       retryMax,
		retryWaitMin:   time.Duration(retryWaitMin) * time.Second,
		retryWaitMax:   time.Duration(retryWaitMax) * time.Second,
		retryDisabled:  retryDisabled,
	})
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("configuring context failed", err.Error()))
		return
	}

	resp.DataSourceData = meta
	resp.ResourceData = meta
}

// Resources returns slice of functions used to instantiate resource implementations
func (p *Provider) Resources(_ context.Context) []func() resource.Resource {
	resources := make([]func() resource.Resource, 0)

	for _, subprovider := range p.subproviders {
		resources = append(resources, subprovider.FrameworkResources()...)
	}

	return resources
}

// DataSources returns slice of functions used to instantiate data source implementations
func (p *Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	dataSources := make([]func() datasource.DataSource, 0)

	for _, subprovider := range p.subproviders {
		dataSources = append(dataSources, subprovider.FrameworkDataSources()...)
	}

	return dataSources
}

func getFrameworkConfigInt(tfValue types.Int64, envKey string) (int, error) {
	ret := int(tfValue.ValueInt64())
	if tfValue.IsNull() {
		if v := os.Getenv(envKey); v != "" {
			vv, err := strconv.Atoi(v)
			if err != nil {
				return 0, err
			}
			ret = vv
		}
	}
	return ret, nil
}

func getFrameworkConfigBool(tfValue types.Bool, envKey string) (bool, error) {
	ret := tfValue.ValueBool()
	if tfValue.IsNull() {
		if v := os.Getenv(envKey); v != "" {
			vv, err := strconv.ParseBool(v)
			if err != nil {
				return false, err
			}
			ret = vv
		}
	}
	return ret, nil
}
