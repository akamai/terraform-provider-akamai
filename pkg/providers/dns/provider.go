package dns

import (
	"context"
	"errors"
	"fmt"
	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"os"
	"strings"
	"sync"
)

type (
	provider struct {
		*schema.Provider
	}
)

var (
	// DeprecatedSectionNotice is returned for schema configurations that are deprecated
	// Terraform now supports section aliases
	// TODO: Add alias example to the examples directory
	DeprecatedSectionNotice = func(n string) string {
		return fmt.Sprintf(`The setting %q has been deprecated. See:
https://www.terraform.io/docs/configuration/providers.html#alias-multiple-provider-configurations`, n)
	}

	once sync.Once

	inst *provider
)

// Subprovider returns a core sub provider
func Subprovider() akamai.Subprovider {
	once.Do(func() {
		inst = &provider{Provider: Provider()}

		// HACK: fixing this up to remove it when we use the subprovider entry
		delete(inst.Provider.Schema, "edgerc")
	})

	return inst
}

// Provider returns the Akamai terraform.Resource provider.
func Provider() *schema.Provider {

	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"edgerc": {
				Description: "The full file path for the edgerc configuration file.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"section": {
				Description: "The section of the edgerc file to use for configuration",
				Optional:    true,
				Type:        schema.TypeString,
				Default:     "default",
			},
			"dns_section": {
				Optional:   true,
				Type:       schema.TypeString,
				Default:    "default",
				Deprecated: DeprecatedSectionNotice("dns_section"),
			},
			"gtm_section": {
				Optional:   true,
				Type:       schema.TypeString,
				Default:    "default",
				Deprecated: DeprecatedSectionNotice("gtm_section"),
			},
			"papi_section": {
				Optional:   true,
				Type:       schema.TypeString,
				Default:    "default",
				Deprecated: DeprecatedSectionNotice("papi_section"),
			},
			"property_section": {
				Optional:   true,
				Type:       schema.TypeString,
				Default:    "default",
				Deprecated: DeprecatedSectionNotice("property_section"),
			},
			"property": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     getConfigOptions("property"),
			},
			"dns": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     getConfigOptions("dns"),
			},
			"gtm": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     getConfigOptions("gtm"),
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
	dnsv2Config, dnsErr := getConfigDNSV2Service(d)

	if dnsErr != nil || dnsv2Config == nil {
		return nil, fmt.Errorf("One or more Akamai Edgegrid provider configurations must be defined")
	}

	return nil, nil
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
	DNSv2Config, err = edgegrid.Init(edgerc, section)
	if err != nil {
		return nil, err
	}

	dnsv2.Init(DNSv2Config)
	edgegrid.SetupLogging()
	return &DNSv2Config, nil
}

func getConfigOptions(section string) *schema.Resource {
	section = strings.ToUpper(section)

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("AKAMAI_" + section + "_HOST"); v != "" {
						return v, nil
					} else if v := os.Getenv("AKAMAI_HOST"); v != "" {
						return v, nil
					}

					return nil, errors.New("host not set")
				},
			},
			"access_token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("AKAMAI_" + section + "_ACCESS_TOKEN"); v != "" {
						return v, nil
					} else if v := os.Getenv("AKAMAI_ACCESS_TOKEN"); v != "" {
						return v, nil
					}

					return nil, errors.New("access_token not set")
				},
			},
			"client_token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("AKAMAI_" + section + "_CLIENT_TOKEN"); v != "" {
						return v, nil
					} else if v := os.Getenv("AKAMAI_CLIENT_TOKEN"); v != "" {
						return v, nil
					}

					return nil, errors.New("client_token not set")
				},
			},
			"client_secret": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("AKAMAI_" + section + "_CLIENT_SECRET"); v != "" {
						return v, nil
					} else if v := os.Getenv("AKAMAI_CLIENT_SECRET"); v != "" {
						return v, nil
					}

					return nil, errors.New("client_secret not set")
				},
			},
			"max_body": {
				Type:     schema.TypeInt,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("AKAMAI_" + section + "_MAX_SIZE"); v != "" {
						return v, nil
					} else if v := os.Getenv("AKAMAI_MAX_SIZE"); v != "" {
						return v, nil
					}

					return 131072, nil
				},
			},
		},
	}
}

func (p *provider) Name() string {
	return "deprecated"
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
