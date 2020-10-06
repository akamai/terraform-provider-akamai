provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_cp_code" "test" {
  name     = "renamed cpcode"
  contract = "ctr1"
  group    = "grp1"
  product  = "prd1"
}
