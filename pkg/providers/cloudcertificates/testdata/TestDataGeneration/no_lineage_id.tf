provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_generation" "test" {
  generation_id = 929
}
