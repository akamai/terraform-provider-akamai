provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}
locals {
  certs = [
    "certificate-data",
    "certificate-data"
  ]
}

resource "random_integer" "rand" {
  min = 0
  max = 1
}

resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
  client_certificate_id = 12345
  version_number        = 1
  signed_certificate    = local.certs[random_integer.rand.result]
  trust_chain           = "trustchain-data"
  wait_for_deployment   = true
}