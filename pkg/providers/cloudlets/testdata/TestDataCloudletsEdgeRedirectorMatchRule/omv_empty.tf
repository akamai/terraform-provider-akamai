provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_edge_redirector_match_rule" "test" {
  match_rules {
    name = "rule1"
    start = 10
    end = 10000
    match_url = "example.com"
    redirect_url = "/abc/sss"
    status_code = 307
    use_relative_url = "copy_scheme_hostname"
    matches {
      match_type = "clientip"
      match_value = "127.0.0.1"
    }
  }
}