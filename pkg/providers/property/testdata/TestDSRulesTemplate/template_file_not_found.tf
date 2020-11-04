provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_rules_template" "test" {
  template_file = "testdata/TestDSRulesTemplate/rules/templates/template_file_not_found.json"
  variables {
    name = "criteriaMustSatisfy"
    value = "all"
    type = "string"
  }
  variables {
    name = "name"
    value = "test"
    type = "string"
  }
  variables {
    name = "enabled"
    value = "true"
    type = "bool"
  }
  variables {
    name = "options"
    value = "{\"enabled\":true}"
    type = "jsonBlock"
  }
}
