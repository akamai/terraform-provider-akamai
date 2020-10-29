provider "akamai" {
  edgerc = "~/.edgerc"
}



data "akamai_appsec_rate_policies" "test" {
    config_id = 43253
    version = 7
}

