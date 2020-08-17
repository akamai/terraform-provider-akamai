package papi

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-akamai/pkg/providers/config"
	"log"
)

type Provider struct{}

func (p Provider) PrepareSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"papi_section": {
			Optional: true,
			Type:     schema.TypeString,
			Default:  "default",
		},
		"property_section": {
			Optional:    true,
			Type:        schema.TypeString,
			DefaultFunc: schema.EnvDefaultFunc("EDGERC_ENVIRONMENT", "default"),
		},
		"property": {
			Optional: true,
			Type:     schema.TypeSet,
			Elem:     config.GetOptions("property"),
		},
	}
}

func (p Provider) PrepareResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_cp_code":             resourceCPCode(),
		"akamai_edge_hostname":       resourceSecureEdgeHostName(),
		"akamai_property":            resourceProperty(),
		"akamai_property_rules":      resourcePropertyRules(),
		"akamai_property_variables":  resourcePropertyVariables(),
		"akamai_property_activation": resourcePropertyActivation(),
	}
}

func (p Provider) PrepareDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_contract":       dataSourcePropertyContract(),
		"akamai_cp_code":        dataSourceCPCode(),
		"akamai_group":          dataSourcePropertyGroups(),
		"akamai_property_rules": dataPropertyRules(),
		"akamai_property":       dataSourceAkamaiProperty(),
	}
}

type set interface {
	List() []interface{}
}

func (p Provider) GetServiceConfig(data *schema.ResourceData) (*edgegrid.Config, error) {
	var papiConfig edgegrid.Config
	if _, ok := data.GetOk("property"); ok {
		log.Printf("[DEBUG] Setting property config via HCL")
		config := data.Get("property").(set).List()[0].(map[string]interface{})

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
	edgerc := data.Get("edgerc").(string)
	if section, ok := data.GetOk("property_section"); ok && section != "default" {
		papiConfig, err = edgegrid.Init(edgerc, section.(string))
	} else if section, ok := data.GetOk("papi_section"); ok && section != "default" {
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
