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
        type                = "object"
        name                = "abcde"
        name_case_sensitive = true
        name_has_wildcard   = false
        options {
          value                = ["test"]
          value_has_wildcard   = true
          value_case_sensitive = true
          value_escaped        = true
        }
      }
    }
  }
}
