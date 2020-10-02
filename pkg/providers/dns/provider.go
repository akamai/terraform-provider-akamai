package dns

import (
	"errors"
	"fmt"
	"sync"

	dns "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configdns"
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

		client dns.DNS
	}

	// Option is a dns provider option
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
			"dns_section": {
				Optional:   true,
				Type:       schema.TypeString,
				Default:    "default",
				Deprecated: akamai.NoticeDeprecatedUseAlias("dns_section"),
			},
			"dns": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     config.Options("dns"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_authorities_set": dataSourceAuthoritiesSet(),
			"akamai_dns_record_set":  dataSourceDNSRecordSet(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_dns_zone":   resourceDNSv2Zone(),
			"akamai_dns_record": resourceDNSv2Record(),
		},
	}
	return provider
}

// WithClient sets the client interface function, used for mocking and testing
func WithClient(c dns.DNS) Option {
	return func(p *provider) {
		p.client = c
	}
}

// Client returns the DNS interface
func (p *provider) Client(meta akamai.OperationMeta) dns.DNS {
	if p.client != nil {
		return p.client
	}
	return dns.Client(meta.Session())
}

func getConfigDNSV2Service(d *schema.ResourceData) (*edgegrid.Config, error) {
	var DNSv2Config edgegrid.Config
	var err error
	dns, err := tools.GetSetValue("dns", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, err
	}
	if err == nil {
		log.Infof("[DEBUG] Setting property config via HCL")
		//cfg := dns.List()[0].(map[string]interface{})
		dnsConfig := dns.List()
		if len(dnsConfig) == 0 {
			return nil, fmt.Errorf("'dns' property in provider must have at least one entry")
		}
		configMap, ok := dnsConfig[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("dns config entry is of invalid type; should be 'map[string]interface{}'")
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
		DNSv2Config = edgegrid.Config{
			Host:         host,
			AccessToken:  accessToken,
			ClientToken:  clientToken,
			ClientSecret: clientSecret,
			MaxBody:      maxBody,
		}
		return &DNSv2Config, nil
	}

	edgerc, err := tools.GetStringValue("edgerc", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, err
	}

	var section string

	for _, s := range tools.FindStringValues(d, "dns_section", "config_section") {
		if s != "default" {
			section = s
			break
		}
	}

	if section != "" {
		d.Set("config_section", section)
	}

	DNSv2Config, err = edgegrid.Init(edgerc, section)
	if err != nil {
		return nil, err
	}

	return &DNSv2Config, nil
}

func (p *provider) Name() string {
	return "dns"
}

// DNSProviderVersion update version string anytime provider adds new features
const DNSProviderVersion string = "v0.8.3"

func (p *provider) Version() string {
	return DNSProviderVersion
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

	_, err := getConfigDNSV2Service(d)
	if err != nil {
		return nil
	}
	return nil
}
