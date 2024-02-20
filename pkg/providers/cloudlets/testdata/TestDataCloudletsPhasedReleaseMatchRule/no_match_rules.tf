provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_phased_release_match_rule" "test" {}