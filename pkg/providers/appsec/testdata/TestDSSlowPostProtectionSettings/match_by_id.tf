provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_slow_post" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
}

output "configsedge_post_output_text" {
  value = data.akamai_appsec_slow_post.test.output_text
}

