provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_certificate" "test" {
  contract_id = "test_contract"
  group_id    = "123"
  key_size    = "2048"
  key_type    = "RSA"
  sans        = ["test.example.com"]
}
