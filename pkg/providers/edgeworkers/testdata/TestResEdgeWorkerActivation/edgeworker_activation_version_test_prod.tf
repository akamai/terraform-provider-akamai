provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_edgeworker_activation" "test" {
  edgeworker_id = 1234
  network       = "PRODUCTION"
  version       = "test"
}