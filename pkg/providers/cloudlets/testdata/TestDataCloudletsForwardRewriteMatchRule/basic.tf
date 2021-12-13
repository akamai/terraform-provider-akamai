provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_forward_rewrite_match_rule" "test" {
  match_rules {
    name = "rule1"
    forward_settings {}
    disabled = true
  }
}