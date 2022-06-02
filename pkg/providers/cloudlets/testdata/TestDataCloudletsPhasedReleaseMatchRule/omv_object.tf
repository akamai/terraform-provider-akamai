provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_phased_release_match_rule" "test" {
  match_rules {
    name = "rule2"
    forward_settings {
      origin_id = "1234"
      percent   = 30
    }
  }
  match_rules {
    name      = "rule1"
    start     = 10
    end       = 10000
    match_url = "example.com"
    matches {
      case_sensitive = true
      match_type     = "header"
      object_match_value {
        type                = "object"
        name                = "Accept"
        name_case_sensitive = true
        name_has_wildcard   = false
        options {
          value                = ["asfas"]
          value_has_wildcard   = false
          value_case_sensitive = true
          value_escaped        = false
        }
      }
    }
    forward_settings {
      origin_id = "33"
      percent   = 20
    }
  }
}