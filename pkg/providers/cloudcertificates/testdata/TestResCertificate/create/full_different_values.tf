provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_certificate" "test" {
  base_name      = "test-name"
  contract_id    = "test_contract"
  group_id       = "123"
  key_size       = "P-256"
  key_type       = "ECDSA"
  secure_network = "ENHANCED_TLS"
  sans           = ["test.example.com", "test.example2.com"]
  subject = {
    state    = "CA"
    locality = "Test City"
  }
}
