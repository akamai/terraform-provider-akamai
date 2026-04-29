provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_certificate" "test" {
  contract_id    = "test_contract"
  group_id       = "123"
  key_size       = "2048"
  key_type       = "RSA"
  secure_network = "STANDARD_TLS"
  geo_class      = "RESERVED_GLOBAL"
  sans           = ["test.example.com"]
}
