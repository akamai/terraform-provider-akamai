provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_mtlstruststore_ca_set" "test" {
  name                = "set-1"
  description         = "Test CA Set for validation"
  allow_insecure_sha1 = false
  version_description = "Second version for testing"

  certificates = [
    {
      certificate_pem = <<EOT
-----BEGIN CERTIFICATE-----
UPDATED
-----END CERTIFICATE-----
EOT
      description     = "Test certificate"
    },
    {
      certificate_pem = <<EOT
-----BEGIN CERTIFICATE-----
FOO
-----END CERTIFICATE-----
EOT
      description     = "second cert"
    }
  ]
  timeouts {
    delete = "5m"
  }
}
