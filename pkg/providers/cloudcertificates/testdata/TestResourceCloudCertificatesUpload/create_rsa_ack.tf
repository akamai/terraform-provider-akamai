provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_upload" "test" {
  lineage_id           = 500001
  acknowledge_warnings = true
  algorithms = {
    RSA = {
      signed_certificate_pem = "-----BEGIN CERTIFICATE-----\nRSACERT\n-----END CERTIFICATE-----\n"
    }
  }
}
