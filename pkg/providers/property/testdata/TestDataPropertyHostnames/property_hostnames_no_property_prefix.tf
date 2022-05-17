provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_property_hostnames" "akaprophosts" {
  group_id    = "grp_test"
  contract_id = "ctr_test"
  property_id = "test"
}