provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cp_codes" "test" {
  group_id = "grp_22"
}
