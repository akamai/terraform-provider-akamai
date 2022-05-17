provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_api_prioritization_match_rule" "test" {
  match_rules {
    pass_through_percent = 0
    matches {
      match_type     = "clientip"
      match_operator = "equals"
      check_ips      = "invalid"
      object_match_value {
        type  = "simple"
        value = ["fghi"]
      }
    }
  }
}