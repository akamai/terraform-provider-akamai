provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_rules_template" "test" {
  template_file = "testdata/TestDSRulesTemplate/nested-pattern-rules/main.json"
  variables {
    name  = "includesBundle"
    value = jsonencode(formatlist("#include:%s", ["nested/include1.json", "nested/include2.json", "nested/include2.json"]))
    type  = "jsonBlock"
  }
}
