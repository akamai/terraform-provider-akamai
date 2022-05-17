provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_audience_segmentation_match_rule" "test" {
  match_rules {
    name      = "complex_simple_rule"
    start     = 1
    end       = 2
    match_url = "example1.com"
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

    forward_settings {
      origin_id                 = "origin_1"
      use_incoming_query_string = true
      path_and_qs               = "path_1"
    }
  }
  match_rules {
    name      = "complex_object_rule"
    start     = 2
    end       = 3
    match_url = "example2.com"
    matches {
      match_type     = "cookie"
      match_operator = "exists"
      case_sensitive = true
      negate         = true
      check_ips      = "CONNECTING_IP XFF_HEADERS"
      object_match_value {
        type                = "object"
        name                = "object name"
        name_case_sensitive = true
        name_has_wildcard   = false
        options {
          value                = ["cookie1=value1", "cookie2=value2"]
          value_has_wildcard   = false
          value_case_sensitive = true
          value_escaped        = false
        }
      }
    }

    forward_settings {
      origin_id                 = "origin_2"
      use_incoming_query_string = false
      path_and_qs               = "path_2"
    }
  }
  match_rules {
    name      = "complex_range_rule"
    start     = 3
    end       = 4
    match_url = "example3.com"
    matches {
      match_type     = "range"
      match_operator = "equals"
      case_sensitive = false
      negate         = false
      check_ips      = "CONNECTING_IP"
      object_match_value {
        type  = "range"
        value = [1, 50]
      }
    }

    forward_settings {
      origin_id                 = "origin_3"
      use_incoming_query_string = true
      path_and_qs               = "path_3"
    }
  }
}