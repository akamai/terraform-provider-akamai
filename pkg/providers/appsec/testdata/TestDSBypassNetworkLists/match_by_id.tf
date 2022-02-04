provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

data "akamai_appsec_bypass_network_lists" "test" {
  config_id = 43253
}

