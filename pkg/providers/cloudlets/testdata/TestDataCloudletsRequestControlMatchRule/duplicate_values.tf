provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_request_control_match_rule" "test" {
  match_rules {
    allow_deny = "allow"
    matches {
      match_type     = "header"
      match_value    = "this is value"
      match_operator = "equals"
      object_match_value {
        type  = "simple"
        value = ["no, this is value"]
      }
    }
  }
}