provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_mtlskeystore_client_certificate_akamai" "test" {
  certificate_name    = "test-certificate-drift"
  contract_id         = "123456789"
  geography           = "CORE"
  group_id            = 987654321
  notification_emails = ["testemail3@example.com", "testemail4@example.com"]
  secure_network      = "STANDARD_TLS"
}