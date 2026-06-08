provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

variable "security_policy_id_v1" {
  type    = string
  default = "2222_333333"
}

resource "akamai_appsec_waf_ruleset" "test" {
  config_id          = 111111
  security_policy_id = var.security_policy_id_v1
}


