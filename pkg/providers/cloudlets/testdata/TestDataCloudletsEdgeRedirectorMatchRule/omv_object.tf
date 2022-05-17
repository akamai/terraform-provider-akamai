provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_edge_redirector_match_rule" "test" {
  match_rules {
    name                      = "rule2"
    start                     = 10
    end                       = 10000
    redirect_url              = "/abc/sss"
    status_code               = 307
    use_incoming_query_string = false
    use_relative_url          = "copy_scheme_hostname"
    matches {
      match_type     = "hostname"
      match_operator = "equals"
      object_match_value {
        type  = "simple"
        value = ["abc"]
      }
    }
    matches {
      case_sensitive = true
      match_type     = "cookie"
      match_operator = "equals"
      negate         = false
      object_match_value {
        type                = "object"
        name                = "abcde"
        name_case_sensitive = true
        name_has_wildcard   = false
        options {
          value                = ["asfas"]
          value_has_wildcard   = false
          value_case_sensitive = true
          value_escaped        = false
        }
      }
    }
  }
}
