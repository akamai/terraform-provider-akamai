provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_rules" "rules" {
  property_id = ""
}
