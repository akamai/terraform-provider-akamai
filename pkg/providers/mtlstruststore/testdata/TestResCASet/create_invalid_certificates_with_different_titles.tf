provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_mtlstruststore_ca_set" "test_ca_set" {
  name                = "DXE-5282"
  description         = "Test CA Set for validation"
  allow_insecure_sha1 = false
  version_description = "Initial version for testing"
  certificates = [
    {
      certificate_pem = "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV123\n-----END CERTIFICATE-----"
      description     = "Incorrect PEM format first group"
    },
    {
      certificate_pem = "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV345\n-----END CERTIFICATE-----"
      description     = "Incorrect PEM format first group"
    },
    {
      certificate_pem = "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV567\n-----END CERTIFICATE-----"
      description     = "Incorrect PEM format second group"
    },
    {
      certificate_pem = "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV789\n-----END CERTIFICATE-----"
      description     = "Incorrect PEM format second group"
    }
  ]
}