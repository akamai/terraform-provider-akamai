provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_edge_redirector_match_rule" "test" {

  match_rules {
    type = "erMatchRule"
    aka_rule_id = "a58392a7a43f19a3"
    end = 0
    id = 0
    location = "/cloudlets/api/v2/policies/276858/versions/6/rules/a58392a7a43f19a3"
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
    use_incoming_scheme_and_host = true
    use_relative_url = "copy_scheme_hostname"
  }
  match_rules {
    type = "erMatchRule"
    aka_rule_id = "e38515c6542d2ed8"
    end = 0
    id = 0
    location = "/cloudlets/api/v2/policies/276858/versions/6/rules/e38515c6542d2ed8"
    match_url = "ddd.aaa"
    name = "rule 2"
    redirect_url = "sss.com"
    start = 0
    status_code = 301
    use_incoming_query_string = true
    use_relative_url = "none"
  }
  match_rules {
    type = "erMatchRule"
    aka_rule_id = "e1969ed65202167f"
    end = 0
    id = 0
    location = "/cloudlets/api/v2/policies/276858/versions/6/rules/e1969ed65202167f"
    match_url = "abc.com"
    name = "r1"
    redirect_url = "/ddd"
    start = 0
    status_code = 301
    use_incoming_query_string = false
    use_incoming_scheme_and_host = true
    use_relative_url = "copy_scheme_hostname"
  }


}
