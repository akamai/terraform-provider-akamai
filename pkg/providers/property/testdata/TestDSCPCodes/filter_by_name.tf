provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cp_codes" "test" {
  filter_by_name = "test_name"
  contract_id    = "ctr_11"
  group_id       = "grp_22"
}
