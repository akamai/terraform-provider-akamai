provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

resource "akamai_appsec_version_notes" "test" {
  config_id     = 43253
  version_notes = "Test Notes"
}

