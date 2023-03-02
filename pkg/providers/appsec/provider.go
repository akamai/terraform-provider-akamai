// Package appsec contains implementation for Akamai Terraform sub-provider responsible for Application Security
package appsec

import (
	"fmt"
	"sync"

	"github.com/apex/log"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/config"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	provider struct {
		*schema.Provider

		client appsec.APPSEC
	}
	// Option is a appsec provider option
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
			"appsec_section": {
				Optional:   true,
				Type:       schema.TypeString,
				Default:    "default",
				Deprecated: akamai.NoticeDeprecatedUseAlias("appsec_section"),
			},
			"appsec": {
				Optional:   true,
				Type:       schema.TypeSet,
				Elem:       config.Options("appsec"),
				Deprecated: akamai.NoticeDeprecatedUseAlias("appsec"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_appsec_advanced_settings_attack_payload_logging": dataSourceAdvancedSettingsAttackPayloadLogging(),
			"akamai_appsec_advanced_settings_evasive_path_match":     dataSourceAdvancedSettingsEvasivePathMatch(),
			"akamai_appsec_advanced_settings_logging":                dataSourceAdvancedSettingsLogging(),
			"akamai_appsec_advanced_settings_pragma_header":          dataSourceAdvancedSettingsPragmaHeader(),
			"akamai_appsec_advanced_settings_prefetch":               dataSourceAdvancedSettingsPrefetch(),
			"akamai_appsec_api_endpoints":                            dataSourceAPIEndpoints(),
			"akamai_appsec_api_request_constraints":                  dataSourceAPIRequestConstraints(),
			"akamai_appsec_attack_groups":                            dataSourceAttackGroups(),
			"akamai_appsec_bypass_network_lists":                     dataSourceBypassNetworkLists(),
			"akamai_appsec_configuration":                            dataSourceConfiguration(),
			"akamai_appsec_configuration_version":                    dataSourceConfigurationVersion(),
			"akamai_appsec_contracts_groups":                         dataSourceContractsGroups(),
			"akamai_appsec_custom_deny":                              dataSourceCustomDeny(),
			"akamai_appsec_custom_rule_actions":                      dataSourceCustomRuleActions(),
			"akamai_appsec_custom_rules":                             dataSourceCustomRules(),
			"akamai_appsec_eval":                                     dataSourceEval(),
			"akamai_appsec_eval_groups":                              dataSourceEvalGroups(),
			"akamai_appsec_eval_penalty_box":                         dataSourceEvalPenaltyBox(),
			"akamai_appsec_eval_rules":                               dataSourceEvalRules(),
			"akamai_appsec_export_configuration":                     dataSourceExportConfiguration(),
			"akamai_appsec_failover_hostnames":                       dataSourceFailoverHostnames(),
			"akamai_appsec_hostname_coverage":                        dataSourceAPIHostnameCoverage(),
			"akamai_appsec_hostname_coverage_match_targets":          dataSourceAPIHostnameCoverageMatchTargets(),
			"akamai_appsec_hostname_coverage_overlapping":            dataSourceAPIHostnameCoverageOverlapping(),
			"akamai_appsec_ip_geo":                                   dataSourceIPGeo(),
			"akamai_appsec_malware_content_types":                    dataSourceMalwareContentTypes(),
			"akamai_appsec_malware_policies":                         dataSourceMalwarePolicies(),
			"akamai_appsec_malware_policy_actions":                   dataSourceMalwarePolicyActions(),
			"akamai_appsec_match_targets":                            dataSourceMatchTargets(),
			"akamai_appsec_penalty_box":                              dataSourcePenaltyBox(),
			"akamai_appsec_rate_policies":                            dataSourceRatePolicies(),
			"akamai_appsec_rate_policy_actions":                      dataSourceRatePolicyActions(),
			"akamai_appsec_reputation_profile_actions":               dataSourceReputationProfileActions(),
			"akamai_appsec_reputation_profile_analysis":              dataSourceReputationAnalysis(),
			"akamai_appsec_reputation_profiles":                      dataSourceReputationProfiles(),
			"akamai_appsec_rule_upgrade_details":                     dataSourceRuleUpgrade(),
			"akamai_appsec_rules":                                    dataSourceRules(),
			"akamai_appsec_security_policy":                          dataSourceSecurityPolicy(),
			"akamai_appsec_security_policy_protections":              dataSourcePolicyProtections(),
			"akamai_appsec_selectable_hostnames":                     dataSourceSelectableHostnames(),
			"akamai_appsec_selected_hostnames":                       dataSourceSelectedHostnames(),
			"akamai_appsec_siem_definitions":                         dataSourceSiemDefinitions(),
			"akamai_appsec_siem_settings":                            dataSourceSiemSettings(),
			"akamai_appsec_slow_post":                                dataSourceSlowPostProtectionSettings(),
			"akamai_appsec_threat_intel":                             dataSourceThreatIntel(),
			"akamai_appsec_tuning_recommendations":                   dataSourceTuningRecommendations(),
			"akamai_appsec_version_notes":                            dataSourceVersionNotes(),
			"akamai_appsec_waf_mode":                                 dataSourceWAFMode(),
			"akamai_appsec_wap_selected_hostnames":                   dataSourceWAPSelectedHostnames(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_appsec_activations":                              resourceActivations(),
			"akamai_appsec_advanced_settings_attack_payload_logging": resourceAdvancedSettingsAttackPayloadLogging(),
			"akamai_appsec_advanced_settings_evasive_path_match":     resourceAdvancedSettingsEvasivePathMatch(),
			"akamai_appsec_advanced_settings_logging":                resourceAdvancedSettingsLogging(),
			"akamai_appsec_advanced_settings_pragma_header":          resourceAdvancedSettingsPragmaHeader(),
			"akamai_appsec_advanced_settings_prefetch":               resourceAdvancedSettingsPrefetch(),
			"akamai_appsec_api_constraints_protection":               resourceAPIConstraintsProtection(),
			"akamai_appsec_api_request_constraints":                  resourceAPIRequestConstraints(),
			"akamai_appsec_attack_group":                             resourceAttackGroup(),
			"akamai_appsec_bypass_network_lists":                     resourceBypassNetworkLists(),
			"akamai_appsec_configuration":                            resourceConfiguration(),
			"akamai_appsec_configuration_rename":                     resourceConfigurationRename(),
			"akamai_appsec_custom_deny":                              resourceCustomDeny(),
			"akamai_appsec_custom_rule":                              resourceCustomRule(),
			"akamai_appsec_custom_rule_action":                       resourceCustomRuleAction(),
			"akamai_appsec_eval":                                     resourceEval(),
			"akamai_appsec_eval_group":                               resourceEvalGroup(),
			"akamai_appsec_eval_penalty_box":                         resourceEvalPenaltyBox(),
			"akamai_appsec_eval_rule":                                resourceEvalRule(),
			"akamai_appsec_ip_geo":                                   resourceIPGeo(),
			"akamai_appsec_ip_geo_protection":                        resourceIPGeoProtection(),
			"akamai_appsec_malware_policy":                           resourceMalwarePolicy(),
			"akamai_appsec_malware_policy_action":                    resourceMalwarePolicyAction(),
			"akamai_appsec_malware_policy_actions":                   resourceMalwarePolicyActions(),
			"akamai_appsec_malware_protection":                       resourceMalwareProtection(),
			"akamai_appsec_match_target":                             resourceMatchTarget(),
			"akamai_appsec_match_target_sequence":                    resourceMatchTargetSequence(),
			"akamai_appsec_penalty_box":                              resourcePenaltyBox(),
			"akamai_appsec_rate_policy":                              resourceRatePolicy(),
			"akamai_appsec_rate_policy_action":                       resourceRatePolicyAction(),
			"akamai_appsec_rate_protection":                          resourceRateProtection(),
			"akamai_appsec_reputation_profile":                       resourceReputationProfile(),
			"akamai_appsec_reputation_profile_action":                resourceReputationProfileAction(),
			"akamai_appsec_reputation_profile_analysis":              resourceReputationAnalysis(),
			"akamai_appsec_reputation_protection":                    resourceReputationProtection(),
			"akamai_appsec_rule":                                     resourceRule(),
			"akamai_appsec_rule_upgrade":                             resourceRuleUpgrade(),
			"akamai_appsec_security_policy":                          resourceSecurityPolicy(),
			"akamai_appsec_security_policy_rename":                   resourceSecurityPolicyRename(),
			"akamai_appsec_selected_hostnames":                       resourceSelectedHostname(),
			"akamai_appsec_siem_settings":                            resourceSiemSettings(),
			"akamai_appsec_slow_post":                                resourceSlowPostProtectionSetting(),
			"akamai_appsec_slowpost_protection":                      resourceSlowPostProtection(),
			"akamai_appsec_threat_intel":                             resourceThreatIntel(),
			"akamai_appsec_version_notes":                            resourceVersionNotes(),
			"akamai_appsec_waf_mode":                                 resourceWAFMode(),
			"akamai_appsec_waf_protection":                           resourceWAFProtection(),
			"akamai_appsec_wap_selected_hostnames":                   resourceWAPSelectedHostnames(),
		},
	}
	return provider
}

// WithClient sets the client interface function, used for mocking and testing
func WithClient(c appsec.APPSEC) Option {
	return func(p *provider) {
		p.client = c
	}
}

// Client returns the PAPI interface
func (p *provider) Client(meta akamai.OperationMeta) appsec.APPSEC {
	if p.client != nil {
		return p.client
	}
	return appsec.Client(meta.Session())
}

func getAPPSECV1Service(d *schema.ResourceData) (interface{}, error) {
	var section string

	for _, s := range tools.FindStringValues(d, "appsec_section", "config_section") {
		if s != "default" {
			section = s
			break
		}
	}

	if section != "" {
		if err := d.Set("config_section", section); err != nil {
			return nil, fmt.Errorf("%s: %s", tools.ErrValueSet, err.Error())
		}
	}

	return nil, nil
}

func (p *provider) Name() string {
	return "appsec"
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

	_, err := getAPPSECV1Service(d)
	if err != nil {
		return nil
	}

	return nil
}
