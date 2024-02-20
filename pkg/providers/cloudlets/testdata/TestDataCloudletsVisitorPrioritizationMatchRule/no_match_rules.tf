provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_visitor_prioritization_match_rule" "test" {}