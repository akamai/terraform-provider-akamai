provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_cp_code" "test" {
  name     = "renamed cpcode"
  contract = "ctr_1"
  group    = "grp_1"
  product  = "prd_1"
}
