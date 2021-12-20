provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_application_load_balancer_match_rule" "test" {
  match_rules {
    name      = "rule1"
    start     = 10
    end       = 10000
    match_url = "example.com"
    matches {
      case_sensitive = true
      match_type     = "cookie"
      object_match_value {
        name                = "abcde"
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
    }
  }
}