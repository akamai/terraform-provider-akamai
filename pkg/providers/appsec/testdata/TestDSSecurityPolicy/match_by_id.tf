provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

data "akamai_appsec_security_policy" "test" {
  security_policy_name = "akamaitools"
  config_id            = 43253
}

output "securitypolicy" {
  value = data.akamai_appsec_security_policy.test.security_policy_id
}

output "securitypolicies" {
  value = data.akamai_appsec_security_policy.test.security_policy_id_list
}

