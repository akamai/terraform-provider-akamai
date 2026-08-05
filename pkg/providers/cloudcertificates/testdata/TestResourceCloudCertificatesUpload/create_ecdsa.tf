provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_upload" "test" {
  lineage_id = 500001
  algorithms = {
    ECDSA = {
      signed_certificate_pem = "-----BEGIN CERTIFICATE-----\nECDSACERT\n-----END CERTIFICATE-----\n"
    }
  }
}
