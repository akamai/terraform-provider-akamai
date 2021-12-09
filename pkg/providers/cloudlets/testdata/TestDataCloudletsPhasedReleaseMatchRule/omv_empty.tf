provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_phased_release_match_rule" "test" {
  match_rules {
    name      = "rule1"
    start     = 10
    end       = 10000
    match_url = "example.com"
    matches {
      match_type  = "clientip"
      match_value = "127.0.0.1"
    }
    forward_settings {
      origin_id = "31"
      percent   = 20
    }
  }
  match_rules {
    name = "rule2"
    forward_settings {
      origin_id = "1234"
      percent   = 10
    }
  }
}