provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_property_rules_template" "test" {
  template_file = "property-snippet/template_invalid_json.json"
}
