provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_upload" "test" {
  lineage_id = 500001
  algorithms = {
    INVALID_KEY_TYPE = {
      signed_certificate_pem = "-----BEGIN CERTIFICATE-----\nRSACERT\n-----END CERTIFICATE-----\n"
    }
  }
}
