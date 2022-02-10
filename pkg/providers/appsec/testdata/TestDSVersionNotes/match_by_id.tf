provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

data "akamai_appsec_version_notes" "test" {
  config_id = 43253
}

