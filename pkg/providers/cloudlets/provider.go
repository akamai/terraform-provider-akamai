// Package cloudlets contains implementation for Akamai Terraform sub-provider responsible for managing Cloudlets applications
package cloudlets

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/terraform-provider-akamai/v10/pkg/subprovider"
)

type (
	// Subprovider gathers cloudlets resources and data sources
	Subprovider struct{}
)

var _ subprovider.Subprovider = &Subprovider{}

const (
	// defaultPollActivationInterval is the default polling interval for activation status checks
	defaultPollActivationInterval = time.Minute
	// defaultPollRetryInterval is the default polling interval for retrying policy activation
	defaultPollRetryInterval = 15 * time.Second
	// defaultRetryTimeout is the default timeout for policy activation retries
	defaultRetryTimeout = 10 * time.Minute
)

// NewSubprovider returns a new cloudlets subprovider
func NewSubprovider() *Subprovider {
	return &Subprovider{}
}

// SDKResources returns the cloudlets resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_cloudlets_application_load_balancer":            resourceCloudletsApplicationLoadBalancer(),
		"akamai_cloudlets_application_load_balancer_activation": resourceCloudletsApplicationLoadBalancerActivation(defaultALBPollActivationInterval, defaultALBRetryTimeout),
		"akamai_cloudlets_policy":                               resourceCloudletsPolicy(defaultDeletionPolicyPollInterval),
		"akamai_cloudlets_policy_activation":                    resourceCloudletsPolicyActivation(defaultPollActivationInterval, defaultPollRetryInterval, defaultRetryTimeout),
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
