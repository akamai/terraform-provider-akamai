data "akamai_property_rules" "rules" {
  rules {
    rule {
        name = "child"
        criteria {
            name = "criteriaAll"
        }
    }
    behavior {
      name = "siteShield"
      option {
        key   = "ssmap"
        value = "mapname.akamai.net"
      }
    }
  }
}
