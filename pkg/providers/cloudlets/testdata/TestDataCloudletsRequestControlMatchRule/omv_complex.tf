provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudlets_request_control_match_rule" "test" {
  match_rules {
    name       = "complex_simple_rule"
    start      = 1
    end        = 2
    allow_deny = "allow"
    matches {
      match_type     = "method"
      match_operator = "contains"
      case_sensitive = true
      negate         = false
      check_ips      = ""
      object_match_value {
        type  = "simple"
        value = ["GET", "POST"]
      }
    }
  }
  match_rules {
    name       = "complex_object_rule"
    start      = 2
    end        = 3
    allow_deny = "allow"
    matches {
      match_type     = "header"
      match_operator = "exists"
      case_sensitive = true
      negate         = true
      check_ips      = "CONNECTING_IP XFF_HEADERS"
      object_match_value {
        type                = "object"
        name                = "Accept"
        name_case_sensitive = true
        name_has_wildcard   = false
        options {
          value                = ["text/html*", "text/css*"]
          value_has_wildcard   = false
          value_case_sensitive = true
          value_escaped        = false
        }
      }
    }
  }
}