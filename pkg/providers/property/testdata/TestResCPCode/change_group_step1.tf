provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_cp_code" "test" {
  name     = "test cpcode"
  contract = "ctr1"
  group    = "grp2"
  product  = "prd1"
}
