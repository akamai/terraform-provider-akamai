provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_certificate" "test" {
  base_name      = "test-name"
  contract_id    = "test_contract"
  group_id       = "123"
  key_size       = "2048"
  key_type       = "RSA"
  secure_network = "ENHANCED_TLS"
  sans           = ["test.example2.com"]
  subject = {
    common_name  = "example..com"
    organization = "Test Org"
    country      = "US"
    state        = "CA"
    locality     = "Test City"
  }
}
