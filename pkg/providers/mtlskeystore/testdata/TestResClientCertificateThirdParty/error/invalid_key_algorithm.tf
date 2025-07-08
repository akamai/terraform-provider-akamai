provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_mtlskeystore_client_certificate_third_party" "test" {
  certificate_name    = "test-certificate"
  contract_id         = "ctr_12345"
  geography           = "CORE"
  group_id            = 1234
  key_algorithm       = "INVALID"
  notification_emails = ["jkowalski@akamai.com", "jsmith@akamai.com"]
  secure_network      = "STANDARD_TLS"
  subject             = "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/"
  versions = {
    "v1" = {},
    "v2" = {},
  }
}