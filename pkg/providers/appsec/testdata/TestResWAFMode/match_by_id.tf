provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_waf_mode" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
  mode               = "AAG"
}

output "configsedge_post_output_text" {
  value = akamai_appsec_waf_mode.test.output_text
}

