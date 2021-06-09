terraform {
  required_version = ">= 0.12"
  required_providers {
    akamai = {
      source = "akamai/akamai/akamai"
      version = "0.9.1"
    }
    local = {
      source = "hashicorp/local"
    }
  }
}

provider "akamai" {
  edgerc = "~/.edgerc"
}

// Example data source and resource definitions: change the configuration name and security policy ID
// below to values that are applicable to your configuration, and edit the supporting JSON files as needed.

data "akamai_appsec_configuration" "appsec_config" {
  name = "Akamai Tools"
}
output "appsec_config_production_version" {
  value = data.akamai_appsec_configuration.appsec_config.production_version
}
output "appsec_config_latest_version" {
  value = data.akamai_appsec_configuration.appsec_config.latest_version
}
output "appsec_config_output_text" {
  value = data.akamai_appsec_configuration.appsec_config.output_text
}

data "akamai_appsec_security_policy" "security_policy" {
  security_policy_name = "akamaitools"
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
}
output "security_policy" {
  value = data.akamai_appsec_security_policy.security_policy.security_policy_id
}
output "securitypolicies" {
  value = data.akamai_appsec_security_policy.security_policy.output_text
}

data "akamai_appsec_selectable_hostnames" "selectable_hostnames" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
}
output "selectable_hostnames_output_text" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames.output_text
}
output "selectable_hostnames_json" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames.hostnames_json
}
output "selectable_hostnames" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames.hostnames
}

data "akamai_appsec_configuration_version" "configuration_version" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
}
output "configuration_version" {
  value = data.akamai_appsec_configuration_version.configuration_version.output_text
}
output "configuration_version_staging_status" {
  value = data.akamai_appsec_configuration_version.configuration_version.staging_status
}
output "configuration_version_production_status" {
  value = data.akamai_appsec_configuration_version.configuration_version.production_status
}
output "configuration_version_output_text" {
  value = data.akamai_appsec_configuration_version.configuration_version.latest_version
}

data "akamai_appsec_export_configuration" "export" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
  version = data.akamai_appsec_configuration.appsec_config.latest_version
  search = ["selectedHosts"]
}
output "export_configuration" {
  value = data.akamai_appsec_export_configuration.export.output_text
}

resource "akamai_appsec_selected_hostnames" "selected_hostnames" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
  hostnames = ["rinaldi.sandbox.akamaideveloper.com", "sujala.sandbox.akamaideveloper.com"]
  mode = "REPLACE"
}

data "akamai_appsec_selected_hostnames" "data_selected_hostnames" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
}
output "selected_hostnames" {
  value = data.akamai_appsec_selected_hostnames.data_selected_hostnames.hostnames
}
output "selected_hostnames_json" {
  value = data.akamai_appsec_selected_hostnames.data_selected_hostnames.hostnames_json
}
output "output_text" {
  value = data.akamai_appsec_selected_hostnames.data_selected_hostnames.output_text
}

resource "akamai_appsec_match_target" "match_targets" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
  match_target = file("${path.module}/match_targets.json")
}
  
resource "akamai_appsec_match_target_sequence" "match_target_sequence" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
  match_target_sequence = file("${path.module}/match_target_sequence.json")
  depends_on = [akamai_appsec_match_target.match_targets]
}

data "akamai_appsec_match_targets" "match_targets" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
}
output "ds_match_targets" {
  value = data.akamai_appsec_match_targets.match_targets.output_text
}

data "local_file" "rules" {
  filename = "${path.module}/custom_rules_simple.json"
}

resource "akamai_appsec_custom_rule" "custom_rule" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
  custom_rule = data.local_file.rules.content
}

resource "akamai_appsec_custom_rule" "custom_rule1" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
  custom_rule = file("${path.module}/custom_rules_simple1.json")
}

data "akamai_appsec_custom_rules" "custom_rule" {
    config_id = data.akamai_appsec_configuration.appsec_config.config_id
}
output "custom_rules" {
  value = data.akamai_appsec_custom_rules.custom_rule.output_text
}

resource "akamai_appsec_activations" "appsecactivations" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
  network = "STAGING"
  notes = "TEST Notes"
  activate = false
  notification_emails = ["plodine@akamai.com"]
}

resource "akamai_appsec_rate_policy" "appsecratepolicy" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
  rate_policy = file("${path.module}/rate_policy.json")
}

data "akamai_appsec_rate_policies" "appsecreatepolicies" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
}
output "ds_rate_policies" {
  value = data.akamai_appsec_rate_policies.appsecreatepolicies.output_text
}

resource  "akamai_appsec_rate_policy_action" "appsecreatepolicysaction" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
  security_policy_id = "AAAA_81230"
  rate_policy_id = akamai_appsec_rate_policy.appsecratepolicy.id
  ipv4_action = "alert"
  ipv6_action = "none"
}

resource "akamai_appsec_slow_post" "appsecslowpostprotectionsettings" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
  security_policy_id = "AAAA_81230"
  slow_rate_action = "alert"
  slow_rate_threshold_rate = 10
  slow_rate_threshold_period = 30
  duration_threshold_timeout = 20
}

data "akamai_appsec_rate_policy_actions" "appsecreatepolicysactions" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
  security_policy_id = "AAAA_81230"
}
output "ds_rate_policy_actions" {
  value = data.akamai_appsec_rate_policy_actions.appsecreatepolicysactions.output_text
}

resource "akamai_appsec_custom_rule_action" "create_custom_rule_action" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
  security_policy_id = "AAAA_81230"
  custom_rule_id = akamai_appsec_custom_rule.custom_rule.custom_rule_id
  custom_rule_action = "alert"
}
output "custom_rule_action" {
  value = akamai_appsec_custom_rule_action.create_custom_rule_action.custom_rule_id
}

data "akamai_appsec_custom_rule_actions" "create_custom_rule_actions" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
  security_policy_id = "AAAA_81230"
}
output "custom_rule_actions" {
  value = data.akamai_appsec_custom_rule_actions.create_custom_rule_actions.output_text
}

resource "akamai_appsec_waf_mode" "waf_mode" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
  security_policy_id = "AAAA_81230"
  mode = "AAG" // KRS
}
output "waf_mode" {
  value = akamai_appsec_waf_mode.waf_mode.output_text
}

data "akamai_appsec_waf_mode" "waf_modes" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
  security_policy_id = "AAAA_81230"
}
output "waf_modes" {
  value = data.akamai_appsec_waf_mode.waf_modes.output_text
}

resource "akamai_appsec_penalty_box" "penalty_box" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
  security_policy_id = "AAAA_81230"
  penalty_box_action = "alert"
  penalty_box_protection = true
}

data "akamai_appsec_penalty_box" "penalty_box" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
  security_policy_id = "AAAA_81230"
}
output "penalty_box" {
  value = data.akamai_appsec_penalty_box.penalty_box.output_text
}

resource "akamai_appsec_waf_protection" "waf_protection" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
  security_policy_id = "AAAA_81230"
  enabled = true
}
output "waf_output" {
  value = akamai_appsec_waf_protection.waf_protection.output_text
}

data "akamai_appsec_security_policy_protections" "security_policy_protections" {
  config_id = data.akamai_appsec_configuration.appsec_config.config_id
  security_policy_id = "AAAA_81230"
}
output "security_policy_protections" {
  value = data.akamai_appsec_security_policy_protections.security_policy_protections.output_text
}
