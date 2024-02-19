provider "akamai" {
  edgerc = "../../test/edgerc"
}


resource "akamai_property_bootstrap" "test" {
  name        = "property_name"
  group_id    = "grp_1"
  contract_id = "ctr_222"
  product_id  = "prd_3"
}