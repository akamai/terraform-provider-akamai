provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_cp_code" "test" {
  name     = "renamed cpcode"
  contract = "ctr_2"
  group    = "grp_2"
  product  = "prd_2"
}
