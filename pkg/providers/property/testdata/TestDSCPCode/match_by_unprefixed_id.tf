provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cp_code" "test" {
  name     = "234"
  contract = "ctr_test"
  group    = "grp_test"
}
