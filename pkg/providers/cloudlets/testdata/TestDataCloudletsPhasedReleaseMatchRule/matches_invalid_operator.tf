provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_phased_release_match_rule" "test" {
  match_rules {
    name = "rule1"
    matches {
      match_type     = "clientip"
      match_value    = "127.0.0.1"
      match_operator = "invalid"
      object_match_value {
        type  = "simple"
        value = ["fghi"]
      }
    }
    forward_settings {
      origin_id = "1234"
      percent   = 30
    }
  }
}