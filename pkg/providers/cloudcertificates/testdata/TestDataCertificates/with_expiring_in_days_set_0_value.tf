provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_certificates" "test" {
  expiring_in_days = 0
}