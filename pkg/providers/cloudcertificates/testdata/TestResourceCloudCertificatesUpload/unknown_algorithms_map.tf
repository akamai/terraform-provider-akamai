provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "terraform_data" "algorithms" {
  input = {
    RSA = {
      signed_certificate_pem = "-----BEGIN CERTIFICATE-----\nRSACERT\n-----END CERTIFICATE-----\n"
    }
  }
}

resource "akamai_cloudcertificates_upload" "test" {
  lineage_id = 500001
  algorithms = terraform_data.algorithms.output
}
