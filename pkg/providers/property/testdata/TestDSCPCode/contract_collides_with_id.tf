provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cp_code" "test" {
  name        = "cpc_test2"
  contract    = "ctr_test"
  contract_id = "ctr_test"
  group       = "grp_test"
}
