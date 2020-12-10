package iam

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type provider struct {
	client iam.IAM
	cache  Cache

	assertMeta func(akamai.OperationMeta)
}

// Schema returns the subprovider's config schema map
func (p *provider) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

// Resources returns the subprovider's resource schema map
func (p *provider) Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_iam_user": p.resUser(),
	}
}

// DataSources returns the subprovider's data source schema map
func (p *provider) DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_iam_roles":              p.dsRoles(),
		"akamai_iam_groups":             p.dsGroups(),
		"akamai_iam_countries":          p.dsCountries(),
		"akamai_iam_contact_types":      p.dsContactTypes(),
		"akamai_iam_supported_langs":    p.dsLanguages(),
		"akamai_iam_notification_prods": p.dsNotificationProds(),
		"akamai_iam_timeout_policies":   p.dsTimeoutPolicies(),
		"akamai_iam_states":             p.dsStates(),
	}
}

// Provider returns a new provider schema instance
func (p *provider) ProviderSchema() *schema.Provider {
	return &schema.Provider{
		Schema:         p.Schema(),
		DataSourcesMap: p.DataSources(),
		ResourcesMap:   p.Resources(),
	}
}

// Configure receives the core provider's config data
func (p *provider) Configure(log log.Interface, d *schema.ResourceData) diag.Diagnostics {
	return nil
}

// Name returns the subprovider's name
func (p *provider) Name() string {
	return "iam"
}

// Version returns the subprovider's version
func (p *provider) Version() string {
	return "v0.0.1"
}

// SetIAM allows injection of an IAM.Client
func (p *provider) SetIAM(c iam.IAM) {
	p.client = c
}

// SetSession allows injection of a session.Session
func (p *provider) SetSession(s session.Session) {
	p.SetIAM(iam.Client(s))
}

// SetCache allows injection of a Cache
func (p *provider) SetCache(c Cache) {
	p.cache = c
}

func (p *provider) log(ctx context.Context) log.Interface {
	return log.FromContext(ctx)
}
