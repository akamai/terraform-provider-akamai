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
      match_type     = "method"
      match_operator = "equals"
    }
    forward_settings {
      origin_id = "33"
    }
    disabled = false
  }
}