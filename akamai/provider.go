package akamai

import (
	"errors"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"os"
	"strings"
)

const (
	Version = "0.1.0"
)

// Config contains the Akamai provider configuration (unused).
type Config struct {
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

// Provider returns the Akamai terraform.Resource provider.
func Provider() terraform.ResourceProvider {
	client.UserAgent = client.UserAgent + " terraform/" + Version

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"edgerc": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
			},
			"dns_section": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				Default:  "default",
			},
			"papi_section": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				Default:  "default",
			},
			"property_section": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				Default:  "default",
			},
			"property": &schema.Schema{
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     getConfigOptions("property"),
			},
			"dns": &schema.Schema{
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     getConfigOptions("dns"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_authorities_set": dataSourceAuthoritiesSet(),
			"akamai_contract":        dataSourcePropertyContract(),
			"akamai_cp_code":         dataSourceCPCode(),
			"akamai_dns_record_set":  dataSourceDNSRecordSet(),
			"akamai_group":           dataSourcePropertyGroups(),
			"akamai_property_rules":  dataPropertyRules(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_cp_code":             resourceCPCode(),
			"akamai_dns_zone":            resourceDNSv2Zone(),
			"akamai_dns_record":          resourceDNSv2Record(),
			"akamai_edge_hostname":       resourceSecureEdgeHostName(),
			"akamai_property":            resourceProperty(),
			"akamai_property_rules":      resourcePropertyRules(),
			"akamai_property_variables":  resourcePropertyVariables(),
			"akamai_property_activation": resourcePropertyActivation(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	dnsv2Config, dnsErr := getConfigDNSV2Service(d)
	papiConfig, papiErr := getPAPIV1Service(d)

	if dnsErr != nil && papiErr != nil || dnsv2Config == nil && papiConfig == nil {
		return nil, fmt.Errorf("at least one configuration must be defined")
	}

	return &Config{}, nil
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
	return &DNSv2Config, nil
}

func getPAPIV1Service(d resourceData) (*edgegrid.Config, error) {
	var papiConfig edgegrid.Config
	if _, ok := d.GetOk("property"); ok {
		log.Printf("[DEBUG] Setting property config via HCL")
		config := d.Get("property").(set).List()[0].(map[string]interface{})

		papiConfig = edgegrid.Config{
			Host:         config["host"].(string),
			AccessToken:  config["access_token"].(string),
			ClientToken:  config["client_token"].(string),
			ClientSecret: config["client_secret"].(string),
			MaxBody:      config["max_body"].(int),
		}

		papi.Init(papiConfig)
		return &papiConfig, nil
	}

	var err error
	edgerc := d.Get("edgerc").(string)
	if section, ok := d.GetOk("property_section"); ok && section != "default" {
		papiConfig, err = edgegrid.Init(edgerc, section.(string))
	} else if section, ok := d.GetOk("papi_section"); ok && section != "default" {
		papiConfig, err = edgegrid.Init(edgerc, section.(string))
	} else {
		papiConfig, err = edgegrid.Init(edgerc, "default")
	}

	if err != nil {
		return nil, err
	}

	papi.Init(papiConfig)
	return &papiConfig, nil
}
