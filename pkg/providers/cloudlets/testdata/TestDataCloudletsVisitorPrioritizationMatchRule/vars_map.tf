provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudlets_visitor_prioritization_match_rule" "test" {

  match_rules {
    end       = 0
    match_url = null
    matches {
      case_sensitive = true
      match_operator = "equals"
      match_type     = "hostname"
      match_value    = "3333.dom"
      negate         = false
    }
    matches {
      case_sensitive = false
      match_operator = "equals"
      match_type     = "cookie"
      match_value    = "cookie=cookievalue"
      negate         = false
    }
    matches {
      case_sensitive = false
      match_operator = "equals"
      match_type     = "extension"
      match_value    = "txt"
      negate         = false
    }
    name                 = "rul3"
    start                = 0
    pass_through_percent = -1
  }
  match_rules {
    end                  = 0
    match_url            = "ddd.aaa"
    name                 = "rule 2"
    start                = 0
    pass_through_percent = 100
    disabled             = true
  }
  match_rules {
    end                  = 0
    match_url            = "abc.com"
    name                 = "r1"
    start                = 0
    pass_through_percent = 50.55
    disabled             = false
  }
}

