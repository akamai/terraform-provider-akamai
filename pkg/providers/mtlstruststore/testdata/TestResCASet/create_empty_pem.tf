provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_mtlstruststore_ca_set" "test" {
  name                = "set-1"
  description         = "Test CA Set for validation"
  allow_insecure_sha1 = false
  version_description = "Initial version for testing"

  certificates = [
    {
      certificate_pem = ""
      description     = "Test certificate"
    }
  ]
}
