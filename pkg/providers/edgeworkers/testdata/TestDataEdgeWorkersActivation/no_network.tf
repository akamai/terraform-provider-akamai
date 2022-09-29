provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_edgeworker_activation" "test" {
  edgewroker_id = 1
}
