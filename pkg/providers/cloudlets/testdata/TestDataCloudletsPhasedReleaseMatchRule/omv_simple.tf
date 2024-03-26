provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudlets_phased_release_match_rule" "test" {
  match_rules {
    name           = "rule2"
    matches_always = true
    forward_settings {
      origin_id = "1234"
      percent   = 10
    }
  }
  match_rules {
    name      = "rule1"
    start     = 10
    end       = 10000
    match_url = "example.com"
    matches {
      match_type     = "method"
      match_operator = "equals"
      object_match_value {
        type  = "simple"
        value = ["GET"]
      }
    }
    forward_settings {
      origin_id = "33"
      percent   = 20
    }
    disabled = false
  }
}