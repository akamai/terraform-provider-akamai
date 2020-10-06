provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_cp_code" "test" {
  name     = "test cpcode"
  contract = "ctr_test"
  group    = "grp_test"
  product  = "prd_test"
}
