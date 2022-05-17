provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_activations" "test" {
  config_id           = 43253
  version             = 7
  network             = "STAGING"
  notes               = "TEST Notes"
  activate            = true
  notification_emails = ["martin@email.io"]
}

