package iam

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type provider struct {
	Log    log.Interface
	Client iam.IAM

	cacheSet  func(string, interface{}) error
	cacheGet  func(string, interface{}) error
	checkMeta func(akamai.OperationMeta)
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
		"akmai_iam_roles":              p.dsRoles(),
		"akmai_iam_groups":             p.dsGroups(),
		"akmai_iam_countries":          p.dsCountries(),
		"akmai_iam_contact_types":      p.dsContactTypes(),
		"akmai_iam_supported_langs":    p.dsSupportedLangs(),
		"akmai_iam_notification_prods": p.dsNotificationProds(),
		"akmai_iam_timeout_policies":   p.dsTimeoutPolicies(),
		"akmai_iam_states":             p.dsStates(),
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

// Accept dependencies from Meta
func (p *provider) handleMeta(m interface{}) {
	if m == nil {
		return
	}

	meta := akamai.Meta(m)

	if p.checkMeta == nil {
		p.checkMeta = makeMetaCheck(meta)
	}

	p.checkMeta(meta)

	p.SetSession(meta.Session())
	p.SetApexLogger(meta.Log())

	p.cacheSet = func(key string, val interface{}) error {
		return meta.CacheSet(p, key, val)
	}

	p.cacheGet = func(key string, out interface{}) error {
		return meta.CacheGet(p, key, out)
	}
}

// Build a function that verifies the assumption that we receive exactly one meta value
func makeMetaCheck(originalMeta akamai.OperationMeta) func(akamai.OperationMeta) {
	if originalMeta == nil {
		panic("BUG: originalMeta can't be nil")
	}

	return func(newMeta akamai.OperationMeta) {
		if newMeta != originalMeta {
			panic("BUG: Received a new and different meta (invariant broken)")
		}
	}
}

func (p *provider) CacheSet(key string, val interface{}) error {
	if p.cacheSet == nil {
		return nil
	}

	return p.cacheSet(key, val)
}

func (p *provider) CacheGet(key string, out interface{}) error {
	if p.cacheGet == nil {
		return akamai.ErrCacheDisabled
	}

	return p.cacheGet(key, out)
}

// SetClient allows injection of an IAM.Client
func (p *provider) SetClient(c iam.IAM) {
	p.Client = c
}

// SetSession allows injection of a session.Session
func (p *provider) SetSession(s session.Session) {
	p.SetClient(iam.Client(s))
}

// SetApexLogger allows injection of a log.Interface
func (p *provider) SetApexLogger(l log.Interface) {
	p.Log = l
}
