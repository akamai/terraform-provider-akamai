provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_hostnames_diff" "diff" {
  contract_id = "ctr_1"
  group_id    = "grp_1"
  property_id = "prp_1"
}