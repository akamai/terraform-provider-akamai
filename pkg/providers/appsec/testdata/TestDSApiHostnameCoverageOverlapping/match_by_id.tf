provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_hostname_coverage_overlapping" "test" {
  config_id = 43253
  hostname  = "example.com"
}

