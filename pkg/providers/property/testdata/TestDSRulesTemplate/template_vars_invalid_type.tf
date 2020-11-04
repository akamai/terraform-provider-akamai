provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_rules_template" "test" {
  template_file = "testdata/TestDSRulesTemplate/rules/templates/template_in.json"
  variables {
    name = "criteriaMustSatisfy"
    value = "all"
    type = "test"
  }
}
