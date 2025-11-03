provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_certificate" "testcert" {
  certificate_id = "12345"
}