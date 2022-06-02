provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_phased_release_match_rule" "test" {
  match_rules {
    name = "rule1"
    forward_settings {
      origin_id = "1234"
      percent   = 123
    }
  }
}