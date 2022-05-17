provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_visitor_prioritization_match_rule" "test" {

  match_rules {
    name                 = "rule 2"
    start                = 0
    end                  = 0
    pass_through_percent = 50.5
    matches {
      match_type     = "hostname"
      match_operator = "equals"
      object_match_value {
        type  = "range"
        value = [1, 50]
      }
    }
  }
}

