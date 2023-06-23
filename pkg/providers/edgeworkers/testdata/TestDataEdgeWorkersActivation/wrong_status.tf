provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_edgeworker_activation" "test" {
  edgeworker_id = 4
  network       = "STAGING"
}
