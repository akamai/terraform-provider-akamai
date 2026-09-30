provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_upload" "test" {
  algorithms = {
    RSA = {
      signed_certificate_pem = "-----BEGIN CERTIFICATE-----\nRSACERT\n-----END CERTIFICATE-----\n"
    }
  }
}
