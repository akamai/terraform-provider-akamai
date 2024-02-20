provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_forward_rewrite_match_rule" "test" {}