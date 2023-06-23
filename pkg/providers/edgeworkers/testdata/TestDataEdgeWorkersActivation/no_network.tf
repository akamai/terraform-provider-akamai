provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_edgeworker_activation" "test" {
  edgewroker_id = 1
}
