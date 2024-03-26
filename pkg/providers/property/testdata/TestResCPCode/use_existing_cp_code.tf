provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cp_code" "test" {
  name        = "test cpcode"
  contract_id = "ctr_test"
  group_id    = "grp_test"
  product_id  = "prd_test"
}
