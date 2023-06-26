// Package botman contains implementation for Akamai Terraform sub-provider responsible for maintaining Bot Manager
package botman

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/providers/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers botman resources and data sources
	Subprovider struct {
		client botman.BotMan
	}

	option func(p *Subprovider)
)

var (
	once sync.Once

	inst *Subprovider

	getLatestConfigVersion     = appsec.GetLatestConfigVersion
	getModifiableConfigVersion = appsec.GetModifiableConfigVersion
)

var _ subprovider.Plugin = &Subprovider{}

// NewSubprovider returns a core sub provider
func NewSubprovider(opts ...option) *Subprovider {
	once.Do(func() {
		inst = &Subprovider{}

		for _, opt := range opts {
			opt(inst)
		}
	})

	return inst
}

func withClient(c botman.BotMan) option {
	return func(p *Subprovider) {
		p.client = c
	}
}

// Client returns the PAPI interface
func (p *Subprovider) Client(meta meta.Meta) botman.BotMan {
	if p.client != nil {
		return p.client
	}
	return botman.Client(meta.Session())
}

// Resources returns terraform resources for botman
func (p *Subprovider) Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
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
	}
}

// DataSources returns terraform data sources for botman
func (p *Subprovider) DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
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
	}
}
