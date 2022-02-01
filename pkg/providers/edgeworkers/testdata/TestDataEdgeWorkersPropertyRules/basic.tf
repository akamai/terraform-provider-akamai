provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_edgeworkers_property_rules" "test" {
  edgeworker_id = 123
}