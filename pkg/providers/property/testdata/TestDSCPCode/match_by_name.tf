provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cp_code" "test" {
  name        = "test cpcode"
  contract_id = "ctr_11"
  group_id    = "grp_22"
}
