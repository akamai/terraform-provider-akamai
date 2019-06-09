package akamai

import (
	"fmt"
	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Config contains the Akamai provider configuration (unused).
type Config struct {
}

// Provider returns the Akamai terraform.Resource provider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"edgerc": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
			},
			"cps_section": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				Default:  "default",
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
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_authorities_set": dataSourceAuthoritiesSet(),
			"akamai_contract":        dataSourcePropertyContract(),
			"akamai_cp_code":         dataSourceCPCode(),
			"akamai_dns_record_set":  dataSourceDNSRecordSet(),
			"akamai_group":           dataSourcePropertyGroups(),
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
	dnsv2Config, err := getConfigDNSV2Service(d)
	if err != nil {
		return nil, err
	}

	papiConfig, err := getPAPIV1Service(d)
	if err != nil {
		return nil, err
	}

	if dnsv2Config == nil && papiConfig == nil {
		return nil, fmt.Errorf("at least one edgerc section must be defined")
	}

	return &Config{}, nil
}

func getConfigDNSV2Service(d *schema.ResourceData) (*edgegrid.Config, error) {
	edgerc := d.Get("edgerc").(string)
	section := d.Get("dns_section").(string)
	DNSv2Config, err := edgegrid.Init(edgerc, section)
	if err != nil {
		return nil, err
	}

	dnsv2.Init(DNSv2Config)

	return &DNSv2Config, nil
}

func getPAPIV1Service(d *schema.ResourceData) (*edgegrid.Config, error) {
	edgerc := d.Get("edgerc").(string)
	section := d.Get("papi_section").(string)

	papiConfig, err := edgegrid.Init(edgerc, section)
	if err != nil {
		return nil, err
	}

	papi.Init(papiConfig)

	return &papiConfig, nil
}
