provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_request_control_match_rule" "test" {
  match_rules {
    allow_deny = "allow"
    matches {
      match_type     = "clientip"
      match_operator = "equals"
      check_ips      = "incorrect_check_ips"
      object_match_value {
        type  = "simple"
        value = ["a"]
      }
    }
  }
}