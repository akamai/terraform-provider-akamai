provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

data "akamai_botman_bot_endpoint_coverage_report" "test" {
  operation_id = "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"
}