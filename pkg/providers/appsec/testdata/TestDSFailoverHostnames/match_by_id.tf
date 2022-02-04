provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

data "akamai_appsec_failover_hostnames" "test" {
  config_id = 43253
}

