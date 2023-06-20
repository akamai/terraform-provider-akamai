provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cp_code" "test" {
  name        = "234"
  contract_id = "ctr_test"
  group_id    = "grp_test"
}
