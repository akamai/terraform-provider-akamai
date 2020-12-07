provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_property_rules" "rules" {
  property_id = ""
}
