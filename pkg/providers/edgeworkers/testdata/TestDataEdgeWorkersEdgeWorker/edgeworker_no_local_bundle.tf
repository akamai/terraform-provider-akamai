provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_edgeworker" "test" {
  edgeworker_id = 1
}
