// Package botman contains implementation for Akamai Terraform sub-provider responsible for maintaining Bot Manager
package botman

import (
	"fmt"
	"sync"

	"github.com/apex/log"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/providers/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	provider struct {
		*schema.Provider

		client botman.BotMan
	}

	// Option is a botman provider option
	Option func(p *provider)
)

var (
	once sync.Once

	inst *provider

	getLatestConfigVersion     = appsec.GetLatestConfigVersion
	getModifiableConfigVersion = appsec.GetModifiableConfigVersion
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
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_botman_akamai_bot_category":               dataSourceAkamaiBotCategory(),
			"akamai_botman_akamai_bot_category_action":        dataSourceAkamaiBotCategoryAction(),
			"akamai_botman_akamai_defined_bot":                dataSourceAkamaiDefinedBot(),
			"akamai_botman_bot_analytics_cookie":              dataSourceBotAnalyticsCookie(),
			"akamai_botman_bot_analytics_cookie_values":       dataSourceBotAnalyticsCookieValues(),
			"akamai_botman_bot_category_exception":            dataSourceBotCategoryException(),
			"akamai_botman_bot_detection":                     dataSourceBotDetection(),
			"akamai_botman_bot_detection_action":              dataSourceBotDetectionAction(),
			"akamai_botman_bot_endpoint_coverage_report":      dataSourceBotEndpointCoverageReport(),
			"akamai_botman_bot_management_settings":           dataSourceBotManagementSettings(),
			"akamai_botman_challenge_action":                  dataSourceChallengeAction(),
			"akamai_botman_challenge_interception_rules":      dataSourceChallengeInterceptionRules(),
			"akamai_botman_client_side_security":              dataSourceClientSideSecurity(),
			"akamai_botman_conditional_action":                dataSourceConditionalAction(),
			"akamai_botman_custom_bot_category":               dataSourceCustomBotCategory(),
			"akamai_botman_custom_bot_category_action":        dataSourceCustomBotCategoryAction(),
			"akamai_botman_custom_bot_category_sequence":      dataSourceCustomBotCategorySequence(),
			"akamai_botman_custom_client":                     dataSourceCustomClient(),
			"akamai_botman_custom_defined_bot":                dataSourceCustomDefinedBot(),
			"akamai_botman_custom_deny_action":                dataSourceCustomDenyAction(),
			"akamai_botman_javascript_injection":              dataSourceJavascriptInjection(),
			"akamai_botman_recategorized_akamai_defined_bot":  dataSourceRecategorizedAkamaiDefinedBot(),
			"akamai_botman_response_action":                   dataSourceResponseAction(),
			"akamai_botman_serve_alternate_action":            dataSourceServeAlternateAction(),
			"akamai_botman_transactional_endpoint":            dataSourceTransactionalEndpoint(),
			"akamai_botman_transactional_endpoint_protection": dataSourceTransactionalEndpointProtection(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_botman_akamai_bot_category_action":        resourceAkamaiBotCategoryAction(),
			"akamai_botman_bot_analytics_cookie":              resourceBotAnalyticsCookie(),
			"akamai_botman_bot_category_exception":            resourceBotCategoryException(),
			"akamai_botman_bot_detection_action":              resourceBotDetectionAction(),
			"akamai_botman_bot_management_settings":           resourceBotManagementSettings(),
			"akamai_botman_challenge_action":                  resourceChallengeAction(),
			"akamai_botman_challenge_interception_rules":      resourceChallengeInterceptionRules(),
			"akamai_botman_client_side_security":              resourceClientSideSecurity(),
			"akamai_botman_conditional_action":                resourceConditionalAction(),
			"akamai_botman_custom_bot_category":               resourceCustomBotCategory(),
			"akamai_botman_custom_bot_category_action":        resourceCustomBotCategoryAction(),
			"akamai_botman_custom_bot_category_sequence":      resourceCustomBotCategorySequence(),
			"akamai_botman_custom_client":                     resourceCustomClient(),
			"akamai_botman_custom_defined_bot":                resourceCustomDefinedBot(),
			"akamai_botman_custom_deny_action":                resourceCustomDenyAction(),
			"akamai_botman_javascript_injection":              resourceJavascriptInjection(),
			"akamai_botman_recategorized_akamai_defined_bot":  resourceRecategorizedAkamaiDefinedBot(),
			"akamai_botman_serve_alternate_action":            resourceServeAlternateAction(),
			"akamai_botman_transactional_endpoint":            resourceTransactionalEndpoint(),
			"akamai_botman_transactional_endpoint_protection": resourceTransactionalEndpointProtection(),
		},
	}
	return provider
}

// WithClient sets the client interface function, used for mocking and testing
func WithClient(c botman.BotMan) Option {
	return func(p *provider) {
		p.client = c
	}
}

// Client returns the PAPI interface
func (p *provider) Client(meta akamai.OperationMeta) botman.BotMan {
	if p.client != nil {
		return p.client
	}
	return botman.Client(meta.Session())
}

func getBotmanService(d *schema.ResourceData) (interface{}, error) {
	var section string

	for _, s := range tf.FindStringValues(d, "botman_section", "config_section") {
		if s != "default" {
			section = s
			break
		}
	}

	if section != "" {
		if err := d.Set("config_section", section); err != nil {
			return nil, fmt.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}

	return nil, nil
}

func (p *provider) Name() string {
	return "botman"
}

// ProviderVersion update version string anytime provider adds new features
const ProviderVersion string = "v1.0.1"

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

	_, err := getBotmanService(d)
	if err != nil {
		return nil
	}

	return nil
}
