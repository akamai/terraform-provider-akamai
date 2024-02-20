provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_edge_redirector_match_rule" "test" {}