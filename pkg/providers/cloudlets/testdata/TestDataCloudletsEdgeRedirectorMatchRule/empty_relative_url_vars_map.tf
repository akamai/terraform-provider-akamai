provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_edge_redirector_match_rule" "test" {

  match_rules {
    end                       = 0
    match_url                 = "ddd.aaa"
    name                      = "rule 2"
    redirect_url              = "sss.com"
    start                     = 0
    status_code               = 301
    use_incoming_query_string = true
  }
  match_rules {
    end                       = 0
    match_url                 = "abc.com"
    name                      = "rule 1"
    redirect_url              = "/ddd"
    start                     = 0
    status_code               = 301
    use_incoming_query_string = false
    use_relative_url          = ""
  }


}
