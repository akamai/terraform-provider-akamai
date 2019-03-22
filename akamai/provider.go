package akamai

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
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
			"dns_section": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				Default:  "default",
			},
			"dnsv2_section": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				Default:  "default",
			},
			"papi_section": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				Default:  "default",
			},
			"cps_section": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				Default:  "default",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_dns_record_set":  dataSourceDnsRecordSet(),
			"akamai_authorities_set": dataSourceAuthoritiesSet(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_cp_code":       resourceCPCode(),
			"akamai_dns_zone":      resourceDNSZone(),
			"akamai_dnsv2_zone":    resourceDNSv2Zone(),
			"akamai_dnsv2_record":  resourceDNSv2Record(),
			"akamai_property":      resourceProperty(),
      "akamai_cps_enrollment": resourceEnrollment(),
			"akamai_property_rule": resourcePropertyRule(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	dnsConfig, err := getConfigDNSV1Service(d)
	if err != nil {
		return nil, err
	}
	dnsv2Config, err := getConfigDNSV2Service(d)
	if err != nil {
		return nil, err
	}
	papiConfig, err := getPAPIV1Service(d)
	if err != nil {
		return nil, err
	}

	cpsConfig, err := getCPSV2Service(d)
	if err != nil {
		return nil, err
	}

	if dnsConfig == nil && dnsv2Config == nil && papiConfig == nil && cpsConfig == nil {
		return nil, fmt.Errorf("at least one edgerc section must be defined")
	}

	return &Config{}, nil
}

func getConfigDNSV1Service(d *schema.ResourceData) (*edgegrid.Config, error) {
	edgerc := d.Get("edgerc").(string)
	section := d.Get("dns_section").(string)

	DNSConfig, err := edgegrid.Init(edgerc, section)
	if err != nil {
		return nil, err
	}

	dns.Init(DNSConfig)

	return &DNSConfig, nil
}

func getConfigDNSV2Service(d *schema.ResourceData) (*edgegrid.Config, error) {
	edgerc := d.Get("edgerc").(string)
	section := d.Get("dnsv2_section").(string)
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
