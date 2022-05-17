provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_property" "test" {
  name        = "test_property"
  contract_id = "ctr_0"
  product_id  = "prd_0"

  group    = "grp1"
  group_id = "grp2"
}
