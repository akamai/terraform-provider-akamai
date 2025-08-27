// Package property contains implementation for Property Provisioning module used to manage properties
package property

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers property resources and data sources
	Subprovider struct{}

	option func(p *Subprovider)
)

var (
	once sync.Once

	inst *Subprovider
)

var (
	_ subprovider.Subprovider = &Subprovider{}
)

var (
	client     papi.PAPI
	hapiClient hapi.HAPI
	iamClient  iam.IAM
)

// NewSubprovider returns a new property subprovider
func NewSubprovider(opts ...option) *Subprovider {
	once.Do(func() {
		inst = &Subprovider{}

		for _, opt := range opts {
			opt(inst)
		}
	})

	return inst
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

// IAMClient returns the IAM interface
func IAMClient(meta meta.Meta) iam.IAM {
	if iamClient != nil {
		return iamClient
	}
	return iam.Client(meta.Session())
}

// SDKResources returns the property resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_cp_code":                     resourceCPCode(),
		"akamai_edge_hostname":               resourceSecureEdgeHostName(),
		"akamai_property":                    resourceProperty(),
		"akamai_property_activation":         resourcePropertyActivation(),
		"akamai_property_include":            resourcePropertyInclude(),
		"akamai_property_include_activation": resourcePropertyIncludeActivation(),
	}
}

// SDKDataSources returns the property data sources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKDataSources() map[string]*schema.Resource {
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

// FrameworkResources returns the property resources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{
		NewBootstrapResource,
		NewHostnameBucketResource,
	}
}

// FrameworkDataSources returns the property data sources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCPCodesDataSource,
		NewHostnameActivationDataSource,
		NewHostnameActivationsDataSource,
		NewHostnamesDiffDataSource,
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
		return str.AddPrefix(given.(string), prefix)
	}
}
