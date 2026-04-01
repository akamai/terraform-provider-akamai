provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_configuration" "test" {
  name = "Akamai Tools"
}

output "configsedge" {
  value = data.akamai_appsec_configuration.test.config_id
}

output "configsedgelatestversion" {
  value = data.akamai_appsec_configuration.test.latest_version
}

output "configsedgeconfiglist" {
  value = data.akamai_appsec_configuration.test.output_text
}

output "host_names" {
  value = data.akamai_appsec_configuration.test.host_names
}

