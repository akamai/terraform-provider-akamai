provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_rules_template" "test" {
  template_file = "testdata/TestDSRulesTemplate/rules/templates/template_in.json"
  var_definition_file = "invalid_path"
  var_values_file = "testdata/TestDSRulesTemplate/rules/variables/variables.json"
}
