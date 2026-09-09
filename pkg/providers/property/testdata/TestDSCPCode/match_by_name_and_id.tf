provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cp_code" "test" {
  cp_code_name = "test cpcode"
  cp_code_id   = "cpc_234"
  contract_id  = "ctr_11"
  group_id     = "grp_22"
}