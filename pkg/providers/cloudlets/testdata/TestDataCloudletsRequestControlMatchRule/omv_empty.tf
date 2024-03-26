provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudlets_request_control_match_rule" "test" {
  match_rules {
    allow_deny = "allow"
    name       = "empty_object_rule"
    matches {
      match_type  = "clientip"
      match_value = "127.0.0.1"
    }
  }
}