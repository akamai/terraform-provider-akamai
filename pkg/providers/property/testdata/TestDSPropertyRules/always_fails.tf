provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_property_rules" "rules" {
  property_id = ""
}
