provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_visitor_prioritization_match_rule" "test" {

  match_rules {
    name                 = "rule 2"
    pass_through_percent = 50.5
    matches {
      match_type     = "hostname"
      match_operator = "equals"
    }
  }
}