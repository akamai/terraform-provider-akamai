provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_forward_rewrite_match_rule" "test" {
  match_rules {
    name = "rule1"
    matches {
      match_type     = "clientip"
      match_operator = "invalid"
      object_match_value {
        type  = "simple"
        value = ["fghi"]
      }
    }
    forward_settings {}
  }
}
