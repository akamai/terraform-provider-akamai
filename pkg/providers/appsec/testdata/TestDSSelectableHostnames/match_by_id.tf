provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_selectable_hostnames" "test" {
  config_id = 43253
}

output "selectablehostnames" {
  value = data.akamai_appsec_selectable_hostnames.test.hostnames
}

