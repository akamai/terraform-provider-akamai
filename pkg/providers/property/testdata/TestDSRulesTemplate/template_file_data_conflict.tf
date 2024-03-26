provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_rules_template" "test" {
  template_file = "testdata/TestDSRulesTemplate/rules/property-snippets/template_in.json"
  template {
    template_data = <<EOT
{
  "rules": {
    "name": "$${env.name}",
    "children": [
      "#include:snippets/some-template.json"
    ]
  }
}
EOT
    template_dir  = "testdata/TestDSRulesTemplate/rules/property-snippets/"
  }
  var_definition_file = "testdata/TestDSRulesTemplate/rules/variables/variableDefinitions.json"
  var_values_file     = "testdata/TestDSRulesTemplate/rules/variables/variables.json"
}
