provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_property_rules_template" "test" {
  template {
    template_dir = "testdata/TestDSRulesTemplate/rules/property-snippets/"
  }
  var_definition_file = "testdata/TestDSRulesTemplate/rules/variables/variableDefinitions.json"
  var_values_file     = "testdata/TestDSRulesTemplate/rules/variables/variables.json"
}
