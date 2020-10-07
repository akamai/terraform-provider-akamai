package appsec

import (
	"errors"
	"sync"

	appsecv1 "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/config"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	provider struct {
		*schema.Provider

		client appsec.APPSEC
	}
	// Option is a appsec provider option
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
		Schema: map[string]*schema.Schema{
			"appsec_section": {
				Optional:   true,
				Type:       schema.TypeString,
				Default:    "default",
				Deprecated: akamai.NoticeDeprecatedUseAlias("appsec_section"),
			},
			"appsec": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     config.Options("appsec"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_appsec_configuration":         dataSourceConfiguration(),
			"akamai_appsec_configuration_version": dataSourceConfigurationVersion(),
			"akamai_appsec_custom_rules":          dataSourceCustomRules(),
			"akamai_appsec_custom_rule_actions":   dataSourceCustomRuleActions(),
			"akamai_appsec_export_configuration":  dataSourceExportConfiguration(),
			"akamai_appsec_match_targets":         dataSourceMatchTargets(),
			"akamai_appsec_rate_policies":         dataSourceRatePolicies(),
			"akamai_appsec_rate_policy_actions":   dataSourceRatePolicyActions(),
			"akamai_appsec_selectable_hostnames":  dataSourceSelectableHostnames(),
			"akamai_appsec_security_policy":       dataSourceSecurityPolicy(),
			"akamai_appsec_selected_hostnames":    dataSourceSelectedHostnames(),
			"akamai_appsec_slow_post":             dataSourceSlowPostProtectionSettings(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_appsec_configuration_version_clone": resourceConfigurationClone(),
			"akamai_appsec_selected_hostnames":          resourceSelectedHostname(),
			"akamai_appsec_security_policy_clone":       resourceSecurityPolicyClone(),
			"akamai_appsec_match_target":                resourceMatchTarget(),
			"akamai_appsec_custom_rule":                 resourceCustomRule(),
			"akamai_appsec_custom_rule_action":          resourceCustomRuleAction(),
			//"akamai_appsec_activations":                 resourceActivations(),
			"akamai_appsec_rate_policy":        resourceRatePolicy(),
			"akamai_appsec_rate_policy_action": resourceRatePolicyAction(),
			"akamai_appsec_slow_post":          resourceSlowPostProtectionSetting(),
		},
	}
	return provider
}

/*
type resourceData interface {
	GetOk(string) (interface{}, bool)
	Get(string) interface{}
}

type set interface {
	List() []interface{}
}
*/
// WithClient sets the client interface function, used for mocking and testing
func WithClient(c appsec.APPSEC) Option {
	return func(p *provider) {
		p.client = c
	}
}

// Client returns the PAPI interface
func (p *provider) Client(meta akamai.OperationMeta) appsec.APPSEC {
	if p.client != nil {
		return p.client
	}
	return appsec.Client(meta.Session())
}

func getAPPSECV1Service(d *schema.ResourceData) (*edgegrid.Config, error) {
	var APPSECv1Config edgegrid.Config
	var err error
	appsec, err := tools.GetSetValue("appsec", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, err
	}
	if err == nil {
		cfg := appsec.List()[0].(map[string]interface{})

		APPSECv1Config = edgegrid.Config{
			Host:         cfg["host"].(string),
			AccessToken:  cfg["access_token"].(string),
			ClientToken:  cfg["client_token"].(string),
			ClientSecret: cfg["client_secret"].(string),
			MaxBody:      cfg["max_body"].(int),
		}

		appsecv1.Init(APPSECv1Config)
		return &APPSECv1Config, nil
	}

	edgerc, err := tools.GetStringValue("edgerc", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, err
	}

	var section string

	for _, s := range tools.FindStringValues(d, "appsec_section", "config_section") {
		if s != "default" {
			section = s
			break
		}
	}

	if section != "" {
		d.Set("config_section", section)
	}

	APPSECv1Config, err = edgegrid.Init(edgerc, section)
	if err != nil {
		return nil, err
	}
	edgegrid.SetupLogging()
	appsecv1.Init(APPSECv1Config)

	return &APPSECv1Config, nil
}

func (p *provider) Name() string {
	return "appsec"
}

// DnsProviderVersion update version string anytime provider adds new features
const AppSecProviderVersion string = "v0.8.3"

func (p *provider) Version() string {
	return AppSecProviderVersion
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

func (p *provider) Configure(log log.Interface, d *schema.ResourceData) diag.Diagnostics {
	log.Debug("START Configure")

	_, err := getAPPSECV1Service(d)
	if err != nil {
		return nil
	}
	return nil
}
