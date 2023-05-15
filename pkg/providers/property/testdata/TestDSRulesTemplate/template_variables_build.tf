provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_property_rules_template" "test" {
  template_file = "testdata/TestDSRulesTemplate/variables-complex-rules/template_simple.json"
  variables {
    name  = "variableName"
    value = "$${env.firstVariable}$${env.secondVariable}"
    type  = "string"
  }
  variables {
    name  = "variableNameWithSpace"
    value = "$${env.firstVariable} $${env.secondVariable}"
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
  variables {
    name  = "trickyBro"
    value = true
    type  = "bool"
  }
  variables {
    name  = "enabled"
    value = true
    type  = "bool"
  }
  variables {
    name  = "enabledString"
    value = "true"
    type  = "string"
  }
  variables {
    name  = "someNumber"
    value = 25
    type  = "number"
  }
  variables {
    name = "simpleJSON"
    value = jsonencode(
      {
        "allowSampling" : true,
        "cookies" : {
          "type" : "all"
        },
        "customHeaders" : {
          "type" : "all"
        }
      }
    )
    type = "jsonBlock"
  }
  variables {
    name = "nestedJSON"
    type = "jsonBlock"
    value = jsonencode(
      {
        "enabled" : "$${env.enabled}",
        "name" : "$${env.variableName}"
        "someNumber" : "$${env.someNumber}",
        "andJSON" : "$${env.simpleJSON}"
      }
    )
  }

}