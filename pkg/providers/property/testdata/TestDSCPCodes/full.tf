provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cp_codes" "test" {
  contract_id = "ctr_11"
  group_id    = "grp_22"
}
