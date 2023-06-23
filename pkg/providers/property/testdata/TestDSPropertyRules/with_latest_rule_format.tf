provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_rules" "rules" {
  contract_id = "ctr_2"
  group_id    = "grp_2"
  property_id = "prp_2"
  rule_format = "latest"
}
