provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_property_rules_template" "test" {
  template_file = "testdata/TestDSRulesTemplate/rules/property-snippets/template_in.json"
  variables {
    name  = "criteriaMustSatisfy"
    value = "all"
    type  = "test"
  }
}
