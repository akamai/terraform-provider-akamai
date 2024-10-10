provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property" "test" {
  contract_id = "ctr_1"
  group_id    = "grp_2"
  product_id  = "prd_3"
}
