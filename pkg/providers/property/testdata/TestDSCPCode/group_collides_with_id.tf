provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cp_code" "test" {
  name     = "cpc_test2"
  contract = "ctr_test"
  group    = "grp_test"
  group_id    = "grp_test_id"
}
