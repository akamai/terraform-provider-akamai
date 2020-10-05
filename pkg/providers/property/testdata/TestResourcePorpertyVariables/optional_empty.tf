provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property_variables" "var" {
  variables {
    variable {
      name = "var_1"
      hidden = true
      sensitive = true
    }
  }
}