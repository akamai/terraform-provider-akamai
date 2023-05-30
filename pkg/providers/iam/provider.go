// Package iam contains implementation for Akamai Terraform sub-provider responsible for managing identities and accesses
package iam

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/akamai"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	provider struct {
		*schema.Provider

		client iam.IAM
	}

	// Option is a iam provider option
	Option func(p *provider)
)

var (
	once sync.Once

	inst *provider
)

// Subprovider returns a core sub provider
func Subprovider() akamai.Subprovider {
	once.Do(func() {
		inst = &provider{Provider: Provider()}
	})

	return inst
}

// Provider returns the Akamai terraform.Resource provider.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_iam_contact_types":    dataSourceIAMContactTypes(),
			"akamai_iam_countries":        dataSourceIAMCountries(),
			"akamai_iam_grantable_roles":  dataSourceIAMGrantableRoles(),
			"akamai_iam_groups":           dataSourceIAMGroups(),
			"akamai_iam_roles":            dataSourceIAMRoles(),
			"akamai_iam_states":           dataSourceIAMStates(),
			"akamai_iam_supported_langs":  dataSourceIAMLanguages(),
			"akamai_iam_timeout_policies": dataSourceIAMTimeoutPolicies(),
			"akamai_iam_timezones":        dataSourceIAMTimezones(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_iam_blocked_user_properties": resourceIAMBlockedUserProperties(),
			"akamai_iam_group":                   resourceIAMGroup(),
			"akamai_iam_role":                    resourceIAMRole(),
			"akamai_iam_user":                    resourceIAMUser(),
		},
	}
	return provider
}

// WithClient sets the client interface function, used for mocking and testing
func WithClient(c iam.IAM) Option {
	return func(p *provider) {
		p.client = c
	}
}

// Client returns the DNS interface
func (p *provider) Client(meta akamai.OperationMeta) iam.IAM {
	if p.client != nil {
		return p.client
	}
	return iam.Client(meta.Session())
}

func (p *provider) Name() string {
	return "iam"
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
