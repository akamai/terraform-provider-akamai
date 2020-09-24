provider "akamai" {
  edgerc = "~/.edgerc"
  alias = "appsec"
  //appsec_section = "global"
}
/*
data "akamai_appsec_security_policy" "appsecsecuritypolicy" {
  name = "akamaitools"
  config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
  version =  data.akamai_appsec_configuration.appsecconfigedge.version
}
output "securitypolicy" {
  value = data.akamai_appsec_security_policy.appsecsecuritypolicy.policy_id
}
output "securitypolicies" {
  value = data.akamai_appsec_security_policy.appsecsecuritypolicy.policy_list
}
*/
data "akamai_appsec_configuration" "appsecconfigedge" {
  name = "Akamai Tools" //Example for EDGE
  version = 3
}
output "configsedge_prodversion" {
  value = data.akamai_appsec_configuration.appsecconfigedge.production_version
}
output "configsedge_latestversion" {
  value = data.akamai_appsec_configuration.appsecconfigedge.latest_version
}
output "configsedg_eoutputtext" {
  value = data.akamai_appsec_configuration.appsecconfigedge.output_text
}
output "configsedge_configversion" {
  value = data.akamai_appsec_configuration.appsecconfigedge.version
}
/*
resource "akamai_appsec_configuration_clone" "appsecconfigurationclone" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    create_from_version = data.akamai_appsec_configuration.appsecconfigedge.version
    rule_update  = true
   }
*/
/*
data "akamai_appsec_selectable_hostnames" "appsecselectablehostnames" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
   // active_in_staging = true
   // active_in_production = false
 //    active_in_staging = false
 //   active_in_production = false
 //   active_in_staging = true
   // active_in_production = true
   active_in_staging = false
    active_in_production = true
}
*/
/*
output "selectablehostnames_json" {
  value = data.akamai_appsec_selectable_hostnames.appsecselectablehostnames.hostnames_json
}
*/
/*
output "selectablehostnames" {
  value = data.akamai_appsec_selectable_hostnames.appsecselectablehostnames.hostnames
}
*/
data "akamai_appsec_configuration_version" "appsecconfigversion" {
  config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
  version = 1
}
output "configversion" {
  value = data.akamai_appsec_configuration_version.appsecconfigversion.output_text
}
output "configversionstagingstatus" {
  value = data.akamai_appsec_configuration_version.appsecconfigversion.staging_status
}
output "configversionproductionstatus" {
  value = data.akamai_appsec_configuration_version.appsecconfigversion.production_status
}
output "configversion_output_text" {
  value = data.akamai_appsec_configuration_version.appsecconfigversion.latest_version
}
data "akamai_appsec_export_configuration" "export" {
  config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
  version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
 // search = ["ruleActions","securityPolicies","selectedHosts.tf","selectedHosts","customRules","rulesets","reputationProfiles","ratePolicies","matchTargets"] //"selectableHosts"
   search = ["ruleActions","customRules","rulesets","reputationProfiles","ratePolicies","matchTargets","selectedHosts.tf","customRuleActions"] //"selectableHosts"
 //   search = ["matchTargets"]
}
output "exportconfig" {
  value = data.akamai_appsec_export_configuration.export.output_text
}
resource "akamai_appsec_selected_hostnames" "appsecselectedhostnames" {
 config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    //hostnames = merge(["rinaldi.sandbox.akamaideveloper.com","sujala.sandbox.akamaideveloper.com"]
     // data.akamai_appsec_selected_hostnames.dataappsecselectedhostnames.hostnames
    //  )
hostnames = ["rinaldi.sandbox.akamaideveloper.com","sujala.sandbox.akamaideveloper.com"]
   // hostnames = data.akamai_appsec_selected_hostnames.dataappsecselectedhostnames.hostnames
   // hostnames = ["rinaldi.sandbox.akamaideveloper.com"]
}
/*
data "akamai_appsec_selected_hostnames" "dataappsecselectedhostnames" {
 config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
}
output "selectedhosts" {
  value = data.akamai_appsec_selected_hostnames.dataappsecselectedhostnames.hostnames
}
output "selectedhosts_json" {
  value = data.akamai_appsec_selected_hostnames.dataappsecselectedhostnames.hostnames_json
}
*/
/*
resource "akamai_appsec_security_policy_clone" "appsecsecuritypolicyclone" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    create_from_security_policy = "LNPD_76189"
    policy_name = "PL Cloned Test for Launchpad"
    policy_prefix = "PL"
    depends_on = ["akamai_appsec_configuration_clone.appsecconfigurationclone"]
   }
output "secpolicyclone" {
  value = akamai_appsec_security_policy_clone.appsecsecuritypolicyclone.policy_id
}
*/
/*
data "akamai_contract" "contract" {
}
data "akamai_group" "group" {
}
*/
resource "akamai_appsec_match_targets" "appsecmatchtargets" {
  config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
  version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
 //   config_id = akamai_appsec_configuration_clone.appsecconfigurationclone.config_id
  //  version = akamai_appsec_configuration_clone.appsecconfigurationclone.version
    json =  file("${path.module}/match_targets.json")
   /* type =  "website"
    is_negative_path_match =  false
    is_negative_file_extension_match =  true
    default_file = "NO_MATCH" //"BASE_MATCH" //NO_MATCH
    hostnames =  ["example.com","www.example.net","n.example.com"]
    file_paths =  ["/sssi/*","/cache/aaabbc*","/price_toy/*"]
    file_extensions = ["wmls","jpeg","pws","carb","pdf","js","hdml","cct","swf","pct"]
    security_policy = "AAAA_81230"
*/
    //bypass_network_lists = ["888518_ACDDCKERS","1304427_AAXXBBLIST"]
}
/*
data "local_file" "rules" {
  filename = "${path.module}/custom_rules_simple.json"
}
resource "akamai_appsec_custom_rule" "appseccustomrule" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    rules = data.local_file.rules.content
}
*/
resource "akamai_appsec_custom_rule" "appseccustomrule1" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    rules = file("${path.module}/custom_rules_simple1.json")
}
/*
data "akamai_appsec_custom_rules" "appseccustomrule" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
}
output "appseccustomrules" {
  value = data.akamai_appsec_custom_rules.appseccustomrule.output_text
}
*/
/*
resource "akamai_appsec_activations" "appsecactivations" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = 4 //data.akamai_appsec_configuration.appsecconfigedge.version
    network = "STAGING"
    notes  = "TEST Notes"
    activate = true
    notification_emails = ["martin@akava.io"]
}*/
resource "akamai_appsec_rate_policy" "appsecratepolicy" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version_number = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    json =  file("${path.module}/rate_policy.json")
}
resource  "akamai_appsec_rate_policy_action" "appsecreatepolicysaction" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
    rate_policy_id = akamai_appsec_rate_policy.appsecratepolicy.id
    ipv4_action = "alert"
    ipv6_action = "none"
}
/*
resource "akamai_appsec_slow_post_protection_settings" "appsecslowpostprotectionsettings" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
    slow_rate_action = "alert"
    slow_rate_threshold_rate = 10
    slow_rate_threshold_period = 30
    duration_threshold_timeout = 20
}
*/
data "akamai_appsec_rate_policy_actions" "appsecreatepolicysactions" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
}
output "ds_rate_policy_actions" {
  value = data.akamai_appsec_rate_policy_actions.appsecreatepolicysactions.output_text
}
resource "akamai_appsec_custom_rule_action" "appsecreatecustomruleaction" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
    rule_id = akamai_appsec_custom_rule.appseccustomrule1.rule_id
    custom_rule_action = "alert"
}
output "customruleaction" {
  value = akamai_appsec_custom_rule_action.appsecreatecustomruleaction.rule_id
}
/*
data "akamai_appsec_custom_rule_actions" "appsecreatecustomruleactions" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
  //  rule_id = akamai_appsec_custom_rule.appseccustomrule1.rule_id
   // custom_rule_action = "alert"
}
output "customruleactions" {
  value = data.akamai_appsec_custom_rule_actions.appsecreatecustomruleactions.output_text
}*/