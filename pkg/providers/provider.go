package providers

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/httpclient"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-akamai/pkg/providers/config"
	"github.com/terraform-providers/terraform-provider-akamai/version"
	"log"
)

// Config contains the Akamai provider configuration (unused).
type Config struct {
	terraformVersion string
}

type AdapterSetup interface {
	PrepareSchema() map[string]*schema.Schema
	PrepareResources() map[string]*schema.Resource
	PrepareDataSources() map[string]*schema.Resource
	GetServiceConfig(*schema.ResourceData) (*edgegrid.Config, error)
}

func Provider(adapters ...AdapterSetup) plugin.ProviderFunc {
	return func() terraform.ResourceProvider {
		provider := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"edgerc": &schema.Schema{
					Optional:    true,
					Type:        schema.TypeString,
					DefaultFunc: schema.EnvDefaultFunc("EDGERC", nil),
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
				"dns": &schema.Schema{
					Optional: true,
					Type:     schema.TypeSet,
					Elem:     config.GetOptions("dns"),
				},
				"gtm": &schema.Schema{
					Optional: true,
					Type:     schema.TypeSet,
					Elem:     config.GetOptions("gtm"),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"akamai_authorities_set":        dataSourceAuthoritiesSet(),
				"akamai_dns_record_set":         dataSourceDNSRecordSet(),
				"akamai_gtm_default_datacenter": dataSourceGTMDefaultDatacenter(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"akamai_dns_zone":       resourceDNSv2Zone(),
				"akamai_dns_record":     resourceDNSv2Record(),
				"akamai_gtm_domain":     resourceGTMv1Domain(),
				"akamai_gtm_datacenter": resourceGTMv1Datacenter(),
				"akamai_gtm_property":   resourceGTMv1Property(),
				"akamai_gtm_resource":   resourceGTMv1Resource(),
				"akamai_gtm_cidrmap":    resourceGTMv1Cidrmap(),
				"akamai_gtm_geomap":     resourceGTMv1Geomap(),
				"akamai_gtm_asmap":      resourceGTMv1ASmap(),
			},
		}
		for _, p := range adapters {
			subSchema, err := mergeSchema(p.PrepareSchema(), provider.Schema)
			if err != nil {
				// log and exit
				panic(err)
			}
			provider.Schema = subSchema
			resources, err := mergeResource(p.PrepareResources(), provider.ResourcesMap)
			if err != nil {
				// log and exit
				panic(err)
			}
			provider.ResourcesMap = resources
			dataSources, err := mergeResource(p.PrepareDataSources(), provider.DataSourcesMap)
			if err != nil {
				// log and exit
				panic(err)
			}
			provider.DataSourcesMap = dataSources
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
}

func providerConfigure(d *schema.ResourceData, terraformVersion string, adapters ...AdapterSetup) (interface{}, error) {
	log.Printf("[DEBUG] START providerConfigure  %s\n", terraformVersion)
	var configs []*edgegrid.Config
	for _, p := range adapters {
		cfg, err := p.GetServiceConfig(d)
		if err != nil {
			return nil, err
		}
		if cfg != nil {
			configs = append(configs, cfg)
		}
	}
	if len(configs) == 0 {
		return nil, fmt.Errorf("One or more Akamai Edgegrid provider configurations must be defined")
	}

	config := Config{
		terraformVersion: terraformVersion,
	}

	tfUserAgent := httpclient.TerraformUserAgent(config.terraformVersion)
	//og.Printf("[DEBUG] tfUserAgent  %s\n", tfUserAgent)
	providerVersion := fmt.Sprintf("terraform-provider-akamai/%s", version.ProviderVersion)
	//log.Printf("[DEBUG] providerVersion  %s\n", providerVersion)
	//userAgent := fmt.Sprintf("%s %s", tfUserAgent, providerVersion)
	client.UserAgent = fmt.Sprintf("%s %s", tfUserAgent, providerVersion)

	//log.Printf("[DEBUG] CLIENT UserAgent  %s\n", client.UserAgent)

	//client.UserAgent = client.UserAgent + " terraform/" + Version

	//return &Config{}, nil
	return &config, nil
}

func mergeSchema(from, to map[string]*schema.Schema) (map[string]*schema.Schema, error) {
	for k, v := range from {
		if _, ok := to[k]; ok {
			return nil, fmt.Errorf("oops, duplicate key")
		}
		to[k] = v
	}
	return to, nil
}

func mergeResource(from, to map[string]*schema.Resource) (map[string]*schema.Resource, error) {
	for k, v := range from {
		if _, ok := to[k]; ok {
			return nil, fmt.Errorf("oops, duplicate key")
		}
		to[k] = v
	}
	return to, nil
}
