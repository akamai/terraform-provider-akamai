provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_upload" "test" {
  lineage_id = 500001
  algorithms = {
    RSA = {
      signed_certificate_pem = "-----BEGIN CERTIFICATE-----\nRSACERTMODIFIED\n-----END CERTIFICATE-----\n"
    }
    ECDSA = {
      signed_certificate_pem = "-----BEGIN CERTIFICATE-----\nECDSACERT\n-----END CERTIFICATE-----\n"
    }
  }
}
