provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_certificate" "test" {
  contract_id    = "test_contract"
  key_size       = "2048"
  key_type       = "RSA"
  secure_network = "ENHANCED_TLS"
  sans           = ["test.example.com"]
}
