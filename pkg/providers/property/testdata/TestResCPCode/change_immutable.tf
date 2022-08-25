provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_cp_code" "test" {
  name     = "renamed cpcode"
  contract = "ctr_2"
  group_id = "grp_2"
  product  = "prd_2"
}
