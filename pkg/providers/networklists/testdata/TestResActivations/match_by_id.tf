provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_networklist_activations" "test" {
  network_list_id     = "86093_AGEOLIST"
  network             = "STAGING"
  notes               = "Test Notes"
  notification_emails = ["user@example.com"]
  sync_point          = 0
}

