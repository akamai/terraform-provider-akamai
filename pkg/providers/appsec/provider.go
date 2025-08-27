// Package appsec contains implementation for Akamai Terraform sub-provider responsible for Application Security
package appsec

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers appsec resources and data sources
	Subprovider struct {
		client appsec.APPSEC
	}

	option func(p *Subprovider)
)

var (
	once sync.Once

	inst *Subprovider
)

var _ subprovider.Subprovider = &Subprovider{}

// NewSubprovider returns a new appsec subprovider
func NewSubprovider(opts ...option) *Subprovider {
	once.Do(func() {
		inst = &Subprovider{}

		for _, opt := range opts {
			opt(inst)
		}
	})

	return inst
}

// Client returns the APPSEC interface
func (p *Subprovider) Client(meta meta.Meta) appsec.APPSEC {
	if p.client != nil {
		return p.client
	}
	return appsec.Client(meta.Session())
}

// SDKResources returns the appsec resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_appsec_aap_selected_hostnames":                   resourceAAPSelectedHostnames(),
		"akamai_appsec_activations":                              resourceActivations(),
		"akamai_appsec_advanced_settings_ase_penalty_box":        resourceAdvancedSettingsAsePenaltyBox(),
		"akamai_appsec_advanced_settings_attack_payload_logging": resourceAdvancedSettingsAttackPayloadLogging(),
		"akamai_appsec_advanced_settings_evasive_path_match":     resourceAdvancedSettingsEvasivePathMatch(),
		"akamai_appsec_advanced_settings_logging":                resourceAdvancedSettingsLogging(),
		"akamai_appsec_advanced_settings_pii_learning":           resourceAdvancedSettingsPIILearning(),
		"akamai_appsec_advanced_settings_pragma_header":          resourceAdvancedSettingsPragmaHeader(),
		"akamai_appsec_advanced_settings_prefetch":               resourceAdvancedSettingsPrefetch(),
		"akamai_appsec_advanced_settings_request_body":           resourceAdvancedSettingsRequestBody(),
		"akamai_appsec_advanced_settings_ja4_fingerprint":        resourceAdvancedSettingsJA4Fingerprint(),
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
		"akamai_appsec_eval_penalty_box_conditions":              resourceEvalPenaltyBoxConditions(),
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
		"akamai_appsec_penalty_box_conditions":                   resourcePenaltyBoxConditions(),
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
		"akamai_appsec_security_policy_default_protections":      resourceSecurityPolicyDefaultProtections(),
		"akamai_appsec_security_policy_rename":                   resourceSecurityPolicyRename(),
		"akamai_appsec_siem_settings":                            resourceSiemSettings(),
		"akamai_appsec_slow_post":                                resourceSlowPostProtectionSetting(),
		"akamai_appsec_slowpost_protection":                      resourceSlowPostProtection(),
		"akamai_appsec_threat_intel":                             resourceThreatIntel(),
		"akamai_appsec_version_notes":                            resourceVersionNotes(),
		"akamai_appsec_waf_mode":                                 resourceWAFMode(),
		"akamai_appsec_waf_protection":                           resourceWAFProtection(),
	}
}

// SDKDataSources returns the appsec data sources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_appsec_aap_selected_hostnames":                   dataSourceAAPSelectedHostnames(),
		"akamai_appsec_advanced_settings_ase_penalty_box":        dataSourceAdvancedSettingsAsePenaltyBox(),
		"akamai_appsec_advanced_settings_attack_payload_logging": dataSourceAdvancedSettingsAttackPayloadLogging(),
		"akamai_appsec_advanced_settings_evasive_path_match":     dataSourceAdvancedSettingsEvasivePathMatch(),
		"akamai_appsec_advanced_settings_logging":                dataSourceAdvancedSettingsLogging(),
		"akamai_appsec_advanced_settings_pii_learning":           dataSourceAdvancedSettingsPIILearning(),
		"akamai_appsec_advanced_settings_pragma_header":          dataSourceAdvancedSettingsPragmaHeader(),
		"akamai_appsec_advanced_settings_prefetch":               dataSourceAdvancedSettingsPrefetch(),
		"akamai_appsec_advanced_settings_request_body":           dataSourceAdvancedSettingsRequestBody(),
		"akamai_appsec_advanced_settings_ja4_fingerprint":        dataSourceAdvancedSettingsJA4Fingerprint(),
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
		"akamai_appsec_eval_penalty_box_conditions":              dataSourceEvalPenaltyBoxConditions(),
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
		"akamai_appsec_penalty_box_conditions":                   dataSourcePenaltyBoxConditions(),
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
		"akamai_appsec_siem_definitions":                         dataSourceSiemDefinitions(),
		"akamai_appsec_siem_settings":                            dataSourceSiemSettings(),
		"akamai_appsec_slow_post":                                dataSourceSlowPostProtectionSettings(),
		"akamai_appsec_threat_intel":                             dataSourceThreatIntel(),
		"akamai_appsec_tuning_recommendations":                   dataSourceTuningRecommendations(),
		"akamai_appsec_version_notes":                            dataSourceVersionNotes(),
		"akamai_appsec_waf_mode":                                 dataSourceWAFMode(),
	}
}

// FrameworkResources returns the appsec resources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{
		NewRapidRulesResource,
	}
}

// FrameworkDataSources returns the appsec data sources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewRapidRulesDataSource,
	}
}
