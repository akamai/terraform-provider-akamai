provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudlets_request_control_match_rule" "test" {
  match_rules {
    name           = "basic_rule"
    allow_deny     = "allow"
    matches_always = true
    disabled       = true
  }
}