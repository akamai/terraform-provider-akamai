provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_phased_release_match_rule" "test" {

  match_rules {
    name = "rule 2"
    matches {
      match_type     = "hostname"
      match_operator = "equals"
    }
    forward_settings {
      origin_id = "1234"
      percent   = 10
    }
  }
}