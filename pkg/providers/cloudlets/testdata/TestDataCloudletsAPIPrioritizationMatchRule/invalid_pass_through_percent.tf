provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_api_prioritization_match_rule" "test" {

  match_rules {
    pass_through_percent = -2
  }
}