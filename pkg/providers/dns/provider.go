package dns

import (
	"context"
	"log"
	"sync"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/config"
	"github.com/hashicorp/go-hclog"
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
	//ConfigureFunc: providerConfigure,
	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}
	return provider
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	log.Printf("[DEBUG] START providerConfigure  %s\n", terraformVersion)
	cfg, err := getConfigDNSV2Service(d)
	if err != nil {
		return nil, nil
	}
	return cfg, nil
}

type resourceData interface {
	GetOk(string) (interface{}, bool)
	Get(string) interface{}
}

type set interface {
	List() []interface{}
}

func getConfigDNSV2Service(d resourceData) (*edgegrid.Config, error) {
	var DNSv2Config edgegrid.Config
	var err error
	if _, ok := d.GetOk("dns"); ok {
		config := d.Get("dns").(set).List()[0].(map[string]interface{})

		DNSv2Config = edgegrid.Config{
			Host:         config["host"].(string),
			AccessToken:  config["access_token"].(string),
			ClientToken:  config["client_token"].(string),
			ClientSecret: config["client_secret"].(string),
			MaxBody:      config["max_body"].(int),
		}

		dnsv2.Init(DNSv2Config)
		return &DNSv2Config, nil
	}

	edgerc := d.Get("edgerc").(string)
	section := d.Get("dns_section").(string)
	if section == "" {
		section = d.Get("config_section").(string)
	}
	DNSv2Config, err = edgegrid.Init(edgerc, section)
	if err != nil {
		return nil, err
	}

	dnsv2.Init(DNSv2Config)
	edgegrid.SetupLogging()
	return &DNSv2Config, nil
}

func (p *provider) Name() string {
	return "dns"
}

func (p *provider) Version() string {
	return "v0.8.3"
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

func (p *provider) Configure(ctx context.Context, log hclog.Logger, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	state, err := p.ConfigureFunc(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return state, nil
}
