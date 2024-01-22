provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_property_include" "test" {
  contract_id = "ctr_123"
  group_id    = "grp_123"
  product_id  = "prd_test"
  name        = "test_include"
  type        = "INVALID_TYPE"
  rule_format = "INVALID_RULE_FORMAT"
  rules       = "INVALID_JSON"
}