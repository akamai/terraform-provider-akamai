provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "terraform_data" "cert" {
  input = "-----BEGIN CERTIFICATE-----\nRSACERT\n-----END CERTIFICATE-----\n"
}

resource "akamai_cloudcertificates_upload" "test" {
  lineage_id = 500001
  algorithms = {
    RSA = {
      signed_certificate_pem = terraform_data.cert.output
    }
  }
}
