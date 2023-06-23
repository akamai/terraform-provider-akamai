provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_siem_definitions" "test" {
  siem_definition_name = "SIEM Version 01"
}

