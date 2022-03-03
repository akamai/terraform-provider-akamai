provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_edgeworkers_activation" "test" {
  edgeworker_id = 1234
  network       = "STAGING"
  version       = "test"
}