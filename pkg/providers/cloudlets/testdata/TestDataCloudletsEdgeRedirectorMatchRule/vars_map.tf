provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_edge_redirector_match_rule" "test" {

  match_rules {
    end = 0
    match_url = null
    matches {
      case_sensitive = true
      match_operator = "equals"
      match_type = "hostname"
      match_value = "3333.dom"
      negate = false
    }
    matches {
      case_sensitive = false
      match_operator = "equals"
      match_type = "cookie"
      match_value = "cookie=cookievalue"
      negate = false
    }
    matches {
      case_sensitive = false
      match_operator = "equals"
      match_type = "extension"
      match_value = "txt"
      negate = false
    }
    name = "rul3"
    redirect_url = "/abc/sss"
    start = 0
    status_code = 307
    use_incoming_query_string = false
    use_relative_url = "copy_scheme_hostname"
  }
  match_rules {
    end = 0
    match_url = "ddd.aaa"
    name = "rule 2"
    redirect_url = "sss.com"
    start = 0
    status_code = 301
    use_incoming_query_string = true
    use_relative_url = "none"
  }
  match_rules {
    end = 0
    match_url = "abc.com"
    name = "r1"
    redirect_url = "/ddd"
    start = 0
    status_code = 301
    use_incoming_query_string = false
    use_relative_url = "copy_scheme_hostname"
  }


}
