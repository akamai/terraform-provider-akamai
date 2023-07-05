// Package property contains implementation for Property Provisioning module used to manage properties
package property

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/subprovider"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// PluginSubprovider gathers property resources and data sources written using terraform-plugin-sdk
	PluginSubprovider struct{}

	// FrameworkSubprovider gathers property resources and data sources written using terraform-plugin-framework
	FrameworkSubprovider struct{}
)

var (
	client     papi.PAPI
	hapiClient hapi.HAPI
)

var _ subprovider.Plugin = &PluginSubprovider{}
var _ subprovider.Framework = &FrameworkSubprovider{}

// NewPluginSubprovider returns a core SDKv2 based sub provider
func NewPluginSubprovider() *PluginSubprovider {
	return &PluginSubprovider{}
}

// NewFrameworkSubprovider returns a core Framework based sub provider
func NewFrameworkSubprovider() *FrameworkSubprovider {
	return &FrameworkSubprovider{}
}

// Client returns the PAPI interface
func Client(meta meta.Meta) papi.PAPI {
	if client != nil {
		return client
	}
	return papi.Client(meta.Session())
}

// HapiClient returns the HAPI interface
func HapiClient(meta meta.Meta) hapi.HAPI {
	if hapiClient != nil {
		return hapiClient
	}
	return hapi.Client(meta.Session())
}

// Resources returns terraform resources for property
func (p *PluginSubprovider) Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_cp_code":                     resourceCPCode(),
		"akamai_edge_hostname":               resourceSecureEdgeHostName(),
		"akamai_property":                    resourceProperty(),
		"akamai_property_activation":         resourcePropertyActivation(),
		"akamai_property_include":            resourcePropertyInclude(),
		"akamai_property_include_activation": resourcePropertyIncludeActivation(),
	}
}

// DataSources returns terraform data sources for property
func (p *PluginSubprovider) DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_contract":                    dataSourcePropertyContract(),
		"akamai_contracts":                   dataSourceContracts(),
		"akamai_cp_code":                     dataSourceCPCode(),
		"akamai_group":                       dataSourcePropertyGroup(),
		"akamai_groups":                      dataSourcePropertyMultipleGroups(),
		"akamai_properties":                  dataSourceProperties(),
		"akamai_properties_search":           dataSourcePropertiesSearch(),
		"akamai_property":                    dataSourceProperty(),
		"akamai_property_activation":         dataSourcePropertyActivation(),
		"akamai_property_hostnames":          dataSourcePropertyHostnames(),
		"akamai_property_include_activation": dataSourcePropertyIncludeActivation(),
		"akamai_property_include_parents":    dataSourcePropertyIncludeParents(),
		"akamai_property_include_rules":      dataSourcePropertyIncludeRules(),
		"akamai_property_includes":           dataSourcePropertyIncludes(),
		"akamai_property_products":           dataSourcePropertyProducts(),
		"akamai_property_rule_formats":       dataSourcePropertyRuleFormats(),
		"akamai_property_rules":              dataSourcePropertyRules(),
		"akamai_property_rules_builder":      dataSourcePropertyRulesBuilder(),
		"akamai_property_rules_template":     dataSourcePropertyRulesTemplate(),
	}
}

// Resources returns terraform resources for property
func (p *FrameworkSubprovider) Resources() []func() resource.Resource {
	return []func() resource.Resource{}
}

// DataSources returns terraform data sources for property
func (p *FrameworkSubprovider) DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewIncludeDataSource,
	}
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
