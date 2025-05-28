// Package cloudlets contains implementation for Akamai Terraform sub-provider responsible for managing Cloudlets applications
package cloudlets

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/cloudlets"
	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/cloudlets/v3"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/subprovider"
)

type (
	// Subprovider gathers cloudlets resources and data sources
	Subprovider struct{}
)

var (
	_ subprovider.Subprovider = &Subprovider{}
)

var (
	client   cloudlets.Cloudlets
	v3Client v3.Cloudlets
)

// NewSubprovider returns a new cloudlets subprovider
func NewSubprovider() *Subprovider {
	return &Subprovider{}
}

// Client returns the cloudlets interface
func Client(meta meta.Meta) cloudlets.Cloudlets {
	if client != nil {
		return client
	}
	return cloudlets.Client(meta.Session())
}

// ClientV3 returns the cloudlets v3 interface
func ClientV3(meta meta.Meta) v3.Cloudlets {
	if v3Client != nil {
		return v3Client
	}
	return v3.Client(meta.Session())
}

// SDKResources returns the cloudlets resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_cloudlets_application_load_balancer":            resourceCloudletsApplicationLoadBalancer(),
		"akamai_cloudlets_application_load_balancer_activation": resourceCloudletsApplicationLoadBalancerActivation(),
		"akamai_cloudlets_policy":                               resourceCloudletsPolicy(),
		"akamai_cloudlets_policy_activation":                    resourceCloudletsPolicyActivation(),
	}
}

// SDKDataSources returns the cloudlets data sources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKDataSources() map[string]*schema.Resource {
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

// FrameworkResources returns the cloudlets resources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{}
}

// FrameworkDataSources returns the cloudlets data sources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPolicyActivationDataSource,
		NewSharedPolicyDataSource,
	}
}
