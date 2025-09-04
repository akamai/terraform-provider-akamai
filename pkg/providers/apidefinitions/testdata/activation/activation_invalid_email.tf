provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_apidefinitions_activation" "a1" {
  api_id                  = 1
  version                 = 1
  network                 = "STAGING"
  notification_recipients = ["user"]
  notes                   = "Notes"
}