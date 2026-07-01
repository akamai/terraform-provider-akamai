provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_security_policy" "main" {
  security_policy_name = "policy"
  config_id            = 111111
}

resource "akamai_appsec_waf_ruleset" "test" {
  config_id          = 111111
  security_policy_id = data.akamai_appsec_security_policy.main.security_policy_id
}