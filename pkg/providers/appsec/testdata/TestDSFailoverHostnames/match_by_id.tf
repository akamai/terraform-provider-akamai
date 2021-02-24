provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_failover_hostnames" "test" {
  config_id = 43253
}


