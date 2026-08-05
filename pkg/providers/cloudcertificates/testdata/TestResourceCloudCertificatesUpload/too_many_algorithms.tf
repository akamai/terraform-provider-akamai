provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_upload" "test" {
  lineage_id = 500001
  algorithms = {
    RSA = {
      signed_certificate_pem = "-----BEGIN CERTIFICATE-----\nRSACERT1-----END CERTIFICATE-----\n"
    }
    ECDSA = {
      signed_certificate_pem = "-----BEGIN CERTIFICATE-----\nECDSACERT-----END CERTIFICATE-----\n"
    }
    DSA = {
      signed_certificate_pem = "-----BEGIN CERTIFICATE-----\nDSACERT-----END CERTIFICATE-----\n"
    }
  }
}
