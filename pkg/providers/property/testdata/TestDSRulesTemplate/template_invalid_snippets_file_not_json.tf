provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_property_rules_template" "test" {
  template_file = "testdata/TestDSRulesTemplate/output/property-snippets/template_invalid_json.json"
}
