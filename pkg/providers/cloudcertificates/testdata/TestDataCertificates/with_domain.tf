provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_certificates" "test" {
  contract_id = "A-123"
  group_id    = "1234"
  domain      = "test2.example.com"
}