provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}


// This resource and locals are used to simulate reading external certificate from some other resource,
// effectively forcing unknown value for the certificate_pem field during the plan phase.
resource "random_integer" "max_idx" {
  min = 0
  max = 1
}

locals {
  certs = [
    "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n",
    "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n"
  ]
}

resource "akamai_mtlstruststore_ca_set" "test" {
  name                = "set-1"
  description         = "Test CA Set for validation"
  allow_insecure_sha1 = false
  version_description = "Initial version for testing"

  certificates = [
    {
      certificate_pem = local.certs[random_integer.max_idx.result]
      description     = "Test certificate"
    }
  ]
}
