data "akamai_property_rules" "rules" {
  rules {
    behavior {
      name = "siteShield"
      option {
        key   = "ssmap"
        value = "mapname.akamai.net"
      }
    }
  }
  variables = <<-EOF
  {
      "variables": [
          {
              "name": "var",
              "description": "desc",
              "value": "val",
              "hidden": true,
              "sensitive": true
          }
      ]
  }
  EOF
}
