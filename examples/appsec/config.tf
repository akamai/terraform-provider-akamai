provider "akamai" {
  edgerc = "~/.edgerc"
  alias = "appsec"
  //appsec_section = "global"
}

/*
data "akamai_appsec_configuration" "appsecconfig" {
  name = "Akamai Tools"
  
}

output "configs" {
  value = data.akamai_appsec_configuration.appsecconfig.config_id
}
*/
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
}*/
/*
output "selectablehostnames" {
  value = data.akamai_appsec_selectable_hostnames.appsecselectablehostnames.hostnames_json
}

output "selectablehostnames" {
  value = data.akamai_appsec_selectable_hostnames.appsecselectablehostnames.hostnames
}
*/
output "configsedge" {
  value = data.akamai_appsec_configuration.appsecconfigedge.config_id
}

output "configsedgelatestversion" {
  value = data.akamai_appsec_configuration.appsecconfigedge.latest_version 
}

output "configsedgeconfiglist" {
  value = data.akamai_appsec_configuration.appsecconfigedge.config_list
}

output "configsedgeconfigversion" {
  value = data.akamai_appsec_configuration.appsecconfigedge.version
}

data "akamai_appsec_configuration_version" "appsecconfigversion" {
  config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
  version = 1
  
}

output "configversion" {
  value = data.akamai_appsec_configuration_version.appsecconfigversion.version_list
}

output "configversionstagingstatus" {
  value = data.akamai_appsec_configuration_version.appsecconfigversion.staging_status
}

output "configversionproductionstatus" {
  value = data.akamai_appsec_configuration_version.appsecconfigversion.production_status
}

/*data "akamai_appsec_export_configuration" "export" {
  config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
  version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
  search = ["ruleActions","securityPolicies","selectedHosts.tf","selectedHosts","customRules","rulesets","reputationProfiles","ratePolicies","matchTargets"] //"selectableHosts"
  
}

output "exportconfig" {
  value = data.akamai_appsec_export_configuration.export.output_text
}
*/

/* 
resource "akamai_appsec_selected_hostnames" "appsecselectedhostnames" {
    config_id = akamai_appsec_configuration_clone.appsecconfigurationclone.config_id
    version = akamai_appsec_configuration_clone.appsecconfigurationclone.version
    hostnames = ["rinaldi.sandbox.akamaideveloper.com","sujala.sandbox.akamaideveloper.com"] 
    //hostnames = data.akamai_appsec_selectable_hostnames.appsecselectablehostnames.hostnames 
   // hostnames = ["rinaldi.sandbox.akamaideveloper.com"]  
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
/*
resource "akamai_appsec_match_targets" "appsecmatchtargets" {
  config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
  version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
  
 //   config_id = akamai_appsec_configuration_clone.appsecconfigurationclone.config_id
  //  version = akamai_appsec_configuration_clone.appsecconfigurationclone.version
   // type =  "website"
    json =  file("${path.module}/match_targets.json")
     sequence =  1
    is_negative_path_match =  false
    is_negative_file_extension_match =  true
    default_file = "NO_MATCH" //"BASE_MATCH" //NO_MATCH
    hostnames =  ["example.com","www.example.net","n.example.com"]
    file_paths =  ["/sssi/*","/cache/aaabbc*","/price_toy/*"]
    file_extensions = ["wmls","jpeg","pws","carb","pdf","js","hdml","cct","swf","pct"]
    security_policy = "f1rQ_106946"

    //bypass_network_lists = ["888518_ACDDCKERS","1304427_AAXXBBLIST"]
}
*/

data "local_file" "rules" {
  filename = "${path.module}/custom_rules_simple.json"
}

resource "akamai_appsec_custom_rule" "appseccustomrule" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    rules = data.local_file.rules.content
}

resource "akamai_appsec_custom_rule" "appseccustomrule1" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    rules = file("${path.module}/custom_rules_simple1.json")
}


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