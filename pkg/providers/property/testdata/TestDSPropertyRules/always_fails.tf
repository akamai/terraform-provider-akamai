provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_property_rules" "rules" {
  contract_id = "ctr_test"
  group_id = "grp_test"
  property_id = "prp_test"
}
