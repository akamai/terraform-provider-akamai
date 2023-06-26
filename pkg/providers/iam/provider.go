// Package iam contains implementation for Akamai Terraform sub-provider responsible for managing identities and accesses
package iam

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers IAM resources and data sources
	Subprovider struct {
		client iam.IAM
	}

	option func(p *Subprovider)
)

var (
	once sync.Once

	inst *Subprovider
)

var _ subprovider.Plugin = &Subprovider{}

// NewSubprovider returns a core sub provider
func NewSubprovider() *Subprovider {
	once.Do(func() {
		inst = &Subprovider{}
	})

	return inst
}

func withClient(c iam.IAM) option {
	return func(p *Subprovider) {
		p.client = c
	}
}

// Client returns the DNS interface
func (p *Subprovider) Client(meta meta.Meta) iam.IAM {
	if p.client != nil {
		return p.client
	}
	return iam.Client(meta.Session())
}

// Resources returns terraform resources for IAM
func (p *Subprovider) Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_iam_blocked_user_properties": resourceIAMBlockedUserProperties(),
		"akamai_iam_group":                   resourceIAMGroup(),
		"akamai_iam_role":                    resourceIAMRole(),
		"akamai_iam_user":                    resourceIAMUser(),
	}
}

// DataSources returns terraform data sources for IAM
func (p *Subprovider) DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_iam_contact_types":    dataSourceIAMContactTypes(),
		"akamai_iam_countries":        dataSourceIAMCountries(),
		"akamai_iam_grantable_roles":  dataSourceIAMGrantableRoles(),
		"akamai_iam_groups":           dataSourceIAMGroups(),
		"akamai_iam_roles":            dataSourceIAMRoles(),
		"akamai_iam_states":           dataSourceIAMStates(),
		"akamai_iam_supported_langs":  dataSourceIAMLanguages(),
		"akamai_iam_timeout_policies": dataSourceIAMTimeoutPolicies(),
		"akamai_iam_timezones":        dataSourceIAMTimezones(),
	}
}
