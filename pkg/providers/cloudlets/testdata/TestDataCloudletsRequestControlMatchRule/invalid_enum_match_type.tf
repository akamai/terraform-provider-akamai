provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_request_control_match_rule" "test" {
  match_rules {
    allow_deny = "allow"
    matches {
      match_type     = "invalid_match_type"
      match_operator = "equals"
      object_match_value {
        type  = "simple"
        value = ["a"]
      }
    }
  }
}