provider "akamai" {
  edgerc = "~/.edgerc"
}




resource "akamai_appsec_siem_settings" "test" {
  config_id = 43253
  version = 7
  enable_siem = true
  enable_for_all_policies = false
  enable_botman_siem = true
  siem_id = 1 //data.akamai_appsec_siem_definitions.siem_definition.id
  security_policy_ids = [12345]//data.akamai_appsec_security_policy.security_policies.policy_ids
}


