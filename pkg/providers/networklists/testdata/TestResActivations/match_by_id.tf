provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_networklist_activations" "test" {
  name                = "Network list test"
  network             = "STAGING"
  notes               = "TEST Notes"
  activate            = true
  notification_emails = ["martin@akava.io"]
}

