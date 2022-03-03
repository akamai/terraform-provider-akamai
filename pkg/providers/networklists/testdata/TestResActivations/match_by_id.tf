provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_networklist_activations" "test" {
  name                = "Network list test"
  network             = "STAGING"
  notes               = "TEST Notes"
  notification_emails = ["user@example.com"]
}

