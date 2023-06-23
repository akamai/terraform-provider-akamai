provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_edgeworkers_property_rules" "test" {
  edgeworker_id = 123
}