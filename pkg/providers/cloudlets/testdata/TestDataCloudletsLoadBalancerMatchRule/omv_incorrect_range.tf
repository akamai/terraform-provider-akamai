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
      match_type  = "clientip"
      match_value = "127.0.0.1"
      object_match_value {
        type  = "range"
        value = ["abc", "asdf"]
      }
    }
    forward_settings {
      origin_id = "12"
    }
  }
}