provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

# Simulates: security_policy_id = data.akamai_appsec_security_policy.new_policy.security_policy_id
# The local value holds the same ID as the inline value used in the first step,
# reproducing the pattern where a user switches from an inline string to a
# data-source / local reference that resolves to the identical known value.
locals {
  security_policy_id = "2222_333333"
}

resource "akamai_appsec_waf_ruleset" "test" {
  config_id          = 111111
  security_policy_id = local.security_policy_id
}

