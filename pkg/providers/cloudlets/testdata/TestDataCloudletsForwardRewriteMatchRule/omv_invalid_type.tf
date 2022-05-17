provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_forward_rewrite_match_rule" "test" {
  match_rules {
    name      = "rule1"
    start     = 10
    end       = 10000
    match_url = "example.com"
    matches {
      match_type     = "clientip"
      match_operator = "equals"
      object_match_value {
        type  = "invalid_type"
        value = ["fghi"]
      }
    }
    forward_settings {
      origin_id = "33"
    }
  }
}