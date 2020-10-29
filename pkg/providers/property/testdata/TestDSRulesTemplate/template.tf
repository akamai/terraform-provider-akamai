provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_rules_template" "test" {
  template_file = "testdata/TestDSRulesTemplate/rules/rules.tmpl"
  variables = {
    "name" = "test"
    "criteriaMustSatisfy" = "test-criteria"
  }
  template_dir = "templates"
}
