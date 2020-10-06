package gtm

import (
	"errors"
	"fmt"
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"
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

		client gtm.GTM
	}

	// Option is a gtm provider option
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
		Schema: map[string]*schema.Schema{
			"gtm_section": {
				Optional:   true,
				Type:       schema.TypeString,
				Default:    "default",
				Deprecated: akamai.NoticeDeprecatedUseAlias("gtm_section"),
			},
			"gtm": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     config.Options("gtm"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_gtm_default_datacenter": dataSourceGTMDefaultDatacenter(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_gtm_domain":     resourceGTMv1Domain(),
			"akamai_gtm_property":   resourceDNSv1Propery(),
			"akamai_gtm_datacenter": resourceDNSv1Datacenter(),
			"akamai_gtm_resource":   resourceDNSv1Resource(),
			"akamai_gtm_asmap":      resourceDNSv1ASmap(),
			"akamai_gtm_geomap":     resourceDNSv1Geomap(),
			"akamai_gtm_cidrmap":    resourceDNSv1Cidrmap(),
		},
	}
	return provider
}

// WithClient sets the client interface function, used for mocking and testing
func WithClient(c gtm.DNS) Option {
	return func(p *provider) {
		p.client = c
	}
}

// Client returns the DNS interface
func (p *provider) Client(meta akamai.OperationMeta) gtm.GTM {
	if p.client != nil {
		return p.client
	}
	return gtm.Client(meta.Session())
}

func getConfigGTMV1Service(d *schema.ResourceData) (*edgegrid.Config, error) {
	var GTMv1Config edgegrid.Config
	var err error
	gtm, err := tools.GetSetValue("gtm", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, err
	}
	if err == nil {
		log.Infof("[DEBUG] Setting property config via HCL")
		//cfg := gtm.List()[0].(map[string]interface{})
		gtmConfig := gtm.List()
		if len(gtmConfig) == 0 {
			return nil, fmt.Errorf("'gtm' property in provider must have at least one entry")
		}
		configMap, ok := gtmConfig[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("gtm config entry is of invalid type; should be 'map[string]interface{}'")
		}
		host, ok := configMap["host"].(string)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "host", "string")
		}
		accessToken, ok := configMap["access_token"].(string)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "access_token", "string")
		}
		clientToken, ok := configMap["client_token"].(string)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "client_token", "string")
		}
		clientSecret, ok := configMap["client_secret"].(string)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "client_secret", "string")
		}
		maxBody, ok := configMap["max_body"].(int)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "max_body", "int")
		}
		GTMv1Config = edgegrid.Config{
			Host:         host,
			AccessToken:  accessToken,
			ClientToken:  clientToken,
			ClientSecret: clientSecret,
			MaxBody:      maxBody,
		}
		return &GTMv1Config, nil
	}

	edgerc, err := tools.GetStringValue("edgerc", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, err
	}

	var section string

	for _, s := range tools.FindStringValues(d, "gtm_section", "config_section") {
		if s != "default" {
			section = s
			break
		}
	}

	if section != "" {
		d.Set("config_section", section)
	}

	GTMv1Config, err = edgegrid.Init(edgerc, section)
	if err != nil {
		return nil, err
	}

	return &GTMv1Config, nil
}

func (p *provider) Name() string {
	return "gtm"
}

// GTMProviderVersion update version string anytime provider adds new features
const GTMProviderVersion string = "v0.8.3"

func (p *provider) Version() string {
	return GTMProviderVersion
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

	_, err := getConfigGTMV1Service(d)
	if err != nil {
		return nil
	}
	return nil
}
