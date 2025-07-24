provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_mtlskeystore_client_certificate_akamai" "test" {
  certificate_name    = "test-certificate"
  contract_id         = "ctr_123456789"
  geography           = "CORE"
  group_id            = 987654321
  notification_emails = ["testemail1@example.com", "testemail2@example.com"]
  secure_network      = "STANDARD_TLS"
  key_algorithm       = "RSA"
}