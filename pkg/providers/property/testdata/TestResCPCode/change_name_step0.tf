provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_cp_code" "test" {
  name        = "test cpcode"
  contract_id = "ctr_1"
  group_id    = "grp_1"
  product_id  = "prd_1"
}
