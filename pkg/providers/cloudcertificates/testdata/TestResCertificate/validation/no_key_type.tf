provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_certificate" "test" {
  contract_id    = "test_contract"
  group_id       = "123"
  key_size       = "2048"
  secure_network = "ENHANCED_TLS"
  sans           = ["test.example.com"]
}
