provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_property_rules" "rules" {
  group_id    = "grp_2"
  contract_id = ""
  property_id = "prp_2"
}
