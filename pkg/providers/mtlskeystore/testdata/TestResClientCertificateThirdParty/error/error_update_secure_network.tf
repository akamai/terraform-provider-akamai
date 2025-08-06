provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_mtlskeystore_client_certificate_third_party" "test" {
  certificate_name    = "test-certificate"
  contract_id         = "ctr_12345"
  geography           = "CORE"
  group_id            = 1234
  notification_emails = ["jkowalski@akamai.com", "jsmith@akamai.com"]
  secure_network      = "ENHANCED_TLS"
  versions = {
    "v1" = {},
  }
}