provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_rules_template" "rules1" {
  template_file = "testdata/TestDSRulesTemplate/multiple-templates/snippet1.json"
  variables {
    name  = "test1"
    value = "abc"
    type  = "string"
  }
}

data "akamai_property_rules_template" "rules2" {
  template_file = "testdata/TestDSRulesTemplate/multiple-templates/snippet2.json"
  variables {
    name  = "test2"
    value = "cba"
    type  = "string"
  }
}



