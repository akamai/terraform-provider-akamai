provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_request_control_match_rule" "test" {
  match_rules {
    matches {
      match_type     = "clientip"
      match_operator = "equals"
      object_match_value {
        type  = "simple"
        value = ["a"]
      }
    }
  }
}