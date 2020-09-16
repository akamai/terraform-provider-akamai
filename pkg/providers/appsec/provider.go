package appsec

import (
	"errors"
	"sync"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
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
	}
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
			"akamai_appsec_export_configuration":  dataSourceExportConfiguration(),
			"akamai_appsec_selectable_hostnames":  dataSourceSelectableHostnames(),
			"akamai_appsec_security_policy":       dataSourceSecurityPolicy(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_appsec_configuration_clone":   resourceConfigurationClone(),
			"akamai_appsec_selected_hostnames":    resourceSelectedHostnames(),
			"akamai_appsec_security_policy_clone": resourceSecurityPolicyClone(),
			"akamai_appsec_match_targets":         resourceMatchTargets(),
			"akamai_appsec_custom_rule":           resourceCustomRule(),
			"akamai_appsec_activations":           resourceActivations(),
			"akamai_appsec_rate_policy":           resourceRatePolicy(),
		},
	}
	return provider
}

type resourceData interface {
	GetOk(string) (interface{}, bool)
	Get(string) interface{}
}

type set interface {
	List() []interface{}
}

func getAPPSECV1Service(d resourceData) (*edgegrid.Config, error) {
	var APPSECv1Config edgegrid.Config
	var err error
	if _, ok := d.GetOk("appsec"); ok {
		config := d.Get("appsec").(set).List()[0].(map[string]interface{})

		APPSECv1Config = edgegrid.Config{
			Host:         config["host"].(string),
			AccessToken:  config["access_token"].(string),
			ClientToken:  config["client_token"].(string),
			ClientSecret: config["client_secret"].(string),
			MaxBody:      config["max_body"].(int),
		}

		appsec.Init(APPSECv1Config)
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
	APPSECv1Config, err = edgegrid.Init(edgerc, section)
	if err != nil {
		return nil, err
	}

	appsec.Init(APPSECv1Config)
	edgegrid.SetupLogging()
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
