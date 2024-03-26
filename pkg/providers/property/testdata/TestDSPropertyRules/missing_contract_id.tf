provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_rules" "rules" {
  group_id    = "grp_2"
  property_id = "prp_2"
}
