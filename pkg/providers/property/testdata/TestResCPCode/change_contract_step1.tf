provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_cp_code" "test" {
  name     = "test cpcode"
  contract = "2"
  group    = "1"
  product  = "prd_1"
}
