provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_penalty_boxes" "test" {
  config_id          = 43253
  version            = 7
  security_policy_id = "AAAA_81230"
}

output "appsecpenaltyboxes" {
  value = data.akamai_appsec_penalty_boxes.test.output_text
}
