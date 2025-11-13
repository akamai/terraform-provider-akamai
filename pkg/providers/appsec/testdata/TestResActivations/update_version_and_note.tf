provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_activations" "test" {
  config_id           = 43253
  version             = 8
  network             = "STAGING"
  note                = "Updated Notes"
  notification_emails = ["user@example.com"]
}