provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_certificates" "test" {
  include_certificate_materials = true
  issuer                        = "Test Org1"
}