provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_certificate" "testcert" {
  certificate_id            = "12345"
  include_hostname_bindings = true
}