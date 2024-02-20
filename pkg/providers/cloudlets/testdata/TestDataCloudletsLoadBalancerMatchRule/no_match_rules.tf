provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_application_load_balancer_match_rule" "test" {}