provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_property_rules_template" "test" {
  template_file = "testdata/TestDSRulesTemplate/cyclic-dependency/cyclic_dependency_rules.json"
  variables {
    name  = "variableName"
    value = "$${env.$${env.innerOne}}"
    type  = "string"
  }
  variables {
    name  = "innerOne"
    value = "moreInnerOne"
    type  = "string"
  }
  variables {
    name  = "moreInnerOne"
    value = "$${env.variableName}"
    type  = "string"
  }
}