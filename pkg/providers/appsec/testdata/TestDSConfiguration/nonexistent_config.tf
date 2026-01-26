provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_configuration" "test" {
  name = "Nonexistent"
}

output "configsedge" {
  value = data.akamai_appsec_configuration.test.config_id
}

output "configsedgelatestversion" {
  value = data.akamai_appsec_configuration.test.latest_version
}

