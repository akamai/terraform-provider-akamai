provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_rules_template" "test" {
  template_file = "testdata/TestDSRulesTemplate/variables-complex-rules-with-include/template_simple.json"
  variables {
    name  = "variableName"
    value = "$${env.firstVariable}$${env.secondVariable}"
    type  = "string"
  }
  variables {
    name  = "firstVariable"
    value = "first"
    type  = "string"
  }
  variables {
    name  = "secondVariable"
    value = "second"
    type  = "string"
  }
  variables {
    name  = "trickier"
    value = "tricky"
    type  = "string"
  }
  variables {
    name  = "tricky"
    value = "wow"
    type  = "string"
  }
}
