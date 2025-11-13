provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_certificates" "test" {
  certificate_name   = "test_certificate1"
  certificate_status = ["ACTIVE", "READY_FOR_USE"]
}