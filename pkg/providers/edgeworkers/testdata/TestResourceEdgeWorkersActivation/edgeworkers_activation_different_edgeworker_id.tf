provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_edgeworkers_activation" "test" {
  edgeworker_id = 4321
  network       = "STAGING"
  version       = "test"
}