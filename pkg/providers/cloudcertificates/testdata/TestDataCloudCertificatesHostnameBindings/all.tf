provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_hostname_bindings" "test" {
  contract_id      = "K-0N7RAK71"
  group_id         = "123456"
  domain           = "example.com"
  network          = "PRODUCTION"
  expiring_in_days = 30
}
