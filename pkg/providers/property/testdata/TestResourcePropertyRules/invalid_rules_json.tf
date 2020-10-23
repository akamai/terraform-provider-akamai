provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property_rules" "rules" {
  contract_id = "1"
  group_id = "1"
  property_id = "1"
  rules = "abc"
}
