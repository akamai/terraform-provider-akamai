provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_api_hostname_coverage_match_targets" "test" {
  config_id = 43253
    version = 7
  hostname = "example.com"
}
