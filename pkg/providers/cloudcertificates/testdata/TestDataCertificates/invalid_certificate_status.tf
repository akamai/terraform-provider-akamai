provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_certificates" "test" {
  certificate_status = ["INVALID_STATUS"]
}
