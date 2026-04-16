provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}
resource "akamai_appsec_penalty_box" "test" {
  config_id              = 43253
  penalty_box_action     = "none"
  penalty_box_protection = false
}
