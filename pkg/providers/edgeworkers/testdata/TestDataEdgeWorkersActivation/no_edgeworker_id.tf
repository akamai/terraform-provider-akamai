provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_edgeworker_activation" "test" {
  network = "PRODUCTION"
}
