package cloudlets

import (
	"sync"

	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/cloudlets"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
)

type (
	provider struct {
		*schema.Provider

		client cloudlets.Cloudlets
	}

	// Option is a cloudlets provider option
	Option func(p *provider)
)

var (
	once sync.Once

	inst *provider
)

// Subprovider returns a core sub provider
func Subprovider(opts ...Option) akamai.Subprovider {
	once.Do(func() {
		inst = &provider{Provider: Provider()}

		for _, opt := range opts {
			opt(inst)
		}
	})

	return inst
}

// Provider returns the Akamai terraform.Resource provider.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_cloudlets_api_prioritization_match_rule":        dataSourceCloudletsAPIPrioritizationMatchRule(),
			"akamai_cloudlets_application_load_balancer":            dataSourceCloudletsApplicationLoadBalancer(),
			"akamai_cloudlets_application_load_balancer_match_rule": dataSourceCloudletsApplicationLoadBalancerMatchRule(),
			"akamai_cloudlets_audience_segmentation_match_rule":     dataSourceCloudletsAudienceSegmentationMatchRule(),
			"akamai_cloudlets_edge_redirector_match_rule":           dataSourceCloudletsEdgeRedirectorMatchRule(),
			"akamai_cloudlets_forward_rewrite_match_rule":           dataSourceCloudletsForwardRewriteMatchRule(),
			"akamai_cloudlets_phased_release_match_rule":            dataSourceCloudletsPhasedReleaseMatchRule(),
			"akamai_cloudlets_request_control_match_rule":           dataSourceCloudletsRequestControlMatchRule(),
			"akamai_cloudlets_visitor_prioritization_match_rule":    dataSourceCloudletsVisitorPrioritizationMatchRule(),
			"akamai_cloudlets_policy":                               dataSourceCloudletsPolicy(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_cloudlets_application_load_balancer":            resourceCloudletsApplicationLoadBalancer(),
			"akamai_cloudlets_application_load_balancer_activation": resourceCloudletsApplicationLoadBalancerActivation(),
			"akamai_cloudlets_policy":                               resourceCloudletsPolicy(),
			"akamai_cloudlets_policy_activation":                    resourceCloudletsPolicyActivation(),
		},
	}
	return provider
}

// WithClient sets the client interface function, used for mocking and testing
func WithClient(c cloudlets.Cloudlets) Option {
	return func(p *provider) {
		p.client = c
	}
}

// Client returns the Cloudlets interface
func (p *provider) Client(meta akamai.OperationMeta) cloudlets.Cloudlets {
	if p.client != nil {
		return p.client
	}
	return cloudlets.Client(meta.Session())
}

func (p *provider) Name() string {
	return "cloudlets"
}

// ProviderVersion update version string anytime provider adds new features
const ProviderVersion string = "v0.0.1"

func (p *provider) Version() string {
	return ProviderVersion
}

func (p *provider) Schema() map[string]*schema.Schema {
	return p.Provider.Schema
}

func (p *provider) Resources() map[string]*schema.Resource {
	return p.Provider.ResourcesMap
}

func (p *provider) DataSources() map[string]*schema.Resource {
	return p.Provider.DataSourcesMap
}

func (p *provider) Configure(_ log.Interface, _ *schema.ResourceData) diag.Diagnostics {
	return nil
}
