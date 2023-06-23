provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_include" "test" {
  contract_id = "ctr_1234"
  group_id    = "grp_1234"
  product_id  = "prd_test2"
  name        = "test_include"
  type        = "MICROSERVICES"
  rule_format = "v2022-06-28"
  rules       = data.akamai_property_rules_template.rules.json
}


data "akamai_property_rules_template" "rules" {
  template_file = "testdata/TestResPropertyInclude/property-snippets/simple_rules.json"
}