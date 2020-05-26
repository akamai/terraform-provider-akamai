package akamai

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/httpclient"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-akamai/version"
)

//const (
//	Version = "0.2.0"
//)

// Config contains the Akamai provider configuration (unused).
type Config struct {
	terraformVersion string
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
	//client.UserAgent = client.UserAgent + " terraform/" + Version

	provider := &schema.Provider{
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
			"gtm_section": &schema.Schema{
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
			"gtm": &schema.Schema{
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     getConfigOptions("gtm"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_authorities_set":        dataSourceAuthoritiesSet(),
			"akamai_contract":               dataSourcePropertyContract(),
			"akamai_cp_code":                dataSourceCPCode(),
			"akamai_dns_record_set":         dataSourceDNSRecordSet(),
			"akamai_group":                  dataSourcePropertyGroups(),
			"akamai_property_rules":         dataPropertyRules(),
			"akamai_property":               dataSourceAkamaiProperty(),
			"akamai_gtm_default_datacenter": dataSourceGTMDefaultDatacenter(),
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
			"akamai_gtm_domain":          resourceGTMv1Domain(),
			"akamai_gtm_datacenter":      resourceGTMv1Datacenter(),
			"akamai_gtm_property":        resourceGTMv1Property(),
			"akamai_gtm_resource":        resourceGTMv1Resource(),
			"akamai_gtm_cidrmap":         resourceGTMv1Cidrmap(),
			"akamai_gtm_geomap":          resourceGTMv1Geomap(),
			"akamai_gtm_asmap":           resourceGTMv1ASmap(),
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
	papiConfig, papiErr := getPAPIV1Service(d)
	gtmConfig, gtmErr := getConfigGTMV1Service(d)

	if dnsErr != nil && papiErr != nil && gtmErr != nil || dnsv2Config == nil && papiConfig == nil && gtmConfig == nil {
		return nil, fmt.Errorf("One or more Akamai Edgegrid provider configurations must be defined")
	}

	config := Config{
		terraformVersion: terraformVersion,
	}

	tfUserAgent := httpclient.TerraformUserAgent(config.terraformVersion)
	log.Printf("[DEBUG] tfUserAgent  %s\n", tfUserAgent)
	providerVersion := fmt.Sprintf("terraform-provider-akamai/%s", version.ProviderVersion)
	log.Printf("[DEBUG] providerVersion  %s\n", providerVersion)
	//userAgent := fmt.Sprintf("%s %s", tfUserAgent, providerVersion)
	client.UserAgent = fmt.Sprintf("%s %s", tfUserAgent, providerVersion)

	log.Printf("[DEBUG] CLIENT UserAgent  %s\n", client.UserAgent)

	//client.UserAgent = client.UserAgent + " terraform/" + Version

	//return &Config{}, nil
	return &config, nil
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

func getConfigGTMV1Service(d resourceData) (*edgegrid.Config, error) {
	var GTMv1Config edgegrid.Config
	var err error
	if _, ok := d.GetOk("gtm"); ok {
		config := d.Get("gtm").(set).List()[0].(map[string]interface{})

		GTMv1Config = edgegrid.Config{
			Host:         config["host"].(string),
			AccessToken:  config["access_token"].(string),
			ClientToken:  config["client_token"].(string),
			ClientSecret: config["client_secret"].(string),
			MaxBody:      config["max_body"].(int),
		}

		gtm.Init(GTMv1Config)
		edgegrid.SetupLogging()
		return &GTMv1Config, nil
	}

	edgerc := d.Get("edgerc").(string)
	section := d.Get("gtm_section").(string)
	GTMv1Config, err = edgegrid.Init(edgerc, section)
	if err != nil {
		return nil, err
	}

	gtm.Init(GTMv1Config)
	return &GTMv1Config, nil
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
