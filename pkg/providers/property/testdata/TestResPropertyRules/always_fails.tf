provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property_rules" "rules" {}
