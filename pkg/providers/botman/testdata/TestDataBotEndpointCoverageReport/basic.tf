provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_botman_bot_endpoint_coverage_report" "test" {
}