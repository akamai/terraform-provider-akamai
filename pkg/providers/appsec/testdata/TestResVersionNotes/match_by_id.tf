provider "akamai" {
  edgerc = "~/.edgerc"
}



data "akamai_appsec_version_notes" "test" {
    config_id = 43253
    version = 7
    notes = "Test Notes"
}


