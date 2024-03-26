provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

variable "a_list" {
  type    = list(string)
  default = ["abc", "def"]
}

data "akamai_property_rules_template" "test" {
  template_file = "testdata/TestDSRulesTemplate/rules-with-list/main.json"
  variables {
    name  = "aList"
    value = jsonencode(var.a_list)
    type  = "jsonBlock"
  }
}
