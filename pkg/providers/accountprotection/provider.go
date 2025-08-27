// Package accountprotection contains implementation for Akamai Terraform sub-provider responsible for maintaining Bot Manager
package accountprotection

import (
	"sync"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/accountprotection"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/providers/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers account protection resources and data sources
	Subprovider struct {
		client apr.AccountProtection
	}

	option func(p *Subprovider)
)

var (
	once sync.Once

	inst *Subprovider

	getLatestConfigVersion     = appsec.GetLatestConfigVersion
	getModifiableConfigVersion = appsec.GetModifiableConfigVersion
)

var _ subprovider.Subprovider = &Subprovider{}

// NewSubprovider returns a new botman subprovider
func NewSubprovider(opts ...option) *Subprovider {
	once.Do(func() {
		inst = &Subprovider{}

		for _, opt := range opts {
			opt(inst)
		}
	})

	return inst
}

// Client returns the BotMan interface
func (p *Subprovider) Client(meta meta.Meta) apr.AccountProtection {
	if p.client != nil {
		return p.client
	}
	return apr.Client(meta.Session())
}

// SDKResources returns the botman resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_apr_protected_operations":        resourceProtectedOperations(),
		"akamai_apr_general_settings":            resourceGeneralSettings(),
		"akamai_apr_user_risk_response_strategy": resourceUserRiskResponseStrategy(),
		"akamai_apr_user_allow_list":             resourceUserAllowList(),
	}
}

// SDKDataSources returns the botman data sources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_apr_protected_operations":        dataSourceProtectedOperations(),
		"akamai_apr_general_settings":            dataSourceGeneralSettings(),
		"akamai_apr_user_risk_response_strategy": dataSourceUserRiskResponseStrategy(),
		"akamai_apr_user_allow_list":             dataSourceUserAllowList(),
	}
}

// FrameworkResources returns the botman resources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{}
}

// FrameworkDataSources returns the botman data sources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
