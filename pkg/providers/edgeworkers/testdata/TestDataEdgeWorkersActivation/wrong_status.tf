provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_edgeworker_activation" "test" {
  edgeworker_id = 4
  network       = "STAGING"
}
