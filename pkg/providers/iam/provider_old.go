package iam

import (
	"context"
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type providerOld struct {
	client iam.IAM
	cache  Cache
	mtx    sync.Mutex
}

// Schema returns the subprovider's config schema map
func (p *providerOld) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

// Resources returns the subprovider's resource schema map
func (p *providerOld) Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// DataSources returns the subprovider's data source schema map
func (p *providerOld) DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_iam_countries":        p.dsCountries(),
		"akamai_iam_contact_types":    p.dsContactTypes(),
		"akamai_iam_supported_langs":  p.dsLanguages(),
		"akamai_iam_timeout_policies": p.dsTimeoutPolicies(),
		"akamai_iam_states":           p.dsStates(),
	}
}

// ProviderSchema returns a new provider schema instance
func (p *providerOld) ProviderSchema() *schema.Provider {
	return &schema.Provider{
		Schema:         p.Schema(),
		DataSourcesMap: p.DataSources(),
		ResourcesMap:   p.Resources(),
	}
}

// Configure receives the core provider's config data
func (p *providerOld) Configure(_ log.Interface, _ *schema.ResourceData) diag.Diagnostics {
	return nil
}

// Name returns the subprovider's name
func (p *providerOld) Name() string {
	return "iam"
}

// Version returns the subprovider's version
func (p *providerOld) Version() string {
	return "v0.0.1"
}

// SetIAM allows injection of an IAM.Client
func (p *providerOld) SetIAM(c iam.IAM) {
	p.client = c
}

// SetSession allows injection of a session.Session
func (p *providerOld) SetSession(s session.Session) {
	p.SetIAM(iam.Client(s))
}

// SetCache allows injection of a Cache
func (p *providerOld) SetCache(c Cache) {
	p.cache = c
}

func (p *providerOld) log(ctx context.Context) log.Interface {
	return log.FromContext(ctx)
}
