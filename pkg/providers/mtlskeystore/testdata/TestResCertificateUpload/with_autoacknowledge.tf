provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
  client_certificate_id     = 12345
  version_number            = 1
  signed_certificate        = "certificate-data"
  trust_chain               = "trustchain-data"
  wait_for_deployment       = true
  auto_acknowledge_warnings = true
}