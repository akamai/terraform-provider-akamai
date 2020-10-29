provider "akamai" {
  edgerc = "~/.edgerc"
}


data "akamai_appsec_security_policy" "test" {
  name = "akamaitools" 
  config_id = 43253
  version =  7
}

output "securitypolicy" {
  value = data.akamai_appsec_security_policy.test.policy_id
}

output "securitypolicies" {
  value = data.akamai_appsec_security_policy.test.policy_list
}
