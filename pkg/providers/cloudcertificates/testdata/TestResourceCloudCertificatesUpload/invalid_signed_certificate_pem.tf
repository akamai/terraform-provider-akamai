provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_upload" "test" {
  lineage_id = 500001
  algorithms = {
    RSA = {
      signed_certificate_pem = "not-a-valid-pem-certificate"
    }
  }
}
