provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property_variables" "var" {
  variables {
    variable {
      name = "var_1_updated"
      hidden = true
      sensitive = true
      description = "desc"
      value = "val"
    }
    variable {
      name = "var_2_updated"
      hidden = false
      sensitive = false
      description = "desc"
      value = "val_2"
    }
  }
}