provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_rules_template" "test" {
  template_file = "testdata/TestDSRulesTemplate/child-with-child-rules/main.json"
  variables {
    name  = "variableName"
    value = "simpleName"
    type  = "string"
  }
  variables {
    name  = "include1Name"
    value = "include1.json"
    type  = "string"
  }
  variables {
    name  = "include2Name"
    value = "include2.json"
    type  = "string"
  }
}
