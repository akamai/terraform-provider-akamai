provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_property_include" "test" {
  contract_id = "ctr_123"
  group_id    = "grp_123"
  product_id  = "prd_test"
  name        = "test_include"
  type        = "MICROSERVICES"
  rule_format = "v2022-06-28"
  rules       = file("testdata/TestResPropertyInclude/property-snippets/simple_rules.json")
}