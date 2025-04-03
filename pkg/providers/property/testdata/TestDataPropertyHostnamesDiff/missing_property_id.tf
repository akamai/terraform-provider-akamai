provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_hostnames_diff" "diff" {
  contract_id = "1"
  group_id    = "1"
}