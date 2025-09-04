provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_apr_transactional_endpoint" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
  operation_id       = "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"
}