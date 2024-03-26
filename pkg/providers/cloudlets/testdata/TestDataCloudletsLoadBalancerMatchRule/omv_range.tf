provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudlets_application_load_balancer_match_rule" "test" {
  match_rules {
    name      = "rule1"
    start     = 10
    end       = 10000
    match_url = "example.com"
    matches {
      match_type = "clientip"
      object_match_value {
        type  = "range"
        value = [1, 50]
      }
    }
    forward_settings {
      origin_id = "12"
    }
  }
}