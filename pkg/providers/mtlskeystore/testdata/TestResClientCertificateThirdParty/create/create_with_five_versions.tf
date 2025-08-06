provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_mtlskeystore_client_certificate_third_party" "test" {
  certificate_name    = "test-certificate"
  contract_id         = "ctr_12345"
  geography           = "CORE"
  group_id            = 1234
  key_algorithm       = "RSA"
  notification_emails = ["jkowalski@akamai.com", "jsmith@akamai.com"]
  secure_network      = "STANDARD_TLS"
  versions = {
    "v1" = {},
    "v2" = {},
    "v3" = {},
    "v4" = {},
    "v5" = {},
  }
}