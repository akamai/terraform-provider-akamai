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
}
