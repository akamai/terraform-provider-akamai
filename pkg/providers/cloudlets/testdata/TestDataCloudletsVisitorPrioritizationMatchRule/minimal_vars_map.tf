provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_visitor_prioritization_match_rule" "test" {

  match_rules {
    pass_through_percent = 0
  }
}

