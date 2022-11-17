provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_property_include_rules" "test" {
  contract_id = "ctr_1"
  group_id    = "grp_2"
  include_id  = "12345"
  version     = 1
}

