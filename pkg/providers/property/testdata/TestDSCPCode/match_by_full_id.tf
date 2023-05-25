provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cp_code" "test" {
  name     = "cpc_234"
  contract = "ctr_test"
  group    = "grp_test"
}
