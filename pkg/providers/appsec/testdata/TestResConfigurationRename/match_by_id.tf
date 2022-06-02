provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_configuration_rename" "test" {
  name        = "Akamai Tools New"
  description = "TF Tools"
  config_id   = 432531
}

