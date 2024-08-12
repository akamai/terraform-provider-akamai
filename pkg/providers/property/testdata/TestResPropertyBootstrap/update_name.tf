provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_bootstrap" "test" {
  name        = "property_name2"
  group_id    = "grp_1"
  contract_id = "ctr_2"
  product_id  = "prd_3"
}