terraform {
  required_version = ">= 0.12"
  required_providers {
    akamai = {
      source = "akava.io/akava/akamai"
      version = "0.9.1"
    }
    local = {
      source = "hashicorp/local"
    }

  }
}


provider "akamai" {
  edgerc = "~/.edgerc"
  //alias  = "appsec"
 // appsec_section = "default"
}


data "akamai_appsec_security_policy" "appsecsecuritypolicy" {
  name = "akamaitools"
  config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
  version =  data.akamai_appsec_configuration.appsecconfigedge.latest_version
}
output "securitypolicy" {
  value = data.akamai_appsec_security_policy.appsecsecuritypolicy.policy_id
}
output "securitypolicies" {
  value = data.akamai_appsec_security_policy.appsecsecuritypolicy.output_text
}

data "akamai_appsec_configuration" "appsecconfigedge" {
  name    = "Akamai Tools" //Example for EDGE

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

/*
resource "akamai_appsec_configuration_version_clone" "appsecconfigurationclone" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    create_from_version = 3
    rule_update  = true
   }

*/
/*
data "akamai_appsec_selectable_hostnames" "appsecselectablehostnames" {
  config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
  version   = data.akamai_appsec_configuration.appsecconfigedge.latest_version
 
}


output "selectablehostnames_output_text" {
  value = data.akamai_appsec_selectable_hostnames.appsecselectablehostnames.output_text
}

output "selectablehostnames_json" {
  value = data.akamai_appsec_selectable_hostnames.appsecselectablehostnames.hostnames_json
}

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
/*
data "akamai_appsec_export_configuration" "export" {
  config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
  version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
 // search = ["ruleActions","securityPolicies","selectedHosts.tf","selectedHosts","customRules","rulesets","reputationProfiles","ratePolicies","matchTargets"] //"selectableHosts"
 //  search = ["ruleActions","customRules","rulesets","reputationProfiles","ratePolicies","matchTargets","selectedHosts.tf","customRuleActions"] //"selectableHosts"
   search = ["selectedHosts"]
}
output "exportconfig" {
  value = data.akamai_appsec_export_configuration.export.output_text
}
*/
/*
resource "akamai_appsec_selected_hostnames" "appsecselectedhostnames" {
  config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
  version   = data.akamai_appsec_configuration.appsecconfigedge.latest_version
  hostnames = ["rinaldi.sandbox.akamaideveloper.com", "sujala.sandbox.akamaideveloper.com"]
  // hostnames = data.akamai_appsec_selected_hostnames.dataappsecselectedhostnames.hostnames
  // hostnames = ["rinaldi.sandbox.akamaideveloper.com"]
  mode = "REPLACE"
  // mode = "APPEND"
  // mode = "REMOVE"
}
data "akamai_appsec_selected_hostnames" "dataappsecselectedhostnames" {
  config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
  version   = data.akamai_appsec_configuration.appsecconfigedge.latest_version
}
output "selectedhosts" {
  value = data.akamai_appsec_selected_hostnames.dataappsecselectedhostnames.hostnames
}
output "selectedhosts_json" {
  value = data.akamai_appsec_selected_hostnames.dataappsecselectedhostnames.hostnames_json
}
output "output_text" {
  value = data.akamai_appsec_selected_hostnames.dataappsecselectedhostnames.output_text
}*/
/*
resource "akamai_appsec_security_policy_clone" "appsecsecuritypolicyclone" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    create_from_security_policy_id = "LNPD_76189"
    security_policy_name = "PLE Cloned Test for Launchpad"
    security_policy_prefix = "PLE"
    //depends_on = ["akamai_appsec_configuration_clone.appsecconfigurationclone"]
   }
output "secpolicyclone" {
  value = akamai_appsec_security_policy_clone.appsecsecuritypolicyclone.policy_id
}

*/
/*
resource "akamai_appsec_match_target" "appsecmatchtargets" {
  config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
  version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
 //   config_id = akamai_appsec_configuration_clone.appsecconfigurationclone.config_id
  //  version = akamai_appsec_configuration_clone.appsecconfigurationclone.version
    match_target =  file("${path.module}/match_targets.json")
   

   
}*/
/*
resource "akamai_appsec_match_target_sequence" "appsecmatchtargetsequence" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = 11 //data.akamai_appsec_configuration.appsecconfigedge.latest_version  
    type = "website"
   // json = file("${path.module}/match_target_sequence.json")
    sequence_map = {
      2971336 = 1
      2052813 = 2
    }  
    depends_on = ["akamai_appsec_match_target.appsecmatchtargets"]
}
*/
/*
data "akamai_appsec_match_targets" "appsecmatchtargets" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
}
output "ds_match_targets" {
  value = data.akamai_appsec_match_targets.appsecmatchtargets.output_text
}
*/

data "local_file" "rules" {
  filename = "${path.module}/custom_rules_simple.json"
}
resource "akamai_appsec_custom_rule" "appseccustomrule" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    custom_rule = data.local_file.rules.content
}


resource "akamai_appsec_custom_rule" "appseccustomrule1" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    custom_rule = file("${path.module}/custom_rules_simple1.json")
}


data "akamai_appsec_custom_rules" "appseccustomrule" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
}
output "appseccustomrules" {
  value = data.akamai_appsec_custom_rules.appseccustomrule.output_text
}

/*
resource "akamai_appsec_activations" "appsecactivations" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = 4 //data.akamai_appsec_configuration.appsecconfigedge.version
    network = "STAGING"
    notes  = "TEST Notes"
    activate = false
    notification_emails = ["martin@akava.io"]
}*/
/*
resource "akamai_appsec_rate_policy" "appsecratepolicy" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version_number = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    json =  file("${path.module}/rate_policy.json")
}

data "akamai_appsec_rate_policies" "appsecreatepolicies" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
}
output "ds_rate_policies" {
  value = data.akamai_appsec_rate_policies.appsecreatepolicies.output_text
}*/
/*
resource  "akamai_appsec_rate_policy_action" "appsecreatepolicysaction" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
    rate_policy_id = akamai_appsec_rate_policy.appsecratepolicy.id
    ipv4_action = "alert"
    ipv6_action = "none"
}
*/
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
/*
data "akamai_appsec_rate_policy_actions" "appsecreatepolicysactions" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
}
output "ds_rate_policy_actions" {
  value = data.akamai_appsec_rate_policy_actions.appsecreatepolicysactions.output_text
}
*/

resource "akamai_appsec_custom_rule_action" "appsecreatecustomruleaction" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
    custom_rule_id = akamai_appsec_custom_rule.appseccustomrule1.custom_rule_id
    custom_rule_action = "alert"
}
output "customruleaction" {
  value = akamai_appsec_custom_rule_action.appsecreatecustomruleaction.custom_rule_id
}
/*
data "akamai_appsec_custom_rule_actions" "appsecreatecustomruleactions" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
}
output "customruleactions" {
  value = data.akamai_appsec_custom_rule_actions.appsecreatecustomruleactions.output_text
}*/

/*
resource "akamai_appsec_waf_mode" "appsecwafmode" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
    mode = "AAG" //KRS
}

output "appsecwafmode" {
  value = akamai_appsec_waf_mode.appsecwafmode.output_text
}*/
/*
data "akamai_appsec_waf_modes" "appsecwafmodes" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
}

output "appsecwafmodes" {
  value = data.akamai_appsec_waf_modes.appsecwafmodes.output_text
}*/

/*
resource "akamai_appsec_penalty_box" "appsecpenaltybox" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
    action = "alert" 
    penalty_box_protection = true
}

output "appsecpenaltybox" {
  value = akamai_appsec_penalty_box.appsecpenaltybox.output_text
}

data "akamai_appsec_penalty_boxes" "appsecpenaltyboxes" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
}

output "appsecpenaltyboxes" {
  value = data.akamai_appsec_penalty_boxes.appsecpenaltyboxes.output_text
}
*/

resource "akamai_appsec_waf_protection" "appsecwafprotection" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
    enabled = true
}

output "appsecwafmode" {
  value = akamai_appsec_waf_protection.appsecwafprotection.output_text
}

data "akamai_appsec_waf_protection" "appsecwafprotection" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
}

output "appsecwafprotection" {
  value = data.akamai_appsec_waf_protection.appsecwafprotection.output_text
}