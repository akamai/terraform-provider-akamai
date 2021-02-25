provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_property_rules_template" "test" {
  template_file = "testdata/TestDSRulesTemplate/rules/property-snippets/template_in.json"
  var_definition_file = "testdata/TestDSRulesTemplate/rules/variables/variableDefinitions.json"
  var_values_file = "testdata/TestDSRulesTemplate/rules/variables/variables.json"
}
