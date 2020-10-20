package property

import (
	"fmt"
	"strings"
	"sync"

	"github.com/apex/log"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/config"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	provider struct {
		*schema.Provider

		client papi.PAPI
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
				Deprecated: akamai.NoticeDeprecatedUseAlias("property"),
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

func getPAPIV1Service(d *schema.ResourceData) (interface{}, error) {
	var section string

	for _, s := range tools.FindStringValues(d, "property_section", "papi_section", "config_section") {
		if s != "default" {
			section = s
			break
		}
	}

	if section != "" {
		d.Set("config_section", section)
	}

	return nil, nil
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

	_, err := getPAPIV1Service(d)
	if err != nil {
		return nil
	}

	return nil
}

// Returns a schema.SchemaStateFunc that ensures ID values always have the given prefix
func statePrefixer(prefix string) schema.SchemaStateFunc {
	return func(given interface{}) string {
		s := given.(string)
		if !strings.HasPrefix(s, prefix) {
			s = fmt.Sprintf("%s%s", prefix, s)
		}

		return s
	}
}
