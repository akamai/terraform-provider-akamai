provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_rules_template" "test" {
  template_file = "testdata/TestDSRulesTemplate/cyclic-dependency/cyclic_dependency_rules.json"
  variables {
    name  = "variableName"
    value = "$${env.fatherchild}"
    type  = "string"
  }
  variables {
    name  = "fatherchild"
    value = "$${env.variableName}"
    type  = "string"
  }
}