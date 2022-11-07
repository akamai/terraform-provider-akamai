package property

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/config"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
)

type (
	provider struct {
		*schema.Provider

		client papi.PAPI

		hapiClient hapi.HAPI
	}

	// Option is a papi provider option
	Option func(p *provider)
)

var (
	once sync.Once

	inst *provider
)

// Subprovider returns a core sub provider
func Subprovider(opts ...Option) akamai.Subprovider {
	once.Do(func() {
		inst = &provider{Provider: Provider()}

		for _, opt := range opts {
			opt(inst)
		}
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
				Optional:   true,
				Type:       schema.TypeSet,
				Elem:       config.Options("property"),
				MaxItems:   1,
				Deprecated: akamai.NoticeDeprecatedUseAlias("property"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_contract":                dataSourcePropertyContract(),
			"akamai_contracts":               dataSourceAkamaiContracts(),
			"akamai_cp_code":                 dataSourceCPCode(),
			"akamai_group":                   dataSourcePropertyGroup(),
			"akamai_groups":                  dataSourcePropertyMultipleGroups(),
			"akamai_property_rules":          dataPropertyRules(),
			"akamai_property_rule_formats":   dataPropertyRuleFormats(),
			"akamai_property":                dataSourceAkamaiProperty(),
			"akamai_property_rules_template": dataSourcePropertyRulesTemplate(),
			"akamai_properties":              dataSourceAkamaiProperties(),
			"akamai_property_products":       dataSourceAkamaiPropertyProducts(),
			"akamai_property_hostnames":      dataSourceAkamaiPropertyHostnames(),
			"akamai_properties_search":       dataSourcePropertiesSearch(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_cp_code":             resourceCPCode(),
			"akamai_edge_hostname":       resourceSecureEdgeHostName(),
			"akamai_property":            resourceProperty(),
			"akamai_property_variables":  resourcePropertyVariables(),
			"akamai_property_activation": resourcePropertyActivation(),
		},
	}
	return provider
}

// WithClient sets the client interface function, used for mocking and testing
func WithClient(c papi.PAPI) Option {
	return func(p *provider) {
		p.client = c
	}
}

// Client returns the PAPI interface
func (p *provider) Client(meta akamai.OperationMeta) papi.PAPI {
	if p.client != nil {
		return p.client
	}
	return papi.Client(meta.Session())
}

// HapiClient returns the HAPI interface
func (p *provider) HapiClient(meta akamai.OperationMeta) hapi.HAPI {
	if p.hapiClient != nil {
		return p.hapiClient
	}
	return hapi.Client(meta.Session())
}

func getPAPIV1Service(d *schema.ResourceData) error {
	var inlineConfig *schema.Set
	for _, key := range []string{"property", "config"} {
		opt, err := tools.GetSetValue(key, d)
		if err != nil {
			if !errors.Is(err, tools.ErrNotFound) {
				return err
			}
			continue
		}
		if inlineConfig != nil {
			return fmt.Errorf("only one inline config section can be defined")
		}
		inlineConfig = opt
	}
	if err := d.Set("config", inlineConfig); err != nil {
		return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	for _, s := range tools.FindStringValues(d, "property_section", "papi_section", "config_section") {
		if s != "default" && s != "" {
			if err := d.Set("config_section", s); err != nil {
				return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
			}
			break
		}
	}

	return nil
}

func (p *provider) Name() string {
	return "property"
}

// ProviderVersion update version string anytime provider adds new features
const ProviderVersion string = "v0.8.3"

func (p *provider) Version() string {
	return ProviderVersion
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

	if err := getPAPIV1Service(d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// compactJSON converts a JSON-encoded byte slice to a compact form (so our JSON fixtures can be readable)
func compactJSON(encoded []byte) string {
	buf := bytes.Buffer{}
	if err := json.Compact(&buf, encoded); err != nil {
		panic(fmt.Sprintf("%s: %s", err, string(encoded)))
	}

	return buf.String()
}

// addPrefixToState returns a function that ensures string values are prefixed correctly
func addPrefixToState(prefix string) schema.SchemaStateFunc {
	return func(given interface{}) string {
		if given.(string) == "" {
			return ""
		}
		return tools.AddPrefix(given.(string), prefix)
	}
}
