provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_rules_builder" "default" {
  rules_v2023_01_05 {
    name = "default"
    behavior {
      restrict_object_caching {}
      origin {}
    }
  }
}
