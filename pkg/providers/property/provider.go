package property

import (
	"context"
	"log"
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
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
			"papi_section": {
				Optional:   true,
				Type:       schema.TypeString,
				Default:    "default",
				Deprecated: akamai.NoticeDeprecatedUseAlias("papi_section"),
			},
			"property_section": {
				Optional:   true,
				Type:       schema.TypeString,
				Default:    "default",
				Deprecated: akamai.NoticeDeprecatedUseAlias("property_section"),
			},
			"property": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     config.Options("property"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_contract":       dataSourcePropertyContract(),
			"akamai_cp_code":        dataSourceCPCode(),
			"akamai_group":          dataSourcePropertyGroups(),
			"akamai_property_rules": dataPropertyRules(),
			"akamai_property":       dataSourceAkamaiProperty(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_cp_code":             resourceCPCode(),
			"akamai_edge_hostname":       resourceSecureEdgeHostName(),
			"akamai_property":            resourceProperty(),
			"akamai_property_rules":      resourcePropertyRules(),
			"akamai_property_variables":  resourcePropertyVariables(),
			"akamai_property_activation": resourcePropertyActivation(),
		},
	}
	return provider
}

type resourceData interface {
	GetOk(string) (interface{}, bool)
	Get(string) interface{}
}

type set interface {
	List() []interface{}
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
		papiConfig, err = edgegrid.Init(edgerc, d.Get("config_section").(string))
	}

	if err != nil {
		return nil, err
	}

	papi.Init(papiConfig)
	return &papiConfig, nil
}

func (p *provider) Name() string {
	return "property"
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
	log.Named(p.Name()).Debug("START Configure")

	cfg, err := getPAPIV1Service(d)
	if err != nil {
		return nil, nil
	}

	return cfg, nil
}
