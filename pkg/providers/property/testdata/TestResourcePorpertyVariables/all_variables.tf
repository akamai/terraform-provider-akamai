provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property_variables" "var" {
  variables {
    variable {
      name = "var_1"
      hidden = true
      sensitive = true
      description = "desc"
      value = "val"
    }
    variable {
      name = "var_2"
      hidden = false
      sensitive = false
      description = "desc"
      value = "val_2"
    }
  }
}