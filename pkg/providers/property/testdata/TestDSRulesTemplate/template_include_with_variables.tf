provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_property_rules_template" "test" {
  template_file = "testdata/TestDSRulesTemplate/variables-with-includes-rules/main.json"
  variables {
    name  = "variableName"
    value = "simpleName"
    type  = "string"
  }
  variables {
    name  = "includeName"
    value = "include1.json"
    type  = "string"
  }
}
