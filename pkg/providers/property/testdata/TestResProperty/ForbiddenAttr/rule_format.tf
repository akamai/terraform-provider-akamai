provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property" "test" {
  name        = "test_property"
  product_id  = "prd_0"
  group_id    = "grp_0"
  contract_id = "ctr_0"

  rule_format = "not empty"
}
