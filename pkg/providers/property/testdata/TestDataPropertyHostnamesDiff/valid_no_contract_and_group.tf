provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_hostnames_diff" "diff" {
  property_id = "prp_1"
}