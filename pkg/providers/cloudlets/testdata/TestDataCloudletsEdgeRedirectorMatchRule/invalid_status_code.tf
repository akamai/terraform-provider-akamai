provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_edge_redirector_match_rule" "test" {

  match_rules {
    redirect_url     = "/ddd"
    status_code      = 111
    use_relative_url = "copy_scheme_hostname"
  }
}
