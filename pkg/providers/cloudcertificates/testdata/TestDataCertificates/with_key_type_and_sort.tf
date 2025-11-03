provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_certificates" "test" {
  key_type = "RSA"
  sort     = "-createdDate"
}