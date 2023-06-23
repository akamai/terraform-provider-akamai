provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_hostnames" "akaprophosts" {
  group_id    = "test"
  contract_id = "ctr_test"
  property_id = "prp_test"
}