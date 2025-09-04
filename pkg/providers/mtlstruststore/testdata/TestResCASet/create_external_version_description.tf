provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_test" "sample" {
  input = "Initial version for testing"
}

resource "akamai_mtlstruststore_ca_set" "test" {
  name                = "set-1"
  description         = "Test CA Set for validation"
  allow_insecure_sha1 = false
  version_description = resource.akamai_test.sample.output

  certificates = [
    {
      certificate_pem = <<EOT
-----BEGIN CERTIFICATE-----
MIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV
-----END CERTIFICATE-----
EOT
      description     = "Test certificate"
    }
  ]
}
