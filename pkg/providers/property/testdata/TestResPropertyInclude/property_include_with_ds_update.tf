provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_property_include" "test" {
  contract_id = "ctr_123"
  group_id    = "grp_123"
  product_id  = "prd_test"
  name        = "test include"
  type        = "MICROSERVICES"
  rule_format = "v2022-06-28"
  rules       = data.akamai_property_rules_template.rules.json
}

data "akamai_property_include_rules" "rules" {
  contract_id = "ctr_123"
  group_id    = "grp_123"
  include_id  = akamai_property_include.test.id
  version     = akamai_property_include.test.latest_version
}


data "akamai_property_rules_template" "rules" {
  template_file = "testdata/TestResPropertyInclude/property-snippets/simple_rules.json"
}