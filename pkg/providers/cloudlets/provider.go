// Package cloudlets contains implementation for Akamai Terraform sub-provider responsible for managing Cloudlets applications
package cloudlets

import (
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudlets"
	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudlets/v3"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/subprovider"
)

type (
	// Subprovider gathers cloudlets resources and data sources
	Subprovider struct {
		client   cloudlets.Cloudlets
		v3Client v3.Cloudlets
	}

	option func(p *Subprovider)
)

var (
	once sync.Once

	inst *Subprovider
)

var _ subprovider.Plugin = &Subprovider{}

// NewSubprovider returns a core sub provider
func NewSubprovider(opts ...option) *Subprovider {
	once.Do(func() {
		inst = &Subprovider{}

		for _, opt := range opts {
			opt(inst)
		}
	})

	return inst
}

func withClient(c cloudlets.Cloudlets) option {
	return func(p *Subprovider) {
		p.client = c
	}
}

// Client returns the Cloudlets interface for v2
func (p *Subprovider) Client(meta meta.Meta) cloudlets.Cloudlets {
	if p.client != nil {
		return p.client
	}
	return cloudlets.Client(meta.Session())
}

// V3Client returns the Cloudlets interface for v3
func (p *Subprovider) V3Client(meta meta.Meta) v3.Cloudlets {
	if p.v3Client != nil {
		return p.v3Client
	}
	return v3.Client(meta.Session())
}

// Resources returns terraform resources for cloudlets
func (p *Subprovider) Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_cloudlets_application_load_balancer":            resourceCloudletsApplicationLoadBalancer(),
		"akamai_cloudlets_application_load_balancer_activation": resourceCloudletsApplicationLoadBalancerActivation(),
		"akamai_cloudlets_policy":                               resourceCloudletsPolicy(),
		"akamai_cloudlets_policy_activation":                    resourceCloudletsPolicyActivation(),
	}
}

// DataSources returns terraform data sources for cloudlets
func (p *Subprovider) DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
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
	}
}
